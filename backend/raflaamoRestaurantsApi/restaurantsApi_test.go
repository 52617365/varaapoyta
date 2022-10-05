package raflaamoRestaurantsApi

import (
	"testing"
)

func TestRaflaamoRestaurantsApi_GetRestaurants(t *testing.T) {
	raflaamoRestaurantsApi, err := GetRaflaamoRestaurantsApiStruct("rovaniemi")
	if err != nil {
		t.Errorf("[TestRaflaamoRestaurantsApi_GetRestaurants] (GetRaflaamoRestaurantsApiStruct) did not expect error but we got one: %s", err)
	}

	restaurantsFromApi, err := raflaamoRestaurantsApi.getRestaurantsFromRaflaamoApi()
	if err != nil {
		t.Errorf("[TestRaflaamoRestaurantsApi_GetRestaurants] (getRestaurantsFromRaflaamoApi) did not expect error but we got one: %s", err)
	}

	if len(restaurantsFromApi) == 0 {
		t.Errorf("[TestRaflaamoRestaurantsApi_GetRestaurants] (restaurantsFromApi) expected length of result to be < 0 but it was 0")
	}
}
