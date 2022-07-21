package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	dd := getCurrentDate()
	dt := getCurrentTime()
	fmt.Println(dd)
	fmt.Println(dt)
}

// TODO: add some format check before sending in request.
type payload struct {
	// restaurantId is the same as slot_id in payload.
	restaurantId string
	date         string // format is 2022-07-20
	time         string // format is 12:30 corresponds to (12%3A00) in payload.
	amount       string // Amount of eaters.
}

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

func makePayload(id int) string {
	payload_struct := payload{
		restaurantId: strconv.Itoa(id),
		date:         "2022-07-21",
		time:         "17:00",
		amount:       "1",
	}
	// example payload https://s-varaukset.fi/online/reserve/availability/fi/357?date=2022-07-20&slot_id=357&time=12%3A00&amount=1&price_code=&check=1
	payloadString := fmt.Sprintf(
		"https://s-varaukset.fi/online/reserve/availability/fi/%s?date=%s&slot_id=%s&time=%s&amount=%s&price_code=&check=1",
		payload_struct.restaurantId,
		payload_struct.date,
		payload_struct.restaurantId,
		strings.Replace(payload_struct.time, ":", "%3A", -1), // We replace the ":" in time with "%3A" to fit request format.
		payload_struct.amount,
	)
	fmt.Println(payloadString)
	return payloadString
}

func generateUrls() []string {
	urls := []string{}
	for i := 1; i < 21; i++ {
		payload_string := makePayload(i)
		urls = append(urls, payload_string)
	}
	return urls
}

func getAvailableTables() []string {
	urls := generateUrls()

	results := make([]string, len(urls))
	for i := 0; i < len(urls); i++ {
		c := make(chan string)
		go func(url string, channel *chan string) {
			res, err := http.Get(url)
			if err != nil {
				return
			}
			if res.StatusCode != 200 {
				return
			}
			c <- url
		}(urls[i], &c)
		results[i] = <-c
	}
	return results
}
