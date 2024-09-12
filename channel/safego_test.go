package channel

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSafeGo(t *testing.T) {
	wg := new(sync.WaitGroup)
	type args struct {
		do_func func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "正常函数",
			args: args{do_func: func() {
				defer wg.Done()
				time.Sleep(500 * time.Millisecond)
			}},
		},
		{
			name: "异常函数",
			args: args{do_func: func() {
				defer wg.Done()
				time.Sleep(500 * time.Millisecond)
				panic("error")
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				assert.Nil(t, recover())
			}()
			wg.Add(1)
			SafeGo(tt.args.do_func)
		})
	}
	wg.Wait()
}

func TestSafeGoWait(t *testing.T) {
	type args struct {
		fs []func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "正常函数+异常函数",
			args: args{fs: []func(){
				func() {
					time.Sleep(500 * time.Millisecond)
				},
				func() {
					panic("error")
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				assert.Nil(t, recover())
			}()
			SafeGoWait(tt.args.fs...)
		})
	}
}
