package main

import (
	"backend/restaurants"
	"fmt"
	"log"
	"net/http"
)

// @Experimental fix is already in place, if it does not work, revisit the problem. the relative raflaamoTime seems to be off when current raflaamoTime is 22:51 and the closing raflaamoTime is 23:30. Closing raflaamoTime points to 2am.
// TODO: Get some auto-completion into the city field on the front end? Figure out how to do this.
// @Performance: we have to make it faster, it's too slow right now but make it faster once everything else works.

var allPossibleCities = [...]string{
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

func setCorrectRequestHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json")
}

func main() {
	raflaamoRestaurants, err := restaurants.GetRestaurantsAndAvailableTables("rovaniemi", 1)
	if err != nil {
		log.Fatalln("err")
	}
	for _, restaurant := range raflaamoRestaurants {
		fmt.Println(restaurant.Name.FiFi)
	}
	//r := mux.NewRouter()
	//r.HandleFunc("/raflaamo/tables/{city}/{amount_of_eaters}", entryPoint).Methods("GET")
	//log.Fatal(http.ListenAndServe(":10000", r))
}

//func entryPoint(w http.ResponseWriter, r *http.Request) {
//	setCorrectRequestHeaders(&w)
//	vars := mux.Vars(r)
//	city := vars["city"]
//	if isNotValidCity(city) {
//		serializedErr, _ := json.Marshal("Sisään syötetyllä kaupungilla ei ole ravintoloita olemassa")
//		w.Write(serializedErr)
//		return
//	}
//
//	amountOfEaters := vars["amount_of_eaters"] //  This is the amount of eaters.
//	amountOfEatersInt := getIntFromAmountOfEaters(amountOfEaters)
//
//	if amountOfEatersInt == -1 {
//		serializedErr, _ := json.Marshal("amount of eaters is unknown")
//		w.Write(serializedErr)
//		return
//	}
//
//	availableTables, err := get_available_tables(city, amountOfEatersInt)
//	if err != nil {
//		errorMessage, _ := json.Marshal(err)
//		_, err2 := w.Write(errorMessage)
//		if err2 != nil {
//			return
//		}
//	}
//	serialize, _ := json.Marshal(availableTables)
//
//	_, err2 := w.Write(serialize)
//	if err2 != nil {
//		return
//	}
//}
//
//func isNotValidCity(city string) bool {
//	return !Contains(allPossibleCities, strings.ToLower(city))
//}
//
//func getIntFromAmountOfEaters(amountOfEaters string) int {
//	if amountOfEaters == "" {
//		return -1
//	}
//	if val, err := strconv.Atoi(amountOfEaters); err == nil {
//		return val
//	}
//	return -1
//
//}
