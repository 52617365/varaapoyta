/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package raflaamoTime

import (
	"reflect"
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
			coveredTimes := &CoveredTimes{
				time:            tt.fields.time,
				timeWindowStart: tt.fields.timeWindowStart,
				timeWindowsEnd:  tt.fields.timeWindowsEnd,
			}
			if got := coveredTimes.ConvertUnixTimeToString(); got != tt.want {
				t.Errorf("ConvertUnixTimeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRaflaamoTimes_GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(t *testing.T) {
	type fields struct {
		TimeAndDate                               *TimeAndDate
		AllFutureRaflaamoReservationTimeIntervals []int64
		AllFutureGraphApiTimeIntervals            []string
	}
	tests := []struct {
		name   string
		fields fields
		args   string
		want   []string
	}{
		{name: "TestRaflaamoTimes_GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward", fields: fields{
			TimeAndDate: &TimeAndDate{
				CurrentTime: 30000,
				CurrentDate: "",
			},
		}, args: "1930", want: []string{"0800", "1400", "2000"}},

		{name: "TestRaflaamoTimes_GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward", fields: fields{
			TimeAndDate: &TimeAndDate{
				CurrentTime: 30000,
				CurrentDate: "",
			},
		}, args: "2100", want: []string{"0800", "1400", "2000"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			times := &RaflaamoTimes{
				TimeAndDate: tt.fields.TimeAndDate,
				AllFutureRaflaamoReservationTimeIntervals: tt.fields.AllFutureRaflaamoReservationTimeIntervals,
				AllFutureGraphApiTimeIntervals:            tt.fields.AllFutureGraphApiTimeIntervals,
			}
			times.GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward(tt.args)
			if !reflect.DeepEqual(times.AllFutureGraphApiTimeIntervals, tt.want) {
				t.Errorf("GetAllGraphApiUnixTimeIntervalsFromCurrentPointForward() = %v, want %v", times.AllFutureGraphApiTimeIntervals, tt.want)
			}
		})
	}
}
