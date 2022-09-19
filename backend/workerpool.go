package main

type job struct {
	slot             *covered_times
	restaurant_id    string
	current_time     *date_and_time
	amount_of_eaters int
}

type Result struct {
	value *parsed_graph_data
	err   error
}

func worker(jobs <-chan job, results chan<- Result) {
	for j := range jobs {
		graph_data, err := interact_with_api(j.slot, j.restaurant_id, j.current_time, j.amount_of_eaters)
		if err != nil {
			results <- Result{nil, err}
		}
		results <- Result{graph_data, nil}
	}
}
