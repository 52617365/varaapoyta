package main

import (
	"testing"

	"golang.org/x/exp/slices"
)

// TestTimeSlotsFromCurrentPointForward | Test to see if the function correctly gets all the graph time windows from current time forward.
func TestTimeSlotsFromCurrentPointForward(t *testing.T) {
	t.Parallel()
	time_windows := [...]covered_times{
		{time: "0200", time_window_start: "0000", time_window_end: "0600"},
		{time: "0800", time_window_start: "0600", time_window_end: "1200"},
		{time: "1400", time_window_start: "1200", time_window_end: "1800"},
		{time: "2000", time_window_start: "1800", time_window_end: "0000"},
	}

	current_time := "0700"
	expected_time_slot_windows := []covered_times{
		{time: "0800", time_window_start: "0600", time_window_end: "1200"},
		{time: "1400", time_window_start: "1200", time_window_end: "1800"},
		{time: "2000", time_window_start: "1800", time_window_end: "0000"},
	}

	second_time_slot_windows := get_time_slots_from_current_point_forward(time_windows, current_time)
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

func TestReturnTimeslotsInbetween(t *testing.T) {
	t.Parallel()
	expected_result_range := [...]string{"0015", "0030", "0045", "0100"}

	start_time := "0015"
	end_time := "0100"
	closing_time := "0200"

	all_available_time_slots := get_all_reservation_times(closing_time)

	time_slots, err := time_slots_in_between(start_time, end_time, all_available_time_slots)

	if err != nil {
		t.Errorf(`TestReturn_time_slots_in_between failed completely with start_time: %s and end_time: %s`, start_time, end_time)
	}

	for index := range time_slots {
		if time_slots[index] != expected_result_range[index] {
			t.Errorf(`expected time slot to be %s but it was %s`, time_slots[index], expected_result_range[index])
		}
	}
}

func TestConvert_uneven_minutes_to_even(t *testing.T) {
	t.Parallel()
	test_uneven_number := "1228"
	expected_even_number := "1230"

	even_number := convert_uneven_minutes_to_even(test_uneven_number)

	if even_number != expected_even_number {
		t.Fatalf(`expected even number to be %s but it was %s`, expected_even_number, even_number)
	}

	test_uneven_number2 := "1938"
	expected_even_number2 := "1945"

	even_number2 := convert_uneven_minutes_to_even(test_uneven_number2)

	if even_number2 != expected_even_number2 {
		t.Fatalf(`expected even number to be %s but it was %s`, expected_even_number2, even_number2)
	}
}

// doesnt pass
func TestGetAllReservationTimes(t *testing.T) {
	times := get_all_reservation_times("0100")
	if len(times) != 62 {
		t.Fatalf(`expected len to be %d but it was %d`, 62, len(times))
	}
}

// doesnt pass
func TestReturnTimeslotsInbetween2(t *testing.T) {
	t.Parallel()
	expected_result_range := []string{"0000", "0015", "0030", "1800", "1815", "1830", "1845",
		"1900", "1915", "1930", "1945",
		"2000", "2015", "2030", "2045",
		"2100", "2115", "2130", "2145",
		"2200", "2215", "2230", "2245",
		"2300", "2315", "2330", "2345",
	}

	start_time := "1800"
	end_time := "0100"
	closing_time := "0115" // last time is therefore 0030
	// seems to be returning times 1h before it should.
	all_available_time_slots := get_all_reservation_times(closing_time)

	time_slots, err := time_slots_in_between(start_time, end_time, all_available_time_slots)

	if err != nil {
		t.Errorf(`TestReturn_time_slots_in_between failed completely with start_time: %s and end_time: %s`, start_time, end_time)
	}

	for _, time_slot := range time_slots {
		if !slices.Contains(expected_result_range, time_slot) {
			t.Errorf(`expected time slot to contain %s but it did not`, time_slot)
		}
	}
}
