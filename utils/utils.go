package utils

import "time"

func GetCurrentTimePtr() *time.Time {
	t := time.Now()
	return &t
}
