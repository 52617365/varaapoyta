package main

// @Problem the getAvailableTables() function gets all tables (1500 sites, x 64 possible time slots).
// This is not sustainable because there are around 160k requests that need to be sent, this is hella slow.

// @Brainstorm to fix problem
// @Solution one
// Get availability when demanded? E.g. user presses button on website, send request with that information.
// In the back end the restaurants' information would be held in a struct, E.g. restaurant url etc.

// @Solution two
// Find a way to reduce the amount of requests. This could be done by finding an alternative to getting open restaurants', E.g. do something with the graph on site.
// Or get better checks into the requests to narrow down the amount of requests (for example, if first request does not have the stuff we want, don't do it on the same page but at a different time).
func main() {
	getAvailableTables()
}
