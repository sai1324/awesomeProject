package main

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	// 创建一个新的 cron 调度器
	c := cron.New()
	// 定义一个定时任务
	task := cron.NewChain(
		// 可选的任务装饰器，用于日志记录
		cron.SkipIfStillRunning(cron.DefaultLogger),
	).Then(cron.FuncJob(func() {
		//执行定时任务
		fmt.Println("定时任务执行:", time.Now())
	}))
	// 使用 cron 表达式设置定时任务的执行时间

	// 添加定时任务到 cron 调度器
	_, err := c.AddJob("@every 1s", task)
	if err != nil {
		log.Fatal("添加定时任务失败:", err)
	}

	c.Start()

	// 等待一段时间，观察定时任务的执行情况
	time.Sleep(1 * time.Minute)

	// 停止 cron 调度器
	c.Stop()
}
