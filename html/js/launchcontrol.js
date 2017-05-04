$(document).ready(function() {
   'use strict';

    var fleshlight = $( "#fleshlight" ).fleshlight();
    console.log("Acquiring websocket");
    var launchSocket = new WebSocket("ws://" + location.host + "/v1/socket");
    launchSocket.onmessage = function(event) {
        console.log(event.data);
        var action = JSON.parse(event.data);
        fleshlight.fleshlight("move", action.pos, action.spd);
    }
}());
