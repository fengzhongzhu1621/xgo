//go:generate mockgen -source=$GOFILE -destination=./mock/$GOFILE -package=mock
//go:generate gotests -w -all $GOFILE
package gotests

import "sync/atomic"

type configImpl struct {
	name  string
	value atomic.Value
}

func newConfigImpl(name string) *configImpl {
	return &configImpl{
		name: name,
	}
}
