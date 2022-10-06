package raflaamoTime

import (
	"backend/graphApiResponseStructure"
	"reflect"
	"regexp"
	"testing"
)

var regexToMatchTime = regexp.MustCompile(`\d{2}:\d{2}`)
var regexToMatchDate = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

func TestGraphApiReservationTimes_GetTimeSlotsInBetweenIntervals(t *testing.T) {

	want := []string{"1615", "1630", "1645", "1700"}

	allNeededRaflaamoTimes := GetAllNeededRaflaamoTimes(regexToMatchTime, regexToMatchDate)
	mockParsedGraphData := graphApiResponseStructure.ParsedGraphData{Name: "something"}
	mockParsedIntervalData := make([]graphApiResponseStructure.ParsedIntervalData, 2, 2)
	mockParsedIntervalData[0] = graphApiResponseStructure.ParsedIntervalData{From: 1665061200, To: 1665064800}
	mockParsedGraphData.Intervals = &mockParsedIntervalData

	graphApiReservationTimes := GetGraphApiReservationTimes(&mockParsedGraphData)

	timeSlotsInBetween := graphApiReservationTimes.GetTimeSlotsInBetweenIntervals(allNeededRaflaamoTimes.AllRaflaamoReservationTimeIntervals)

	if !reflect.DeepEqual(want, timeSlotsInBetween) {
		t.Errorf("[raflaamoReservingTimes_test.go] (TestGraphApiReservationTimes_GetTimeSlotsInBetweenIntervals) expected %s but got %s", want, timeSlotsInBetween)
	}
}
