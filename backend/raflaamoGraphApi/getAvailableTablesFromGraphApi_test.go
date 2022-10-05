package raflaamoGraphApi

import (
	"testing"
)

func TestRaflaamoGraphApi_getGraphApiResponseFromTimeSlot(t *testing.T) {
	graphApi := GetRaflaamoGraphApi()

	slot, err := graphApi.getGraphApiResponseFromTimeSlot("https://s-varaukset.fi/api/recommendations/slot/281/2022-08-12/2145/1")
	if err != nil {
		t.Errorf("[TestRaflaamoGraphApi_getGraphApiResponseFromTimeSlot] (getGraphApiResponseFromTimeSlot) did not expect error")
	}

	if slot == nil || slot.Name == "" {
		t.Errorf("[TestRaflaamoGraphApi_getGraphApiResponseFromTimeSlot] (getGraphApiResponseFromTimeSlot) did not return valid data.")
	}
}
