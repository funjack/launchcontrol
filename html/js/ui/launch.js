/*

Lauchcontrol UI Fleshlight
https://github.com/funjack/launchcontrol

Copyright 2017 Funjack

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
1. Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.
3. Neither the name of the copyright holder nor the names of its contributors
may be used to endorse or promote products derived from this software without
specific prior written permission.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/

;(function($, doc, win) {
    "use strict";

    /*
     * fleshlight widget
     *
     * Simulates the movements of a Fleshlight mounted in a Launch.
     *
     * Will create blocks elements in a given element, so must be used on a
     * block element like <div>.
     *
     */
    $.widget( "launchcontrol.fleshlight", {
        options: {
            ruler: true,
            range: 0.22, // there is 22% of movement with the default images.
            rulerImg: "img/ruler.png",
            fleshlightImg: "img/fleshlight.png",
        },

        _create: function() {
            var cDiv = doc.createElement("div");
            cDiv.className = "c-fleshlight";

            var ruler = doc.createElement("img");
            ruler.className = "o-ruler";
            ruler.src = this.options.rulerImg;
	    this.ruler = $(ruler)

            var fleshlight = doc.createElement("img");
            fleshlight.className = "o-fleshlight";
            fleshlight.src = this.options.fleshlightImg;
	    this.fleshlight = $(fleshlight)

            cDiv.appendChild(ruler);
            cDiv.appendChild(fleshlight);
            this.element.append(cDiv);

            if (!this.options.ruler) {
                this.ruler.fadeOut(0);
            }
        },

        // move to to a given position at a given speed.
        //
        // position: place to move to in percent (0-100).
        // speed: speed to move at in percent (20-100).
        move: function(position, speed) {
            position = position < 0 ? 0 : position;
            position = position > 100 ? 100 : position;
            speed = speed < 20 ? 20 : speed;
            speed = speed > 100 ? 100 : speed;

            var p = this.options.range * position;
            var duration = this.moveDuration(position, speed);
            this.fleshlight.stop(true);
            this.fleshlight.animate({bottom: [p+"%", "linear"]}, duration);
        },

        // positions returns the current position in percent (0-100).
        position: function() {
            var bottomPx = parseFloat(this.fleshlight.css("bottom"));
            var widgetHeightPx = parseFloat(this.element.css("height"));
            var percentValue = bottomPx / widgetHeightPx * 100;
            return Math.round(percentValue / this.options.range);
        },

        // moveDuration returns the time in milliseconds it will take to move
        // to position at speed.
        //
        // position: position in percent (0-100).
        // speed:    speed in percent (20-100).
        moveDuration: function(position, speed) {
            var distance = Math.abs(position-this.position());
            return this._calcDuration(distance, speed);
        },

        // _calcDuration returns duration of a move in milliseconds for a given
        // distance/speed.
        //
        // distance: amount to move percent (0-100).
        // speed: speed to move at in percent (20-100).
        _calcDuration: function(distance, speed) {
	    var mil = Math.pow(speed/25000, -0.95);
	    return mil/(90/distance);
	    /*
            // TODO figure out the real timings and best way to calculate them.
            //
            // Rough measurements 'felt' like a log10... :) and putting 100% at
            // 300ms given a full stroke for now:
            var delayPerPercent = (Math.log10(100/speed) * 8.5 + 3);
            return distance * delayPerPercent;
	    */
        },

        // stroke continues to move between two points at a set speed.
        //
        // start: start position of the stroke in percent (0-100).
        // stop : end position of the stroke in percent (0-100).
        // speed: speed to move at in percent (20-100).
        stroke: function(start, stop, speed) {
            win.clearTimeout(this.timeout);
            var fl = this;
            fl.move(start, speed);
            fl.timeout = win.setTimeout(function() {
                fl.move(stop, speed);
                win.setTimeout(function() {
                    fl.stroke(start, stop, speed);
                }, fl.moveDuration(stop, speed));
            }, fl.moveDuration(start, speed));
        },

        // play runs a Launchcontrol style script starting at a given position.
        //
        // script: array of actions objects:
        //   action.at  = time in ms.
        //   action.pos = position in percent (0-100).
        //   action.spd = speed in percent (20-100).
        // i: number in the array to start at.
        play: function(script, i) {
            i = typeof i !== 'undefined' ? i : 0;
            var fl = this;
            if (i >= script.length) { return; }
            var delay = i === 0 ? script[i].at : script[i].at-script[i-1].at;
            win.setTimeout(function(){
                fl.move(script[i].pos, script[i].spd);
                fl.play(script, i+1);
            }, delay);
        },

        // stop halts any active script / stroke.
        stop: function() {
            win.clearTimeout(this.timeout);
        },

        // toggleRuler displays or hides the ruler image.
        toggleRuler: function() {
            this.ruler.fadeToggle();
        },

    });
})(jQuery, document, window);
