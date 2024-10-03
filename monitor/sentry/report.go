package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
)

var sentryOn bool

// SetSentryCaptureSwitch init the sentryEnabled var
func SetSentryCapture(sentryOn bool) {
	sentryOn = sentryOn
}

func EnableSentry() {
	sentryOn = true
}

func DisableSentry() {
	sentryOn = false
}

// CaptureEvent is a shortcut to build and send an event to sentry
func CaptureEvent(message string, extra map[string]interface{}) {
	event := sentry.NewEvent()
	event.Message = message
	event.Level = "error"
	event.Timestamp = time.Now()
	event.Extra = extra
	if sentryOn {
		sentry.CaptureEvent(event)
	}
}
