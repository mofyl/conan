package leaky

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func TestLeaky(t *testing.T) {

	l := NewLeaky(nil)

	l.SetEnable(OPEN_LEAKY)

	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < 10; i++ {
		go func() {

			for {
				select {
				case <-ctx.Done():
					return
				default:
					if l.TryGetToken() {
						fmt.Println("Get Token")
					}
				}
			}

		}()
	}

	time.Sleep(5 * time.Second)

	l.Close()
	cancel()
}

func TestInterfacePoint(t *testing.T) {

	num1 := 2
	num2 := 3

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		initTime := time.Now()

		var v interface{}
		var tmp int
		for i := 0; i < 1000; i++ {
			v = num1
			tmp = v.(int)
		}
		fmt.Printf(" interface %d\n", time.Since(initTime).Nanoseconds())
		wg.Done()
		fmt.Printf("interface Tmp %d\n", tmp)
	}()

	wg.Add(1)
	go func() {
		initTime := time.Now()

		var v unsafe.Pointer
		var tmp *int
		for i := 0; i < 1000; i++ {
			v = unsafe.Pointer(&num2)
			tmp = (*int)(v)
		}
		fmt.Printf(" unsafePoint %d\n", time.Since(initTime).Nanoseconds())
		wg.Done()
		fmt.Printf("unsafePoint Tmp %d\n", *tmp)
	}()

	wg.Wait()
}

func BenchmarkLeaky(b *testing.B) {

	l := NewLeaky(nil)

	l.SetEnable(OPEN_LEAKY)

	for i := 0; i < b.N; i++ {
		l.TryGetToken()
	}

	l.Close()

}
