# 通道关闭原则

一般原则上使用通道是不允许接收方关闭通道和 不能关闭一个有多个并发发送者的通道。 换而言之， 你只能在发送方的 goroutine
中关闭只有该发送方的通道。

* 一个发送者，一个接收者：发送者关闭 channel，接收者使用 select 或 for range 判断 channel 是否关闭。
* 一个发送者，多个接收者：发送者关闭 channel，同上。
* 多个发送者，一个接收者：接收者接收完毕后，使用专用的 stop channel 关闭；发送者使用 select 监听 stop channel 是否关闭。
* 多个发送者，多个接收者：任意一方使用专用的 stop channel 关闭；发送者、接收者都使用 select 监听 stop channel 是否关闭。


* [如何优雅地关闭 channel?](https://learnku.com/go/t/23459/how-to-close-the-channel-gracefully)

## 什么情况下关闭 channel 会造成 panic ？

【知识点】在下述 4 种情况关闭 channel 会引发 panic：

* 未初始化时关闭
* 重复关闭
* 关闭后发送
* 发送时关闭

另外，从 golang 的报错中我们可以知道，golang 认为第3种和第4种情况属于一种情况。

## 如何判断 channel 是否关闭？

【知识点】go channel 关闭后，读取该 channel 永远不会阻塞，且只会输出对应类型的零值。

* [go channel 关闭的那些事儿](https://juejin.cn/post/7033671944587182087)