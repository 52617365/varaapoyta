package main

import (
	"fmt"
	"net/http"
)

// TODO: get data first and then make the endpoints.

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func hellobroski(w http.ResponseWriter, r *http.Request) {
	name := 1 + 1
	fmt.Println(name)
	fmt.Println("serious gang shit")
}

func main() {
}
