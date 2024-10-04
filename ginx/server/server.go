package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fengzhongzhu1621/xgo/config"
	"github.com/fengzhongzhu1621/xgo/logging"
	log "github.com/sirupsen/logrus"
)

const (
	defaultGraceTimeout = 30 * time.Second

	defaultIdleTimeout  = 180 * time.Second
	defaultReadTimeout  = 60 * time.Second
	defaultWriteTimeout = 60 * time.Second
)

// Server 提供了一个基本的HTTP服务器启动、停止和等待关闭的框架，允许服务器在接收到停止信号时优雅地处理现有的连接，并在超时后强制关闭。
type Server struct {
	// 服务器地址
	addr string
	// 实际的HTTP服务器对象
	server *http.Server
	// 用于优雅地停止服务器的通道
	stopChan chan struct{}
	// 服务器的配置信息
	config *config.Config
}

// NewServer 提供的配置信息创建并初始化一个HTTP服务器，并注册相应的路由规则
func NewServer(cfg *config.Config) *Server {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	log.Infof("the server addr: %s", addr)

	// 设置超时时间
	readTimeout := defaultReadTimeout
	if cfg.Server.ReadTimeout > 0 {
		readTimeout = time.Duration(cfg.Server.ReadTimeout) * time.Second
	}
	writeTimeout := defaultWriteTimeout
	if cfg.Server.WriteTimeout > 0 {
		writeTimeout = time.Duration(cfg.Server.WriteTimeout) * time.Second
	}
	idleTimeout := defaultIdleTimeout
	if cfg.Server.IdleTimeout > 0 {
		idleTimeout = time.Duration(cfg.Server.IdleTimeout) * time.Second
	}

	// 注册路由
	router := NewRouter(cfg)

	// 定义HTTP服务器
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	// 返回自定义服务器结构体
	return &Server{
		addr:     addr,
		server:   server,
		stopChan: make(chan struct{}, 1),
		config:   cfg,
	}
}

// Run 启动HTTP服务器，并且设置了一个协程来监听上下文的取消信号（如SIGINT或SIGTERM），以便优雅地停止服务器
func (s *Server) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		// 当接收到停止信号时，调用Stop()方法优雅地停止服务器
		// SIGINT or SIGTERM
		log.Info("I have to go...")
		log.Info("Stopping server gracefully")
		s.Stop()
	}()

	// 启动服务
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			log.Info(err.Error())
		}
	}()

	// 等待服务关闭，阻塞当前协程，直到服务器停止
	s.Wait()
	log.Info("Shutting down")
}

// Stop 优雅地停止HTTP服务器，它设置了超时时间，并在超时后强制关闭服务器
func (s *Server) Stop() {
	defer log.Info("Server stopped")

	// 设置优雅关闭的超时时间
	graceTimeout := defaultGraceTimeout
	if s.config.Server.GraceTimeout > 0 {
		graceTimeout = time.Duration(s.config.Server.GraceTimeout) * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), graceTimeout)
	defer cancel()
	log.Infof("Waiting %s seconds before killing connections...", graceTimeout)

	// 禁用keep-alive连接
	s.server.SetKeepAlivesEnabled(false)
	if err := s.server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Wait is over due to error")
		s.server.Close()
	}

	// 刷新缓存的日志
	logging.GetApiLogger().Sync()
	logging.GetWebLogger().Sync()

	// 发送信号，表示服务器已经停止
	s.stopChan <- struct{}{}
}

// Wait 阻塞当前协程，直到服务器停止
func (s *Server) Wait() {
	<-s.stopChan
}
