package main

import (
	"os"
	"sync"
	"time"

	"math/rand"

	"github.com/abates/log"
)

var wg sync.WaitGroup

func testLog(l log.Logger, num int) {
	wg.Add(1)
	go func(num int) {
		ll := l.Printf("LL%d", num)
		for i := 0; i < rand.Intn(10); i++ {
			ll.Printf("LL%d: message %d", num, i)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		wg.Done()
	}(num)
}

func main() {
	l := log.New(os.Stdout)
	l.Printf("One")
	testLog(l, 1)

	l.Printf("Two")
	testLog(l, 2)

	l.Printf("Three")
	//go testLog(3)
	time.Sleep(time.Second)
	l.Printf("Four")
	time.Sleep(time.Second)
	l.Printf("Five")
	l.Printf("Six")
	time.Sleep(time.Second)
	wg.Wait()
}
