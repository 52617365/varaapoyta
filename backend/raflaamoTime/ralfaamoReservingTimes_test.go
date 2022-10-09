/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTime

//var regexToMatchTime = regexp.MustCompile(`\d{2}:\d{2}`)
//var regexToMatchDate = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

// Test passed at 15:44, try again later.
//func TestGraphApiReservationTimes_GetTimeSlotsInBetweenIntervals(t *testing.T) {
//
//	want := []string{"1615", "1630", "1645", "1700"}
//
//	allNeededRaflaamoTimes := GetAllNeededRaflaamoTimes(regexToMatchTime, regexToMatchDate)
//	mockParsedGraphData := graphApiResponseStructure.ParsedGraphData{Name: "something"}
//	mockParsedIntervalData := make([]graphApiResponseStructure.ParsedIntervalData, 1)
//	mockParsedIntervalData[0] = graphApiResponseStructure.ParsedIntervalData{From: 1665144157000 /*3pm*/, To: 1665147757000 /*5pm*/}
//	mockParsedGraphData.Intervals = &mockParsedIntervalData
//
//	graphApiReservationTimes := GetGraphApiReservationTimes(&mockParsedGraphData)
//
//	timeSlotsInBetween := graphApiReservationTimes.GetTimeSlotsInBetweenUnixIntervals(allNeededRaflaamoTimes.AllFutureRaflaamoReservationTimeIntervals)
//
//	if !reflect.DeepEqual(want, timeSlotsInBetween) {
//		t.Errorf("[raflaamoReservingTimes_test.go] (TestGraphApiReservationTimes_GetTimeSlotsInBetweenIntervals) expected %s but got %s", want, timeSlotsInBetween)
//	}
//}
