Launchcontrol VLC extension
===========================

**Notice:** please use Kodi if possible, the VLC extension has less features
and might be a little buggy.

VLC Lua extension that can send scripts and commands to a Launchcontrol server.

What is working:
 - Script loading from `file://` sources on input change (eg open video, next
   item in playlist)
 - Handling pause/resume
 - Seeking/jumping to different time when playback is paused

Gotchas:
 - Seeking/jumping while playing does not work. You will have to pause/resume
   to sync a new position to Launchcontrol.
 - Loading script from other sources like `http://`.

Install
-------

Place `launchcontrol.lua` in the `extensions` directory, create it if it doesn't
already exist:
- Linux: `~/.local/share/vlc/lua/extensions/`
- Mac: `/Users/<NAME>/Library/Application Support/org.videolan.vlc/lua/extensions/`
- Windows: `C:\Users\<NAME>\AppData\Roaming\vlc\lua\extensions`

Usage
-----

- Run Launchcontrol server on `http://localhost:6969`
- In VLC enable extension `Launchcontrol` in the `View` menu.
- Play a video that has a paired `.kiiroo` file.
