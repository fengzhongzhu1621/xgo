package log

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
)

type AppHook struct {
	AppName string
}

func (h *AppHook) Levels() []Level {
	return AllLevels
}

func (h *AppHook) Fire(entry *Entry) error {
	entry.Data["app"] = h.AppName
	return nil
}

func TestLogrusInfo(t *testing.T) {
	LogrusSetLevel(TraceLevel)
	LogrusSetReportCaller(true)
	LogrusInfo("info msg")
	LogrusWithFields(Fields{
		"name": "dj",
		"age":  18,
	}).Info("info msg")

	writer1 := &bytes.Buffer{}
	writer2 := os.Stdout
	writer3, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	LogrusSetOutput(io.MultiWriter(writer1, writer2, writer3))
	LogrusInfo("info msg")

	LogrusSetFormatter(&JSONFormatter{})
	LogrusInfo("info msg")

	h := &AppHook{AppName: "awesome-web"}
	LogrusAddHook(h)
	LogrusInfo("info msg")
}
