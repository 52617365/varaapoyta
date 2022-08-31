package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var all_possible_cities = [...]string{
	"helsinki",
	"espoo",
	"vantaa",
	"nurmijärvi",
	"kerava",
	"järvenpää",
	"vihti",
	"porvoo",
	"lohja",
	"hyvinkää",
	"karkkila",
	"riihimäki",
	"tallinna",
	"hämeenlinna",
	"lahti",
	"forssa",
	"salo",
	"kotka",
	"kouvola",
	"akaa",
	"loimaa",
	"heinola",
	"hamina",
	"kaarina",
	"turku",
	"kangasala",
	"raisio",
	"tampere",
	"nokia",
	"luumäki",
	"laitila",
	"lappeenranta",
	"mikkeli",
	"rauma",
	"ulvila",
	"pori",
	"jyväskylä",
	"imatra",
	"pieksämäki",
	"savonlinna",
	"varkaus",
	"seinäjoki",
	"kuopio",
	"joensuu",
	"kitee",
	"vaasa",
	"iisalmi",
	"lieksa",
	"kokkola",
	"ylivieska",
	"nurmes",
	"kajaani",
	"sotkamo",
	"muhos",
	"kempele",
	"oulu",
	"rovaniemi",
	"kittilä"}

func set_correct_request_headers(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json")
}

func Contains[T comparable](arr [58]T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tables/{city}/{amount_of_eaters}", entry_point).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", r))
}
func entry_point(w http.ResponseWriter, r *http.Request) {
	// TODO: set this to explicit urls once it works.
	set_correct_request_headers(&w)
	vars := mux.Vars(r)
	city := vars["city"]
	if is_not_valid_city(city) {
		w.Write([]byte("no restaurants with that city"))
		//_, _ = fmt.Fprintf(w, "no restaurants with that city")
		return
	}

	amount_of_eaters := vars["amount_of_eaters"] //  This is the amount of eaters.
	amount_of_eaters_int := get_int_from_amount_of_eaters(amount_of_eaters)

	if amount_of_eaters_int == -1 {
		w.Write([]byte("amount of eaters is unknown"))
		//_, _ = fmt.Fprintf(w, "amount of eaters is unknown")
		return
	}

	// TODO: why is times not included in the serialize_string?
	available_tables := get_available_tables(city, amount_of_eaters_int)
	serialize, _ := json.Marshal(available_tables)

	_, err2 := w.Write(serialize)
	if err2 != nil {
		return
	}
}

func is_not_valid_city(city string) bool {
	return !Contains(all_possible_cities, city)
}

func get_int_from_amount_of_eaters(amount_of_eaters string) int {
	if amount_of_eaters == "" {
		return -1
	}
	if val, err := strconv.Atoi(amount_of_eaters); err == nil {
		return val
	}
	return -1

}
