# 通道关闭原则
一般原则上使用通道是不允许接收方关闭通道和 不能关闭一个有多个并发发送者的通道。 换而言之， 你只能在发送方的 goroutine 中关闭只有该发送方的通道。

* [如何优雅地关闭 channel?](https://learnku.com/go/t/23459/how-to-close-the-channel-gracefully)