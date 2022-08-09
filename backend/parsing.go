package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// This file handles everything related to parsing shit.

// 2022-07-21 12:45:54.1414084 +0300 EEST m=+0.001537301
func getCurrentDate() *string {
	re, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	dt := time.Now().String()
	string_formatted := re.FindString(dt)
	return &string_formatted
}

func getCurrentTime() string {
	var re, _ = regexp.Compile(`\d{2}:\d{2}`)
	dt := time.Now().String()
	return re.FindString(dt)
}

// Numbers after 1000 are 4 digits so check if number is under 1000, if so, add trailing zero.
func formatTimesToString(times []int) []string {
	formattedStrings := make([]string, len(times))
	// This for loop first makes sure that everything is the same length (converts 1000< to 4 digits E.g. 800 is 0800)
	// This is done to handle them all the same way.
	for i, t := range times {
		var formattedString string
		// AKA if number length is only 3
		if t < 1000 {
			formattedString = fmt.Sprintf("0%d", t)
			// Here formattedString will look like "1230" so we can assume it will always have 3 indices.
		} else {
			// AKA if number length is already 4
			formattedString = strconv.Itoa(t)
		}
		hour := formattedString[:2]
		minutes := formattedString[2:]
		formattedString = fmt.Sprintf("%s:%s", hour, minutes)
		formattedStrings[i] = formattedString
	}
	return formattedStrings
}

func getAllPossibleTimes() *[]string {
	// Here we have all the possible times when you can reserve a table.
	times := []int{
		800, 815, 830, 845, 900, 915, 930, 945, 1000, 1015, 1030,
		1045, 1100, 1115, 1130, 1145, 1200, 1215, 1230, 1245, 1300,
		1315, 1330, 1345, 1400, 1415, 1430, 1445, 1500, 1515, 1530,
		1545, 1600, 1615, 1630, 1645, 1700, 1715, 1730, 1745, 1800,
		1815, 1830, 1845, 1900, 1915, 1930, 1945, 2000, 2015, 2030,
		2045, 2100, 2115, 2130, 2145, 2200, 2215, 2230, 2245, 2300,
		2315, 2330, 2345,
	}

	// Won't be an error since getCurrentTime returns right value everytime.
	currentTime, _ := strconv.Atoi(strings.ReplaceAll(getCurrentTime(), ":", ""))
	var timesWeWant []string

	// Get all the times we want. (List is sorted, so we can assume that if a number is larger, everything after it will be too, so we don't need a branch for everything after that)
	for i := 0; i < len(times); i++ {
		if times[i] > currentTime {
			timesWeWant = formatTimesToString(times[i:])
			break
		}
	}
	return &timesWeWant
}

// Gets the restaurants from the passed in argument. Returns error if nothing is found.
func getRestaurantsFromCity(city *string) (*[]response_fields, error) {
	restaurants := getAllRestaurantsFromRaflaamoApi()
	captured_restaurants := make([]response_fields, 0, len(*restaurants))

	for _, restaurant := range *restaurants {
		if strings.Contains(strings.ToLower(*restaurant.Address.Municipality.Fi_FI), strings.ToLower(*city)) {
			captured_restaurants = append(captured_restaurants, restaurant)
		}
	}
	if len(captured_restaurants) == 0 {
		return &captured_restaurants, errors.New("no restaurants found")
	}
	return &captured_restaurants, nil
}
