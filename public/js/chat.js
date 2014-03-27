var server = "wss://zephyr.gamingrobot.net/ws"
if (window.location.search === '?debug') {
    server = 'ws://localhost:3000/ws';
}
var c=new ReconnectingWebSocket(server);
c.onopen = function(){
  c.onmessage = function(message){
    var steamevent = JSON.parse(message.data)
    console.log(steamevent);
    var newMessage = $('<li>').text(JSON.stringify(steamevent));
    $('#messages').append(newMessage);
  };
  $('form').submit(function(){
    var msg = $('#message').val();
    if(msg.length > 0){
        var chatmsg = {"Event": "SendMessageEvent", "Body": {"SteamId": $('#steamid').val(), "ChatEntryType": EChatEntryType.ChatMsg, "Message": msg}}
        console.log(chatmsg)
        console.log(JSON.stringify(chatmsg))
        c.send(JSON.stringify(chatmsg));
        $('#message').val('');
        return false;
    }
    return false;
  });
}

