package main

import "fmt"

func main() {
	type UserVipLevel int64

	const (
		V0 UserVipLevel = 0
		V1 UserVipLevel = 1
		V2 UserVipLevel = 2
		V3 UserVipLevel = 3
		V4 UserVipLevel = 4
		V5 UserVipLevel = 5
		V6 UserVipLevel = 6
		V7 UserVipLevel = 7
		V8 UserVipLevel = 8
		V9 UserVipLevel = 9
	)

	fmt.Printf("================================================================:%#v",UserVipLevel(11))
}
