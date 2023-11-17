package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
)
func RunTimer(value string)  {
	fmt.Println(value)
	fmt.Println("每5秒执行一次")
}
func main() {
	c := cron.New()
	c.AddFunc("@every 5s", func() {
		RunTimer("传入参数")
	})
	c.AddFunc("@every 5s", func() {
		RunTimer("第二条任务传入参数")
	})
	c.Start()
	select {}
}