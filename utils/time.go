package utils

import (
	"time"
)

const (
	TZ_FROMAT = "2006-01-02T15:04:05.000000000Z"
	//HW_FROFROMAT  = "2006-01-02T15:04:05.0000000+08:00"
	HW_FROFROMAT = "2006/01/02 15:04:05 GMT+08:00"
	DATE_FROMAT  = "2006-01-02"
	TIME_FROMAT  = "2006-01-02 15:04:05"
)

func GetNowDayStr() string {
	return time.Now().Format(DATE_FROMAT)
}

func GetDaysAgoStr(days int) string {
	return time.Now().AddDate(0, 0, -days).Format(DATE_FROMAT)
}

func GetFullTime() string {
	return time.Now().Format(TIME_FROMAT)
}
