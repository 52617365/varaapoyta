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

func makePayLoad(payload payload) string {
	// TODO: make payload from struct
	// TODO: Turn time from 12:30 to 12%3A00 etc (: is %3A) (might not need to tbh)
	// example payload https://s-varaukset.fi/online/reserve/availability/fi/357?date=2022-07-20&slot_id=357&time=12%3A00&amount=1&price_code=&check=1
	payloadString := fmt.Sprintf(
		"https://s-varaukset.fi/online/reserve/availability/fi/%s?date=%s&slot_id=%s&time=%s&amount=%s&price_code=&check=1",
		payload.restaurantId,
		payload.date,
		payload.restaurantId,
		strings.Replace(payload.time, ":", "%3A", -1), // We replace the ":" in time with "%3A" to fit request format.
		payload.amount,
	)
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
