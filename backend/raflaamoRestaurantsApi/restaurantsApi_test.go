/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoRestaurantsApi

import (
	"testing"
)

func TestRaflaamoRestaurantsApi_GetRestaurants(t *testing.T) {
	raflaamoRestaurantsApi, err := GetRaflaamoRestaurantsApi("rovaniemi")
	if err != nil {
		t.Errorf("[TestRaflaamoRestaurantsApi_GetRestaurants] (GetRaflaamoRestaurantsApi) did not expect error but we got one: %s", err)
	}

	var currentTime int64 = 43200
	restaurantsFromApi, err := raflaamoRestaurantsApi.getRestaurantsFromRaflaamoApi(currentTime)
	if err != nil {
		t.Errorf("[TestRaflaamoRestaurantsApi_GetRestaurants] (getRestaurantsFromRaflaamoApi) did not expect error but we got one: %s", err)
	}

	if len(restaurantsFromApi) == 0 {
		t.Errorf("[TestRaflaamoRestaurantsApi_GetRestaurants] (restaurantsFromApi) expected length of result to be < 0 but it was 0")
	}
}
