var server = "wss://zephyr.gamingrobot.net/ws"
if (window.location.search === '?debug') {
    server = 'ws://localhost:3000/ws';
}
var c = new ReconnectingWebSocket(server);
c.onopen = function() {
    c.onmessage = function(message) {
        var steamevent = JSON.parse(message.data)
        console.log(steamevent);
        var newMessage = $('<li>').text(JSON.stringify(steamevent));
        $('#content-friends >.chat >.messages').append(newMessage);
    };
    $('form').submit(function() {
        var msg = $('#message').val();
        if (msg.length > 0) {
            var chatmsg = {
                "Event": "SendMessageEvent",
                "Body": {
                    "SteamId": $('#steamid').val(),
                    "ChatEntryType": EChatEntryType.ChatMsg,
                    "Message": msg
                }
            }
            console.log(chatmsg)
            console.log(JSON.stringify(chatmsg))
            c.send(JSON.stringify(chatmsg));
            $('#message').val('');
            return false;
        }
        return false;
    });
}


$(document).ready(function() {
    $("#left-bar-list ul li").each(function() {
        $(this).click(function() {
            left_bar_click($(this))
        });
    });
});

function left_bar_click(element) {
    var chatid = element.attr("id");
    chatid = chatid.split("-")
    var type = chatid[0];
    var id = chatid[1];
    console.log(type, id);
    if (type === "friends" || type === "groups") {
        if (!$("#chats-" + id).length) { //check that the element doesnt already exist
            //add the chat to sidebar
            var $newside = $("#chats-default").clone();
            $newside.attr("id", "chats-" + id)
            $newside.removeAttr("style") //show but without the display:block
            $newside.append(element.clone().children());
            $newside.click(function() {
                left_bar_click($newside)
            });
            $newside.appendTo("#chats");
        }
        if (!$("#chat-" + id).length) { //check that the element doesnt already exist
            //show the new chat
            var $newchat = $("#chat-default").clone();
            $newchat.attr("id", "chat-" + id)
            $el = element.clone().unwrap()
            $newchat.find(".chat-header").empty().append($el.children());
            $newchat.removeAttr("style") //show but without the display:block
            //hide all other chats
            $(".chat-bar").each(function() {
                $(this).hide();
            });
            //add our chat
            $newchat.appendTo("#container");
        }

    } else if (type === "chats") {
        var $chat = $("#chat-" + id);
        //hide all other chats
        $(".chat-bar").each(function() {
            $(this).hide();
        });
        $chat.removeAttr("style") //show but without the display:block
    }
    //hide the chat sidebar
    if (type === "friends") {
        var $chat = $("#chat-" + id + " .right-bar");
        $chat.hide();
    }
}