/*

Lauchcontrol UI Controls
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

;(function($, doc) {
    "use strict";

    /*
     * controls widget
     *
     * Player control buttons.
     *
     */
    $.widget( "launchcontrol.controls", {
        options: {
            play: null,
            pause: null,
            stop: null,
            back: null,
            forward: null,
            load: null,
        },

        _create: function() {
            var bGroup = doc.createElement("div");
            bGroup.className = "btn-group";
            bGroup.setAttribute("role", "group");
            bGroup.setAttribute("aria-label", "controls");

            var input = doc.createElement("input");
            input.setAttribute("type", "file");
            input.style.display = "none";
            bGroup.appendChild(input);

            var file = $( input );
            var loadFunc = this.options.load;
            file.change(function() { loadFunc(file.prop('files')[0]); });
            bGroup.appendChild(this._createButton("glyphicon-open-file",
                        function() { file.trigger("click"); }));

            //bGroup.appendChild(this._createButton("glyphicon-backward",
            //            this.options.back));
            bGroup.appendChild(this._createButton("glyphicon-play",
                        this.options.play));
            bGroup.appendChild(this._createButton("glyphicon-pause",
                        this.options.pause));
            bGroup.appendChild(this._createButton("glyphicon-stop",
                        this.options.stop));
            //bGroup.appendChild(this._createButton("glyphicon-forward",
            //            this.options.forward));
            this.element.append(bGroup);
        },

        _createButton: function(glyph, callback) {
            var button = doc.createElement("button");
            button.className = "btn btn-default";
            button.setAttribute("type", "button");
            button.onclick = callback;

            var span = doc.createElement("span");
            span.className = "glyphicon " + glyph;
            span.setAttribute("aria-hidden", "true");

            button.appendChild(span);
            return button
        },
    });

})(jQuery, document);
