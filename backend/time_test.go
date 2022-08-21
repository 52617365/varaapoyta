package main

import (
	"fmt"
	"golang.org/x/exp/slices"
	"reflect"
	"testing"
)

func Fuzz_get_string_time_from_unix(f *testing.F) {
	var unix_time int64 = 90900
	f.Add(unix_time)
	f.Fuzz(func(t *testing.T, unix_time int64) {
		string_time := get_string_time_from_unix(unix_time)

		if len(string_time) < 4 {
			t.Errorf("could not get string_time correctly")
		}
	})
}
func Fuzz_get_unix_from_time(f *testing.F) {
	f.Add("1230")
	f.Fuzz(func(t *testing.T, time string) {
		unix_time := get_unix_from_time(time)

		if unix_time == -1 {
			t.Errorf("fuzzing resulted in -1")
		}
	})
}

func Fuzz_time_slots_inbetween(f *testing.F) {
	var number_one int64 = 9882912392
	var number_two int64 = 98812398123
	f.Add(number_one, number_two)
	f.Fuzz(func(t *testing.T, current_time int64, end_time int64) {
		reservation_times, _ := get_all_reservation_times(get_unix_from_time("1500"), get_unix_from_time("1800"))
		results, err := time_slots_in_between(current_time, end_time, reservation_times)

		if current_time == -1 || end_time == -1 && err == nil {
			t.Errorf("expected an error but got none")
		}
		if current_time == end_time && err == nil {
			t.Errorf("expected an error but got none")
		}
		if current_time > end_time && err == nil {
			t.Errorf("expected an error but got none")
		}
		if results == nil && err == nil {
			t.Errorf("uncaught error in time_slots_inbetween")
		}
	})
}

func Fuzz_get_all_reservation_times(f *testing.F) {
	var number_one int64 = 18281991
	var number_two int64 = 92288128
	f.Add(number_one, number_two)
	f.Fuzz(func(t *testing.T, starting_time int64, closing_time int64) {
		reservation_times, err := get_all_reservation_times(starting_time, closing_time)

		if starting_time > closing_time && err == nil {
			t.Errorf("expected an error with starting_time: %d and closing_time: %d", starting_time, closing_time)
		}
		if !(starting_time > closing_time) && err != nil {
			t.Errorf("unexpected error with starting_time: %d and closing_time: %d", starting_time, closing_time)
		}
		if reservation_times == nil && err == nil {
			t.Errorf("unexpected nil value in Fuzz_get_all_reservation_times")
		}
	})
}

// TestTimeSlotsFromCurrentPointForward | Test to see if the function correctly gets all the graph time windows from current time forward.
func TestTimeSlotsFromCurrentPointForward(t *testing.T) {
	t.Parallel()

	current_time := "0700"
	expected_time_slot_windows := []covered_times{
		{time: "0800", time_window_start: "0600", time_window_end: "1200"},
		{time: "1400", time_window_start: "1200", time_window_end: "1800"},
		{time: "2000", time_window_start: "1800", time_window_end: "0000"},
	}

	second_time_slot_windows := get_time_slots_from_current_point_forward(current_time)
	for time_slot_window_index, time_slot_window := range second_time_slot_windows {
		if time_slot_window.time != expected_time_slot_windows[time_slot_window_index].time {
			t.Fatalf("Expected window time to be %s but it was %s", expected_time_slot_windows[time_slot_window_index].time, time_slot_window.time)
		}
		if time_slot_window.time_window_start != expected_time_slot_windows[time_slot_window_index].time_window_start {
			t.Fatalf("Expected time window start time to be %s but it was %s", expected_time_slot_windows[time_slot_window_index].time_window_start, time_slot_window.time_window_start)
		}
		if time_slot_window.time_window_end != expected_time_slot_windows[time_slot_window_index].time_window_end {
			t.Fatalf("Expected time window end time to be %s but it was %s", expected_time_slot_windows[time_slot_window_index].time_window_end, time_slot_window.time_window_end)
		}
	}
}

func Fuzz_times_from_current_point_forward(f *testing.F) {
	f.Add("1500")
	f.Fuzz(func(t *testing.T, current_time string) {
		time_slot_windows := get_time_slots_from_current_point_forward(current_time)
		for _, time_slot := range time_slot_windows {
			if time_slot.time_window_end < current_time {
				t.Errorf(`Did not expect %s to be in the time_slot but it was.`, current_time)
			}
		}
	})
}

