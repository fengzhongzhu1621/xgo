package uuid

import (
	"encoding/hex"

	"github.com/pborman/uuid"
	"github.com/rogpeppe/fastuuid"
)

var generator *fastuuid.Generator

func initUUID() (err error) {
	generator, err = fastuuid.NewGenerator()
	return
}

func NextID() string {
	id := generator.Next()
	return hex.EncodeToString(id[:])
}

func Uuid() string {
	return uuid.New()
}
