package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	tickers := []string{"AAPL", "GOOG", "INFY"}
	ch := make(chan string)

	for _, t := range tickers {
		sym := t
		go func() {
			price := 200 + rand.Float64()
			for {
				price += (rand.Float64() - 0.5) * 2
				ch <- fmt.Sprintf("[%s] %s: %.2f", time.Now().Format("15:04:05"), sym, price)
				time.Sleep(1 * time.Second)
			}
		}()
	}

	go func() {
		time.Sleep(10 * time.Second)
		close(ch)
	}()

	for msg := range ch {
		fmt.Println(msg)
	}
}
