package main

import (
	"fmt"
	"sort"
)

func main() {
	// 结构体排序
	sorts()
	return
	m := map[int]int{
		8:  10,
		6:  8,
		1:  2,
		2:  3,
		3:  5,
		30: 30,
	}
	fmt.Printf("======:%#v", mapAscSort(m))
}
// TODO map 排序比较绕，没有现成的，必须循环
func mapAscSort(m map[int]int) map[int]int {
	r := make(map[int]int)
	if len(m) == 0 {
		return r
	}
	var mArr []int
	for k := range m {
		mArr = append(mArr, k)
	}
 // 数组排序
	sort.Ints(mArr)
	for _, v := range mArr {
		fmt.Printf("===v==:%v\n", v)
		if val, ok := m[v]; ok {
			r[v] = val
		}
	}
	return r
}

func sorts() {
	var map1 = make(map[string]string)
	var map2 = make(map[string]string)
	var map3 = make(map[string]string)
	var map4 = make(map[string]string)
	map1["sort"] = "10"
	map2["sort"] = "20"
	map3["sort"] = "30"
	map4["sort"] = "40"

	type ss struct {
		Id    int64             `json:"id"`
		Attrs map[string]string `json:"map_as"`
	}

	var arr []*ss
	arr = append(arr, &ss{
		Id:    3,
		Attrs: map3,
	})
	arr = append(arr, &ss{
		Id:    4,
		Attrs: map4,
	})
	arr = append(arr, &ss{
		Id:    1,
		Attrs: map1,
	})
	arr = append(arr, &ss{
		Id:    2,
		Attrs: map2,
	})

	sort.Slice(arr, func(i, j int) bool {
		// 大到小
		return arr[i].Id > arr[j].Id
	})
	sort.Slice(arr, func(i, j int) bool {
		// 小到大
		return arr[i].Id < arr[j].Id
	})
	for _, v := range arr {
		fmt.Printf("===:%#v\n", v)
	}

}
