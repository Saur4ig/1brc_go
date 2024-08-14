package types

type TempData struct {
	City        string
	Temperature float32
}

type Result struct {
	Min     float32
	Max     float32
	Mean    float32
	Visited int
}
