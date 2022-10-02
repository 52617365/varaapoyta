package main

import (
	"backend/timeUtils"
	"net/http"
)

type job struct {
	slot           *timeUtils.CoveredTimes
	restaurantId   string
	currentTime    *timeUtils.DateAndTime
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
