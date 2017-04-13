package kiiroo

import "time"

// DefaultAlgorithm implements the Algorithm interface trying to mimik the
// Kiiroo apps.
type DefaultAlgorithm struct{}

// Actions converts Kiiroo events into Actions that can be send to a Launch.
func (da DefaultAlgorithm) Actions(es Events) []TimedAction {
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

	var speed int
	if rawSpd > prvSpd {
		// Go faster
		speed = prvSpd + (rawSpd-prvSpd)/6
	} else {
		// Go slower
		speed = prvSpd - rawSpd/2
	}

	if speed < 20 {
		return 20
	}
	return speed
}
