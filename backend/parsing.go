package main

import (
	"regexp"
	"time"
)

// This file handles everything related to parsing shit.

// 2022-07-21 12:45:54.1414084 +0300 EEST m=+0.001537301
func getCurrentDate() string {
	re, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	dt := time.Now().String()
	return re.FindString(dt)
}

func getAllPossibleTimes() []string {
	return []string{
		"8:00", "8:15", "8:30", "8:45",
		"9:00", "9:15", "9:30", "9:45",
		"10:00", "10:15", "10:30", "10:45",
		"11:00", "11:15", "11:30", "11:45",
		"12:00", "12:15", "12:30", "12:45",
		"13:00", "13:15", "13:30", "13:45",
		"14:00", "14:15", "14:30", "14:45",
		"15:00", "15:15", "15:30", "15:45",
		"16:00", "16:15", "16:30", "16:45",
		"17:00", "17:15", "17:30", "17:45",
		"18:00", "18:15", "18:30", "18:45",
		"19:00", "19:15", "19:30", "19:45",
		"20:00", "20:15", "20:30", "20:45",
		"21:00", "21:15", "21:30", "21:45",
		"22:00", "22:15", "22:30", "22:45",
		"23:00", "23:15", "23:30", "23:45",
	}
}

// todo: convert to upper
// 0 15 30 45
func getCurrentTime() string {
	var re, _ = regexp.Compile(`\d{2}:\d{2}`)
	dt := time.Now().String()
	return re.FindString(dt)
}
