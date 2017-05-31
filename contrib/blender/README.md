Funscripting Blender Addon
==========================

Funscripting is an addon for Blender to 'script' Launch movements for a video.
With this addon you can create your own scripts (in Funscript format) that can
be played back by Launchcontrol for any movie you like.

[Blender](https://www.blender.org/) is a very powerful free and open-source 3D
editor suite, that comes with a video editor. Scrubbing through a video is very
fast and easy. 

This addon consists of a panel for the sequencer with percentage buttons
representing Launch positions. This will insert the buttons value in a custom
property on the selected strip and marks it as a keyframe. The export button
will save this as a Funscript that can be played back using Launchcontrol.

Installation
------------

1. Download the `funscripting.py` file
2. Start Blender
3. Open User Preferences (Ctrl+Alt+U)
4. Click the Add-ons tab
5. Click Install from File.
6. Select the downloaded `funscripting.py`
7. Mark the checkbox of `Funscripting Addon`

If the above steps are not clear, this
[stackexchange](https://blender.stackexchange.com/questions/1688/installing-an-addon)
answer explains the process with screenshots.

Usage
-----

Don't Panic. Blender can be very overwhelming but you don't need to know what
everything does to create a Funscript :)

### Switch to video editing layout
TODO
### Import a movie file
TODO
### Script movie
TODO
### Export Funscript
TODO

Tips
----

When scripting, keep in mind that the Launch can only move up and down. That
means the next position is **relative** to the previous one. This works great
for penetration moves, but get tricky when you try to script moments like
touching and licking. Then end and begin position do not necessarily match. So
you need to be creative if you want to script that :)
