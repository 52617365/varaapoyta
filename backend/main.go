/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package main

import (
	"backend/endpoint"
	"backend/restaurants"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	config := cors.DefaultConfig()
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"https://raflaamo.rasmusmaki.com"},
		AllowMethods:  []string{"GET"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length"},
	}))
	r.GET("/raflaamo/tables/:city/:amountOfEaters", func(c *gin.Context) {
		e := &endpoint.Endpoint{
			C:    c,
			Cors: config,
		}
		e.UserParameters = e.GetUserRaflaamoParameters()

		if e.UserInputIsInvalid() {
			c.JSON(http.StatusBadRequest, "no results found with that city")
			return
		}

		init := restaurants.GetInitializeProgram(e.UserParameters.City, e.UserParameters.AmountOfEaters)
		collectedRestaurants, err := init.GetRestaurantsAndAvailableTables()
		if err != nil {
			c.JSON(http.StatusBadRequest, "there was problems getting restaurants and/or available tables")
			return
		}
		c.JSON(http.StatusOK, collectedRestaurants)
	})
	log.Fatalln(r.Run(":10000"))
}
