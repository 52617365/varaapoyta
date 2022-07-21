package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	payload := payload{
		restaurantId: "2",
		date:         "2022-07-20",
		time:         "12:00",
		amount:       "2",
	}

	payloadString := makePayLoad(payload)
	fmt.Println(payloadString)
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
func getAvailableTables(time time.Duration) {
	fmt.Println("The time is:", time)
	URL := "http://dummy.restapiexample.com/api/v1/employee/1"
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatalln("Oopsie mudafuka")
	}
	fmt.Println(resp.StatusCode)
}
