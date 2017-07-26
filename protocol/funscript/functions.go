package funscript

import (
	"errors"
	"math"
	"time"
)

var ErrOutOfRange = errors.New("out of range")

// Speed returns the speed (in percent) to move the given distance (in percent)
// in the given duration.
func Speed(dist int, dur time.Duration) (speed int) {
	if dist <= 0 {
		return 0
	} else if dist > 100 {
		return 100
	}
	mil := float64(dur.Nanoseconds()/1e6) * 90 / float64(dist)
	speed = int(25000 * math.Pow(float64(mil), -1.05))
	return speed
}

// Duration returns the time it will take to move the given distance (in
// percent) at the given speed (in percent.)
func Duration(dist int, spd int) (dur time.Duration) {
	if dist <= 0 {
		return 0
	}
	mil := math.Pow(float64(spd)/25000, -0.95)
	dur = time.Duration(mil/(90/float64(dist))) * time.Millisecond
	return dur
}

// Distance returns the distance (in percent) that will be moved with the given
// speed (in percent) and duration.
func Distance(spd int, dur time.Duration) (dist int) {
	if spd <= 0 {
		return 0
	}
	mil := math.Pow(float64(spd)/25000, -0.95)
	diff := mil - float64(dur.Nanoseconds()/1e6)
	dist = 90 - int(diff/mil*90)
	return dist
}
