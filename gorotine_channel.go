package main

//测试协程和管道，利用协程和管道实现一个简单的并发编程
//协程用go+函数名启动可以理解成轻量版的线程，他的创建和销毁不需要cpu干预

import (
	"fmt"
	"math/rand"
	"time"
)

// 生产者模型，不断将随机数发送到通道中
func producer(ch chan<- int) {
	for {
		// 生成一个随机数并将其发送到通道中
		num := rand.Intn(100)
		ch <- num
		// 随机等待一段时间再继续发送下一个数
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

// 消费者模型不断从通道中取数
func consumer(ch <-chan int) {
	for {
		// 从通道中读取一个数并且输出出来
		num := <-ch
		fmt.Println("Received number:", num)

		// 随机等待一段时间再继续读取下一个数
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func main() {
	// 创建一个带有缓冲区大小为 10 的整数通道
	ch := make(chan int, 10)

	// 启动一个生产者 协程和一个消费者 协程
	go producer(ch)
	go consumer(ch)

	// 让主线程暂停一段时间
	time.Sleep(10 * time.Second)
}
