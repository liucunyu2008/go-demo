package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"strconv"
	"strings"
	"time"
	defUtils "yasf.com/backend/playground/pg_def/util"
	pbUserExtensionMonthlyWelfare "yasf.com/backend/playground/user_extension/api/monthly_welfare/v1"

)

const (
	ACTIVE_DONE     = "1"
	ACTIVE_NOT_DONE = "0"
)

const (
	CYCLE_NOT_UNCLOCKED = 1 //未解锁
	CYCLE_CLOCKED       = 2 //已解锁
	CYCLE_ENDED         = 3 //已结束
	CYCLE_CLAIMED       = 4 //已领取
	CYCLE_HIDE          = 5 //隐藏
)

// 1:可领取, 2:已领取, 3:暂未解锁, 4:提升活跃值, 5:隐藏
const (
	AVAILABLE_TO_CLAIM      = 1
	HAVE_CLAIMED            = 2
	NOT_UNLOCKED_YET        = 3
	INCREASE_ACTICITY_VALUE = 4
	HIDE_BUTTON             = 5
)

const USER_IS_VIEW_TYPE = 1 //气泡标识入库数据

type MsgPartCommonRsp struct {
	ErrCode int32  `json:"err_code"`
	ErrDesc string `json:"err_desc"`
}

type MsgGetUserMonthlyWelfareInfoRsp struct {
	MsgPartCommonRsp
	Name                  string            `json:"name"`
	RemainderReceiveTimes int64             `json:"remainder_receive_times"`
	VipLevelConfig        []*vipLevelConfig `json:"vip_level_config"`
}

type vipLevelConfig struct {
	VipLevel           int64          `json:"vip_level"`            //vip等级
	UserIsVipLevel     bool           `json:"user_is_vip_level"`    //当前等级
	Number             int64          `json:"number"`               //可领取次数
	RewardsRefreshInfo []*RefreshItem `json:"rewards_refresh_info"` //刷新日信息
}

type RefreshItem struct {
	RewardsRefreshDay int64          `json:"rewards_refresh_day"` //刷新日期
	Img               string         `json:"img"`                 //加赠图片
	Status            int64          `json:"status"`              // 1 //未解锁 2 //已解锁 3 //已结束 4 //已领取
	ReceiveStatus     int64          `json:"receive_status"`      // 1:可领取, 2:已领取, 3:暂未解锁, 4:提升活跃值, 5:隐藏
	Rewards           []*rewardsItem `json:"rewards"`             //奖励物品
}

type rewardsItem struct {
	Id     string  `json:"id"`     //奖励id
	Type   int64   `json:"type"`   //奖励类型 1：道具，2：饰品，6：碎片
	Name   string  `json:"name"`   //奖励名称
	Desc   string  `json:"desc"`   //奖励描述
	Number int64   `json:"number"` //奖励数量
	Value  float64 `json:"value"`
	High   bool    `json:"high"` //高价值开关
	Img    string  `json:"img"`  //图片
}



func main() {

	var rewardsRefreshList []*pbUserExtensionMonthlyWelfare.RewardsRefreshItem
	rewardsRefreshList = append(rewardsRefreshList, &pbUserExtensionMonthlyWelfare.RewardsRefreshItem{
		RewardsRefreshDays: 1,
		Img:                "",
	})

	rewardsRefreshList = append(rewardsRefreshList, &pbUserExtensionMonthlyWelfare.RewardsRefreshItem{
		RewardsRefreshDays: 22,
		Img:                "",
	})


	var rewardList []*pbUserExtensionMonthlyWelfare.RewardItem
	rewardList = append(rewardList, &pbUserExtensionMonthlyWelfare.RewardItem{
		MetaId:      "aaa",
		Type:        1,
		TreasureNum: 1,
		RewardsDays: "1",
		High:        true,
	})
	rewardList = append(rewardList, &pbUserExtensionMonthlyWelfare.RewardItem{
		MetaId:      "aaa",
		Type:        1,
		TreasureNum: 1,
		RewardsDays: "22",
		High:        true,
	})
	var levelWelfareList []*pbUserExtensionMonthlyWelfare.LevelWelfare
	levelWelfareList = append(levelWelfareList, &pbUserExtensionMonthlyWelfare.LevelWelfare{
		VipLevel:           1,
		RewardsRefreshItem: rewardsRefreshList,
		RewardItem:         rewardList,
	})
	levelWelfareList = append(levelWelfareList, &pbUserExtensionMonthlyWelfare.LevelWelfare{
		VipLevel:           2,
		RewardsRefreshItem: rewardsRefreshList,
		RewardItem:         rewardList,
	})
	welfareReply:= &pbUserExtensionMonthlyWelfare.UbeShelfMonthlyWelfareReply{
		 ReceiveTimes:1701277200,
		 Item: &pbUserExtensionMonthlyWelfare.MonthlyWelfareItem{
			 Id:           1,
			 Name:         "aaa",
			 StartAt:      "2023-12-01 00:00:00",
			 EndAt:        "2023-12-31 23:59:59",
			 Status:       1,
			 LevelWelfare: levelWelfareList,
		 },
	}

	uid := "15459460"
	rsp := getUserMonthlyWelfareInfo(uid, welfareReply)
	fmt.Printf("================================================================:%#v", rsp)
}

