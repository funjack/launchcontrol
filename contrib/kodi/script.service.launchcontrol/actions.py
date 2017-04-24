import launchcontrol
import xbmcaddon
import xbmcgui

__addon__        = xbmcaddon.Addon()
__addonid__      = __addon__.getAddonInfo('id')
__addonname__    = __addon__.getAddonInfo('name')
__addonversion__ = __addon__.getAddonInfo('version')

def TestConnection():
    """TestConnection sends a small script to Launchcontrol."""
    progressDialog = xbmcgui.DialogProgress()
    progressDialog.create("Launchcontrol", "Testing connection...")
    try:
        l = launchcontrol.Client(
                url=__addon__.getSetting("address"),
                positionmin=__addon__.getSetting("positionmin"),
                positionmax=__addon__.getSetting("positionmax"),
                speedmin=__addon__.getSetting("speedmin"),
                speedmax=__addon__.getSetting("speedmax"))
        l.Play("{0.50:4,1.00:0,2.50:4,3.00:0}", "x-text/kiiroo")
    except Exception as e:
        progressDialog.close()
        xbmcgui.Dialog().ok("Launchcontrol connection test" , "Failed:", e.message)
    else:
        progressDialog.close()
        xbmcgui.Dialog().ok("Launchcontrol connection test" , "Success")

if __name__ == "__main__":
    TestConnection()
