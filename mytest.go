package main

import (
	"fmt"
	"sync"
	"github.com/panjf2000/ants"
	"time"
)
var sum int32
func demoFunc() error {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
	return nil
}

func main111() {
	defer ants.Release()

	p, _ := ants.NewPool(3)
	runTimes := 1000

	// use the common pool
	var wg sync.WaitGroup
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		fmt.Println("running is :",p.Running(),"cap is:",p.Cap(),"free is ",p.Free())

		p.Submit(func() error {
			demoFunc()
			wg.Done()
			return nil
		})
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", ants.Running())
	fmt.Printf("finish all tasks.\n")
}
