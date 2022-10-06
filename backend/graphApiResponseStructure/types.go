package graphApiResponseStructure

type ParsedGraphData struct {
	Name      string                `json:"name"`
	Intervals *[]ParsedIntervalData `json:"intervals"` // were only interested in the first index.
	Id        int                   `json:"id"`
}

type ParsedIntervalData struct {
	From  int64  `json:"from"`  // From is a unix timestamp in ms.
	To    int64  `json:"to"`    // To is a unix timestamp in ms.
	Color string `json:"color"` // Optional field, we can match this to see if the restaurant has available tables. (if not nil it does.)
}
