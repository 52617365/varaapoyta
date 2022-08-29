package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tables/{city}", entry_point).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", r))
}
func entry_point(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city := vars["city"]
	if city == "" {
		// TODO: Get all the cities and then match if it even exists on raflaamo site here.
		_, _ = fmt.Fprintf(w, "no city provided")
		return
	}
	restaurants, err := filter_valid_restaurants_from_city(city)
	if err != nil {
		// if error we return this from the endpoint.
		log.Fatalln(err)
	}
	if len(restaurants) == 0 {
		log.Fatalln("no restaurants found")
	}
	available_tables := get_available_tables(restaurants, 1)
	serialize, _ := json.Marshal(available_tables)
	_, _ = fmt.Fprintf(w, string(serialize))
}
