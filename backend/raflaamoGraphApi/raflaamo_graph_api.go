package raflaamoGraphApi

import (
	"backend/timeUtils"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"
)

type RaflaamoGraphApi struct {
	httpClient *http.Client
	requestUrl string
}
type TimeSlot struct {
	currentTimeInUnix          int64
	reservationTimeUnix        int64 // TODO: give better name
	graphEndTimeUnix           int64 // TODO: give better name
	restaurantStartingTimeUnix int64 // TODO: figure out if we even need this
	restaurantClosingTimeUnix  int64 // TODO: figure out if we even need this
}

// We determine if there is a timeUtils slot with open tables by looking at the "color" field in the response.
// The color field will contain "transparent" if it does not contain a graph (open times), else it contains nil (meaning there are open tables)
func timeSlotDoesNotContainOpenTables(data *parsed_graph_data) bool {
	return (*data.Intervals)[0].Color == "transparent"
}
func (graphApi *RaflaamoGraphApi) constructPayload(idFromReservationPageUrl string, currentDate string, time *timeUtils.CoveredTimes, amountOfEaters int) {
	timeSlotString := timeUtils.GetStringTimeFromUnix(time.time)

	// replacing the 17(:)00 to match the format in url.
	timeSlotString = strings.Replace(timeSlotString, ":", "", -1)
	requestUrl := fmt.Sprintf("https://s-varaukset.fi/api/recommendations/slot/%s/%s/%s/%d", idFromReservationPageUrl, currentDate, timeSlotString, amountOfEaters)
	graphApi.requestUrl = requestUrl
}

