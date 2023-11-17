package main

import (
	"fmt"
	"strings"
	"yasf.com/backend/playground/pg_def/util"
)


func main() {
	infoStartAt, err := util.GetShortTsTime(getTsCstToTsString("2023-11-10T16:26:23+0800 CST"))
	if err != nil {
		fmt.Printf("==err=:%#v",err.Error())
		return
	}
fmt.Println("=====================infoStartAt===========\n",infoStartAt)
}


func  getTsCstToTsString(ts string) string {
	tsArr := strings.Split(ts, "+")
	if len(tsArr) == 2 {
		return strings.Trim(strings.Replace(tsArr[0], "T", " ", -1)," ")
	}
	return ts
}