func getUserMonthlyWelfareInfo(uid string, welfareReply *pbUserExtensionMonthlyWelfare.UbeShelfMonthlyWelfareReply) MsgGetUserMonthlyWelfareInfoRsp {
	rsp := MsgGetUserMonthlyWelfareInfoRsp{}

	//todo 获取用户本月已领取过奖励的时间
	intUid := int64(1701363600)
	ctx := context.Background()
	userAllReceiveTimes, err := getUserCurrentMonthReceiveDays(ctx, intUid)
	if err != nil {
		rsp.ErrCode, rsp.ErrDesc = int32(defUtils.ErrCode(err)), err.Error()
		return rsp
	}

	userIntVipLevel := 1
	pgActiveStatus := "1"
	//todo 根据用户当前的权益激活状态，判断礼包领取状态
	name, remainderTimes, vipLevelConfigItems, err := getMonthlyWelfareLevelReceiveStatus(uid, welfareReply, userIntVipLevel, userAllReceiveTimes, pgActiveStatus)
	if err != nil {
		rsp.ErrCode, rsp.ErrDesc = int32(defUtils.ErrCode(err)), err.Error()
		return rsp
	}

	vipLevelConfigList := supplementVipLevelConfig(uid, vipLevelConfigItems, int64(userIntVipLevel))

	if len(vipLevelConfigList) > 0 {
		sort.Slice(vipLevelConfigList, func(i, j int) bool {
			return vipLevelConfigList[i].VipLevel < vipLevelConfigList[j].VipLevel
		})
	}
	
	for _,v:=range vipLevelConfigList{
		if v.UserIsVipLevel ==true{
			//fmt.Printf("=============vv===================================================:%#v\n",v)
			for _,val := range v.RewardsRefreshInfo{
				fmt.Printf("===========vip:%v;==val===================================================:%#v\n",v.VipLevel,val)
				continue
			}
			continue
		}
	}
	
	
	rsp.Name = name
	rsp.RemainderReceiveTimes = remainderTimes
	rsp.VipLevelConfig = vipLevelConfigList
	return rsp
}

func getUserCurrentMonthReceiveDays(ctx context.Context, uid int64) (string, error) {
	receiveTimes := "1701277200"
	return receiveTimes, nil
}

