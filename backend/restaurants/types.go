package restaurants

import (
	"regexp"
)

var regexToMatchRestaurantId = regexp.MustCompile(`[^fi/]\d+`)
var regexToMatchTime = regexp.MustCompile(`\d{2}:\d{2}`)
var regexToMatchDate = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
