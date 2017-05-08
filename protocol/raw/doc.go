/*
Package raw manages the raw launch protocol.

The raw format are timed position,speed commands using the same values as the
Launch's BLE protocol.

The file format is a JSON encoded array of timed actions:

	[
	  {
	    "at": <time>,
	    "pos": <position>,
	    "spd": <speed>
	  },
	  ...
	]

	time    : integer, time in ms when the action is executed
	position: integer, position to move to 0-100 (bottom ... top)
	speed   : integer, speed to move at 0-100 (slow ... fast)

The raw format uses the same values as the BLE protocol, giving the script the
biggest amount of control but with great power comes great responsibility ;-)

Tips

Canceling previous actions before they finish will often result in a more
smooth transition between the moves.

At speed 20% a full stroke takes ~900ms

*/
package raw
