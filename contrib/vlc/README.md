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
    - Make sure you can connect to the url before using the extension, VLC might
      hang if the server is not reachable.
- In VLC enable extension `Launchcontrol` in the `View` menu.
- Play a video that has a paired `.kiiroo` file.

Configuration
-------------

The extension does not (yet) store the config, until it does you have to edit
the `launchcontrol.lua` file in order to make the changes persistent:

```lua
--[[ Config ]]--
local clientConfig = {
  url = "http://127.0.0.1:6969",
  latency = 0,
  positionMin = 0,
  positionMax = 100,
  speedMin = 20,
  speedMax = 100,
}
```
