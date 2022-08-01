package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// TODO: make post request work again.
// TODO: pass city here.
func getRestaurants() interface{} {
	data := generate_and_serialize_payload()
	r, err := http.NewRequest("POST", "https://api.raflaamo.fi/query", bytes.NewBuffer(data))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("client_id", "jNAWMvWD9rp637RaR")
	r.Header.Add("Origin", "https://raflaamo.fi")
	r.Header.Add("Referer", "https://raflaamo.fi")
	r.Header.Add("Sec-Fetch-Site", "same-site")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.63 Safari/537.36")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Body)
	decoded := deserialize_response(&res)
	defer res.Body.Close()

	fmt.Println(decoded)
	return decoded
}
