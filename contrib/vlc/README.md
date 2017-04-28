Launchcontrol VLC extension
===========================

**Notice:** please use Kodi if possible, the VLC extension has less features
and might be a little buggy.

VLC Lua extension that can send scripts and commands to a Launchcontrol server.

What is working:
 - Script loading on input change (eg open video, next item in playlist)
 - Handling pause/resume
 - Seeking/jumping to different time when playback is paused

Gotchas:
 - Displays errors for all the scripts it could not find. You can get rid of
   these errors by selecting `Hide future errors`.
 - Seeking/jumping while playing does not work. You will have to pause/resume
   to sync a new position to Launchcontrol.

Install
-------

Place `launchcontrol.lua` in the `extensions` directory, create it if it doesn't
already exist:
- Linux: `~/.local/share/vlc/lua/extensions/`
- MacOS: `/Users/<NAME>/Library/Application Support/org.videolan.vlc/lua/extensions/`
- Windows: `C:\Users\<NAME>\AppData\Roaming\vlc\lua\extensions`

Usage
-----

- Run Launchcontrol server on `http://localhost:6969`
- In VLC enable extension `Launchcontrol` in the `View` menu.
- Play a video that has a paired `.kiiroo` file.