func getMonthlyWelfareLevelReceiveStatus(uid string, shelfMonthlyWelfare *pbUserExtensionMonthlyWelfare.UbeShelfMonthlyWelfareReply, userVipLevel int, userAllReceiveTimes, activeStatus string) (string, int64, []*vipLevelConfig, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"uid":                      uid,
		"shelfMonthlyWelfareItems": shelfMonthlyWelfare.Item,
		"userVipLevel":             userVipLevel,
		"userAllReceiveTimes":      userAllReceiveTimes,
		"activeStatus":             activeStatus,
	})
	logEntry.Debug("<ubPlayground.getMonthlyWelfareLevelReceiveStatus> debug")

	if shelfMonthlyWelfare == nil || shelfMonthlyWelfare.Item == nil {
		logEntry.Error("<ubPlayground.getMonthlyWelfareLevelReceiveStatus> shelfMonthlyWelfare == nil")
		return "", 0, nil, errors.New("shelfMonthlyWelfare is nil")
	}

	monthlyWelfareConfigItem := shelfMonthlyWelfare.Item
	name := monthlyWelfareConfigItem.Name
	newTime := time.Now()
	startAt, err := defUtils.GetShortTsTime(monthlyWelfareConfigItem.StartAt)
	if err != nil {
		logEntry.WithFields(logrus.Fields{
			"StartAt": monthlyWelfareConfigItem.StartAt,
		}).WithError(err).Error("<ubPlayground.getMonthlyWelfareLevelReceiveStatus> startAt GetShortTsTime  error")
		return "", 0, nil, err
	}

	endAt, err := defUtils.GetShortTsTime(monthlyWelfareConfigItem.EndAt)
	if err != nil {
		logEntry.WithFields(logrus.Fields{
			"endAt": monthlyWelfareConfigItem.EndAt,
		}).WithError(err).Error("<ubPlayground.getMonthlyWelfareLevelReceiveStatus> endAt GetShortTsTime  error")
		return "", 0, nil, err

	}

	if startAt.After(newTime) || endAt.Before(newTime) {
		logEntry.WithFields(logrus.Fields{
			"startAt": startAt,
			"endAt":   endAt,
			"newTime": newTime,
		}).Error("<ubPlayground.getMonthlyWelfareLevelReceiveStatus> After or Before  error")
		return "", 0, nil, errors.New("StartAt or endAt  Out of range")
	}

	remainderTimes := int64(0)
	var list []*vipLevelConfig
	for _, v := range shelfMonthlyWelfare.Item.LevelWelfare {
		rewardsRefreshList, number, err := refreshItemsHandle(uid, userVipLevel, v, shelfMonthlyWelfare.ReceiveTimes, userAllReceiveTimes, activeStatus)
		if err != nil {
			logEntry.WithFields(logrus.Fields{
				"v": v,
			}).WithError(err).Error("<ubPlayground.getMonthlyWelfareLevelReceiveStatus> RefreshItemsHandle  error")
			return "", 0, nil, err
		}
		logEntry.WithFields(logrus.Fields{
			"userVipLevel": userVipLevel,
			"v":            userVipLevel,
			"ReceiveTimes": shelfMonthlyWelfare.ReceiveTimes,
		}).Debug("<ubPlayground.getMonthlyWelfareLevelReceiveStatus>  LevelWelfare info ")

		var userIsVipLevel bool
		if int64(userVipLevel) == v.VipLevel {
			userIsVipLevel = true
			remainderTimes = int64(number)
		}

		list = append(list, &vipLevelConfig{
			VipLevel:           v.VipLevel,
			UserIsVipLevel:     userIsVipLevel,
			RewardsRefreshInfo: rewardsRefreshList,
			Number:             int64(number),
		})
	}

	return name, remainderTimes, list, nil
}

func supplementVipLevelConfig(uid string, vipLevelConfigItems []*vipLevelConfig, userIntVipLevel int64) []*vipLevelConfig {
	logEntry := logrus.WithFields(logrus.Fields{
		"uid":                 uid,
		"vipLevelConfigItems": vipLevelConfigItems,
		"userIntVipLevel":     userIntVipLevel,
	})
	logEntry.Debug("<ubPlayground.supplementVipLevelConfig> debug")

	currVipMap := make(map[int64]*vipLevelConfig)
	for _, val := range vipLevelConfigItems {
		currVipMap[val.VipLevel] = val
	}

	var vlItems []*vipLevelConfig
	for i := 0; i <= 9; i++ {
		if _, ok := currVipMap[int64(i)]; ok {
			vlItems = append(vlItems, currVipMap[int64(i)])
			continue
		}
		var userIsVipLevel bool
		if userIntVipLevel == int64(i) {
			userIsVipLevel = true
		}
		vlItems = append(vlItems, &vipLevelConfig{
			VipLevel:           int64(i),
			UserIsVipLevel:     userIsVipLevel,
			Number:             0,
			RewardsRefreshInfo: []*RefreshItem{},
		})
	}

	logEntry.WithFields(logrus.Fields{
		"vlItems": vlItems,
	}).Debug("<ubPlayground.supplementVipLevelConfig> result")

	return vlItems
}

