package raflaamoGraphApi

import (
	"testing"
)

func TestRaflaamoGraphApi_getGraphApiResponseFromTimeSlot(t *testing.T) {
	graphApi := GetRaflaamoGraphApi()

	slot, err := graphApi.GetGraphApiResponseFromTimeSlot("https://s-varaukset.fi/api/recommendations/slot/281/2022-08-12/2145/1")
	if err != nil {
		t.Errorf("[TestRaflaamoGraphApi_getGraphApiResponseFromTimeSlot] (GetGraphApiResponseFromTimeSlot) did not expect error")
	}

	if slot == nil || slot.Name == "" {
		t.Errorf("[TestRaflaamoGraphApi_getGraphApiResponseFromTimeSlot] (GetGraphApiResponseFromTimeSlot) did not return valid data.")
	}
}
