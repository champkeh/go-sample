package main

import (
	"time"
	"fmt"
)

// 使用一个单独的 go routine 来执行这个函数
// delay参数用于控制转动的快慢
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

	// 暂停10秒，模拟程序的执行
	time.Sleep(10 * time.Second)

	fmt.Print("\r")
	fmt.Println("program end.")
}
