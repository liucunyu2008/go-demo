package main

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
	defUtils "yasf.com/backend/playground/pg_def/util"
)
var activityMatchScoreCache                                               *cache.Cache
func main() {


	startAt, err := defUtils.GetShortTsTime("2022-06-08 12:23:23")
	if err != nil {

	}

	endAt, err := defUtils.GetShortTsTime("2023-11-08 12:23:23")
	if err != nil {

	}

	if startAt.After(endAt){
		fmt.Println("=============大于===================")
		return
	}

	if startAt.Before(endAt){
		fmt.Println("=============小于===================")
		return
	}


	if startAt.Before(time.Now()) && endAt.After(time.Now()) {
		fmt.Printf("=================================actAttrsTimeBeforeOk===============================:%#v;%#v", startAt, endAt)
	}
	fmt.Printf("=================================actAttrsTime")
	ss:=(endAt.Unix()-time.Now().Unix())*1000
	fmt.Printf("\n=tt=====%v\n",time.Duration(ss))

	fmt.Printf("\n=tt=55555====%v\n",endAt.Unix()-time.Now().Unix())
	actId:="act_id"
	SetActivityDropCache(actId, time.Duration(ss))


}
func  SetActivityDropCache(actId string,ets time.Duration) {
	activityMatchScoreCache = cache.New(time.Minute*15, time.Minute*15)
	activityMatchScoreCache.Set(actId, "1212121",ets)
}