package backbone

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/fengzhongzhu1621/xgo/db/zookeeper/zkclient"
	tlsutil "github.com/fengzhongzhu1621/xgo/ginx/utils/tls"
	"github.com/fengzhongzhu1621/xgo/monitor/opentelemetry"
	"github.com/fengzhongzhu1621/xgo/network/ssl"
	log "github.com/sirupsen/logrus"
)

// NewClient create a new http client
func NewClient(c *ssl.TLSClientConfig, conf ...tlsutil.ExtraClientConfig) (*http.Client, error) {
	tlsConf := new(tls.Config)
	if c != nil && len(c.CAFile) != 0 && len(c.CertFile) != 0 && len(c.KeyFile) != 0 {
		var err error
		tlsConf, err = ssl.ClientTslConfVerity(c.CAFile, c.CertFile, c.KeyFile, c.Password)
		if err != nil {
			return nil, err
		}
	}

	if c != nil {
		tlsConf.InsecureSkipVerify = c.InsecureSkipVerify
	}

	// set api request timeout to 25s, so that we can stop the long request like searching all hosts
	responseHeaderTimeout := 25 * time.Second
	if len(conf) > 0 {
		if timeout := conf[0].ResponseHeaderTimeout; timeout != 0 {
			responseHeaderTimeout = timeout
		}
	}
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig:     tlsConf,
		Dial: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		MaxIdleConnsPerHost:   100,
		ResponseHeaderTimeout: responseHeaderTimeout,
	}

	client := &http.Client{
		Transport: transport,
	}

	opentelemetry.WrapperTraceClient(client)

	return client, nil
}

func ListenAndServe(c Server, svcDisc ServiceRegisterInterface, cancel context.CancelFunc) error {
	handler := c.Handler
	if c.PProfEnabled {
		rootMux := http.NewServeMux()
		rootMux.HandleFunc("/", c.Handler.ServeHTTP)
		rootMux.Handle("/debug/", http.DefaultServeMux)
		handler = rootMux
	}
	server := &http.Server{
		Addr:    net.JoinHostPort(c.ListenAddr, strconv.FormatUint(uint64(c.ListenPort), 10)),
		Handler: handler,
	}
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM)
	go func() {
		for {
			select {
			case sig := <-exit:
				log.Infof("receive signal %v, begin to shutdown", sig)
				svcDisc.Cancel()
				if err := svcDisc.ClearRegisterPath(); err != nil && err != zkclient.ErrNoNode {
					break
				}
				time.Sleep(time.Second * 5)
				server.SetKeepAlivesEnabled(false)
				err := server.Shutdown(context.Background())
				if err != nil {
					log.Fatalf("Could not gracefully shutdown the server: %v \n", err)
				}
				log.Info("server shutdown done")
				cancel()
				return
			}
		}
	}()

	if !tlsutil.IsTLS(c.TLS) {
		log.Infof("start insecure server on %s:%d", c.ListenAddr, c.ListenPort)
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen and serve failed, err: %v", err)
			}
		}()
		return nil
	}

	tlsC, err := ssl.ClientTslConfVerity(c.TLS.CAFile,
		c.TLS.CertFile,
		c.TLS.KeyFile,
		c.TLS.Password)
	if err != nil {
		return fmt.Errorf("generate tls config failed. err: %v", err)
	}

	server.TLSConfig = tlsC
	log.Infof("start secure server on %s:%d", c.ListenAddr, c.ListenPort)
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve failed, err: %v", err)
		}
	}()

	return nil
}
