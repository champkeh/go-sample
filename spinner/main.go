package main

import (
	"time"
	"fmt"
)

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func main() {
	go spinner(200 * time.Millisecond)
	time.Sleep(10 * time.Second)
	fmt.Print("\r")
	fmt.Println("program end.")
}
