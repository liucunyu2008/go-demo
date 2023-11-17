//实现部分
package main

import (
	"fmt"
	"sync"
	"testing"
)

// 包私有,该结构体不能直接使用
type goldaxe struct{}

func (axe goldaxe) TellTruth() {
	fmt.Printf("==========================\n")
}

// 全局变量
var (
	axe  *goldaxe
	once sync.Once
)

// 由于单例类型不能在包外直接使用，用一个接口类型带出去
type GoldAxe interface {
	TellTruth()
}

// 用于获取单例模式对象,大家都是一样的斧子
func GetGoldAxe() GoldAxe {
	once.Do(func() {
		axe = &goldaxe{}
	})

	return axe
}

//测试部分


const axeCounts = 100

func main()  {
	var t *testing.T
	//Test1(t)
	fmt.Printf("====:%#v",t)
	return
	Test2(t)
}

func Test1(t *testing.T) {
	ins1 := GetGoldAxe()
	ins1.TellTruth()
	ins2 := GetGoldAxe()
	if ins1 != ins2 {
		t.Fatal("instance is not equal")
	}
}

func Test2(t *testing.T) {
	start := make(chan struct{})
	//信号量初始化
	wg := sync.WaitGroup{}
	//信号量搞个100个
	wg.Add(axeCounts)
	//金斧子数组，并进行了列表初始化
	//这么写你紫定get了:var float_array = [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
	instances := [axeCounts]GoldAxe{}
	for i := 0; i < axeCounts; i++ {
		//开启100个协程
		go func(index int) {
			//获取channel的值，由于没有协程只能阻塞
			<-start
			instances[index] = GetGoldAxe()
			wg.Done()
		}(i)
	}
	//关闭channel，所有协程同时GetGoldAxe，达到了并发创建实例的情况
	close(start)
	//等待大家都获得自己金斧子后，才到下一步
	wg.Wait()
	for i := 1; i < axeCounts; i++ {
		if instances[i] != instances[i-1] {
			t.Fatal("instance is not equal")
		}
	}
}