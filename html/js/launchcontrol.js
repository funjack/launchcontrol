$(document).ready(function(loc, client) {
    'use strict';

    var controls = $( "#controls" ).controls({
        play: client.resume,
        pause: client.pause,
        stop: client.stop,
        load: client.play,
    });

    var fleshlight = $( "#fleshlight" ).fleshlight();
    console.log("Acquiring websocket");
    var launchSocket = new WebSocket("ws://" + loc.host + "/v1/socket");
    launchSocket.onmessage = function(event) {
        console.log(event.data);
        var action = JSON.parse(event.data);
        fleshlight.fleshlight("move", action.pos, action.spd);
    }
}(location, launchcontrolClient));
