package main

import (
	"fmt"
	"math/rand"
)

func main() {


	for i := 0; i <10;i++{
		s:=randInit(1,100)
		fmt.Println("===\n",s)
	}


}

func  randInit(min, max int) int{
	return rand.Intn(max-min-10)
}

