package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// TODO: add some format check before sending in request.
type payload struct {
	// restaurantId is the same as slot_id in payload.
	restaurantId string
	date         string // format is 2022-07-20
	time         string // format is 12:30 corresponds to (12%3A00) in payload.
	amount       string // Amount of eaters.
}

func makePayload(id int) string {
	// TODO: make time floor into the latest time e.g (17:15 goes to 17:30) etc.
	payloadStruct := payload{
		restaurantId: strconv.Itoa(id),
		date:         getCurrentDate(),
		time:         getCurrentTime(),
		amount:       "1",
	}
	// example payload https://s-varaukset.fi/online/reserve/availability/fi/357?date=2022-07-20&slot_id=357&time=12%3A00&amount=1&price_code=&check=1
	payloadString := fmt.Sprintf(
		"https://s-varaukset.fi/online/reserve/availability/fi/%s?date=%s&slot_id=%s&time=%s&amount=%s&price_code=&check=1",
		payloadStruct.restaurantId,
		payloadStruct.date,
		payloadStruct.restaurantId,
		strings.Replace(payloadStruct.time, ":", "%3A", -1), // We replace the ":" in time with "%3A" to fit request format.
		payloadStruct.amount,
	)
	return payloadString
}

// This file handles everything related to getting information from restaurants.
func generateUrls() []string {
	var urls []string
	for i := 1; i < 5; i++ {
		payloadString := makePayload(i)
		urls = append(urls, payloadString)
	}
	return urls
}

// This will determine if the body contains something and it will return bool
func stringContains(res *string, substr string) bool {
	if !strings.Contains(*res, substr) {
		return false
	}
	return true
}

func getAvailableTables() []string {
	urls := generateUrls()
	results := make([]string, len(urls))
	for i := 0; i < len(urls); i++ {
		c := make(chan string)
		go func(url string, channel *chan string) {
			body, err := getRequestBody(&url)
			if err != nil {
				log.Fatalln("error with get request.")
				return
			}

			if stringContains(&body, "Jatka varausta") {
				c <- url
			}
			// else use it cuz its valid.
		}(urls[i], &c)
		results[i] = <-c
	}
	return results
}
