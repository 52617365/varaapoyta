/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package main

import (
	"backend/raflaamoRestaurantsApi"
	"backend/restaurants"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
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

/*
* Stuff to fix:
GET http://localhost:10000/raflaamo/tables/tampere/1
- TODO: some tableReservationLocalized fields are empty for some reason.
- TODO: some available_time_slots are null for some reason, if they're null they should not even be included.
*/

// Starting from here everything should go into own file.

// In other words, if graph API response had the "transparent" field set.
func graphApiResponseHadNoTimeSlots(timeSlotResult string) bool {
	// This exists because some time slots might have "transparent" field set aka no time slots found.
	// And we are forced to send something down the channel or else it will keep waiting forever expecting n items to iterate.
	return timeSlotResult == ""
}
func GetRestaurantsAndCollectResults(city string, amountOfEaters string) []raflaamoRestaurantsApi.ResponseFields {
	restaurantsInstance := restaurants.GetRestaurants(city, amountOfEaters)
	raflaamoRestaurants := restaurantsInstance.GetRestaurantsAndAvailableTables()
	iterateRestaurants(raflaamoRestaurants)
	return raflaamoRestaurants
}

func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func iterateRestaurants(raflaamoRestaurants []raflaamoRestaurantsApi.ResponseFields) {
	for index := range raflaamoRestaurants {
		restaurant := &raflaamoRestaurants[index]
		err := <-restaurant.GraphApiResults.Err
		if err != nil {
			continue
		}
		timeSlotsForRestaurant := iterateAndCaptureRestaurantTimeSlots(restaurant)
		// Making sure that we don't get restaurants that don't have any time slots.
		if len(timeSlotsForRestaurant) == 0 {
			remove(raflaamoRestaurants, index)
			continue
		}
		restaurant.AvailableTimeSlots = timeSlotsForRestaurant

		slices.Sort(restaurant.AvailableTimeSlots)
	}
}

// Refactor stops here (to own file).

func iterateAndCaptureRestaurantTimeSlots(restaurant *raflaamoRestaurantsApi.ResponseFields) []string {
	availableTimeSlots := make([]string, 0, 50)
	for result := range restaurant.GraphApiResults.AvailableTimeSlotsBuffer {
		if graphApiResponseHadNoTimeSlots(result) {
			continue
		}
		availableTimeSlots = append(availableTimeSlots, result)
	}
	return availableTimeSlots
}

type Endpoint struct {
	c              *gin.Context
	cors           cors.Config
	userParameters *UserParameters
}

type UserParameters struct {
	city           string
	amountOfEaters string
}

func (endpoint *Endpoint) getUserRaflaamoParameters() *UserParameters {
	return &UserParameters{
		city:           endpoint.c.Param("city"),
		amountOfEaters: endpoint.c.Param("amountOfEaters"),
	}
}

/*
* Contract for validateUserInput:
* userParameters has to be populated.
 */

// TODO: add a struct and methods inside of it for endpoint related stuff.
func (endpoint *Endpoint) userInputIsValid() bool {
	usersAmountOfEaters := endpoint.userParameters.amountOfEaters
	usersCityToCheck := endpoint.userParameters.city
	if usersAmountOfEaters == "" || usersCityToCheck == "" {
		return false
	}

	if endpoint.usersAmountOfEatersIsNotNumber() {
		return false
	}

	if endpoint.raflaamoDoesNotContainCity() {
		return false
	}
	return true
}

func (endpoint *Endpoint) raflaamoDoesNotContainCity() bool {
	if !slices.Contains(allPossibleCities, endpoint.userParameters.city) {
		return true
	}
	return false
}

func (endpoint *Endpoint) usersAmountOfEatersIsNotNumber() bool {
	if _, err := strconv.Atoi(endpoint.userParameters.amountOfEaters); err != nil {
		return true
	}
	return false
}

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://raflaamo.rasmusmaki.com"}
	r.GET("/raflaamo/tables/:city/:amountOfEaters", func(c *gin.Context) {
		userParameters := getUserParameters(c)
		if checkIfCityIsInvalid(userParameters.city) {
			c.JSON(http.StatusBadRequest, "no results found with that city")
			return
		}
		collectedRestaurants := GetRestaurantsAndCollectResults(userParameters.city, userParameters.amountOfEaters)
		c.JSON(http.StatusOK, collectedRestaurants)
	})
	log.Fatalln(r.Run(":10000"))
}

func checkIfCityIsInvalid(city string) bool {
	if !slices.Contains(allPossibleCities, strings.ToLower(city)) {
		return true
	}
	return false
}
