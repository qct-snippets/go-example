package main

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// 【知识点】go channel 关闭后，读取该 channel 永远不会阻塞，且只会输出对应类型的零值。
// 以下代码为例，nil 可能也是需要 channel传输的值之一，通常我们无法通过判断是否为类型的零值确定 channel 是否关闭。
// 所以为了避免输出无意义的值，我们需要一种合理的方式判断 channel 是否关闭。
// golang 官方为我们提供了两种方式。
func TestReadFromClosedChan(t *testing.T) {
	var errCh = make(chan error)

	go func() {
		defer close(errCh)
		errCh <- errors.New("sample error")
	}()

	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(i, <-errCh)
		}
	}()

	time.Sleep(time.Second)

	// Output:
	// 0 chan error
	// 1 <nil>
	// 2 <nil>
}

// 解决方案一：使用 channel 的多重返回值（如 err, ok := <-errCh ）
// err, ok := <-errCh 的第二个返回值 ok 表示 errCh 是否关闭。如果已关闭，则返回 false。
func TestReadFromClosedChan2(t *testing.T) {
	var errCh = make(chan error)
	go func() {
		defer close(errCh)
		errCh <- errors.New("chan error")
	}()

	go func() {
		for i := 0; i < 3; i++ {
			if err, ok := <-errCh; ok {
				fmt.Println(i, err)
			} else {
				fmt.Println(i, "channel closed")
			}
		}
	}()

	time.Sleep(time.Second)

	// Output:
	// 0 chan error
}

// 解决方案二：使用 for range 简化语法
// for range 语法会自动判断 channel 是否结束，如果结束则自动退出 for 循环。
func TestReadFromClosedChanRange(t *testing.T) {
	var errCh = make(chan error)
	go func() {
		defer close(errCh)
		errCh <- errors.New("chan error")
	}()

	go func() {
		i := 0
		for err := range errCh {
			fmt.Println(i, err)
			i++
		}
	}()

	time.Sleep(time.Second)

	// Output:
	// 0 chan error
}
