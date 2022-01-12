package main

import (
	"os"
	"sync"
	"time"

	"math/rand"

	"github.com/abates/log"
)

var logger log.Logger

var wg sync.WaitGroup

func testLog(num int) {
	wg.Add(1)
	ll := logger.Printf("LL%d", num)
	formatter := ll.Formatter()
	ll.SetFormatter(log.FormatChain(log.Indent(5), ll.Formatter()))
	go func(num int) {
		for i := 0; i < rand.Intn(12); i++ {
			ll.Printf("LL%d: message %d", num, i)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		ll.SetFormatter(formatter)
		if num%2 == 0 {
			ll.Setf(log.SuccessMessage, "LL%d: Succeeded", num)
		} else {
			ll.Setf(log.FailMessage, "LL%d: Failed", num)
		}
		wg.Done()
	}(num)
}

func main() {
	logger = log.New(os.Stderr, log.Annotate(), log.Colorize())

	logger.Printf("One")
	testLog(1)

	logger.Printf("Two")
	testLog(2)

	logger.Printf("Three")
	testLog(3)
	time.Sleep(time.Second)
	logger.Printf("Four")
	time.Sleep(time.Second)
	logger.Printf("Five")
	logger.Printf("Six")
	time.Sleep(time.Second)
	wg.Wait()
}
