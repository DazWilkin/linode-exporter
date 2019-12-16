package collector

import (
	"testing"
)

var (
	sampleData = [][]float64{
		[]float64{0, 0},
		[]float64{0, 1},
		[]float64{0, 2},
		[]float64{0, 3},
		[]float64{0, 4},
		[]float64{0, 5},
		[]float64{0, 6},
		[]float64{0, 7},
		[]float64{0, 8},
		[]float64{0, 9},
	}
)

func TestAvg(t *testing.T) {
	ts := NewTimeSeries(sampleData)
	got := ts.Avg()
	want := 4.5
	if got != want {
		t.Errorf("got=%f; want=%f", got, want)
	}
}
func TestMin(t *testing.T) {
	ts := NewTimeSeries(sampleData)
	got := ts.Min()
	want := 0.0
	if got != want {
		t.Errorf("got=%f; want=%f", got, want)
	}
}
func TestMax(t *testing.T) {
	ts := NewTimeSeries(sampleData)
	got := ts.Max()
	want := 9.0
	if got != want {
		t.Errorf("got=%f; want=%f", got, want)
	}
}
func TestSum(t *testing.T) {
	ts := NewTimeSeries(sampleData)
	got := ts.Sum()
	want := 45.0
	if got != want {
		t.Errorf("got=%f; want=%f", got, want)
	}
}
