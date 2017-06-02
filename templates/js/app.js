$(document).foundation()

var ws;
//var url = "ws://localhost:5000/ws";
//var url = "ws://192.168.137.1:5000/ws";
var url = "ws://gwordapi.herokuapp.com/ws";
var enabled = false;
var stop = false;
var inProgress = false;
var scrolled = false;
var initialScroll = 0;
var $theForm = $("#the-form");
var $letters = $("#letters");
var $limit = $("#limit");
var $generate = $(".generate");
var $response = $("#response");
var $temp_stop = $("#temp-stop");
var $footer = $("#footer");

$(document).ready(function () {
    if(window.WebSocket === undefined){
        var notification = $("<div>")
            .attr({
                class : "notification"
            })
            .html("<i class='fa fa-exclamation-circle'></i>&nbsp; dangit! your browser doesn't support websockets :(")
            .hide()
            .velocity("transition.slideDownIn", {
                duration : 200,
                easing : "ease-in-out"
            })
            .velocity("transition.slideUpOut", {
                delay : 2000,
                duration : 200,
                easing : "ease-in-out"
            });
        $("body").append(notification);
        console.log("WebSocket unsupported");
    }else{
        ws = initWebSocket();
        console.log("WebSocket initiated");
        enabled = true;
    }
    if($(window).width() <= (39.9375 * 16)){
        $footer.children().eq(0).removeClass("small-4");
        $footer.children().eq(0).addClass("small-12");
        $footer.children().eq(2).css({
            display : "none"
        });
        $footer.children().eq(1).css({
            display : "none"
        });
    }
});

$(window).resize(function () {
    if($(window).width() <= (39.9375 * 16)){
        if($footer.children().length > 1){
            $footer.children().eq(0).removeClass("small-4");
            $footer.children().eq(0).addClass("small-12");
            $footer.children().eq(2).css({
                display : "none"
            });
            $footer.children().eq(1).css({
                display : "none"
            });
        }
    }else{
        $footer.children().eq(0).removeClass("small-12");
        $footer.children().eq(0).addClass("small-4");
        $footer.children().eq(2).css({
            display : "flex"
        });
        $footer.children().eq(1).css({
            display : "flex"
        });
    }
});

$(window).scroll(function () {
    var currentScroll = $(this).scrollTop();
    if (inProgress){
        if (!scrolled){
            if (currentScroll < initialScroll){
                //Scrolling Up
                scrolled = true;
                var notification = $("<div>")
                    .attr({
                        class : "notification"
                    })
                    .html("<i class='fa fa-info-circle'></i>&nbsp; auto-scroll disabled")
                    .hide()
                    .velocity("transition.slideDownIn", {
                        duration : 200,
                        easing : "ease-in-out"
                    })
                    .velocity("transition.slideUpOut", {
                        delay : 2000,
                        duration : 200,
                        easing : "ease-in-out"
                    });
                $("body").append(notification);
            }
        }
    }
    initialScroll = currentScroll;
});

