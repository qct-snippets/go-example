package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// 情况一：channel 的发送次数等于接收次数
// channel 的发送次数等于接收次数时，发送者 go routine 和接收者 go routine 分别都会在发送或接收结束时结束各自的 go routine。
// 而代码中的 ich 会由于没有代码使用被垃圾收集器回收。
// 因此这种情况下，不关闭 channel，没有任何副作用。
func TestIsCloseChannelNecessary_on_equal(t *testing.T) {
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())
	ich := make(chan int)

	// sender
	go func() {
		for i := 0; i < 3; i++ {
			ich <- i
		}
	}()

	// receiver
	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(<-ich)
		}
	}()

	time.Sleep(time.Second)
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())

	// Output:
	// NumGoroutine: 2
	// 0
	// 1
	// 2
	// NumGoroutine: 2
}

// 情况二：channel 的发送次数大于/小于接收次数
// channel 的发送次数小于接收次数时，接收者 go routine 由于等待发送者发送一直阻塞。
// 因此接收者 go routine 一直未退出，ich 也由于一直被接收者使用无法被垃圾回收。
// 未退出的 go routine 和未被回收的 channel 都造成了内存泄漏的问题。
func TestIsCloseChannelNecessary_on_less_sender(t *testing.T) {
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())
	ich := make(chan int)

	// sender
	go func() {
		//defer close(ich)
		for i := 0; i < 2; i++ {
			ich <- i
		}
	}()

	// receiver
	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(<-ich)
		}
	}()

	time.Sleep(time.Second)
	fmt.Println("NumGoroutine:", runtime.NumGoroutine())

	// Output:
	// NumGoroutine: 2
	// 0
	// 1
	// NumGoroutine: 3
}
