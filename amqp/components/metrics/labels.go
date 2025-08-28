package metrics

import (
	"context"

	"github.com/fengzhongzhu1621/xgo/amqp/router"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelKeyHandlerName    = "handler_name"
	labelKeyPublisherName  = "publisher_name"
	labelKeySubscriberName = "subscriber_name"
	labelSuccess           = "success"
	labelAcked             = "acked"

	labelValueNoHandler = "<no handler>"
)

var labelGetters = map[string]func(context.Context) string{
	labelKeyHandlerName:    router.HandlerNameFromCtx,
	labelKeyPublisherName:  router.PublisherNameFromCtx,
	labelKeySubscriberName: router.SubscriberNameFromCtx,
}

func labelsFromCtx(ctx context.Context, labels ...string) prometheus.Labels {
	ctxLabels := map[string]string{}

	for _, l := range labels {
		k := l
		ctxLabels[l] = ""

		getter, ok := labelGetters[k]
		if !ok {
			continue
		}

		v := getter(ctx)
		if v != "" {
			ctxLabels[l] = v
		}
	}

	return ctxLabels
}
