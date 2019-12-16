package collector

import (
	"log"
)

// TimeSeries represents Linode's (Go) SDK's untyped [][]float64
type TimeSeries struct {
	// 1st dimension represents the series and appears to always contain 64 values
	// 2nd dimension (!) [0] = epoch ms; [1] = value
	data [][]float64

	// Pointers to float64 to permit choice between nil or set
	// If unset, compute is run and this calculates all of them at once
	avg *float64
	min *float64
	max *float64
	sum *float64
}

func NewTimeSeries(data [][]float64) *TimeSeries {
	log.Printf("[NewTimeSeries] length: %d", len(data))
	return &TimeSeries{
		data: data,
	}
}
func (ts *TimeSeries) compute() {
	log.Printf("[TimeSeries:compute] Entered")
	// Initiailize statistics
	ts.avg = new(float64)
	ts.min = new(float64)
	ts.max = new(float64)
	ts.sum = new(float64)
	if len(ts.data) > 0 {
		for _, d := range ts.data {
			*ts.sum += d[1]
			if d[1] < *ts.min {
				*ts.min = d[1]
			}
			if d[1] > *ts.max {
				*ts.max = d[1]
			}
		}
		*ts.avg = *ts.sum / float64(len(ts.data))
	}
}
func (ts *TimeSeries) Avg() float64 {
	if ts.avg == nil {
		ts.compute()
	}
	return *ts.avg
}
func (ts *TimeSeries) Min() float64 {
	if ts.min == nil {
		ts.compute()
	}
	return *ts.min
}
func (ts *TimeSeries) Max() float64 {
	if ts.max == nil {
		ts.compute()
	}
	return *ts.max
}
func (ts *TimeSeries) Sum() float64 {
	if ts.sum == nil {
		ts.compute()
	}
	return *ts.sum
}

// func (ts *TimeSeries) Percentile(i uint8) float64 {}
