package restaurants

import (
	"regexp"
)

var regexToMatchRestaurantId = regexp.MustCompile(`[^fi/]\d+`)
var RegexToMatchTime = regexp.MustCompile(`\d{2}:\d{2}`)
var RegexToMatchDate = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
