package rolling

type Metric interface {
	Add(float64)
	Value() float64
}

type Aggregation interface {
	Min() float64
	Max() float64
	Avg() float64
	Sum() float64
}
