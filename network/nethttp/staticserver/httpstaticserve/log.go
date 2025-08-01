package httpstaticserver

import (
	"log"

	"github.com/codeskyblue/go-accesslog"
)

type httpLogger struct{}

func (l httpLogger) Log(record accesslog.LogRecord) {
	log.Printf("%s - %s %d %s", record.Ip, record.Method, record.Status, record.Uri)
}
