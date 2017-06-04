/*
Package funscript manages the Funscript format.

Funscript is Funjack's haptic script format. It's basically JSON encoded timed
positions:

	{
		"version": "1.0",
		"inverted": false,
		"range": 90,
		"actions": [
			{"pos": 0, "at": 100},
			{"pos": 100, "at": 500},
			...
		]
	}

	version: funscript version (optional, default="1.0")
	inverted: positions are inverted (0=100,100=0) (optional, default=false)
	range: range of moment to use in percent (0-100) (optional, default=90)
	actions: script for a Launch
	  pos: position in percent (0-100)
	  at : time to be at position in milliseconds

Movement range

Implementations may override the range value specified in the script.

Define min (bottom) and max (top) positions for the strokes. Defaults are:
min=5 and max=95. The values for min and max must:

	(max - min) == range
	(max - min) <= 90

The defaults of 5/95 are based on the reverse engineering efforts of the Kiiroo
protocol. It's not certain if 0/100 are safe to use, so for now better be safe
then sorry.

Speed algorithm

The "Magic Launch Formula" is used to determine the speed of the movement
commands:

	Speed = 25000 * (Duration * 90/Distance)^(-1.05)

	Speed: Launch speed to send to the Launch
	Duration: Time in milliseconds the move should take
	Distance: Amount of space to move in Launch percent (1-100)

Speed must always be a value between 20-80. As slow commands crash the Launch,
and fast commands cause weird noises and will likely damage the Launch.

Movement algorithm

Scripts always starts on time 0 at the bottom (min position) unless inverted is
true, then start at the top (max position). This is also the implicit previous
action for the first action in the script and does not have to be present in
the actions list of the script.

For each action in the script:

If the position value is equal to the previous action position do not move.
Note: do not completely ignore this action but still continue to use it as the
'previous action' for the next one.

When inverted is true, flip the position value (0=100, 25=75, 70=30, etc.)

Scale the position value with the percentage in range and add the min position
value. Eg:

	min=0,  range=50 and pos=75 : 50*75%+0  = 37
	min=10, range=90 and pos=80 : 80*90%+10 = 72+10=82
	min=20, range=80 and pos=30 : 80*30%+20 = 24+20=44

Calculate the duration (ms) and distance (percent) since the previous action.
Eg with a 100% range:

	{"pos": 25, "at": 100}, {"pos": 100, "at": 500}

	time = 500-100 = 400
	distance = (100-25)*100% = 75

Use the duration and distance to calculate the speed.

Send the calculated position with the calculated speed at the time specified by
the previous action.

Limitations

The Launch obviously has it's speed limitations. There are two type of
impossible moves. Too fast, causing the stroke to be shorter then the script
requested. Too slow, causing the move to be finished earlier then requested.
Because Funscript is based on relative moves this can cause portions of a
script the get out of sync. Limiting the range can help making fast moves
possible again, but can cause slower moves to become impossible.

Launchcontrol uses a 100ms threshold when sending message to the Launch,
scripts trying to fire actions faster than that will be out of sync for those
sections.
*/
package funscript
