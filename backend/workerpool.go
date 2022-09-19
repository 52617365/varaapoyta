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
