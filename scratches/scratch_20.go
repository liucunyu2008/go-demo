package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/runtime/protoimpl"
	"sort"
	"time"
	defUtils "yasf.com/backend/playground/pg_def/util"
	pbUserExtensionMonthlyWelfare "yasf.com/backend/playground/user_extension/api/monthly_welfare/v1"
)


const (
	NOT_IS_USER_MONTHLY_REWARD = 0
	IS_USER_MONTHLY_REWARD     = 1
)
const (
	IS_USER_INFO_MONTHLY_REWARD_FIELD  = "user_monthly_welfare_status" // 0 未激活 1 可领取
	USER_MONTHLY_REWARD_ISSUE_TS       = "user_monthly_welfare_issue_ts"
	ERROR_USER_MONTHLY_REWARD_ISSUE_TS = "2099-01-01 00:00:00" // // 错误补充
)


type RewardsRefreshItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RewardsRefreshDays int64  `protobuf:"varint,1,opt,name=rewards_refresh_days,json=rewardsRefreshDays,proto3" json:"rewards_refresh_days,omitempty"`
	Img                string `protobuf:"bytes,2,opt,name=img,proto3" json:"img,omitempty"`
}


func main() {
	// 获取vip 下的领取奖励时间
	var list []*pbUserExtensionMonthlyWelfare.RewardsRefreshItem
	list = append(list, &pbUserExtensionMonthlyWelfare.RewardsRefreshItem{
		RewardsRefreshDays: 1,
	})
	list = append(list, &pbUserExtensionMonthlyWelfare.RewardsRefreshItem{
		RewardsRefreshDays: 8,
	})
	list = append(list, &pbUserExtensionMonthlyWelfare.RewardsRefreshItem{
		RewardsRefreshDays: 15,
	})
	list = append(list, &pbUserExtensionMonthlyWelfare.RewardsRefreshItem{
		RewardsRefreshDays: 22,
	})


	//nts:=UserNextIssueRewardsRefreshTime(list)
	//fmt.Printf("%v\n",nts)
	//return
	//// TODO 加一
	//mapTs:=MonthlyIssueRewardsTsHandle(list)
	//fmt.Printf("================================:%#v\n",mapTs)
	//
	//s,e:=monthlyIssueRewardsTs(mapTs)
	//fmt.Printf("================================s:%#v;e:%v\n",s,e)
	//
	//return
	userRewardTsUnix := int64(1699348636)
	status,ts:=GetUserInfoMonthlyIssueRewardsHandle(list, userRewardTsUnix,1)
	fmt.Printf("--status---:%v;---ts:%v-\n",status,ts)
}

func GetUserInfoMonthlyIssueRewardsHandle(list []*pbUserExtensionMonthlyWelfare.RewardsRefreshItem, userRewardTsUnix int64, intPgActiveStatus int64) (string, string) {
	logEntry := logrus.WithFields(logrus.Fields{
		"list":             list,
		"userRewardTsUnix": userRewardTsUnix,
	})
	logEntry.Debug("<ubPlayground.GetUserInfoMonthlyIssueRewardsHandle> info")
	monthlyIssueRewardsTsMap := MonthlyIssueRewardsTsHandle(list)
	nextTs := UserNextIssueRewardsRefreshTime(list)
	if intPgActiveStatus == 0 {
		logEntry.WithFields(logrus.Fields{
			"userRewardTsUnix":  userRewardTsUnix,
			"intPgActiveStatus": intPgActiveStatus,
		}).Debug("<ubPlayground.GetUserInfoMonthlyIssueRewardsHandle> intPgActiveStatus  info")
		return fmt.Sprintf("%d", NOT_IS_USER_MONTHLY_REWARD), defUtils.FormatShortTs(nextTs)
	}
	// TODO 挪过来的
	startDay, endDay := monthlyIssueRewardsTs(monthlyIssueRewardsTsMap)
	if startDay == 0 {
		logEntry.WithFields(logrus.Fields{
			"startDay": startDay,
			"endDay":   endDay,
		}).Debug("<ubPlayground.GetUserInfoMonthlyIssueRewardsHandle> MonthlyIssueRewardsTs info")
		return fmt.Sprintf("%d", NOT_IS_USER_MONTHLY_REWARD), defUtils.FormatShortTs(nextTs)
	}
	t := time.Now()
	if userRewardTsUnix == 0 {
		// TODO 这块可以省略
		//startDay, _ := monthlyIssueRewardsTs(monthlyIssueRewardsTsMap)
		//t := time.Now()
		nextTs = time.Date(t.Year(), t.Month(), startDay, 0, 0, 0, 0, t.Location())
		logEntry.WithFields(logrus.Fields{
			"nextTs": nextTs,
		}).Debug("<ubPlayground.GetUserInfoMonthlyIssueRewardsHandle> UserNextIssueRewardsRefreshTime info")
		return fmt.Sprintf("%d", IS_USER_MONTHLY_REWARD), defUtils.FormatShortTs(nextTs)
	}

	if userRewardTsUnix > time.Now().Unix() {
		logEntry.WithFields(logrus.Fields{
			"userRewardTsUnix": userRewardTsUnix,
		}).Debug("<ubPlayground.GetUserInfoMonthlyIssueRewardsHandle> userRewardTsUnix>new time info")
		return fmt.Sprintf("%d", NOT_IS_USER_MONTHLY_REWARD), defUtils.FormatShortTs(nextTs)
	}

	isIssueReward := UserIsIssueRewards(userRewardTsUnix)
	if isIssueReward {
		//TODO 可以注释掉
		//startDay, _ := monthlyIssueRewardsTs(monthlyIssueRewardsTsMap)
		//t := time.Now()
		nextTs = time.Date(t.Year(), t.Month(), startDay, 0, 0, 0, 0, t.Location())
		logEntry.WithFields(logrus.Fields{
			"userRewardTsUnix": userRewardTsUnix,
			"isIssueReward":    isIssueReward,
			"nextTs":           nextTs,
		}).Debug("<ubPlayground.GetUserInfoMonthlyIssueRewardsHandle> UserIsIssueRewards info")
		return fmt.Sprintf("%d", IS_USER_MONTHLY_REWARD), defUtils.FormatShortTs(nextTs)
	}

	userReceiveDay := time.Unix(userRewardTsUnix, 0).Day()
	//TODO 挪到最上面
	//startDay, endDay := monthlyIssueRewardsTs(monthlyIssueRewardsTsMap)
	//if startDay == 0 {
	//	logEntry.WithFields(logrus.Fields{
	//		"startDay": startDay,
	//		"endDay":   endDay,
	//	}).Debug("<ubPlayground.GetUserInfoMonthlyIssueRewardsHandle> MonthlyIssueRewardsTs info")
	//	return fmt.Sprintf("%d", NOT_IS_USER_MONTHLY_REWARD), defUtils.FormatShortTs(nextTs)
	//}

	status := IS_USER_MONTHLY_REWARD
	//TODO 这块 不能等于 重点
	//if userReceiveDay >= startDay && userReceiveDay <= endDay {
	if userReceiveDay >= startDay && userReceiveDay < endDay {
		status = NOT_IS_USER_MONTHLY_REWARD
	}
	logEntry.WithFields(logrus.Fields{
		"startTs":        startDay,
		"endTs":          endDay,
		"userReceiveDay": userReceiveDay,
		"nextTs":         nextTs,
		"status":         status,
	}).Debug("<ubPlayground.GetUserInfoMonthlyIssueRewardsHandle> info")
	if status == IS_USER_MONTHLY_REWARD {
		//TODO 可以注释掉
		//t := time.Now()
		nextTs = time.Date(t.Year(), t.Month(), startDay, 0, 0, 0, 0, t.Location())
	}
	return fmt.Sprintf("%d", status), defUtils.FormatShortTs(nextTs)
}

