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

// Actions converts Kiiroo events into Actions that can be send to a Launch.
func (es Events) Actions() []TimedAction {
	var (
		prevEvent           Event
		prevAction          TimedAction
		actionCount         int
		position            togglePosition
		speed, limitedSpeed int
	)

	// count(actions) <= count(events)
	actions := make([]TimedAction, len(es))

	for _, e := range es {
		// Move only when value is different from previous event
		if e.Value == prevEvent.Value {
			continue
		}

		// Calculate speed for non-limited actions
		speed = calcSpeed(e.Time-prevEvent.Time, speed)

		// Event came in earlier than the limit allows
		if e.Time-prevAction.Time <= limiterTime {
			if limitedSpeed == 0 {
				// Set the speed that will be used while lmiter
				// is active
				limitedSpeed = speed
			}
			// Only trigger if the last executed action longer ago
			// than the current event.
			if prevAction.Time < e.Time {
				// Create a new event in at rate limit in the
				// future.
				actions[actionCount].Position = position.Toggle()
				actions[actionCount].Speed = limitedSpeed
				actions[actionCount].Time = prevAction.Time + limiterTime
				prevAction = actions[actionCount]
				actionCount++
			}
		} else {
			limitedSpeed = 0 // Reset rate limit
			actions[actionCount].Position = position.Toggle()
			actions[actionCount].Speed = speed
			actions[actionCount].Time = e.Time
			prevAction = actions[actionCount]
			actionCount++
		}

		prevEvent = e
	}
	return actions[:actionCount]
}

// calcSpeed calculates the new speed based on the time difference and speed of
// the previous event.
func calcSpeed(t time.Duration, prvSpd int) int {
	if t >= time.Second*2 {
		return 50
	} else if t >= time.Second {
		return 20
	}

	rawSpd := int(100 - int64(t)*110/(1000*1000*1000))
	if rawSpd < 0 {
		rawSpd = 0
	}

	// Go faster
	if rawSpd > prvSpd {
		return prvSpd + (rawSpd-prvSpd)/6
	}

	// Go slower
	spd := prvSpd - rawSpd/2
	if spd < 20 {
		return 20
	}
	return spd
}
