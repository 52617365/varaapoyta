/*
 * Copyright (c) 2022. Rasmus Mäki
 */

package regex

import "regexp"

var RegexToMatchRestaurantId = regexp.MustCompile(`[^fi/]\d+`)
var RegexToMatchTime = regexp.MustCompile(`\d{2}:\d{2}`)
var RegexToMatchDate = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
var TimeRegex = regexp.MustCompile(`\d{2}:\d{2}`)
