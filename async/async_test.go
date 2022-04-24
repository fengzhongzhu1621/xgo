package async

import (
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func fib(p, c int) (int, int) {
	return c, p + c
}

func TestAsync(t *testing.T) {
	var (
		e        error
		multiRes Results
		res      []interface{}
		keyRes   = "two"
	)

	fmt.Println("Testing `Waterfall`")

	res, e = Waterfall(Tasks{
		fib, fib, fib,
		func(p, c int) int {
			return c
		},
	}, 0, 1)

	if e != nil {
		t.Errorf("Error executing a Waterfall (%s)", e.Error())
	}

	fmt.Println("Waterfall result :", res[0].(int))

	fmt.Printf("\nTesting `Parallel` with `runtime.GOMAXPROCS(2)`\n")

	runtime.GOMAXPROCS(2)

	multiRes, e = Parallel(MapTasks{
		"one": func() error {
			for i := 'a'; i < 'a'+26; i++ {
				fmt.Printf("%c ", i)
			}
			return nil // fmt.Errorf("Error in one function")
		},
		"two": func() (int, string, error) {
			time.Sleep(2 * time.Microsecond)
			for i := 0; i < 27; i++ {
				fmt.Printf("%d ", i)
			}

			return 2, "test", nil
		},
		"three": func() int {
			for i := 'z'; i >= 'a'; i-- {
				fmt.Printf("%c ", i)
			}
			return 3
		},
	})

	if e != nil {
		t.Errorf("Error executing a Parallel (%s)", e.Error())
	}

	fmt.Printf("Parallel Result key %s: %+v\n", keyRes, multiRes.Key(keyRes))

	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Printf("\nTesting `Concurrent`\n")

	multiRes, e = Concurrent(Tasks{
		func() int {
			for i := 'a'; i < 'a'+26; i++ {
				fmt.Printf("%c ", i)
			}

			return 1
		},
		func() bool {
			time.Sleep(3 * time.Microsecond)
			for i := 0; i < 27; i++ {
				fmt.Printf("%d ", i)
			}

			return false
		},
		func() {
			for i := 'z'; i >= 'a'; i-- {
				fmt.Printf("%c ", i)
			}
		},
	})

	fmt.Println("Concurrent Result Index: 1", multiRes.Index(1))

	if e != nil {
		t.Errorf("Error executing a Concurrent (%s)", e.Error())
	}
}

func TestAsyncError(t *testing.T) {
	fmt.Printf("\nTesting `Waterfall` with error\n")

	res, e := Waterfall(Tasks{
		func() (int, error) {
			return 1, nil
		},
		func(n int) error {
			fmt.Printf("if %d > 0 then error\n", n)
			if n > 0 {
				return errors.New("Error on second function")
			}
			return nil
		},
		func() error {
			fmt.Println("Function never reached")
			return nil
		},
	})

	if e != nil {
		fmt.Printf("Error executing a Waterfall (%q)\n", e)
	}

	// should be empty
	fmt.Printf("Waterfall result with error should be empty: %+v\n", res)
}
