package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// This file handles everything related to parsing shit.

// 2022-07-21 12:45:54.1414084 +0300 EEST m=+0.001537301
func getCurrentDate() string {
	re, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	dt := time.Now().String()
	return re.FindString(dt)
}

type timeStruct struct {
	hour    uint8
	minutes uint8
}

// TODO: get all times from the current time onwards.
func getAllPossibleTimes() []string {
	// structs so that we can actually work with the times better (can't use larger than etc. on times that are normally treated as strings)
	times := []timeStruct{
		{hour: 8, minutes: 00}, {hour: 8, minutes: 15},
		{hour: 8, minutes: 30}, {hour: 8, minutes: 45},
		{hour: 9, minutes: 00}, {hour: 9, minutes: 15},
		{hour: 9, minutes: 30}, {hour: 9, minutes: 45},
		{hour: 10, minutes: 00}, {hour: 10, minutes: 15},
		{hour: 10, minutes: 30}, {hour: 10, minutes: 45},
		{hour: 11, minutes: 00}, {hour: 11, minutes: 15},
		{hour: 11, minutes: 30}, {hour: 11, minutes: 45},
		{hour: 12, minutes: 00}, {hour: 12, minutes: 15},
		{hour: 12, minutes: 30}, {hour: 12, minutes: 45},
		{hour: 13, minutes: 00}, {hour: 13, minutes: 15},
		{hour: 13, minutes: 30}, {hour: 13, minutes: 45},
		{hour: 14, minutes: 00}, {hour: 14, minutes: 15},
		{hour: 14, minutes: 30}, {hour: 14, minutes: 45},
		{hour: 15, minutes: 00}, {hour: 15, minutes: 15},
		{hour: 15, minutes: 30}, {hour: 15, minutes: 45},
		{hour: 16, minutes: 00}, {hour: 16, minutes: 15},
		{hour: 16, minutes: 30}, {hour: 16, minutes: 45},
		{hour: 17, minutes: 00}, {hour: 17, minutes: 15},
		{hour: 17, minutes: 30}, {hour: 17, minutes: 45},
		{hour: 18, minutes: 00}, {hour: 18, minutes: 15},
		{hour: 18, minutes: 30}, {hour: 18, minutes: 45},
		{hour: 19, minutes: 00}, {hour: 19, minutes: 15},
		{hour: 19, minutes: 30}, {hour: 19, minutes: 45},
		{hour: 20, minutes: 00}, {hour: 20, minutes: 15},
		{hour: 20, minutes: 30}, {hour: 20, minutes: 45},
		{hour: 21, minutes: 00}, {hour: 21, minutes: 15},
		{hour: 21, minutes: 30}, {hour: 21, minutes: 45},
		{hour: 22, minutes: 00}, {hour: 22, minutes: 15},
		{hour: 22, minutes: 00}, {hour: 22, minutes: 15},
		{hour: 23, minutes: 30}, {hour: 23, minutes: 45},
		{hour: 23, minutes: 30}, {hour: 23, minutes: 45},
	}

	currentTime := strings.Split(getCurrentTime(), ":")
	currentTimeHours, _ := strconv.Atoi(currentTime[0])
	currentTimeMinutes, _ := strconv.Atoi(currentTime[1])

	currentTimeStruct := timeStruct{
		hour:    uint8(currentTimeHours),
		minutes: uint8(currentTimeMinutes),
	}

	timesWeWant := make([]string, len(times))
	for _, t := range times {
		if t.hour > currentTimeStruct.hour && t.minutes > currentTimeStruct.minutes {
			// TODO append to array.
			timeWeWant := fmt.Sprintf("%d:%d", t.hour, t.minutes)
			timesWeWant = append(timesWeWant, timeWeWant)
		}

	}
	return timesWeWant
}

// todo: convert to upper
// 0 15 30 45
func getCurrentTime() string {
	var re, _ = regexp.Compile(`\d{2}:\d{2}`)
	dt := time.Now().String()
	return re.FindString(dt)
}
