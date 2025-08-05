package degrade

import (
	"math/rand"
	"time"

	"trpc.group/trpc-go/trpc-go/plugin"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	plugin.Register(pluginName, &Degrade{})
}
