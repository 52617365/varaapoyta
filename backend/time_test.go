package main

import (
	"strings"
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
		time = strings.Replace(time, ":", "", -1)

		unix_time := get_unix_from_time(time)
		if is_invalid_format(time) && unix_time != -1 {
			t.Errorf("expected error")
		}
		if unix_time == -1 && !is_invalid_format(time) {
			t.Errorf("we wanted an error but it wasnt one")
		}
	})
}

//func Fuzz_time_slots_in_between(f *testing.F) {
//	var number_one int64 = 9882912392
//	var number_two int64 = 98812398123
//	f.Add(number_one, number_two)
//	f.Fuzz(func(t *testing.T, current_time int64, end_time int64) {
//		all_time_intervals := get_all_raflaamo_time_intervals()
//		reservation_times, _ := get_time_intervals_in_between_office_hours(get_unix_from_time("1500"), get_unix_from_time("1800"), all_time_intervals)
//		_, err := time_slots_in_between(current_time, end_time, reservation_times)
//
//		if current_time == -1 || end_time == -1 && err == nil {
//			t.Errorf("expected an error but got none")
//		}
//		if current_time == end_time && err == nil {
//			t.Errorf("expected an error but got none")
//		}
//		if current_time > end_time && err == nil {
//			t.Errorf("expected an error but got none")
//		}
//	})
//}

//func Fuzz_get_time_intervals_in_between_office_hours(f *testing.F) {
//	var number_one int64 = 18281991
//	var number_two int64 = 92288128
//	f.Add(number_one, number_two)
//	f.Fuzz(func(t *testing.T, starting_time int64, closing_time int64) {
//		all_time_intervals := get_all_raflaamo_time_intervals()
//		reservation_times, err := get_time_intervals_in_between_office_hours(starting_time, closing_time, all_time_intervals)
//
//		if starting_time > closing_time && err == nil {
//			t.Errorf("expected an error with starting_time: %d and closing_time: %d", starting_time, closing_time)
//		}
//		if !(starting_time > closing_time) && err != nil {
//			t.Errorf("unexpected error with starting_time: %d and closing_time: %d", starting_time, closing_time)
//		}
//		if reservation_times == nil && err == nil {
//			t.Errorf("unexpected nil value in Fuzz_get_all_reservation_times")
//		}
//	})
//}

// TestTimeSlotsFromCurrentPointForward | Test to see if the function correctly gets all the graph time windows from current time forward.
func TestTimeSlotsFromCurrentPointForward(t *testing.T) {
	t.Parallel()

	current_time := "0700"
	expected_time_slot_windows := []covered_times{
		{time: 28800, time_window_start: 21600, time_window_end: 43200},
		{time: 50400, time_window_start: 43200, time_window_end: 64800},
		{time: 72000, time_window_start: 64800, time_window_end: 86400},
	}

	current_time_unix := get_unix_from_time(current_time)
	second_time_slot_windows := get_graph_time_slots_from_current_point_forward(current_time_unix)
	for time_slot_window_index, time_slot_window := range second_time_slot_windows {
		if time_slot_window.time != expected_time_slot_windows[time_slot_window_index].time {
			t.Fatalf("Expected window time to be %d but it was %d", expected_time_slot_windows[time_slot_window_index].time, time_slot_window.time)
		}
		if time_slot_window.time_window_start != expected_time_slot_windows[time_slot_window_index].time_window_start {
			t.Fatalf("Expected time window start time to be %d but it was %d", expected_time_slot_windows[time_slot_window_index].time_window_start, time_slot_window.time_window_start)
		}
		if time_slot_window.time_window_end != expected_time_slot_windows[time_slot_window_index].time_window_end {
			t.Fatalf("Expected time window end time to be %d but it was %d", expected_time_slot_windows[time_slot_window_index].time_window_end, time_slot_window.time_window_end)
		}
	}
}

func Fuzz_times_from_current_point_forward(f *testing.F) {
	var number int64 = 889282828
	f.Add(number)
	f.Fuzz(func(t *testing.T, current_time int64) {
		time_slot_windows := get_graph_time_slots_from_current_point_forward(current_time)
		for _, time_slot := range time_slot_windows {
			if time_slot.time_window_end < current_time {
				t.Errorf(`Did not expect %d to be in the time_slot but it was.`, current_time)
			}
		}
	})
}

//func TestGetAllReservationTimes(t *testing.T) {
//	t.Parallel()
//	tests := []struct {
//		restaurant_starting_time string
//		restaurant_closing_time  string
//		want                     []int64
//	}{
//		// 15 minutes is 900 unix
//		{"2300", "0000", []int64{83700, 84600, 85500, 86400}},
//		{"0000", "0115", []int64{87300, 88200, 89100, 90000, 90900}},
//		{"0700", "0900", []int64{26100, 27000, 27900, 28800, 29700, 30600, 31500, 32400}},
//	}
//
//	for _, test := range tests {
//		testname := fmt.Sprintf("restaurant_starting_time: %s, restaurant_closing_time: %s", test.restaurant_starting_time, test.restaurant_closing_time)
//		t.Run(testname, func(t *testing.T) {
//			all_time_intervals := get_all_raflaamo_time_intervals()
//			times, err := get_time_intervals_in_between_office_hours(get_unix_from_time(test.restaurant_starting_time), get_unix_from_time(test.restaurant_closing_time), all_time_intervals)
//			if err != nil {
//				t.Fatalf(`unexpected error in TestGetAllReservationTimes: %s`, err)
//			}
//			slices.Sort(test.want)
//			slices.Sort(times)
//			if !reflect.DeepEqual(test.want, times) {
//				t.Errorf("length of wrong results is: %d and we wanted: %d", len(times), len(test.want))
//			}
//		})
//	}
//}
