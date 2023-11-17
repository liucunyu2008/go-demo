package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		fmt.Printf("=====:%v\n",time.Now())
	}
}