function initWebSocket() {
    var socket = new WebSocket(url);
    var count = 1;
    var scroll_count = 1;

    socket.onopen = function () {
        console.log("WebSocket ready");
    };

    socket.onmessage = function (e) {
        var res = JSON.parse(e.data);
        if(res['Done'] === true){
            var word_item = $("<div>")
                .attr({
                    class : "word-container small-6 medium-4 large-4 column"
                })
                .html("")
                .hide()
                .velocity("fadeIn", {
                    duration : 200,
                    easing : "ease-in-out"
                });
            if(stop){
                var notification = $("<div>")
                    .attr({
                        class : "notification"
                    })
                    .html("<i class='fa fa-exclamation-circle'></i>&nbsp; Stopped")
                    .hide()
                    .velocity("transition.slideDownIn", {
                        duration : 200,
                        easing : "ease-in-out"
                    })
                    .velocity("transition.slideUpOut", {
                        delay : 2000,
                        duration : 200,
                        easing : "ease-in-out"
                    });href="https://www.cognition.co.ke" class="cognition" target="_blank"
            }else{
                notification = $("<div>")
                    .attr({
                        class : "notification"
                    })
                    .html("<i class='fa fa-check-circle'></i>&nbsp; Done!")
                    .hide()
                    .velocity("transition.slideDownIn", {
                        duration : 200,
                        easing : "ease-in-out"
                    })
                    .velocity("transition.slideUpOut", {
                        delay : 2000,
                        duration : 200,
                        easing : "ease-in-out"
                    });
            }

            $response.append(word_item);
            $("body").append(notification);

            $letters.attr({
                disabled : false
            });
            $limit.attr({
                disabled : false
            });
            $generate.removeClass("stop");
            $generate.addClass("generate");
            $generate.html("<i class='fa fa-cogs'></i>&nbsp;generate!");
            $temp_stop.children().remove();

            count = 1;
            scroll_count = 1;
            scrolled = false;
            inProgress = false;
        }else{
            word_item = $("<div>")
                .attr({
                    class : "word-container small-6 medium-4 large-4 column"
                })
                .html(
                    "<div class='word-item small-12 column'>" +
                        count + ". " + res['Word'] +
                    "</div>"
                )
                .hide()
                .velocity("fadeIn", {
                    duration : 200,
                    easing : "ease-in-out"
                });
            $response.append(word_item);
            count++;
            if(($response.height() + $response.offset().top) >= $(window).height()){
                if(scroll_count <= 15) {
                    scroll_count++;
                }else{
                    if (!scrolled){
                        $response.velocity("scroll", {
                            duration: 500,
                            offset: $response.height() - 20,
                            easing: "ease-in-out"
                        });
                        scroll_count = 1;
                    }
                }
            }
        }
        return false;
    };

    socket.onclose = function () {
        console.log("WebSocket closed");

        /*var notification = $("<div>")
            .attr({
                class : "notification"
            })
            .html("<i class='fa fa-info-circle'></i>&nbsp; websocket closed")
            .hide()
            .velocity("transition.slideDownIn", {
                duration : 200,
                easing : "ease-in-out"
            })
            .velocity("transition.slideUpOut", {
                delay : 2000,
                duration : 200,
                easing : "ease-in-out"
            });
        $("body").append(notification);*/

        $letters.attr({
            disabled : false
        });
        $limit.attr({
            disabled : false
        });
        $generate.removeClass("stop");
        $generate.addClass("generate");
        $generate.html("<i class='fa fa-cogs'></i>&nbsp;generate!");
        $temp_stop.children().remove();

        count = 1;
        scroll_count = 1;
        scrolled = false;
        inProgress = false;

        socket = new WebSocket(url);
    };

    return socket;
}

$theForm.submit(function (e) {
   e.preventDefault();

   if(enabled) {
       if (inProgress) {
           stop = true;
           ws.send(JSON.stringify({
               Letters: "stop",
               Limit: 0,
               Stop: true
           }));
       } else {
           stop = false;
           $response.contents().velocity("fadeOut", {
               duration: 200,
               easing: "ease-in-out"
           });
           $response.html("");

           $letters.attr({
               disabled: true
           });
           $limit.attr({
               disabled: true
           });
           $generate.removeClass("generate");
           $generate.addClass("stop");
           $generate.html("<i class='fa fa-stop-circle'></i>&nbsp;stop");
           $generate.clone().appendTo("#temp-stop");

           ws.send(JSON.stringify({
               Letters: $("#letters").val(),
               Limit: parseInt($("#limit").val()),
               Stop: false
           }));

           inProgress = true;
       }
   }else{
       var notification = $("<div>")
           .attr({
               class : "notification"
           })
           .html("<i class='fa fa-exclamation-circle'></i>&nbsp; dangit! your browser doesn't support websockets :(")
           .hide()
           .velocity("transition.slideUpIn", {
               duration : 200,
               easing : "ease-in-out"
           })
           .velocity("transition.slideDownOut", {
               delay : 2000,
               duration : 200,
               easing : "ease-in-out"
           });
       console.log("WebSocket unsupported");
   }
});