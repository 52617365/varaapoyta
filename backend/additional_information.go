package main

import (
	"errors"
	"regexp"
)

type AdditionalInformationAboutRestaurant struct {
	restaurant   response_fields
	timeSlots    chan Result
	kitchenTimes restaurant_time
}

func getAdditionalInformation(restaurant response_fields, timeSlotsToCheckLength int) AdditionalInformationAboutRestaurant {
	kitchenOfficeHours := getOpeningAndClosimgTimeFromKitchenTime(&restaurant)
	return AdditionalInformationAboutRestaurant{
		restaurant:   restaurant,
		kitchenTimes: kitchenOfficeHours,
		timeSlots:    make(chan Result, timeSlotsToCheckLength),
	}
}

/*
This struct exists because it lets us filter out the restaurants we're not interested in (E.g. from city we didn't want)
whilst associating the response with the correct restaurant.

also,

Stores certain additional information about the restaurant that we might be interested in.
For example:
- The id from the reservation page url
- Time till restaurant closes.
- Time till restaurants kitchen closes.
*/
func (add *AdditionalInformationAboutRestaurant) add() error {
	restaurantId, err := add.getIdFromReservationPageUrl()
	if err != nil {
		return err
	}
	add.restaurant.Links.TableReservationLocalizedId = restaurantId
	currentTime := get_current_date_and_time()
	time := main2.TimeUtils{
		currentTime:        currentTime,
		closingTime:        add.kitchenTimes.closing,
		timeLeftTillClosed: add.kitchenTimes.closing - currentTime,
	}

	// TODO: add this functionality
	time.currentTime = time.getCurrentTime()

	timeTillRestaurantClosed := time.getTimeTillRestaurantClosingTime()
	timeTillRestaurantsKitchenClosed := time.getTimeLeftToReserve()

	add.restaurant.Openingtime.Time_till_restaurant_closed_hours = timeTillRestaurantClosed.hour
	add.restaurant.Openingtime.Time_till_restaurant_closed_minutes = timeTillRestaurantClosed.minutes

	add.restaurant.Openingtime.Time_left_to_reserve_hours = timeTillRestaurantsKitchenClosed.hour
	add.restaurant.Openingtime.Time_left_to_reserve_minutes = timeTillRestaurantsKitchenClosed.minutes
	return nil
}

/*
Gets the id from a restaurants reservation url.
We have already checked here that the string we're matching contains the string "https://s-varaukset.fi/online/reservation/fi", so it should return something.
We return an error in case but in reality it really should not return error.
*/
func (add *AdditionalInformationAboutRestaurant) getIdFromReservationPageUrl() (string, error) {
	restaurant := add.restaurant
	re, _ := regexp.Compile(`[^fi/]\d+`) // This regex gets the first number match from the TableReservationLocalized JSON field which is the id we want. https://regex101.com/r/NtFMrz/1
	reservationPageUrl := restaurant.Links.TableReservationLocalized.Fi_FI

	idFromReservationPageUrl := re.FindString(reservationPageUrl)

	// If regex could not match or if url was invalid (happens sometimes cuz API is weird).
	if idFromReservationPageUrl == "" {
		return "", errors.New("regex did not match anything, something wrong with reservation_page_url")
	}
	return idFromReservationPageUrl, nil
}
