package main

import (
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

	//Link__Anchor-sc-z5lyog-0 hFnksm ButtonStyles__LinkButton-sc-l1rosc-1 dBCuXR ListItemStyles__ButtonLink-sc-1xaojw6-6 kvKMVN
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