func refreshItemsHandle(uid string, userVipLevel int, info *pbUserExtensionMonthlyWelfare.LevelWelfare, receiveTimes int64, userAllReceiveTimes, activeStatus string) ([]*RefreshItem, int, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"uid":                 uid,
		"info":                info,
		"receiveTimes":        receiveTimes,
		"userVipLevel":        userVipLevel,
		"userAllReceiveTimes": userAllReceiveTimes,
		"activeStatus":        activeStatus,
	})
	logEntry.Debug("<ubPlayground.RefreshItemsHandle> info")

	if info == nil || len(info.RewardsRefreshItem) == 0 || len(info.RewardItem) == 0 {
		logEntry.Error("<ubPlayground.RefreshItemsHandle> info is nil or  RewardsRefreshItem len is o RewardItem is o error ")
		return nil, 0, errors.New("RewardsRefreshItem is nil ")
	}

	var cycleItems []int
	var receiveDay int
	for _, v := range info.RewardsRefreshItem {
		cycleItems = append(cycleItems, int(v.RewardsRefreshDays))
	}

	if receiveTimes > 0 {
		isNotCurrMonth := UserIsIssueRewards(receiveTimes)
		if !isNotCurrMonth {
			receiveDay = time.Unix(receiveTimes, 0).Day()
		}
	}

	retrieveStatus, number := getCycleStatusAndRemainderTimes(uid, cycleItems, receiveDay)
	if userVipLevel != int(info.VipLevel) {
		number = 0
	}
	//todo 获取用户活跃状态
	receiveButtonStatusMap := getReceiveDayStatusMap(uid, cycleItems, retrieveStatus, receiveDay, activeStatus, userAllReceiveTimes)

	logEntry.WithFields(logrus.Fields{
		"retrieveStatus":         retrieveStatus,
		"receiveButtonStatusMap": receiveButtonStatusMap,
	}).Debug("<ubPlayground.RefreshItemsHandle> retrieveStatus and receiveButtonStatusMap result")

	var list []*RefreshItem
	for _, rr := range info.RewardsRefreshItem {
		rewardsItems, err := rewardsItemHandle(rr.RewardsRefreshDays, info.RewardItem)
		if err != nil {
			logEntry.WithError(err).Error("<ubPlayground.RefreshItemsHandle> RewardsRefreshDaysHandle  rewardsItemHandle is  error ")
			return nil, 0, err
		}

		list = append(list, &RefreshItem{
			RewardsRefreshDay: rr.RewardsRefreshDays,
			Img:               rr.Img,
			Status:            getRewardsRefreshDaysStatus(userVipLevel, int(rr.RewardsRefreshDays), int(info.VipLevel), retrieveStatus),
			ReceiveStatus:     getReceiveButtonStatus(userVipLevel, int(rr.RewardsRefreshDays), int(info.VipLevel), receiveButtonStatusMap),
			Rewards:           rewardsItems,
		})
	}
	return list, number, nil
}

func UserIsIssueRewards(userReward int64) bool {
	logEntry := logrus.WithFields(logrus.Fields{
		"userReward": userReward,
	})
	logEntry.Debug("<ubPlayground.UserIsIssueRewards>  info")

	userRewardTime := time.Unix(userReward, 0)
	t := time.Now()
	//TODO 改
	if t.Year() != userRewardTime.Year() || t.Month()+1 != userRewardTime.Month() {
		return true
	}
	return false
}

