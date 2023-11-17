package main

import (
	"fmt"
	"time"
	"yasf.com/backend/playground/pg_def/util"
)

func main() {
	//	TODO 时间大小对比
	s := "2023-09-20 18:49:00"
	tsStr, _ := util.GetShortTsTime(s)
	fmt.Printf("===:%v;str:%v\n", time.Now(), tsStr)

	if tsStr.Before(time.Now()) {
		fmt.Printf("小于当前时间===:%v;str:%v\n", time.Now(), tsStr)
		return
	}
	fmt.Printf("大于当前时间===:%v;str:%v\n", time.Now(), tsStr)

	//TODO 获取当月第一天和最后一天
	startTime := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	endTime := time.Date(startTime.Year(), startTime.Month()+1, startTime.Day(), 0, 0, -1, 0, time.Local)
	fmt.Printf("===:%v======:%v", startTime, endTime)
}
