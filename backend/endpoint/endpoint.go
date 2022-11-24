/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package endpoint

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
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
	C              *gin.Context
	Cors           cors.Config
	UserParameters *UserParameters
}

type UserParameters struct {
	City           string
	AmountOfEaters string
}

func (endpoint *Endpoint) GetUserRaflaamoParameters() *UserParameters {
	return &UserParameters{
		City:           endpoint.C.Param("city"),
		AmountOfEaters: endpoint.C.Param("amountOfEaters"),
	}
}

/*
* Contract for userInputIsInvalid:
* endpoint.userParameters has to be populated.
 */

func (endpoint *Endpoint) UserInputIsInvalid() bool {
	usersAmountOfEaters := endpoint.UserParameters.AmountOfEaters
	usersCityToCheck := endpoint.UserParameters.City

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
	return !slices.Contains(allPossibleCities, strings.ToLower(endpoint.UserParameters.City))
}

func (endpoint *Endpoint) usersAmountOfEatersIsNotNumber() bool {
	if _, err := strconv.Atoi(endpoint.UserParameters.AmountOfEaters); err != nil {
		return true
	}
	return false
}
