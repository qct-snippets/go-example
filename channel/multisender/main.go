package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

// 2. 1 个接收者，N 个发送者，接收者通过关闭一个信号通道说 「请不要再发送数据了」
// 我们不能让接收者关闭数据通道，for 不然就会违反了 通道关闭原则。但是我们可以让接收者关闭一个信号通道去通知接收者停止发送数据：
func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	log.SetFlags(0)

	const MaxRandomNumber = 100000
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(1)

	// 通道 dataCh 不曾被关闭过。 是的，通道没有必要关闭。 如果一个通道不会再有 goroutine 去使用它，它最终会被垃圾回收， 不管它是否被关闭。 所以在这里优雅的关闭通道就是不要去关闭通道。
	dataCh := make(chan int, 100)
	stopCh := make(chan struct{})
	// stopCh 是一个信号通道。
	// 它的发送者是 dataCh 的接收者。
	// 它的接收者是 dataCh 的发送者。

	// 发送者
	for i := 0; i < NumSenders; i++ {
		go func() {
			for {
				// 第一个 select 是为了尽可能早的尝试退出 goroutine。
				//  事实上，在这个特殊的例子中，这不是必要的，所以它能省略。
				select {
				case <-stopCh:
					return
				default:
				}

				// 即使 stopCh 已经关闭，如果发送给 dataCh 没有阻塞，那么在第二个 select 中第一个分支可能会在一些循环中不会执行。
				// 但是在这里例子中是可接受的， 所以上面的第一个 select 代码块可以被省略。
				select {
				case <-stopCh:
					return
				case dataCh <- r.Intn(MaxRandomNumber):
				}
			}
		}()
	}

	// 接收者
	go func() {
		defer wgReceivers.Done()

		for value := range dataCh {
			if value == MaxRandomNumber-1 {
				//  dataCh 通道的接收者也是 stopCh 通道的发送者。
				// 在这里关闭停止通道是安全的。
				close(stopCh)
				return
			}
			log.Println(value)
		}
	}()

	wgReceivers.Wait()
}
