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

	want := []string{"1645", "1700", "1715", "1730", "1745"}

	allNeededRaflaamoTimes := GetAllNeededRaflaamoTimes(regexToMatchTime, regexToMatchDate)
	mockParsedGraphData := graphApiResponseStructure.ParsedGraphData{Name: "something"}
	var mockParsedIntervalData []graphApiResponseStructure.ParsedIntervalData
	mockParsedIntervalData[0] = graphApiResponseStructure.ParsedIntervalData{From: 1660322700000, To: 1660322707200} //2h between these intervals
	mockParsedGraphData.Intervals = &mockParsedIntervalData

	graphApiReservationTimes := GetGraphApiReservationTimes(&mockParsedGraphData)

	timeSlotsInBetween := graphApiReservationTimes.GetTimeSlotsInBetweenIntervals(allNeededRaflaamoTimes.AllRaflaamoReservationTimeIntervals)

	if !reflect.DeepEqual(want, timeSlotsInBetween) {
		t.Errorf("[raflaamoReservingTimes_test.go] (TestGraphApiReservationTimes_GetTimeSlotsInBetweenIntervals) expected %s but got %s", want, timeSlotsInBetween)
	}
}
