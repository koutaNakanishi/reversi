<html>
  <head>
 </head>

 <body>



   ゲーム画面
    <form id="backButton">
       <input type="button" value="タイトルに戻る"/>
  </form>

  <img id="test" src="/home/mohutarou/work/reversi/templates/black.png" alt="表示エラー">


<script src="//ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js"> </script>



<script>
var in_room=false;




$(function(){

  $("#backButton").on('click',function(){
    console.log("backButton")
    setTimeout(function(){
      window.location.href = '/';
    }, 1000);
  });

  $.get("http://localhost:8081/createOrJoinRoom",function(data){
    room_url=data;
    if(in_room==false){
      socket = new WebSocket(room_url);
      console.log("新しい接続:"+room_url)
    }
    socket.onopen=function(e){
      in_room=true;
    }
    socket.onclose=function(e){
      in_room=false;
    }
    socket.onerror=function(e){
      alert("err"+e)
    }
  });

});
</script>
 </body>
</html>
