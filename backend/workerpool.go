package main

type job struct {
	slot             covered_times
	restaurant_id    string
	current_time     date_and_time
	amount_of_eaters int
}

// TODO: pass in the response fields into results parameter.
func worker(jobs <-chan job, results chan<- parsed_graph_data) {
	for j := range jobs {
		graph_data, err := interact_with_api(j.slot, j.restaurant_id, j.current_time, j.amount_of_eaters)
		if err != nil {
			if err.Error() == "error connecting to api" {
				return
			}
			if err.Error() == "error deserializing" {
				return
			} else {
				continue
			}
		}
		results <- graph_data
	}
}

func spawn_jobs(jobs chan<- job, restaurant_id string, kitchen_office_hours restaurant_time, time_slots []covered_times, restaurant response_fields, amount_of_eaters int, current_time date_and_time, city string) {
	if restaurant_format_is_incorrect(city, restaurant) {
		return
	}
	restaurant_additional_information := additional_information{
		restaurant:           restaurant,
		kitchen_office_hours: kitchen_office_hours,
	}

	restaurant_additional_information.add()

	for _, time_slot := range time_slots {
		job := job{
			slot:             time_slot,
			restaurant_id:    restaurant_id,
			current_time:     current_time,
			amount_of_eaters: amount_of_eaters,
		}
		jobs <- job
	}
	close(jobs)
}
