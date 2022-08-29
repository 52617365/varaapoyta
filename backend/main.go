package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var possible_cities = [...]string{
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

// TODO: Make endpoints.
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tables/{city}", entry_point).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", r))
}
func entry_point(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city := vars["city"]
	_, _ = fmt.Fprintf(w, city)
	if city == "" {
		// TODO: Get all the cities and then match if it even exists on raflaamo site here.
		_, _ = fmt.Fprintf(w, "no city provided")
		return
	}

	restaurants, err := get_all_restaurants_from_raflaamo_api()
	if err != nil {
		log.Fatalln(err)
	}
	if len(restaurants) == 0 {
		log.Fatalln("no restaurants found")
	}
	//available_tables := get_available_tables(city, restaurants, 1)
	//serialize, _ := json.Marshal(available_tables)
	//_, _ = fmt.Fprintf(w, string(serialize))
}