func  UserIsIssueRewards(userReward int64) bool {
	logEntry := logrus.WithFields(logrus.Fields{
		"userReward": userReward,
	})
	logEntry.Debug("<ubPlayground.UserIsIssueRewards>  info")

	userRewardTime := time.Unix(userReward, 0)
	t := time.Now()
	if t.Year() != userRewardTime.Year() || t.Month() != userRewardTime.Month() {
		return true
	}
	return false
}

func  MonthlyIssueRewardsTsHandle(items []*pbUserExtensionMonthlyWelfare.RewardsRefreshItem) map[int]int {
	logEntry := logrus.WithFields(logrus.Fields{
		"items": items,
	})
	logEntry.Debug("<ubPlayground.MonthlyIssueRewardsTsHandle> info")

	var cycleItems []int
	for _, v := range items {
		cycleItems = append(cycleItems, int(v.RewardsRefreshDays))
	}

	t := time.Now()
	firstOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	// TODO 31 天加1 重点
	cycleItems = append(cycleItems, lastOfMonth.Day()+1)
	sort.Ints(cycleItems)

	issueTsMap := make(map[int]int)
	for _, v := range cycleItems {
		for _, val := range cycleItems {
			if _, ok := issueTsMap[v]; !ok {
				if v < val {
					issueTsMap[v] = val
				}
			}
		}
		//TODO 这块注释掉 重点
		//if v == lastOfMonth.Day() {
		//	issueTsMap[v] = lastOfMonth.Day()
		//}
	}
	return issueTsMap
}

func  monthlyIssueRewardsTs(monthlyIssueRewardsTsMap map[int]int) (int, int) {
	logEntry := logrus.WithFields(logrus.Fields{
		"monthlyIssueRewardsTsMap": monthlyIssueRewardsTsMap,
	})
	logEntry.Debug("<ubPlayground.monthlyIssueRewardsTs>")
	t := time.Now()
	toDay := t.Day()
	// 获取当月的第一天
	firstOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	// 获取当月的最后一天
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	for k, v := range monthlyIssueRewardsTsMap {
		if lastOfMonth.Day() == toDay && k == lastOfMonth.Day() {
			return k, v
		} else {
			if toDay >= k && toDay < v {
				return k, v
			}
		}
	}

	return 0, 0
}

// 下次领奖时间
// 小于当前时间就是下个月的时间 下次领奖日期
func  UserNextIssueRewardsRefreshTime(items []*pbUserExtensionMonthlyWelfare.RewardsRefreshItem) time.Time {
	logEntry := logrus.WithFields(logrus.Fields{
		"items": items,
	})
	logEntry.Debug("<ubPlayground.UserNextIssueRewardsRefreshTime> info")

	t := time.Now()
	year := t.Year()
	month := t.Month()
	day :=t.Day()
	var cycleItems []int
	for _, v := range items {
		cycleItems = append(cycleItems, int(v.RewardsRefreshDays))
	}
	sort.Ints(cycleItems)
	nexTs := 0
	for _, v := range cycleItems {
		if v > day {
			nexTs = v
			break
		}
	}

	if nexTs == 0 {
		for _, v := range cycleItems {
			nexTs = v
			break
		}

		month += 1
		if month > 12 {
			month = 1
			year += 1
		}
	}
	return time.Date(year, month, nexTs, 0, 0, 0, 0, t.Location())
}
