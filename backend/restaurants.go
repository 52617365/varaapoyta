package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

type payloadResult struct {
	restaurant string
	available  bool
}

func makePayload(id int, amount int, time string) string {
	// example payload https://s-varaukset.fi/online/reserve/availability/fi/357?date=2022-07-20&slot_id=357&time=12%3A00&amount=1&price_code=&check=1
	payloadString := fmt.Sprintf(
		"https://s-varaukset.fi/online/reserve/availability/fi/%d?date=%s&slot_id=%d&time=%s&amount=%d&price_code=&check=1",
		id,
		getCurrentDate(),
		id,
		strings.Replace(time, ":", "%3A", -1), // We replace the ":" in time with "%3A" to fit request format.
		amount,
	)
	return payloadString
}

// Hakee kaikki vapaat ravintolat.
func workerRequest(jobs chan string, resultsChan chan payloadResult) {
	for job := range jobs {
		body, err := getRequestBody(&job)
		if err != nil {
			log.Fatal("error sending get request")
		}
		// "Jatka varausta" on sivulla, jos paikka on vapaana.
		if strings.Contains(body, "Jatka varausta") {
			result := payloadResult{
				restaurant: job,
				available:  true,
			}
			resultsChan <- result
		} else {
			result := payloadResult{
				restaurant: job,
				available:  false,
			}
			resultsChan <- result
		}
	}
}

func getAvailableTables() {
	fmt.Println("starting")
	// 4920 = 82 (pages) x 64 (amount of payloads per page, times)
	jobs := make(chan string, 160000)
	resultsChan := make(chan payloadResult, len(jobs))
	times := getAllPossibleTimes()

	// TODO: recover from this some how.
	if times == nil {
		log.Fatal("no times available.")
	}

	// Launch 20 workers
	for i := 0; i < 20; i++ {
		go workerRequest(jobs, resultsChan)
	}

	// Spawning all the jobs.
	for i := 1; i <= 1500; i++ {
		for _, time := range times {
			payloadString := makePayload(i, 1, time)
			jobs <- payloadString
		}
	}
	close(jobs)

	for a := 1; a <= 160000; a++ {
		result := <-resultsChan
		if result.available {
			// TODO: do something with the information.
			fmt.Println(result.restaurant)
		}
	}
}

// This function will scrape all the possible restaurants that are possible to be reserved from a certain city with colly.

// Colly will first look at the page, store the urls from all of the tables that you're able to reserve.
// li > div > section > a(contains link)

type restaurant struct {
	reservingUrl string
	name         string
	canReserve   bool
}

// The <li>'s which contain the stuff, the list items all have the same class on every page so we can just check that class name.
// ListItemStyles__RestaurantListItemWrapper-sc-1xaojw6-15 jDtLYn
func scrapeRestaurantLocations( /*city string*/ ) {
	//URL := fmt.Sprintf(
	//	"https://raflaamo.fi/fi/ravintolat/%s/kaikki", city,
	//)
	res, err := http.Get("https://raflaamo.fi/fi/ravintolat/rovaniemi/kaikki")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Link__Anchor-sc-z5lyog-0 hFnksm ButtonStyles__LinkButton-sc-l1rosc-1 dBCuXR ListItemStyles__ButtonLink-sc-1xaojw6-6 kvKMVN
	// Find the review items
	doc.Find(".ListItemStyles__RestaurantListItemWrapper-sc-1xaojw6-15.jDtLYn").Each(func(i int, s *goquery.Selection) {
		// Here we are getting the name of the restaurant. For some reason, there is two classes with the same name
		// so to avoid duplicates we only get the first node.
		name := s.Find(".Link__Anchor-sc-z5lyog-0.hFnksm.ListItemStyles__RestaurantNameLink-sc-1xaojw6-11.hkqmvx").First().Text()

		linkClass := s.Find(".Link__Anchor-sc-z5lyog-0.hFnksm.ButtonStyles__LinkButton-sc-l1rosc-1.dBCuXR.ListItemStyles__ButtonLink-sc-1xaojw6-6.deMCLr.anchor")
		link, ok := linkClass.Attr("href")
		if ok {
			restaurantStruct := restaurant{
				reservingUrl: link,
				name:         name,
				canReserve:   true,
			}
			fmt.Println("Name of restaurant is: ", restaurantStruct.name)
			fmt.Println("You're able to book this at: ", restaurantStruct.reservingUrl)
			return
		} else {
			restaurantStruct := restaurant{
				reservingUrl: "",
				name:         name,
				canReserve:   false,
			}
			fmt.Println("You can not book", restaurantStruct.name)
			return
		}
	})
}

