/*
Package kiiroo manages the Kiiroo haptics protocol.

Kiiroo protocol is really wonkey so don't expect anything magical from it.

Protocol format is:

	{<time>:<value>,<time>:<value>,...}

	time : x.xx event time in sec
	value: 0-4 position/intensity of event (but not for the Launch)

Using the Kirroo protocol, the positions sends to the Launch are always 5% and
95% (alternating.) Commands send to the Launch are interrupted/canceled by new
ones, giving only the illusion of precision.

The value is used to determine if the Launch should move. The rate at which
events are received determines the speed parameter of the movement command
issued to the Launch.

Speed algorithm

Move with "default" speed of 50% when it has been >=2sec since last event.

Move at the "slowest" speed of 20% when is has been >=1sec since last event.

Raw speed value between events is 100 - (hundredth sec + 10%) as %. Eg:

	0.20 sec = 100-22 = 78%
	0.50 sec = 100-55 = 45%

If raw speed is bigger then previous speed, then increase speed with 1/6 of the
difference. Eg:

	Previous 50 and Raw 78 = 50+78/6 = 63%

If raw speed is smaller then previous speed, then decrease speed with 1/2 of
the difference. Eg:

	Previous 63 and Raw 45 = 63-45/2 = 41%

Speed can not go below 20%.

Movement algorithm

Stop moving if no signal has been received in 150ms.

Move only when value is different from previous event.

When events are coming in faster then 150ms, send command at 150ms intervals
using last calculated speed. Reset the limiter when the receive window >150ms.

The speed is still calculated when the limiter is active, but the last send
speed is used until the limiter is stopped.

*/
package kiiroo
