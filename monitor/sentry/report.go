package report

import (
	"time"

	"github.com/getsentry/sentry-go"
)

var sentryOn bool

// InitErrorReport init the sentryEnabled var
func InitErrorReport(sentryEnabled bool) {
	sentryOn = sentryEnabled
}

// ReportToSentry is a shortcut to build and send an event to sentry
func ReportToSentry(message string, extra map[string]interface{}) {
	event := sentry.NewEvent()
	event.Message = message
	event.Level = "error"
	event.Timestamp = time.Now()
	event.Extra = extra
	if sentryOn {
		sentry.CaptureEvent(event)
	}
}
