/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package main

import (
	"backend/raflaamoRestaurantsApi"
	"backend/restaurants"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"strings"
)

var allPossibleCities = []string{
	"helsinki",
	"espoo",
	"vantaa",
	"nurmijärvi",
	"kerava",
	"järvenpää",
	"vihti",
	"porvoo",
	"lohja",
	"hyvinkää",
	"karkkila",
	"riihimäki",
	"tallinna",
	"hämeenlinna",
	"lahti",
	"forssa",
	"salo",
	"kotka",
	"kouvola",
	"akaa",
	"loimaa",
	"heinola",
	"hamina",
	"kaarina",
	"turku",
	"kangasala",
	"raisio",
	"tampere",
	"nokia",
	"luumäki",
	"laitila",
	"lappeenranta",
	"mikkeli",
	"rauma",
	"ulvila",
	"pori",
	"jyväskylä",
	"imatra",
	"pieksämäki",
	"savonlinna",
	"varkaus",
	"seinäjoki",
	"kuopio",
	"joensuu",
	"kitee",
	"vaasa",
	"iisalmi",
	"lieksa",
	"kokkola",
	"ylivieska",
	"nurmes",
	"kajaani",
	"sotkamo",
	"muhos",
	"kempele",
	"oulu",
	"rovaniemi",
	"kittilä"}

//	func setCorrectRequestHeaders(w *http.ResponseWriter) {
//		(*w).Header().Set("Access-Control-Allow-Origin", "*")
//		(*w).Header().Set("Content-Type", "application/json")
//	}

// In other words, if graph API response had the "transparent" field set.
func graphApiResponseHadNoTimeSlots(timeSlotResult string) bool {
	return timeSlotResult == ""
}
func GetRestaurantsAndCollectResults(city string, amountOfEaters string) []raflaamoRestaurantsApi.ResponseFields {
	restaurantsInstance, err := restaurants.GetRestaurants(city, amountOfEaters)
	if err != nil {
		fmt.Println(err)
	}
	raflaamoRestaurants, err := restaurantsInstance.GetRestaurantsAndAvailableTables()
	if err != nil {
		log.Fatalln(err)
	}
	iterateRestaurants(raflaamoRestaurants)
	return raflaamoRestaurants
}

func iterateRestaurants(raflaamoRestaurants []raflaamoRestaurantsApi.ResponseFields) {
	for index := range raflaamoRestaurants {
		restaurant := &raflaamoRestaurants[index]
		err := <-restaurant.GraphApiResults.Err
		if err != nil {
			continue
		}
		restaurant.AvailableTimeSlots = make([]string, 0, 50)
		iterateRestaurantTimeSlots(restaurant)
		slices.Sort(restaurant.AvailableTimeSlots)
	}
}

func iterateRestaurantTimeSlots(restaurant *raflaamoRestaurantsApi.ResponseFields) {
	for result := range restaurant.GraphApiResults.AvailableTimeSlotsBuffer {
		// This exists because some time slots might have "transparent" field set aka no time slots found.
		// And we are forced to send something down the channel or else it will keep waiting forever expecting n items to iterate.
		if graphApiResponseHadNoTimeSlots(result) {
			continue
		}
		restaurant.AvailableTimeSlots = append(restaurant.AvailableTimeSlots, result)
	}
}

// TODO: figure out why response body is 500 even though the response seems to be fine.
func main() {
	r := gin.Default()
	r.GET("raflaamo/tables/:city/:amountOfEaters", func(c *gin.Context) {
		city := checkIfCityIsInvalid(c)
		amountOfEaters := c.Param("amountOfEaters")
		collectedRestaurants := GetRestaurantsAndCollectResults(city, amountOfEaters)
		c.JSON(http.StatusOK, collectedRestaurants)
	})
	err := r.Run(":10000")
	if err != nil {
		fmt.Println(err)
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func checkIfCityIsInvalid(c *gin.Context) string {
	city := c.Param("city")
	if !slices.Contains(allPossibleCities, strings.ToLower(city)) {
		c.JSON(http.StatusBadRequest, "no results found with that city")
	}
	return city
}
