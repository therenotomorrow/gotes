package entities

type Metrics struct {
	Requests int64
	Errors   int64
}

type Statistics struct {
	Total     int64
	ErrorRate float64
}
