package logger

import (
	"context"
	"fmt"
	"log"
	"os"
)

type SimpleLogging interface {
	Printf(ctx context.Context, format string, v ...interface{})
}

type simpleLogger struct {
	log *log.Logger
}

func (l *simpleLogger) Printf(ctx context.Context, format string, v ...interface{}) {
	_ = l.log.Output(2, fmt.Sprintf(format, v...))
}

// Logger calls Output to print to the stderr.
// Arguments are handled in the manner of fmt.Print.
var SimpleLogger SimpleLogging = &simpleLogger{
	log: log.New(os.Stderr, "redis: ", log.LstdFlags|log.Lshortfile),
}
