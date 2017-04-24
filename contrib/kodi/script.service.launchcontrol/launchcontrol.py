"""Launchcontrol client

The module exposes the Launchcontrol API as a Client object.

Data:

scripttypes -- list of a dictionaries containing supported script formats.
"""

import urllib2

scripttypes = [
        {
            "name"      : "kiiroo", 
            "extensions": ["kiiroo"],
            "mediaType" : "x-text/kiiroo",
        },
        {
            "name"      : "realtouch", 
            "extensions": ["realtouch", "ott"],
            "mediaType" : "x-text/realtouch",
        },
        {
            "name"      : "vorze", 
            "extensions": ["vorze", "csv"],
            "mediaType" : "x-text/vorze",
        },
]

class NotNowException(Exception):
    """Raise when an operation it not compatible with current state"""

class NotSupportedException(Exception):
    """Raise when the specified type is not supported"""

class Client() :
    """Client communicates with a Launchcontrol server.

    Args:
        url: Launchcontrol server url
        latency: Time adjustment in milliseconds
        positionmin: Lowest position in percent the Launch should move to
        positionmax: Highest position in percent the Launch should move to
        speedmin: Slowest speed in percent the Launch should move at
        speedmax: Highest speed in percent the Launch should move to
    """

    def __init__ (self, url="http://127.0.0.1:6969", latency=0,
            positionmin=0, positionmax=100, speedmin=20, speedmax=100):
        self._url = url
        self.latency = int(latency)
        self.positionMin = int(positionmin)
        self.positionMax = int(positionmax)
        self.speedMin = int(speedmin)
        self.speedMax = int(speedmax)

    def Play(self, data, mediaType):
        """Play by sending data as specified mediatype.

        Args:
            data: Raw script data in bytes
            mediaType: Mimetype of the script in data

        Raises:
            NotSupportedException: The script and or mediaType is not
                supported.
        """
        if mediaType != "":
            params = [ "latency=%d" % self.latency,
                    "positionmin=%d" % self.positionMin,
                    "positionmax=%d" % self.positionMax,
                    "speedmin=%d" % self.speedMin,
                    "speedmax=%d" % self.speedMax ]
            req = urllib2.Request(self._url+'/v1/play?%s' % "&".join(params),
                    data=data, headers={'Content-Type': mediaType})
            try:
                r = urllib2.urlopen(req)
            except urllib2.HTTPError as e:
                # Unsupported Media Type (415): Can't handle script.
                if e.code == 415:
                    raise NotSupportedException("script is not supported")
                else:
                    raise e

    def Stop(self):
        """Stop playback.

        Raises:
            NotNowException: Stop can not be performed now, eg because there
                is no script loaded.
        """

        req = urllib2.Request(self._url+'/v1/stop')
        try:
            r = urllib2.urlopen(req)
        except urllib2.HTTPError as e:
            # Conflict (409): Player not in a state that can pause.
            if e.code == 409:
                raise NotNowException("cannot stop script now")
            else:
                raise e

    def Pause(self):
        """Pause playback.

        Raises:
            NotNowException: Pause can not be performed now, eg because there
                is no script loaded.
        """
        req = urllib2.Request(self._url+'/v1/pause')
        try:
            r = urllib2.urlopen(req)
        except urllib2.HTTPError as e:
            # Conflict (409): Player not in a state that can pause.
            if e.code == 409:
                raise NotNowException("cannot pause script now")
            else:
                raise e

    def Resume(self):
        """Resume playback.

        Raises:
            NotNowException: Pause can not be performed now, eg because there
                is no script loaded.
        """
        req = urllib2.Request(self._url+'/v1/resume')
        try:
            r = urllib2.urlopen(req)
        except urllib2.HTTPError as e:
            # Conflict (409): Player not in a state that can resume.
            if e.code == 409:
                raise NotNowException("cannot resume now")
            else:
                raise e

    def Skip(self, time):
        """Skip jumps to a timecode.
        
        Raises:
            NotNowException: Skip can not be performed now, eg because there
                is no script loaded.
        """
        req = urllib2.Request(self._url+'/v1/skip?p=%dms' % time)
        try:
            r = urllib2.urlopen(req)
        except urllib2.HTTPError as e:
            # Conflict (409): Player not in a state that can skip.
            if e.code == 409:
                raise NotNowException("cannot skip now")
            else:
                raise e
