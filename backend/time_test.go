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

	current_time := "0700"
	expected_time_slot_windows := []time_slot_window{
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

// // TODO: make this test pass.
// func TestGetLastPossibleTimeBeforeClosing(t *testing.T) {
// 	current_time := "2300"
// 	end_time := "0100"

// 	reservation_times, _ := time_slots_in_between(current_time, end_time)

// 	closing_time := "00:15"
// 	expected_results := []string{"2300", "2315", "2330"}

// 	times := get_last_possible_time_slot_before_closing(reservation_times, closing_time)

// 	for _, v := range times {
// 		fmt.Println(v)
// 	}
// 	for index, expected_result := range expected_results {
// 		if expected_result != times[index] {
// 			t.Errorf("Expected result to be %s but it was %s", expected_result, times[index])
// 		}
// 	}
// }
