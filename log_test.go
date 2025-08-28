package xgo

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

type User struct {
	Name string
	Age  int
}

func TestLog(t *testing.T) {
	u := User{
		Name: "dj",
		Age:  18,
	}
	log.SetPrefix("Login: ")
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	log.Printf("%s login, age:%d", u.Name, u.Age)

	buf := &bytes.Buffer{}
	logger := log.New(buf, "", log.Lshortfile|log.LstdFlags)
	logger.Printf("%s login, age:%d", u.Name, u.Age)
	fmt.Print(buf.String())

	writer1 := &bytes.Buffer{}
	writer2 := os.Stdout
	writer3, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0o755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}
	logger = log.New(io.MultiWriter(writer1, writer2, writer3), "", log.Lshortfile|log.LstdFlags)
	logger.Printf("%s login, age:%d", u.Name, u.Age)
}
