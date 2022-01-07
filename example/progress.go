package main

import (
	"sync"
	"time"

	"math/rand"

	"github.com/abates/log"
)

var wg sync.WaitGroup

func testLog(num int) {
	ll := log.LineLogger()
	for i := 0; i < 25; i++ {
		ll.Printf("LL%d: message %d", num, i)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
	wg.Done()
}

func main() {
	wg.Add(3)
	log.Printf("One")
	go testLog(1)

	log.Printf("Two")
	go testLog(2)

	log.Printf("Three")
	go testLog(3)
	time.Sleep(time.Second)
	log.Printf("Four")
	time.Sleep(time.Second)
	log.Printf("Five")
	log.Printf("Six")
	time.Sleep(time.Second)
	wg.Wait()
	log.Finish()
}
