package main

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/exp/slices"
)

func Fuzz_get_string_time_from_unix(f *testing.F) {
	var unix_time int64 = 90900
	f.Add(unix_time)
	f.Fuzz(func(t *testing.T, unix_time int64) {
		string_time := get_string_time_from_unix(unix_time)

		if string_time == "" {
			t.Errorf("could not get string_time correctly")
		}
	})
}
func Fuzz_get_unix_from_time(f *testing.F) {
	f.Add(2, 30)
	f.Fuzz(func(t *testing.T, hour int, minutes int) {
		unix_time := get_unix_from_time(hour, minutes)

		if unix_time == -1 {
			t.Errorf("fuzzing resulted in -1")
		}
	})
}

func Fuzz_time_slots_inbetween(f *testing.F) {
	f.Add("1200", "0900")
	f.Fuzz(func(t *testing.T, current_time string, end_time string) {
		reservation_times := get_all_reservation_times("1800")
		results, err := time_slots_in_between(current_time, end_time, reservation_times)

		if len(results) == 0 && err == nil {
			t.Errorf("uncaught error in time_slots_inbetween")
		}
	})
}

func Fuzz_get_all_reservation_times(f *testing.F) {
	f.Add("2000")
	f.Fuzz(func(t *testing.T, closing_time string) {
		reservation_times := get_all_reservation_times(closing_time)

		if reservation_times == nil {
			t.Errorf("unexpected nil value.")
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
			t.Errorf("Expected window time to be %s but it was %s", expected_time_slot_windows[time_slot_window_index].time, time_slot_window.time)
		}
		if time_slot_window.time_window_start != expected_time_slot_windows[time_slot_window_index].time_window_start {
			t.Errorf("Expected time window start time to be %s but it was %s", expected_time_slot_windows[time_slot_window_index].time_window_start, time_slot_window.time_window_start)
		}
		if time_slot_window.time_window_end != expected_time_slot_windows[time_slot_window_index].time_window_end {
			t.Errorf("Expected time window end time to be %s but it was %s", expected_time_slot_windows[time_slot_window_index].time_window_end, time_slot_window.time_window_end)
		}
	}
}

func Fuzz_times_from_current_point_forward(f *testing.F) {
	f.Add("1500")
	f.Fuzz(func(t *testing.T, current_time string) {
		time_slot_windows := get_time_slots_from_current_point_forward(current_time)
		for _, time_slot := range time_slot_windows {
			if time_slot.time == current_time || time_slot.time_window_start == current_time || time_slot.time_window_end == current_time {
				t.Errorf(`Did not expect %s to be in the time_slot but it was.`, current_time)
			}
		}
	})
}
func TestConvert_uneven_minutes_to_even(t *testing.T) {
	t.Parallel()
	tests := []struct {
		time string
		want string
	}{
		{"1938", "1945"},
		{"1228", "1230"},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("time: %s, wanted: %s", test.time, test.want)
		t.Run(testname, func(t *testing.T) {
			result := convert_uneven_minutes_to_even(test.time)
			if result != test.want {
				t.Errorf("got %s, want %s", result, test.want)
			}
		})
	}
}

func Fuzz_convert_uneven_minutes_to_even(f *testing.F) {
	f.Add("1438")
	even_minutes := []string{"15", "30", "45", "00"}
	f.Fuzz(func(t *testing.T, number string) {
		result := convert_uneven_minutes_to_even(number)
		if result == "" {
			return
		}
		result_minutes := result[len(result)-2:]
		if !slices.Contains(even_minutes, result_minutes) {
			t.Errorf(`Expected minutes to be 15, 30, 45 or 00 but it was %s`, result_minutes)
		}
	})
}

func TestGetAllReservationTimes(t *testing.T) {
	times := get_all_reservation_times("0100")
	if len(times) != 78 {
		t.Fatalf(`expected len to be %d but it was %d`, 78, len(times))
	}
}

func TestReturnTimeslotsInbetween(t *testing.T) {
	t.Parallel()
	tests := []struct {
		start_time   string
		end_time     string
		closing_time string
		want         []string
	}{
		{"2300", "0100", "0115", []string{"0000", "0015", "0030", "2315", "2330", "2345"}},
		{"1800", "0200", "0000", []string{"1815", "1830", "1845", "1900", "1915", "1930", "1945", "2000", "2015", "2030", "2045", "2100", "2115", "2130", "2145", "2200", "2215", "2230", "2245", "2300", "2315"}},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("start_time %s,end_time %s,closing_time %s", test.start_time, test.end_time, test.closing_time)
		t.Run(testname, func(t *testing.T) {
			all_available_time_slots := get_all_reservation_times(test.closing_time)
			result, err := time_slots_in_between(test.start_time, test.end_time, all_available_time_slots)
			if err != nil {
				t.Errorf(`time_slots had err: %s`, err)
			}
			if !reflect.DeepEqual(result, test.want) {
				t.Errorf(`result len: %d, expected len: %d`, len(result), len(test.want))
			}
		})
	}
}
