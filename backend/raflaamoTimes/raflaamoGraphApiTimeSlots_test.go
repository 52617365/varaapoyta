/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTimes

import (
	"backend/raflaamoGraphApiTimes"
	"testing"
)

func TestCoveredTimes_ConvertUnixTimeToString(t *testing.T) {
	type fields struct {
		time            int64
		timeWindowStart int64
		timeWindowsEnd  int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "ConvertUnixTimeSoString", fields: fields{
			time:            20000,
			timeWindowStart: 0,
			timeWindowsEnd:  0,
		}, want: "0533"},
		{name: "ConvertUnixTimeSoString", fields: fields{
			time:            80000,
			timeWindowStart: 0,
			timeWindowsEnd:  0,
		}, want: "2213"},
		{name: "ConvertUnixTimeSoString", fields: fields{
			time:            120000,
			timeWindowStart: 0,
			timeWindowsEnd:  0,
		}, want: "0920"},
		{name: "ConvertUnixTimeSoString", fields: fields{
			time:            100000,
			timeWindowStart: 0,
			timeWindowsEnd:  0,
		}, want: "0346"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coveredTimes := &raflaamoGraphApiTimes.CoveredTimes{
				Time:            tt.fields.time,
				TimeWindowStart: tt.fields.timeWindowStart,
				TimeWindowsEnd:  tt.fields.timeWindowsEnd,
			}
			if got := coveredTimes.ConvertUnixTimeToString(); got != tt.want {
				t.Errorf("ConvertUnixTimeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestRaflaamoTimes_GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(t *testing.T) {
//	type fields struct {
//		TimeAndDate                               *TimeAndDate
//		AllFutureRaflaamoReservationTimeIntervals []int64
//		AllFutureGraphApiTimeIntervals            []string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   string
//		want   []string
//	}{
//		{name: "TestRaflaamoTimes_GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward", fields: fields{
//			TimeAndDate: &TimeAndDate{
//				CurrentTime: 30000,
//				CurrentDate: "",
//			},
//		}, args: "1930", want: []string{"0800", "1400", "2000"}},
//
//		{name: "TestRaflaamoTimes_GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward", fields: fields{
//			TimeAndDate: &TimeAndDate{
//				CurrentTime: 30000,
//				CurrentDate: "",
//			},
//		}, args: "2100", want: []string{"0800", "1400", "2000"}},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			times := &RaflaamoTimes{
//				TimeAndDate: tt.fields.TimeAndDate,
//				AllFutureRaflaamoReservationTimeIntervals: tt.fields.AllFutureRaflaamoReservationTimeIntervals,
//			}
//			times.GetAllFutureGraphApiTimeSlots(tt.args)
//			if !reflect.DeepEqual(times.AllFutureGraphApiTimeIntervals, tt.want) {
//				t.Errorf("GetAllFutureGraphApiTimeSlots() = %v, want %v", times.AllFutureGraphApiTimeIntervals, tt.want)
//			}
//		})
//	}
//}
