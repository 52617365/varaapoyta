package main

import (
	"testing"
)

// TestTimeSlotsFromCurrentPointForward | Test to see if the function correctly gets all the graph time windows from current time forward.
func TestTimeSlotsFromCurrentPointForward(t *testing.T) {
	time_windows := [...]time_slot_window{
		{time: "0200", time_window_start: "0000", time_window_end: "0600"},
		{time: "0800", time_window_start: "0600", time_window_end: "1200"},
		{time: "1400", time_window_start: "1200", time_window_end: "1800"},
		{time: "2000", time_window_start: "1800", time_window_end: "0000"},
	}
	test_current_time := "1900"

	first_expected_time_slot_window := time_slot_window{time: "2000", time_window_start: "1800", time_window_end: "0000"}

	time_slot_windows := get_time_slots_from_current_point_forward(time_windows, test_current_time)

	for _, time_slot_window := range time_slot_windows {
		if time_slot_window.time != first_expected_time_slot_window.time {
			t.Errorf("Expected window time to be %s but it was %s", first_expected_time_slot_window.time, time_slot_window.time)
		}
		if time_slot_window.time_window_start != first_expected_time_slot_window.time_window_start {
			t.Errorf("Expected time window start time to be %s but it was %s", first_expected_time_slot_window.time_window_start, time_slot_window.time_window_start)
		}
		if time_slot_window.time_window_end != first_expected_time_slot_window.time_window_end {
			t.Errorf("Expected time window end time to be %s but it was %s", first_expected_time_slot_window.time_window_end, time_slot_window.time_window_end)
		}
	}
}