func TestGetAllReservationTimes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		restaurant_starting_time string
		restaurant_closing_time  string
		want                     []int64
	}{
		// 15 minutes is 900 unix
		{"2300", "0000", []int64{83700, 84600, 85500, 86400}},
		{"0000", "0115", []int64{87300, 88200, 89100, 90000, 90900}},
		{"0700", "0900", []int64{26100, 27000, 27900, 28800, 29700, 30600, 31500, 32400}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("restaurant_starting_time: %s, restaurant_closing_time: %s", test.restaurant_starting_time, test.restaurant_closing_time)
		t.Run(testname, func(t *testing.T) {
			times, err := get_all_reservation_times(get_unix_from_time(test.restaurant_starting_time), get_unix_from_time(test.restaurant_closing_time))
			if err != nil {
				t.Fatalf(`unexpected error in TestGetAllReservationTimes: %s`, err)
			}
			slices.Sort(test.want)
			slices.Sort(times)
			if !reflect.DeepEqual(test.want, times) {
				t.Errorf("length of wrong results is: %d and we wanted: %d", len(times), len(test.want))
			}
		})
	}
}

func TestReturnTimeslotsInbetween(t *testing.T) {
	t.Parallel()
	tests := []struct {
		start_time              string
		end_time                string
		restaurant_opening_time string
		restaurant_closing_time string
		want                    []string
	}{
		{"2300", "0100", "", "0115", []string{"23:15", "23:30", "23:45", "00:00", "00:15", "00:30", "00:45", "01:00"}},
		{"1800", "0200", "", "0000", []string{"18:15", "18:30", "18:45", "19:00", "19:15", "19:30", "19:45", "20:00", "20:15", "20:30", "20:45", "21:00", "21:15", "21:30", "21:45", "22:00", "22:15", "22:30", "22:45", "23:00", "23:15", "23:30", "23:45", "00:00"}},
		{"1800", "0200", "1900", "0000", []string{"19:15", "19:30", "19:45", "20:00", "20:15", "20:30", "20:45", "21:00", "21:15", "21:30", "21:45", "22:00", "22:15", "22:30", "22:45", "23:00", "23:15", "23:30", "23:45", "00:00"}},
		{"1500", "0300", "1700", "2300", []string{"17:15", "17:30", "17:45", "18:00", "18:15", "18:30", "18:45", "19:00", "19:15", "19:30", "19:45", "20:00", "20:15", "20:30", "20:45", "21:00", "21:15", "21:30", "21:45", "22:00", "22:15", "22:30", "22:45", "23:00"}},
		{"1700", "0300", "1300", "2300", []string{"17:15", "17:30", "17:45", "18:00", "18:15", "18:30", "18:45", "19:00", "19:15", "19:30", "19:45", "20:00", "20:15", "20:30", "20:45", "21:00", "21:15", "21:30", "21:45", "22:00", "22:15", "22:30", "22:45", "23:00"}},
		{"", "", "", "", nil},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("start_time %s,end_time %s,restaurant_opening_time %s restaurant_closing_time %s", test.start_time, test.end_time, test.restaurant_opening_time, test.restaurant_closing_time)
		t.Run(testname, func(t *testing.T) {
			all_available_time_slots, _ := get_all_reservation_times(get_unix_from_time(test.restaurant_opening_time), get_unix_from_time(test.restaurant_closing_time))
			result, err := time_slots_in_between(get_unix_from_time(test.start_time), get_unix_from_time(test.end_time), all_available_time_slots)
			if test.start_time == "" && result == nil && err == nil {
				t.Errorf(`expected an error with start_time: %s`, test.start_time)
			}
			slices.Sort(result)
			if len(result) == len(test.want) {
				for _, item := range result {
					if !slices.Contains(test.want, item) {
						t.Errorf(`expected our list to contain %s but it did not.`, item)
					}
				}
			}
			if len(result) != len(test.want) {
				t.Errorf(`result len: %d, expected len: %d`, len(result), len(test.want))
			}
		})
	}
}
