package snowflake

import (
	"testing"
	"os"
	"sync"
	"time"
	"fmt"
)

// integration testing
func TestWorker_NextId(t *testing.T) {
	worker, err := NewWorker(0)
	if err != nil {
		os.Exit(1)
	}
	count, m := 0, make(map[int64]struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan struct{}, 1)
	go func() {
	FOR:
		for {
			id, err := worker.NextId()
			//fmt.Println(id)
			if err != nil {
				os.Exit(1)
			}
			m[id], count = struct{}{}, count + 1

			select {
			case <-ch:
				break FOR
			default:
				continue
			}
		}
		wg.Done()
	}()
	go func() {
		time.Sleep(time.Millisecond)
		ch <- struct{}{}
	}()
	wg.Wait()
	fmt.Println("count:", count)
	fmt.Println("map size:", len(m))
}