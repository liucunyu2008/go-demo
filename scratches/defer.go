package main

import "fmt"

func main() {
	//defer 调用的函数参数的值 defer 被定义时就确定了
	//小结：需要强调的时, defer 调用的函数参数的值在 defer 定义时就确定了, 而 defer 函数内部所使用的变量的值需要在这个函数运行时才确定
	//i := 1
	//defer fmt.Println("Deferred print:",i)
	//i++
	//fmt.Println("2 print:", i)
	//i++
	//fmt.Println("3 print:", i)

	f3()
}

func f1() (r int) {
	r = 1
	//
	defer func(r int) {
		r++
		fmt.Println(" r value =",r)
	}(r)
	fmt.Printf("====================rrrrrrrrrr============:%v\n",r)
	r = 3
	return
}

func f2() (r int) {
	r = 1
	//
	defer func() {
		r++
		fmt.Println(" 2r value =",r)
	}()
	fmt.Printf("====================rrrrrrrrrr============:%v\n",r)
	r = 3
	return
}

func f3() (r int) {
	r=1
	go func(r int) {
		fmt.Printf("=======r=========：%v\n",r)
	}(r)
	r++
	if r==4{
		fmt.Printf("========return========：%v\n",r)
		return
	}
	return
}