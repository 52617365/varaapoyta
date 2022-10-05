package raflaamoGraphApi

type parsedGraphData struct {
	Name      string                `json:"name"`
	Intervals *[]parsedIntervalData `json:"intervals"` // were only interested in the first index.
	Id        int                   `json:"id"`
}

type parsedIntervalData struct {
	From  int64  `json:"from"`  // From is a unix timestamp in ms.
	To    int64  `json:"to"`    // To is a unix timestamp in ms.
	Color string `json:"color"` // Optional field, we can match this to see if the restaurant has available tables. (if not nil it does.)
}