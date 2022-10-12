/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package main

import (
	"backend/restaurants"
	"fmt"
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
- TODO: when we get an error, everything fucks up, channel waits forever. Happens with helsinki, tampere etc. but not with rovaniemi.
*/

// Starting from here everything should go into own file.

// In other words, if graph API response had the "transparent" field set.
func graphApiResponseHadNoTimeSlots(timeSlotResult string) bool {
	// This exists because some time slots might have "transparent" field set aka no time slots found.
	// And we are forced to send something down the channel or else it will keep waiting forever expecting n items to iterate.
	return timeSlotResult == ""
}

func removeIndexFromSlice[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
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
* Contract for userInputIsInvalid:
* endpoint.userParameters has to be populated.
 */

// TODO: add a struct and methods inside of it for endpoint related stuff.
func (endpoint *Endpoint) userInputIsInvalid() bool {
	usersAmountOfEaters := endpoint.userParameters.amountOfEaters
	usersCityToCheck := endpoint.userParameters.city

	if usersAmountOfEaters == "" || usersCityToCheck == "" {
		return true
	}

	if endpoint.usersAmountOfEatersIsNotNumber() {
		return true
	}

	if endpoint.raflaamoDoesNotContainCity() {
		return true
	}
	return false
}

func (endpoint *Endpoint) raflaamoDoesNotContainCity() bool {
	return !slices.Contains(allPossibleCities, strings.ToLower(endpoint.userParameters.city))
}

func (endpoint *Endpoint) usersAmountOfEatersIsNotNumber() bool {
	if _, err := strconv.Atoi(endpoint.userParameters.amountOfEaters); err != nil {
		return true
	}
	return false
}

func main() {
	init := restaurants.GetInitializeProgram("helsinki", "1")
	collectedRestaurants, err := init.GetRestaurantsAndAvailableTables()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(collectedRestaurants)
	//r := gin.Default()
	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"https://raflaamo.rasmusmaki.com"}
	//r.GET("/raflaamo/tables/:city/:amountOfEaters", func(c *gin.Context) {
	//	endpoint := &Endpoint{
	//		c:    c,
	//		cors: config,
	//	}
	//	endpoint.userParameters = endpoint.getUserRaflaamoParameters()
	//
	//	if endpoint.userInputIsInvalid() {
	//		c.JSON(http.StatusBadRequest, "no results found with that city")
	//		return
	//	}
	//
	//	collectedRestaurants := restaurants.GetRestaurantsAndCollectResults(endpoint.userParameters.city, endpoint.userParameters.amountOfEaters)
	//	c.JSON(http.StatusOK, collectedRestaurants)
	//})
	//log.Fatalln(r.Run(":10000"))
}
