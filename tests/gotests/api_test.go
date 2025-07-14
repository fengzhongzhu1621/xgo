//go:generate mockgen -source=$GOFILE -destination=./mock/$GOFILE -package=mock
//go:generate gotests -w -all $GOFILE
package gotests

import (
	"reflect"
	"testing"
)

func Test_newConfigImpl(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *configImpl
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newConfigImpl(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newConfigImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}
