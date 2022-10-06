package raflaamoTime

import (
	"backend/raflaamoGraphApi"
	"reflect"
	"regexp"
	"testing"
)

var regexToMatchTime = regexp.MustCompile(`\d{2}:\d{2}`)
var regexToMatchDate = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

type ParsedIntervalData = raflaamoGraphApi.ParsedIntervalData

func TestGraphApiReservationTimes_GetTimeSlotsInBetweenIntervals(t *testing.T) {

	want := []string{"1645", "1700", "1715", "1730", "1745"}

	allNeededRaflaamoTimes := GetAllNeededRaflaamoTimes(regexToMatchTime, regexToMatchDate)
	mockParsedGraphData := ParsedGraphData{Name: "something"}
	var mockParsedIntervalData []ParsedIntervalData
	mockParsedIntervalData[0] = ParsedIntervalData{From: 1660322700000, To: 1660322707200} //2h between these intervals
	mockParsedGraphData.Intervals = &mockParsedIntervalData

	graphApiReservationTimes := GetGraphApiReservationTimes(&mockParsedGraphData)

	timeSlotsInBetween := graphApiReservationTimes.GetTimeSlotsInBetweenIntervals(allNeededRaflaamoTimes.AllRaflaamoReservationTimeIntervals)

	if !reflect.DeepEqual(want, timeSlotsInBetween) {
		t.Errorf("[raflaamoReservingTimes_test.go] (TestGraphApiReservationTimes_GetTimeSlotsInBetweenIntervals) expected %s but got %s", want, timeSlotsInBetween)
	}
}
