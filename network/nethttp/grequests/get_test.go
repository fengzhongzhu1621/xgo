package grequests

import (
	"context"
	"log"
	"testing"

	"github.com/levigross/grequests/v2"
)

func TestGet(*testing.T) {
	resp, err := grequests.Get(context.Background(), "https://httpbin.org/get",
		grequests.UserAgent("MyAgent"))
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]any
	if err := resp.JSON(&data); err != nil {
		log.Fatal(err)
	}
	log.Println(data)
}
