package main

import "fmt"

func main() {
	for i := 0; i <= 9; i++ {
		fmt.Printf("%d\n", i)
	}

	type test struct {
		Id string
	}
	var list []*test
	list = append(list, &test{
		Id: "test1",
	})
	list = append(list, &test{
		Id: "test2",
	})
	for _, v := range list {
		fmt.Printf("%#v\n", v)
	}
	var expr = 1
	switch expr {
	case 1:
		fmt.Printf("%#v\n", 1)
	case 2:
		fmt.Printf("%#v\n", 2)
	default:
		fmt.Printf("%#v\n", 0)
	}

}
