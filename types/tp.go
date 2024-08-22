package types

type TempData struct {
	City        string
	Temperature float64
}

type Result struct {
	Min     float64
	Max     float64
	Sum     float64
	Station string
	Visited int
}

type Temperature struct {
	Station string
	Temp    float64
}

type Res struct {
	Min     int64
	Max     int64
	Sum     int64
	Station string
	Visited int32
}