type structure struct {
	operationName string
	variables     map[string]interface{}
	query         string
}

func getPayload() structure {
	data := structure{
		operationName: "getRestaurantsByLocation",
		variables: map[string]interface{}{
			"first": 470,
			"input": map[string]interface{}{
				"restaurantType": "ALL",
				"locationName":   "Helsinki",
				"feature": map[string]interface{}{
					"rentableVenues": false,
				},
			},
			"after": "eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LCJvIjoyNC45MzI1OH19",
		},
		query: "fragment Locales on LocalizedString {fi_FI\n }\n\nfragment Restaurant on Restaurant {\n  id\n  name {\n    ...Locales\n    }\n  urlPath {\n    ...Locales\n     }\n    address {\n    municipality {\n      ...Locales\n       }\n        street {\n      ...Locales\n       }\n       zipCode\n     }\n    features {\n    accessible\n     }\n  openingTime {\n    restaurantTime {\n      ranges {\n        start\n        end\n        endNextDay\n         }\n             }\n    kitchenTime {\n      ranges {\n        start\n        end\n        endNextDay\n              }\n             }\n    }\n  links {\n    tableReservationLocalized {\n      ...Locales\n        }\n    homepageLocalized {\n      ...Locales\n          }\n   }\n     \n}\n\nquery getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {\n  listRestaurantsByLocation(first: $first, after: $after, input: $input) {\n    totalCount\n      edges {\n      ...Restaurant\n        }\n     }\n}",
	}
	return data
}

func getRestaurants() {
	data := getPayload()
	dataEncoded, _ := json.Marshal(data)
	resp, err := http.Post("https://api.raflaamo.fi/query", "application/json", bytes.NewBuffer(dataEncoded))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}

//    const data = {
//        operationName: "getRestaurantsByLocation",
//        variables: {
//            first: 470,
//            input: {
//                restaurantType: "ALL",
//                locationName: "Helsinki",
//                feature: {
//                    rentableVenues: false
//                }
//            },
//            after: "eyJmIjowLCJnIjp7ImEiOjYwLjE3MTE2LCJvIjoyNC45MzI1OH19"
//        },
//        query: `fragment Locales on LocalizedString {fi_FI\n }\n\nfragment Restaurant on Restaurant {\n  id\n  name {\n    ...Locales\n    }\n  urlPath {\n    ...Locales\n     }\n    address {\n    municipality {\n      ...Locales\n       }\n        street {\n      ...Locales\n       }\n       zipCode\n     }\n    features {\n    accessible\n     }\n  openingTime {\n    restaurantTime {\n      ranges {\n        start\n        end\n        endNextDay\n         }\n             }\n    kitchenTime {\n      ranges {\n        start\n        end\n        endNextDay\n              }\n             }\n    }\n  links {\n    tableReservationLocalized {\n      ...Locales\n        }\n    homepageLocalized {\n      ...Locales\n          }\n   }\n     \n}\n\nquery getRestaurantsByLocation($first: Int, $after: String, $input: ListRestaurantsByLocationInput!) {\n  listRestaurantsByLocation(first: $first, after: $after, input: $input) {\n    totalCount\n      edges {\n      ...Restaurant\n        }\n     }\n}`
//    }
//
//    // TODO: figure out a way to get through cors.
//    fetch("https://api.raflaamo.fi/query", {
//        method: 'POST',
//        headers: {
//            "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:103.0) Gecko/20100101 Firefox/103.0",
//            "content-type": "application/json",
//            "client_id": "jNAWMvWD9rp637RaR",
//        },
//        body: JSON.stringify(data)
//    }).then(res => console.log(res))
//}
