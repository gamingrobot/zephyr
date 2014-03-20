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


//tab pages stuff

  var pageImages = [];
  var pageNum = 1;
/**
* Reset numbering on tab buttons
*/
function reNumberPages() {
    pageNum = 1;
    var tabCount = $('#pageTab > li').length;
    $('#pageTab > li').each(function() {
        var pageId = $(this).children('a').attr('href');
        if (pageId == "#page1") {
            return true;
        }
        pageNum++;
        $(this).children('a').html('Page ' + pageNum +
            '<button class="close" type="button" ' +
            'title="Remove this page">×</button>');
    });
}
  
$(document).ready(function() {
  /**
   * Add a Tab
   */
  $('#btnAddPage').click(function() {
  pageNum++;
  $('#pageTab').append(
    $('<li><a href="#page' + pageNum + '">' +
    'Page ' + pageNum +
    '<button class="close" type="button" ' +
    'title="Remove this page">×</button>' +
    '</a></li>'));

  $('#pageTabContent').append(
    $('<div class="tab-pane" id="page' + pageNum +
    '">Content page' + pageNum + '</div>'));

  $('#page' + pageNum).tab('show');
  });


  /**
  * Remove a Tab
  */
  $('#pageTab').on('click', ' li a .close', function() {
  var tabId = $(this).parents('li').children('a').attr('href');
  $(this).parents('li').remove('li');
  $(tabId).remove();
  reNumberPages();
  $('#pageTab a:first').tab('show');
  });

  /**
   * Click Tab to show its content 
   */
  $("#pageTab").on("click", "a", function(e) {
  e.preventDefault();
  $(this).tab('show');
  });
});