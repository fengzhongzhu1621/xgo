package main

import (
	"context"

	"github.com/silenceper/log"
)

func handleLocalTimer(_ context.Context) error {
	log.Info("do local timer processing...")
	return nil
}
