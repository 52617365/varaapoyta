package main

import (
	"backend/raflaamoTime"
	"net/http"
)

type job struct {
	slot           *raflaamoTime.CoveredTimes
	restaurantId   string
	currentTime    *raflaamoTime.DateAndTime
	amountOfEaters int
}

type Result struct {
	value *parsedGraphData
	err   error
}

func worker(jobs <-chan job, results chan<- Result) {
	client := &http.Client{}
	for j := range jobs {
		graphApi := RaflaamoGraphApi{
			httpClient: client,
		}
		requestUrl := graphApi.construct_payload(j.restaurantId, j.currentTime.date, j.slot, j.amountOfEaters)
		graphApi.requestUrl = requestUrl

		graphData, err := graphApi.interact_with_api()
		if err != nil {
			results <- Result{nil, err}
		}
		results <- Result{graphData, nil}
	}
}
