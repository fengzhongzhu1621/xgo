package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"sync"
	"time"
)

func examplePool() {
	ch := make(chan *Email, 10)
	p, err := NewPool(
		"smtp.example.com:25",
		4,
		smtp.PlainAuth("", "my@example.com", "123456", "smtp.example.com"),
	)

	if err != nil {
		log.Fatal("failed to create pool:", err)
	}

	var wg sync.WaitGroup
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			for e := range ch {
				err := p.Send(e, 10*time.Second)
				if err != nil {
					fmt.Fprintf(os.Stderr, "email:%v sent error:%v\n", e, err)
				}
			}
		}()
	}

	for i := 0; i < 10; i++ {
		e := NewEmail()
		e.From = "dj <my@example.com>"
		e.To = []string{"you@example.com"}
		e.Subject = "hello world"
		e.Text = []byte(fmt.Sprintf("email %d", i+1))
		ch <- e
	}

	close(ch)
	wg.Wait()
}
