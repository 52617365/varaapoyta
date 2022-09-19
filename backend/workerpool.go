package main

type job struct {
	slot             covered_times
	restaurant_id    string
	current_time     date_and_time
	amount_of_eaters int
}

// Result TODO: use Result instead of parsed_graph_data in the results channel.
type Result struct {
	value parsed_graph_data
	err   error
}

// TODO: pass in the response fields into results parameter.
func worker(jobs <-chan job, results chan<- Result) {
	for j := range jobs {
		graph_data, err := interact_with_api(j.slot, j.restaurant_id, j.current_time, j.amount_of_eaters)
		if err != nil {
			results <- Result{parsed_graph_data{}, err}
		}
		results <- Result{graph_data, nil}
	}
}
