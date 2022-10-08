package restaurants

import (
	"backend/raflaamoGraphApi"
	"backend/raflaamoRestaurantsApi"
	"backend/raflaamoTime"
	"errors"
	"fmt"
	"sync"
)

type Restaurants struct {
	City                   string
	AmountOfEaters         int
	AllNeededRaflaamoTimes *raflaamoTime.RaflaamoTimes
	GraphApi               *raflaamoGraphApi.RaflaamoGraphApi
	RestaurantsApi         *raflaamoRestaurantsApi.RaflaamoRestaurantsApi
}

func GetRestaurants(city string, amountOfEaters int) (*Restaurants, error) {
	allNeededRaflaamoTimes := raflaamoTime.GetAllNeededRaflaamoTimes(RegexToMatchTime, RegexToMatchDate)
	graphApi := raflaamoGraphApi.GetRaflaamoGraphApi()
	initializedRaflaamoRestaurantsApi, err := raflaamoRestaurantsApi.GetRaflaamoRestaurantsApi(city)
	if err != nil {
		return nil, fmt.Errorf("[GetRestaurants] - %w", errors.New("error making restaurants api"))
	}
	return &Restaurants{
		City:                   city,
		AmountOfEaters:         amountOfEaters,
		AllNeededRaflaamoTimes: allNeededRaflaamoTimes,
		GraphApi:               graphApi,
		RestaurantsApi:         initializedRaflaamoRestaurantsApi,
	}, nil
}

// GetRestaurantsAndAvailableTables This is the entry point to the functionality.
func (restaurants *Restaurants) GetRestaurantsAndAvailableTables() ([]raflaamoRestaurantsApi.ResponseFields, error) {
	allRestaurantsFromRaflaamoRestaurantsApi, err := restaurants.RestaurantsApi.GetAllRestaurantsFromRaflaamoRestaurantsApi()
	if err != nil {
		return nil, err
	}

	for index := range allRestaurantsFromRaflaamoRestaurantsApi {
		restaurant := &allRestaurantsFromRaflaamoRestaurantsApi[index]
		go restaurants.getAvailableTablesForRestaurant(restaurant)
	}
	return allRestaurantsFromRaflaamoRestaurantsApi, nil
}

func (restaurants *Restaurants) getAvailableTablesForRestaurant(restaurant *raflaamoRestaurantsApi.ResponseFields) {
	restaurantsKitchenClosingTime := restaurant.Openingtime.Kitchentime.Ranges[0].End
	restaurants.AllNeededRaflaamoTimes.GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(restaurantsKitchenClosingTime)

	raflaamoGraphApiRequestUrlStruct := raflaamoGraphApi.GetRaflaamoGraphApiRequestUrl(restaurant.Links.TableReservationLocalized.FiFi, restaurants.AmountOfEaters, restaurants.AllNeededRaflaamoTimes.TimeAndDate.CurrentDate, regexToMatchRestaurantId)
	restaurant.Links.TableReservationLocalizedId = raflaamoGraphApiRequestUrlStruct.IdFromReservationPageUrl // Storing the id for the front end.

	restaurantGraphApiRequestUrls := raflaamoGraphApiRequestUrlStruct.GenerateGraphApiRequestUrlsForRestaurant(restaurants.AllNeededRaflaamoTimes)

	restaurants.getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel(restaurant, restaurantGraphApiRequestUrls)
}

// TODO: make the channel stuff work here.
func (restaurants *Restaurants) getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel(restaurant *raflaamoRestaurantsApi.ResponseFields, restaurantGraphApiRequestUrls []string) {
	var wg sync.WaitGroup
	for _, restaurantGraphApiRequestUrl := range restaurantGraphApiRequestUrls {
		restaurantGraphApiRequestUrl := restaurantGraphApiRequestUrl
		wg.Add(1)
		go func() {
			defer wg.Done()
			graphApiResponseFromRequestUrl, err := restaurants.GraphApi.GetGraphApiResponseFromTimeSlot(restaurantGraphApiRequestUrl)
			if err != nil {
				restaurant.GraphApiResults.Err <- err
				return
			}
			// getAvailableTableTimesFromRestaurantRequestUrlsIntoRestaurantsChannel TODO: capture restaurants time till kitchen and restaurant closes.

			graphApiReservationTimes := raflaamoTime.GetGraphApiReservationTimes(graphApiResponseFromRequestUrl)

			graphApiReservationTimes.GetTimeSlotsInBetweenIntervals(restaurant.GraphApiResults.AvailableTimeSlotsBuffer, restaurants.AllNeededRaflaamoTimes.AllRaflaamoReservationTimeIntervals)
		}()
	}
	wg.Wait()
	close(restaurant.GraphApiResults.Err)
	close(restaurant.GraphApiResults.AvailableTimeSlotsBuffer)
}
