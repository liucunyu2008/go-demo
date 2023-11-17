//1>.什么是读写锁
//
//　　互斥锁的本质是当一个goroutine访问的时候，其他goroutine都不能访问。这样在资源同步，避免竞争的同时也降低了程序的并发性能。程序由原来的并行执行变成了串行执行。其实，当我们对一个不会变化的数据只做"读"操作的话，是不存在资源竞争的问题的。因为数据是不变的，不管怎么读取，多少goroutine同时读取，都是可以的。
//
//　　所以问题不是出在"读"上，主要是修改，也就是"写"。修改的数据要同步，这样其他goroutine才可以感知到。所以真正的互斥应该是读取和修改、修改和修改之间，读和读是没有互斥操作的必要的。因此，衍生出另外一种锁，叫做读写锁。
//
//　　读写锁可以让多个读操作并发，同时读取，但是对于写操作是完全互斥的。也就是说，当一个goroutine进行写操作的时候，其他goroutine既不能进行读操作，也不能进行写操作。
//
//　　GO中的读写锁由结构体类型sync.RWMutex表示。此类型的方法集合中包含两对方法：
//　　　　一组是对写操作的锁定和解锁，简称"写锁定"和"写解锁"：
//　　　　　　func (*RWMutex)Lock()
//　　　　　　func (*RWMutex)Unlock()
//
//　　　　另一组表示对读操作的锁定和解锁，简称为"读锁定"与"读解锁"：
//　　　　　　func (*RWMutex)RLock()
//　　　　　　func (*RWMutex)RUnlock()
//2>.读写锁的案例

package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	number int
	rwlock sync.RWMutex //定义读写锁
)

func MyRead(n int) {
	rwlock.RLock()         //添加读锁
	defer rwlock.RUnlock() //使用结束时自动解锁
	fmt.Printf("[%d] Goroutine读取数据为: %d\n", n, number)
}

func MyWrite(n int) {
	rwlock.Lock()         //添加写锁
	defer rwlock.Unlock() //使用结束时自动解锁
	//number = rand.Intn(100)
	number = n
	fmt.Printf("%d Goroutine写入数据为: %d\n", n, number)
}

func main() {

	//创建写端
	for index := 1; index <= 10; index++ {
		go MyWrite(index)
	}

	//创建读端
	for index := 1; index <= 20; index++ {
		go MyRead(index)
	}

	for {
		time.Sleep(time.Second)
	}
}
