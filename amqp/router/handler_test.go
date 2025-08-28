package router

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleClose(t *testing.T) {
	num_chan := make(chan int)
	routersCloseCh := make(chan struct{})

	go func() {
		select {
		case <-routersCloseCh:
			fmt.Println("num_chan <- 1")
			num_chan <- 1
		}
	}()

	go func() {
		select {
		case <-routersCloseCh:
			fmt.Println("num_chan <- 2")
			num_chan <- 2
		}
	}()

	close(routersCloseCh)

	num1 := <-num_chan
	num2 := <-num_chan

	assert.Equal(t, 3, num1+num2)
}