// Gets timeslots from raflaamo API that is responsible for returning graph data.
// In the end, instead of drawing a graph with it, we convert it into timeUtils to determine which table is open or not.
// This one sends requests, so we use goroutines in it.
func (graphApi *RaflaamoGraphApi) interactWithApi() (*parsedGraphData, error) {
	client := graphApi.httpClient

	r, err := http.NewRequest("GET", graphApi.requestUrl, nil)

	if err != nil {
		return nil, errors.New("error connecting to api")
	}

	r.Header.Add("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	res, err := client.Do(r)

	// Will throw if we call deserialize_graph_response with a status code other than 200, so we handle it here.
	if err != nil || res.StatusCode != 200 {
		return nil, errors.New("error connecting to api")
	}

	deserializedGraphData, err := deserializeGraphResponse(res)

	// most likely won't jump into this branch but check regardless.
	if err != nil {
		return nil, errors.New("error deserializing")
	}

	if timeSlotDoesNotContainOpenTables(deserializedGraphData) {
		return nil, errors.New("no open tables found")
	}
	// Adding timezone difference into the unix timeUtils. (three hours).
	graphEndUnix := (*deserializedGraphData.Intervals)[0].To
	// Adding 10800000(ms) to the timeUtils to match UTC +2 or +3 (finnish timeUtils) (10800000 ms corresponds to 3h)
	// because graph unix timeUtils fields "to" and "from" come in UTC +0

	graphEndUnix += 10800000
	return deserializedGraphData, nil
}

// This function interacts with the raflaamo graph API and returns the timeUtils slots that we get from that API.
// Function will return error if the provided timestamps were in invalid form (current timeUtils is bigger or equal to the last possible timeUtils interval returned from the API) it's an error because if that is the case, we don't have any times to check.
// TODO: This does not work correctly.
func (graphData *parsedGraphData) extractAvailableTimeIntervalsFromResponse(currentTime *dateAndTime, kitchenOfficeHours *restaurantTime, allReservationTimes []int64) ([]string, error) {
	graphEndUnix := (*graphData.Intervals)[0].To
	if restaurantAlreadyClosed(currentTime, graphEndUnix) {
		return nil, errors.New("restaurant is already closed")
	}

	restaurantAvailableTimeSlots := make([]string, 0, len(allReservationTimes))

	// Here we capture all the available timeUtils intervals into an array by matching against all the possible timeUtils intervals.
	for _, reservationTime := range allReservationTimes {
		// We check if the timestamps are valid here.
		if validGraphTimeSlot(reservationTime, currentTime.time, graphEndUnix) && timeSlotInRestaurantOpeningHours(reservationTime, kitchenOfficeHours.opening, kitchenOfficeHours.closing) {
			slot := getStringTimeFromUnix(reservationTime)

			// Avoiding storing duplicate timeUtils slots because without this, it will.
			if !slices.Contains(restaurantAvailableTimeSlots, slot) {
				restaurantAvailableTimeSlots = append(restaurantAvailableTimeSlots, slot)
			}
		}
	}
	if len(restaurantAvailableTimeSlots) == 0 {
		return nil, errors.New("no timeUtils slots found for restaurant")
	}
	return restaurantAvailableTimeSlots, nil
}

func restaurantAlreadyClosed(currentTime *interface{}, graphEndUnix int64) bool {
	return currentTime.time >= graphEndUnix
}

// Checks to see if the reservation timeUtils checked is larger than the current timeUtils and smaller or equal to the last possible timeUtils to reserve (graph_end_unix)
func (timeSlot *TimeSlot) validGraphTimeSlot( /*reservationTime int64, currentTime int64, graphEndUnix int64*/ ) bool {
	if timeSlot.reservationTimeUnix > timeSlot.currentTimeInUnix && timeSlot.reservationTimeUnix <= timeSlot.graphEndTimeUnix {
		return true
	}
	return false
}

func (timeSlot *TimeSlot) timeSlotInRestaurantOpeningHours( /*reservationTime int64, restaurantStartingTimeUnix int64, restaurantClosingTimeUnix int64*/ ) bool {
	if timeSlot.reservationTimeUnix > timeSlot.restaurantStartingTimeUnix && timeSlot.reservationTimeUnix <= timeSlot.restaurantClosingTimeUnix {
		return true
	}
	return false
}

// Returns all possible timeUtils intervals that can be reserved in the raflaamo reservation page.
// 11:00, 11:15, 11:30 and so on.
func getAllRaflaamoTimeIntervals() []int64 {
	allTimes := make([]int64, 0, 96)
	for hour := 0; hour < 24; hour++ {
		if hour < 10 {
			formattedTimeSlotOne := timeUtils.ConvertStringTimeToUnix(fmt.Sprintf("0%d00", hour))
			formattedTimeSlotTwo := timeUtils.ConvertStringTimeToUnix(fmt.Sprintf("0%d15", hour))
			formattedTimeSlotThree := timeUtils.ConvertStringTimeToUnix(fmt.Sprintf("0%d30", hour))
			formattedTimeSlotFour := timeUtils.ConvertStringTimeToUnix(fmt.Sprintf("0%d45", hour))
			allTimes = append(allTimes, formattedTimeSlotOne)
			allTimes = append(allTimes, formattedTimeSlotTwo)
			allTimes = append(allTimes, formattedTimeSlotThree)
			allTimes = append(allTimes, formattedTimeSlotFour)
		}
		if hour >= 10 {
			formattedTimeSlotOne := timeUtils.ConvertStringTimeToUnix(fmt.Sprintf("%d00", hour))
			formattedTimeSlotTwo := timeUtils.ConvertStringTimeToUnix(fmt.Sprintf("%d15", hour))
			formattedTimeSlotThree := timeUtils.ConvertStringTimeToUnix(fmt.Sprintf("%d30", hour))
			formattedTimeSlotFour := timeUtils.ConvertStringTimeToUnix(fmt.Sprintf("%d45", hour))
			allTimes = append(allTimes, formattedTimeSlotOne)
			allTimes = append(allTimes, formattedTimeSlotTwo)
			allTimes = append(allTimes, formattedTimeSlotThree)
			allTimes = append(allTimes, formattedTimeSlotFour)
		}
	}
	return allTimes
}
