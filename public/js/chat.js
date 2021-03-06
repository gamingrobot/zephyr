var server = "wss://zephyr.gamingrobot.net/ws"
if (window.location.search === '?debug') {
    server = 'ws://localhost:3000/ws';
}
var current_chat = "";

var websocket = new ReconnectingWebSocket(server);
websocket.onopen = function() {
    websocket.onmessage = function(message) {
        var steamevent = JSON.parse(message.data)
        console.log(steamevent);
        var newMessage = $('<li>').text(JSON.stringify(steamevent));
        $('#content-friends >.chat >.messages').append(newMessage);
    };
}


$(document).ready(function() {
    $("#left-bar-list ul li").each(function() {
        $(this).click(function() {
            left_bar_click($(this))
        });
    });
});

function chat_message_submit(element, id) {
    var msgbox = element.find(".message-box")
    var msg = msgbox.val();
    if (msg.length > 0) {
        var body = {
            "SteamId": id,
            "ChatEntryType": EChatEntryType.ChatMsg,
            "Message": msg
        }
        sendEvent("SendMessageEvent", body)
        msgbox.val('');
    }
}

function left_bar_click(element) {
    var chatid = element.attr("id");
    chatid = chatid.split("-")
    var type = chatid[0];
    var id = chatid[1];
    console.log(type, id);
    if (type === "friends" || type === "groups") {
        //left sidebar chats
        if (!$("#chats-" + id).length) { //check that the element doesnt already exist
            //add the chat to sidebar
            var $newside = $("#chats-default").clone();
            $newside.attr("id", "chats-" + id)
            flexShow($newside);
            $newside.append(element.clone().children());
            $newside.click(function() {
                left_bar_click($newside)
            });
            $newside.appendTo("#chats");
        }
        //middle bar chats
        if (!$("#chat-" + id).length) { //check that the element doesnt already exist
            //show the new chat
            var $newchat = $("#chat-default").clone();
            $newchat.attr("id", "chat-" + id)
            $el = element.clone();
            $newchat.find(".chat-header").empty().append($el.children());
            var $chat_form = $newchat.find(".chat-form")
            $chat_form.submit(function() {
                chat_message_submit($chat_form, id)
                return false
            });
            flexShow($newchat);
            hideChats();
            //add our chat
            $newchat.appendTo("#container");
            joinChat(id);
        }

    } else if (type === "chats") {
        var $chat = $("#chat-" + id);
        hideChats();
        current_chat = id;
        flexShow($chat);
    }
    //hide the chat sidebar
    if (type === "friends") {
        var $chat = $("#chat-" + id + " .right-bar");
        $chat.hide();
    }
}

function hideChats() {
    //hide all other chats
    $(".chat-bar").each(function() {
        $(this).hide();
    });
}

function flexShow(element) {
    //show but without the display:block
    element.removeAttr("style")
}

function joinChat(id) {
    current_chat = id;
    var body = {
        "SteamId": id,
    }
    sendEvent("JoinChatEvent", body)
}

function sendEvent(eve, body) {
    var msg = {
        "Event": eve,
        "Body": body
    }
    websocket.send(JSON.stringify(msg));
}