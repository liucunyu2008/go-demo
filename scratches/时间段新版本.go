package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"time"
	defUtils "yasf.com/backend/playground/pg_def/util"
	pbUserExtensionMonthlyWelfare "yasf.com/backend/playground/user_extension/api/monthly_welfare/v1"
)

func main() {
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
	list = append(list, &pbUserExtensionMonthlyWelfare.RewardsRefreshItem{
		RewardsRefreshDays: 30,
	})
	userRewardTsUnix := int64(1698854399)
	status, ts := GetUserInfoMonthlyIssueRewardsHandle(list, userRewardTsUnix, 1)
	fmt.Printf("--status---:%v;---ts:%v-\n", status, ts)

}

const (
	NOT_IS_USER_MONTHLY_REWARD = 0
	IS_USER_MONTHLY_REWARD     = 1
)

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

	if userRewardTsUnix == 0 {
		startDayUnix, _ := monthlyIssueRewardsTs(monthlyIssueRewardsTsMap)
		nextTs = time.Unix(startDayUnix, 0)
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
		startDay, _ := monthlyIssueRewardsTs(monthlyIssueRewardsTsMap)
		nextTs = time.Unix(startDay, 0)
		logEntry.WithFields(logrus.Fields{
			"userRewardTsUnix": userRewardTsUnix,
			"isIssueReward":    isIssueReward,
		}).Debug("<ubPlayground.GetUserInfoMonthlyIssueRewardsHandle> UserIsIssueRewards info")
		return fmt.Sprintf("%d", IS_USER_MONTHLY_REWARD), defUtils.FormatShortTs(nextTs)
	}

	userReceiveDay := time.Unix(userRewardTsUnix, 0).Day()
	startDay, endDay := monthlyIssueRewardsTs(monthlyIssueRewardsTsMap)
	if startDay == 0 {
		logEntry.WithFields(logrus.Fields{
			"startDay": startDay,
			"endDay":   endDay,
		}).Debug("<ubPlayground.GetUserInfoMonthlyIssueRewardsHandle> MonthlyIssueRewardsTs info")
		return fmt.Sprintf("%d", NOT_IS_USER_MONTHLY_REWARD), defUtils.FormatShortTs(nextTs)
	}

	status := IS_USER_MONTHLY_REWARD
	if userRewardTsUnix >= startDay && userRewardTsUnix <= endDay {
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
		nextTs = time.Unix(startDay, 0)
	}
	return fmt.Sprintf("%d", status), defUtils.FormatShortTs(nextTs)
}

func UserIsIssueRewards(userReward int64) bool {
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

func MonthlyIssueRewardsTsHandle(items []*pbUserExtensionMonthlyWelfare.RewardsRefreshItem) map[int64]int64 {
	logEntry := logrus.WithFields(logrus.Fields{
		"items": items,
	})
	logEntry.Debug("<ubPlayground.MonthlyIssueRewardsTsHandle> info")
	t := time.Now()
	var cycleItems []int
	for _, v := range items {
		cycleItems = append(cycleItems, int(time.Date(t.Year(), t.Month(), int(v.RewardsRefreshDays), 0, 0, 0, 0, t.Location()).Unix()))
	}
	
	firstOfMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	cycleItems = append(cycleItems, int(firstOfMonth.Unix()))
	sort.Ints(cycleItems)
	issueTsMap := make(map[int64]int64)
	for _, v := range cycleItems {
		for _, val := range cycleItems {
			if _, ok := issueTsMap[int64(v)]; !ok {
				if v < val {
				if 	time.Unix(int64(val),0).Month()>t.Month(){
					issueTsMap[int64(v)] = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, -1, t.Location()).Unix()
					continue
					}
					issueTsMap[int64(v)] = time.Date(t.Year(), t.Month(), time.Unix(int64(val),0).Day(), 0, 0, 0, -1, t.Location()).Unix()
				}
			}
		}
	}
	return issueTsMap
}

func monthlyIssueRewardsTs(monthlyIssueRewardsTsMap map[int64]int64) (int64, int64) {
	logEntry := logrus.WithFields(logrus.Fields{
		"monthlyIssueRewardsTsMap": monthlyIssueRewardsTsMap,
	})
	logEntry.Debug("<ubPlayground.monthlyIssueRewardsTs>")
	t:=time.Now()
	toDayUnix :=t.Unix()
	for k, v := range monthlyIssueRewardsTsMap {
		if toDayUnix >=k && toDayUnix<=v{
			return k, v
		}
	}
	return 0, 0
}

// 下次领奖时间
// 小于当前时间就是下个月的时间 下次领奖日期
func UserNextIssueRewardsRefreshTime(items []*pbUserExtensionMonthlyWelfare.RewardsRefreshItem) time.Time {
	logEntry := logrus.WithFields(logrus.Fields{
		"items": items,
	})
	logEntry.Debug("<ubPlayground.UserNextIssueRewardsRefreshTime> info")
	t := time.Now()
	year := t.Year()
	month := t.Month()
	day := t.Day()
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
