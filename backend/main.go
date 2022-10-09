/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package main

import (
	"backend/raflaamoRestaurantsApi"
	"backend/restaurants"
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
	restaurantsInstance := restaurants.GetRestaurants(city, amountOfEaters)
	raflaamoRestaurants := restaurantsInstance.GetRestaurantsAndAvailableTables()
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
		iterateAndCaptureRestaurantTimeSlots(restaurant)
		slices.Sort(restaurant.AvailableTimeSlots)
	}
}

func iterateAndCaptureRestaurantTimeSlots(restaurant *raflaamoRestaurantsApi.ResponseFields) {
	for result := range restaurant.GraphApiResults.AvailableTimeSlotsBuffer {
		// This exists because some time slots might have "transparent" field set aka no time slots found.
		// And we are forced to send something down the channel or else it will keep waiting forever expecting n items to iterate.
		if graphApiResponseHadNoTimeSlots(result) {
			continue
		}
		restaurant.AvailableTimeSlots = append(restaurant.AvailableTimeSlots, result)
	}
}

func main() {
	r := gin.Default()
	r.GET("/raflaamo/tables/:city/:amountOfEaters", func(c *gin.Context) {
		city := checkIfCityIsInvalid(c)
		amountOfEaters := c.Param("amountOfEaters")
		collectedRestaurants := GetRestaurantsAndCollectResults(city, amountOfEaters)
		c.JSON(http.StatusOK, collectedRestaurants)
	})
	log.Fatalln(r.Run(":10000"))
} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

func checkIfCityIsInvalid(c *gin.Context) string {
	city := c.Param("city")
	if !slices.Contains(allPossibleCities, strings.ToLower(city)) {
		c.JSON(http.StatusBadRequest, "no results found with that city")
	}
	return city
}