func getCycleStatusAndRemainderTimes(uid string, cycleItems []int, receiveDay int) (map[int]int, int) {
	logEntry := logrus.WithFields(logrus.Fields{
		"uid":        uid,
		"cycleItems": cycleItems,
		"receiveDay": receiveDay,
	})
	logEntry.Debug("<ubPlayground.getCycleStatusAndRemainderTimes> debug")
	// TODO 改
	//today := time.Now().Day()
	today := 1
	currCycle := getCurrentMonthCycleMapNum(uid, cycleItems, today)

	cycleStatusMap := make(map[int]int, 0)
	sort.Sort(sort.IntSlice(cycleItems))
	for _, val := range cycleItems {
		if val < currCycle {
			cycleStatusMap[val] = CYCLE_ENDED
		} else if val == currCycle {
			if receiveDay >= val {
				cycleStatusMap[val] = CYCLE_CLAIMED
			} else {
				cycleStatusMap[val] = CYCLE_CLOCKED
			}
		} else if val > currCycle {
			cycleStatusMap[val] = CYCLE_NOT_UNCLOCKED
		}
	}
	var remainderTimes int
	for _, va := range cycleStatusMap {
		if va == CYCLE_NOT_UNCLOCKED || va == CYCLE_CLOCKED {
			remainderTimes++
		}
		continue
	}
	logEntry.WithFields(logrus.Fields{
		"cycleStatusMap": cycleStatusMap,
		"remainderTimes": remainderTimes,
	}).Debug("<ubPlayground.getCycleStatusAndRemainderTimes> result")

	return cycleStatusMap, remainderTimes
}

func getReceiveDayStatusMap(uid string, cycleItems []int, cycleStatusMap map[int]int, receiveDay int, activeStatus string, userAllReceiveTimes string) map[int]int {
	logEntry := logrus.WithFields(logrus.Fields{
		"uid":                 uid,
		"cycleItems":          cycleItems,
		"cycleStatusMap":      cycleStatusMap,
		"receiveDay":          receiveDay,
		"activeStatus":        activeStatus,
		"userAllReceiveTimes": userAllReceiveTimes,
	})
	logEntry.Debug("<ubPlayground.getReceiveDayStatusMap> debug")
	//currCycle := x.getCurrentMonthCycleMapNum(cycleItems, receiveDay)

	rCycleMap := make(map[int]int)
	if len(userAllReceiveTimes) > 0 {
		rCycleArr := strings.Split(userAllReceiveTimes, ",")
		for _, value := range rCycleArr {
			intValue, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				logrus.WithError(err).Errorf("<ubPlayground.getReceiveDayStatusMap> strconv.ParseInt intValue:%s", value)
			}
			//if intValue == int64(receiveDay) {
			//	continue
			//}
			preCycle := getCurrentMonthCycleMapNum(uid, cycleItems, int(intValue))
			rCycleMap[preCycle] = int(intValue)
		}
	}

	logEntry.WithFields(logrus.Fields{
		"rCycleMap": rCycleMap,
	}).Debug("<ubPlayground.getReceiveDayStatusMap> rCycleMap")

	//todo在map中就已领取
	//todo 获取当月所有的领取日期，并计算所有已领取周期map
	receiveStatusMap := make(map[int]int, 0)
	for key, val := range cycleStatusMap {
		if val == CYCLE_ENDED {
			//todo 周期已过未领取 隐藏
			if _, ok := rCycleMap[key]; ok {
				receiveStatusMap[key] = HAVE_CLAIMED
				continue
			}
			receiveStatusMap[key] = HIDE_BUTTON
		} else if val == CYCLE_CLAIMED {
			//todo 已领取
			receiveStatusMap[key] = HAVE_CLAIMED
		} else if val == CYCLE_CLOCKED {
			//todo 可领取 提升活跃值
			if activeStatus == ACTIVE_DONE {
				receiveStatusMap[key] = AVAILABLE_TO_CLAIM
			} else {
				receiveStatusMap[key] = INCREASE_ACTICITY_VALUE
			}
		} else if val == CYCLE_NOT_UNCLOCKED {
			//todo 暂未解锁
			receiveStatusMap[key] = NOT_UNLOCKED_YET
		}
	}
	logEntry.WithFields(logrus.Fields{
		"receiveStatusMap": receiveStatusMap,
	}).Debug("<ubPlayground.getReceiveDayStatusMap> receiveStatusMap")

	return receiveStatusMap
}

