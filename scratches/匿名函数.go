package main

import (
	"fmt"
	"reflect"
)

var (
	func1 = func(n1, n2 int) int {
		return n1 + n2
	}
)

func main() {
	//n := func(a, b int) int {
	//	return a + b
	//}
	//fmt.Println(n(1,2))

	n := func1(1,2)
	fmt.Println("==========",n)
	fmt.Println(reflect.TypeOf(func1))  //func(int, int) int
}

