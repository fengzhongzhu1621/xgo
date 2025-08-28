package registerdiscover

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

func NewSvcManagerClient(ctx context.Context, svcManagerAddr string) (*ZkClient, error) {
	var err error
	for retry := 0; retry < maxRetry; retry++ {
		client := NewZkClient(svcManagerAddr, 40*time.Second)
		if err = client.Start(); err != nil {
			log.Errorf("connect regdiscv [%s] failed: %v", svcManagerAddr, err)
			time.Sleep(time.Second * 2)
			continue
		}

		if err = client.Ping(); err != nil {
			log.Errorf("connect regdiscv [%s] failed: %v", svcManagerAddr, err)
			time.Sleep(time.Second * 2)
			continue
		}

		return client, nil
	}

	return nil, err
}
