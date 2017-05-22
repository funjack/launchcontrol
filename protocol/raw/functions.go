package raw

import (
	"math"
	"time"
)

// Speed returns the speed (in percent) to move the given distance (in percent)
// in the given duration.
func Speed(dist int, dur time.Duration) (speed int) {
	mil := dur.Nanoseconds() / 1e6 * int64(90/dist)
	speed = int(25000 * math.Pow(float64(mil), -1.05))
	return speed
}

// Duration returns the time it will take to move the given distance (in
// percent) at the given speed (in percent.)
func Duration(dist int, spd int) (dur time.Duration) {
	mil := int64(math.Pow(float64(spd)/25000, -0.95))
	dur = time.Duration(mil/int64(90/dist)) * time.Millisecond
	return dur
}

// Distance returns the distance (in percent) that will be moved with the given
// speed (in percent) and duration.
func Distance(spd int, dur time.Duration) (dist int) {
	mil := int64(math.Pow(float64(spd)/25000, -0.95))
	diff := mil - (dur.Nanoseconds() / 1e6)
	dist = 90 - int(float64(diff)/float64(mil)*90)
	return dist
}
