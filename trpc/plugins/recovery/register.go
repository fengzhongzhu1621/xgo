package recovery

import "trpc.group/trpc-go/trpc-go/filter"

func init() {
	filter.Register("recovery", ServerFilter(), nil)
}
