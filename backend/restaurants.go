package main

import (
	"fmt"
	"log"
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

type payloadResult struct {
	restaurant string
	available  bool
}

func makePayload(id int, amount int, time string) string {
	// TODO: make time floor into the latest time e.g (17:15 goes to 17:30) etc.
	// example payload https://s-varaukset.fi/online/reserve/availability/fi/357?date=2022-07-20&slot_id=357&time=12%3A00&amount=1&price_code=&check=1
	payloadString := fmt.Sprintf(
		"https://s-varaukset.fi/online/reserve/availability/fi/%d?date=%s&slot_id=%d&time=%s&amount=%d&price_code=&check=1",
		id,
		getCurrentDate(),
		id,
		strings.Replace(time, ":", "%3A", -1), // We replace the ":" in time with "%3A" to fit request format.
		//strings.Replace(getCurrentTime(), ":", "%3A", -1), // We replace the ":" in time with "%3A" to fit request format.
		amount,
	)
	return payloadString
}

// Hakee kaikki vapaat ravintolat.
func workerRequest(jobs chan string, resultsChan chan payloadResult) {
	for job := range jobs {
		body, err := getRequestBody(&job)
		if err != nil {
			log.Fatal("error sending get request")
		}

		// TODO: parse body.
		// "Jatka varausta" on sivulla, jos paikka on vapaana.
		if strings.Contains(body, "Jatka varausta") {
			result := payloadResult{
				restaurant: job,
				available:  true,
			}
			resultsChan <- result
		} else {
			result := payloadResult{
				restaurant: job,
				available:  false,
			}
			resultsChan <- result
		}
	}
}

func getAvailableTables() {
	fmt.Println("starting")
	// 4920 = 82 (pages) x 64 (amount of payloads per page, times)
	jobs := make(chan string, 5248)
	resultsChan := make(chan payloadResult, len(jobs))
	times := getAllPossibleTimes()
	// Launch 8 workers
	for i := 0; i < 20; i++ {
		go workerRequest(jobs, resultsChan)
	}

	// Spawning all the jobs.
	for i := 1; i <= 82; i++ {
		for _, time := range times {
			payloadString := makePayload(i, 1, time)
			jobs <- payloadString
		}
	}
	close(jobs)
	//results := make([]string, len(jobs))
	for a := 1; a <= len(jobs); a++ {
		fmt.Println(<-resultsChan)
	}
}
