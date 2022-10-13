/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package main

import (
	"backend/restaurants"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"strconv"
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
	//init := restaurants.GetInitializeProgram("helsinki", "1")
	//collectedRestaurants, err := init.GetRestaurantsAndAvailableTables()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(collectedRestaurants)
	r := gin.Default()
	config := cors.DefaultConfig()
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"https://raflaamo.rasmusmaki.com"},
		AllowMethods:  []string{"GET"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	r.GET("/raflaamo/tables/:city/:amountOfEaters", func(c *gin.Context) {
		endpoint := &Endpoint{
			c:    c,
			cors: config,
		}
		endpoint.userParameters = endpoint.getUserRaflaamoParameters()

		if endpoint.userInputIsInvalid() {
			c.JSON(http.StatusBadRequest, "no results found with that city")
			return
		}

		init := restaurants.GetInitializeProgram(endpoint.userParameters.city, endpoint.userParameters.amountOfEaters)
		collectedRestaurants, err := init.GetRestaurantsAndAvailableTables()
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.Unwrap(err))
			return
		}
		c.JSON(http.StatusOK, collectedRestaurants)
	})
	log.Fatalln(r.Run(":10000"))
}
