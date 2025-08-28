package metrics

import (
	"context"

	"github.com/fengzhongzhu1621/xgo/trpc/plugins/slime/view"
)

// Noop empty implementation.
type Noop struct{}

// Report does nothing.
func (Noop) Report(context.Context, view.IStat) {}

// Inc implements Emitter and does nothing.
func (Noop) Inc(name string, cnt int, tagPairs ...string) {}

// Observe implements Emitter and does nothing.
func (Noop) Observe(name string, v float64, tagPairs ...string) {}
