var launchcontrolClient = (function() {
    "use strict";
    var url = "";
    var scriptTypes = [
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
    ];

    var setUrl = function(u) {
        url = u;
    };

    var play = function(file, callback) {
        var extension = file.name.split('.').pop();
        var mediaType = getMediaType(extension);
        httpPost(url+"/v1/play", file, mediaType, callback);
    };

    var stop = function(callback) {
        httpGet(url+"/v1/stop", callback);
    };

    var pause = function(callback) {
        httpGet(url+"/v1/pause", callback);
    };

    var resume = function(callback) {
        httpGet(url+"/v1/resume", callback);
    };

    var skip = function(time, callback) {
        httpGet(url+"/v1/skip?p="+time+"ms", callback);
    };

    var httpGet = function(url, callback) {
        var xmlHttp = new XMLHttpRequest();
        xmlHttp.onreadystatechange = function() { 
            if (xmlHttp.readyState == 4)
                if(callback && typeof callback == "function") {
                    callback(xmlHttp.responseText, xmlHttp.status);
                }
        };
        xmlHttp.open("GET", url, true);
        xmlHttp.send(null);
    };

    var httpPost = function(url, data, contentType, callback) {
        var xmlHttp = new XMLHttpRequest();
        xmlHttp.onreadystatechange = function() { 
            if (xmlHttp.readyState == 4)
                if(callback && typeof callback == "function") {
                    callback(xmlHttp.responseText, xmlHttp.status);
                }
        };
        xmlHttp.open("POST", url, true);
        if (contentType !== null) {
            xmlHttp.setRequestHeader("Content-Type", contentType);
        }
        xmlHttp.send(data);
    };

    var getMediaType = function(extension) {
        for (var i = 0; i < scriptTypes.length; i++) {
            var scriptType = scriptTypes[i];
            for (var j = 0; j < scriptType.extensions; j++) {
                if (scriptType.extensions[j] == extension) {
                    return scriptType.mediaType;
                }
            }
        }
        return null;
    };

    return {
        url: setUrl,
        play: play,
        stop: stop,
        pause: pause,
        resume: resume,
        skip: skip,
    };
})();
