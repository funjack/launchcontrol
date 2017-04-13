package kiiroo

import "time"

var (
	limiterTime  = time.Millisecond * 151 // maximum event rate
	upPosition   = 95                     // up position %
	downPosition = 5                      // down position %
)

type togglePosition int

func (a *togglePosition) Toggle() int {
	var new togglePosition
	if *a > togglePosition(downPosition) {
		new = togglePosition(downPosition)
	} else {
		new = togglePosition(upPosition)
	}
	*a = new
	return int(*a)
}
