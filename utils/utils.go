package utils

import "time"

func GetStartToday() int64 {
	loc, _ := time.LoadLocation("Local")
	now := time.Now().In(loc)
	startDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Unix()
	return startDay
}
