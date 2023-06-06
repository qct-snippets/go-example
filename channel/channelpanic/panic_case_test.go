package main

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 1.未初始化时关闭
func TestCloseNilChan(t *testing.T) {
	var errCh chan error
	close(errCh)

	// Output:
	// panic: close of nil channel
}

// 2.重复关闭
func TestRepeatClosingChan(t *testing.T) {
	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		close(errCh)
		close(errCh)
	}()

	wg.Wait()

	// Output:
	// panic: close of closed channel
}

// 3.关闭后发送
func TestSendOnClosingChan(t *testing.T) {
	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		close(errCh)
		errCh <- errors.New("chan error")
	}()

	wg.Wait()

	// Output:
	// panic: send on closed channel
}

// 4.发送时关闭
func TestCloseOnSendingToChan(t *testing.T) {
	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(errCh)

		fmt.Println("1")
		go func() {
			fmt.Println("2")
			errCh <- errors.New("chan error") // 由于 chan 没有缓冲队列，代码会一直在此处阻塞
			fmt.Println("3")
		}()
		fmt.Println("4")
		time.Sleep(time.Second) // 等待向 errCh 发送数据
	}()

	wg.Wait()

	// Output:
	// panic: send on closed channel
}
