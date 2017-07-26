package funscript

import (
	"testing"
	"time"
)

type Measurement struct {
	Distance int
	Speed    int
	Duration time.Duration
}

var MeasurementsTestTable = []Measurement{
	{
		// Case 0
		Distance: 90,
		Speed:    20,
		Duration: time.Millisecond * 850,
	},
	{
		// Case 1
		Distance: 90,
		Speed:    30,
		Duration: time.Millisecond * 600,
	},
	{
		// Case 2
		Distance: 90,
		Speed:    40,
		Duration: time.Millisecond * 450,
	},
	{
		// Case 3
		Distance: 90,
		Speed:    50,
		Duration: time.Millisecond * 375,
	},
	{
		// Case 4
		Distance: 90,
		Speed:    60,
		Duration: time.Millisecond * 300,
	},
	{
		// Case 5
		Distance: 90,
		Speed:    70,
		Duration: time.Millisecond * 250,
	},
	{
		// Case 6
		Distance: 90,
		Speed:    80,
		Duration: time.Millisecond * 225,
	},
	{
		// Case 7
		Distance: 45,
		Speed:    20,
		Duration: time.Millisecond * 425,
	},
	{
		// Case 8
		Distance: 45,
		Speed:    30,
		Duration: time.Millisecond * 300,
	},
	{
		// Case 9
		Distance: 45,
		Speed:    40,
		Duration: time.Millisecond * 225,
	},
	{
		// Case 10
		Distance: 45,
		Speed:    50,
		Duration: time.Millisecond * 190,
	},
	{
		// Case 11
		Distance: 45,
		Speed:    60,
		Duration: time.Millisecond * 165,
	},
	{
		// Case 12
		Distance: 0,
		Speed:    0,
		Duration: 0,
	},
	{
		// Case 13
		Distance: 100,
		Speed:    100,
		Duration: time.Millisecond * 210,
	},
}

type timeTolerance time.Duration

func (p timeTolerance) roughlyEqual(a time.Duration, b time.Duration) bool {
	if a > b+time.Duration(p) {
		return false
	}
	if a < b-time.Duration(p) {
		return false
	}
	return true
}

type precentTolerance int

func (p precentTolerance) roughlyEqual(a int, b int) bool {
	if a > b+int(p) {
		return false
	}
	if a < b-int(p) {
		return false
	}
	return true
}

var defaultTimeTolerance = timeTolerance(time.Millisecond * 25)
var defaultPercentTolerance = precentTolerance(5)

func TestSpeed(t *testing.T) {
	for i, c := range MeasurementsTestTable {
		speed := Speed(c.Distance, c.Duration)
		if !defaultPercentTolerance.roughlyEqual(c.Speed, speed) {
			t.Errorf("case %d: not (roughly) equal, want %d, got %d",
				i, c.Speed, speed)
		}
	}
}

func TestDuration(t *testing.T) {
	for i, c := range MeasurementsTestTable {
		duration := Duration(c.Distance, c.Speed)
		if !defaultTimeTolerance.roughlyEqual(c.Duration, duration) {
			t.Errorf("case %d: not (roughly) equal, want %s, got %s",
				i, c.Duration, duration)
		}
	}
}

func TestDistance(t *testing.T) {
	for i, c := range MeasurementsTestTable {
		distance := Distance(c.Speed, c.Duration)
		if !defaultPercentTolerance.roughlyEqual(c.Distance, distance) {
			t.Errorf("case %d: not (roughly) equal, want %d, got %d",
				i, c.Distance, distance)
		}
	}
}
