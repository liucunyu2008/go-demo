//waitGroup
//
//1>.什么是waitGroup
//
//　　WaitGroup用于等待一组Go程的结束。父线程调用Add方法来设定应等待的Go程的数量。每个被等待的Go程在结束时应调用Done方法。同时，主Go程里可以调用Wait方法阻塞至所有Go程结束。
//
//　　实现大致步骤如下：
//　　　　1>.创建 waitGroup对象。
//　　　　　　var wg sync.WaitGroup
//
//　　　　2>.添加 主go程等待的子go程个数。
//　　　　　　wg.Add(数量)
//
//　　　　3>.在各个子go程结束时，调用defer wg.Done()。
//　　　　　　将主go等待的数量-1。注意：实名子go程需传地址。
//
//　　　　4>.在主go程中等待。
//　　　　　　wg.wait()
//2>.waitGroup案例

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	waitGroup()

}


func waitGroup() {
	/**
	  创建 waitGroup对象。
	*/
	var wg sync.WaitGroup
	/**
	  添加 主go程等待的子go程个数。该数量有三种情况:
	      1>.当主Go程添加的子Go程个数和实际子Go程数量相等时,需要等待所有的子Go程执行完毕后主Go程才能正常退出;
	      2>.当主Go程添加的子Go程个数和实际子Go程数量不等时有以下2种情况:
	          a)小于的情况:只需要等待指定的子Go程数量执行完毕后主Go程就会退出，尽管还有其它的子Go程没有运行完成;
	          b)大于的情况:最终会抛出异常"fatal error: all goroutines are asleep - deadlock!"
	*/
	wg.Add(3)

	/**
	  执行子Go程
	*/
	go son1(&wg)
	go son2(&wg)
	go son3(&wg)

	/**
	  在主go程中等待,即主Go程阻塞状态
	*/
	wg.Wait()
}



func son1(group *sync.WaitGroup) {
	/**
	  在各个子go程结束时,一定要调用Done方法，它会通知WaitGroup该子Go程执行完毕哟~
	*/
	defer group.Done()
	time.Sleep(time.Second * 1)
	fmt.Println("son1子Go程结束...")
}

func son2(group *sync.WaitGroup) {
	/**
	  在各个子go程结束时,一定要调用Done方法，它会通知WaitGroup该子Go程执行完毕哟~
	*/
	defer group.Done()
	time.Sleep(time.Second * 2)
	fmt.Println("son2子Go程结束")
}

func son3(group *sync.WaitGroup) {
	/**
	  在各个子go程结束时,一定要调用Done方法，它会通知WaitGroup该子Go程执行完毕哟~
	*/
	defer group.Done()
	time.Sleep(time.Second * 3)
	fmt.Println("son3子Go程结束~~~")
}
