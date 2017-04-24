import xbmc
import xbmcaddon
import xbmcvfs
import launchcontrol

__addon__        = xbmcaddon.Addon()
__addonid__      = __addon__.getAddonInfo('id')
__addonname__    = __addon__.getAddonInfo('name')
__addonversion__ = __addon__.getAddonInfo('version')

class PlayerMonitor(xbmc.Player) :
    """PlayerMonitor sends launch commands based on Player events."""

    def __init__(self):
        self.loadConfig()
        xbmc.Player.__init__(self)

    def loadConfig(self):
        """loadConfig configures the monitor with the Kodi addon settings."""
        self._launch = launchcontrol.Client(
                url=__addon__.getSetting("address"),
                latency=__addon__.getSetting("latency"),
                positionmin=__addon__.getSetting("positionmin"),
                positionmax=__addon__.getSetting("positionmax"),
                speedmin=__addon__.getSetting("speedmin"),
                speedmax=__addon__.getSetting("speedmax"))

    def onPlayBackStarted(self):
        try:
            fileName = self.getPlayingFile()
            data, mediaType = ReadScript(fileName)
            if mediaType != "":
                self._launch.Play(data, mediaType)

                # I hate this sleep but AFAIK there is no way to figure out
                # where playback is resumed, and it takes a while for the
                # player to actually update getTime to the resumed location.
                waitTimeSec = 1
                for i in xrange(3):
                    eventTime = self.getTime()
                    xbmc.sleep(waitTimeSec*1000)
                    # Skip if the difference is larger then waitTimeSec
                    if round(self.getTime()) > round(eventTime+waitTimeSec*i):
                        self.SkipToCurrentTime()
                        break
        except launchcontrol.NotNowException:
            pass
        except launchcontrol.NotSupportedException:
            log("Script for \"%s\" with mediaType \"%s\" is not supported" %
                    (fileName, mediaType), xbmc.LOGNOTICE)
        except Exception as e:
            log("Unhandled exception in onPlayBackStarted: %s" % e,
                    xbmc.LOGERROR)

    def onPlayBackEnded(self):
        try:
            self._launch.Stop()
        except launchcontrol.NotNowException:
            pass
        except Exception as e:
            log("Unhandled exception in onPlayBackEnded: %s" % e,
                    xbmc.LOGERROR)

    def onPlayBackPaused(self):
        try:
            self._launch.Pause()
        except launchcontrol.NotNowException:
            pass
        except Exception as e:
            log("Unhandled exception in onPlayBackPaused: %s" % e,
                    xbmc.LOGERROR)

    def onPlayBackResumed(self):
        try:
            self._launch.Resume()
        except launchcontrol.NotNowException:
            pass
        except Exception as e:
            log("Unhandled exception in onPlayBackResumed: %s" % e,
                    xbmc.LOGERROR)

    def onPlayBackStopped(self):
        try:
            self._launch.Stop()
        except launchcontrol.NotNowException:
            pass
        except Exception as e:
            log("Unhandled exception in onPlayBackStopped: %s" % e,
                    xbmc.LOGERROR)

    def onPlayBackSeek(self, time, seekOffset):
        try:
            self._launch.Skip(time)
        except launchcontrol.NotNowException:
            pass
        except Exception as e:
            log("Unhandled exception in onPlayBackSeek: %s" % e,
                    xbmc.LOGERROR)

    def onPlayBackSeekChapter(self, chapter):
        try:
            self.SkipToCurrentTime()
        except launchcontrol.NotNowException:
            pass
        except Exception as e:
            log("Unhandled exception in onPlayBackSeekChapter: %s" % e,
                    xbmc.LOGERROR)

    def onPlayBackSpeedChanged(self, speed):
        try:
            if speed == 1:
                self.SkipToCurrentTime()
                self._launch.Resume()
            else:
                self._launch.Pause()
        except launchcontrol.NotNowException:
            pass
        except Exception as e:
            log("Unhandled exception in onPlayBackSpeedChanged: %s" % e,
                    xbmc.LOGERROR)

    def SkipToCurrentTime(self):
        """SkipToCurrentTime skips the script to the current player time."""
        # Make sure the player state catched up with the last request.
        xbmc.sleep(1000)
        self._launch.Skip(self.getTime() * 1000)

def ReadScript(filename):
    """ReadScript uses Kodi's vfs to detect and read a scrip file.

    Returns:
        (data, mediaType): Tuple containing the with raw script data and its
            mediatype.
    """
    # filename can be url, path or something in between, that's why the string
    # rsplit function is used instead of url/path specific functions.
    dotSplit = filename.rsplit(".")
    if len(dotSplit) > 1:
        del dotSplit[-1]
    baseFilename = ".".join(dotSplit)

    for scripttype in launchcontrol.scripttypes:
        for extention in scripttype["extensions"]:
            scriptFile = baseFilename + "." + extention
        if xbmcvfs.exists(scriptFile):
            f = xbmcvfs.File(scriptFile)
            data = f.read()
            f.close()
            return (data, scripttype["mediaType"])
    return ("","")

class SettingsMonitor(xbmc.Monitor):
    """SettingsMonitor reloads the player monitors config when the addon settigs change."""

    def __init__(self, player):
        self._player = player
        xbmc.Monitor.__init__(self)

    def onSettingsChanged(self):
        self._player.loadConfig()

def log(txt, level=xbmc.LOGDEBUG):
    """Log to the XBMC/Kodi logfile.
    
    Args:
        txt: Logmessage
        level: Loglevel xbmc style 0-7, ranging from debug...none.
    """
    message = '%s: %s' % (__addonname__, txt.encode('ascii', 'ignore'))
    xbmc.log(msg=message, level=level)

if __name__ == '__main__':
    log('Version %s started' % __addonversion__, xbmc.LOGNOTICE)
    # FIXME: sometimes the addon cannot be uninstalled, it seems the
    # PlayerMonitor can't always be stopped after it has been used. A Kodi
    # restart is then required. :-(
    player = PlayerMonitor()
    monitor = SettingsMonitor(player)
    while not monitor.abortRequested():
        if monitor.waitForAbort(10):
            break
    del monitor
    del player
    log('Version %s stopped' % __addonversion__, xbmc.LOGNOTICE)