func rewardsItemHandle(rewardsRefreshDay int64, rewardsItems []*pbUserExtensionMonthlyWelfare.RewardItem) ([]*rewardsItem, error) {
	logEntry := logrus.WithFields(logrus.Fields{
		"rewardsRefreshDay": rewardsRefreshDay,
		"rewardsItems":      rewardsItems,
	})
	logEntry.Debug("<ubPlayground.rewardsItemHandle> info")

	if len(rewardsItems) == 0 || rewardsRefreshDay == 0 {
		logEntry.Error("<ubPlayground.rewardsItemHandle> rewardsItems or rewardsRefreshDay is empty error")
		return nil, errors.New("rewardsItems is empty")
	}
	var list []*rewardsItem
	for _, v := range rewardsItems {
		if len(v.RewardsDays) == 0 {
			logEntry.Error("<ubPlayground.rewardsItemHandle> RewardsDays len is nil error")
			return nil, errors.New("RewardsDays is empty")
		}
		rewardsDaysArr := strings.Split(v.RewardsDays, ",")
		if len(rewardsDaysArr) == 0 {
			logEntry.Error("<ubPlayground.rewardsItemHandle> rewardsDaysArr len is nil error")
			return nil, errors.New("rewardsDaysArr is empty")
		}

		if !IsRewardsRefreshDayHandle(rewardsRefreshDay, v.RewardsDays) {
			logEntry.Debug("<ubPlayground.rewardsItemHandle> IsRewardsRefreshDayHandle is false")
			continue
		}

		list = append(list, &rewardsItem{
			Type:   v.Type,
			Number: v.TreasureNum,
			High:   v.High,
		})
	}
	logEntry.WithFields(logrus.Fields{
		"list": list,
	}).Debug("<ubPlayground.rewardsItemHandle> list ")
	return list, nil
}

func IsRewardsRefreshDayHandle(rewardsRefreshDay int64, rewardsDays string) bool {
	if len(rewardsDays) == 0 || rewardsRefreshDay == 0 {
		logrus.Error("<ubPlayground.IsRewardsRefreshDayHandle> len(rewardsDays) == 0")
		return false
	}

	rewardsDaysArr := strings.Split(rewardsDays, ",")
	if len(rewardsDaysArr) == 0 {
		logrus.Error("<ubPlayground.IsRewardsRefreshDayHandle> len(rewardsDaysArr) == 0")
		return false
	}

	for _, v := range rewardsDaysArr {
		if v == fmt.Sprintf("%d", rewardsRefreshDay) {
			return true
		}
	}
	return false
}

func getCurrentMonthCycleMapNum(uid string, cycleItems []int, day int) int {
	logEntry := logrus.WithFields(logrus.Fields{
		"uid":        uid,
		"cycleItems": cycleItems,
		"day":        day,
	})
	logEntry.Debug("<ubPlayground.getCurrentMonthCycleMap> debug")

	if len(cycleItems) == 0 {
		return 0
	}

	sort.Slice(cycleItems, func(i, j int) bool {
		return i < j
	})

	var currCycle int
	for _, v := range cycleItems {
		if day >= v {
			currCycle = v
		}
	}

	logEntry.WithFields(logrus.Fields{
		"currCycle": currCycle,
	}).Debug("<ubPlayground.getCurrentMonthCycleMap> timeCycleMap")

	return currCycle
}

func getRewardsRefreshDaysStatus(userVipLevel, rewardsRefreshDays, vipLevel int, retrieveStatus map[int]int) int64 {
	status := int64(CYCLE_NOT_UNCLOCKED)
	if userVipLevel == vipLevel {
		if rewardsRefreshDaysStatus, ok := retrieveStatus[rewardsRefreshDays]; ok {
			status = int64(rewardsRefreshDaysStatus)
		}
	}
	if userVipLevel > vipLevel {
		status = int64(CYCLE_HIDE)
	}
	return status
}

func getReceiveButtonStatus(userVipLevel, rewardsRefreshDays, vipLevel int, retrieveStatus map[int]int) int64 {
	status := int64(NOT_UNLOCKED_YET)
	if userVipLevel == vipLevel {
		if rewardsRefreshDaysStatus, ok := retrieveStatus[rewardsRefreshDays]; ok {
			status = int64(rewardsRefreshDaysStatus)
		}
	}
	if userVipLevel > vipLevel {
		status = int64(HIDE_BUTTON)
	}
	return status
}
