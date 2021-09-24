package crontab

import (
	"fmt"
	"github.com/robfig/cron"
)

func Setup() {
	c := cron.New()
	// 每天午夜跑一次 等式(0 0 0 * * *)0 0/5 * * * ? 更新每天卡可用余额缓存
	c.AddFunc("@daily", testCron)
	// 启动服务
	c.Start()
}

func testCron() {
	fmt.Println("test cron")
}
