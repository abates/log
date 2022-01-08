package main

import (
	"os"
	"sync"
	"time"

	"math/rand"

	"github.com/abates/log"
)

var logger *log.ProgressLogger

var wg sync.WaitGroup

func testLog(num int) {
	ll := logger.LineLogger()
	for i := 0; i < 25; i++ {
		ll.Printf("LL%d: message %d", num, i)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
	if num%2 == 0 {
		ll.Successf("LL%d Succeeded!", num)
	} else {
		ll.Failf("LL%d Failed!", num)
	}
	wg.Done()
}

func main() {
	logger = log.New(os.Stderr, log.SuccessFormatter(), log.PrefixFormatter("", log.LstdFlags), log.ColorFormatter())

	wg.Add(3)
	logger.Printf("One")
	go testLog(1)

	logger.Printf("Two")
	go testLog(2)

	logger.Printf("Three")
	go testLog(3)
	time.Sleep(time.Second)
	logger.Printf("Four")
	time.Sleep(time.Second)
	logger.Printf("Five")
	logger.Printf("Six")
	time.Sleep(time.Second)
	wg.Wait()
}
