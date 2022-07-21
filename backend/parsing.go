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

func getCurrentTime() string {
	re, _ := regexp.Compile(`\d{2}:\d{2}`)
	dt := time.Now().String()
	return re.FindString(dt)
}
