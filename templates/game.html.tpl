<html>
  <head>
    <script>
      //preload
      var board = new Image();    //画像オブジェクト作成
      var white = new Image();
      var black = new Image();
      board.src = "http://localhost:8081/templates/board.png";  //写真のパスを指定する
      black.src = "http://localhost:8081/templates/black.png";
      white.src = "http://localhost:8081/templates/white.png";
    </script>
 </head>

 <body onload="draw1()">



   ゲーム画面
    <form id="backButton">
       <input type="button" value="タイトルに戻る"/>
  </form>

<script src="//ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js"> </script>
<canvas id="cv1" width="800" height="800"></canvas>



<script>
//描画系
function draw1(){
  var ctx = document.getElementById("cv1").getContext("2d");
  var boardX=100,boardY=100;
  var komaX=100,komaY=100;
  for(var y=0;y<8;y++){
    for(var x=0;x<8;x++){
      ctx.drawImage(board,x*boardX,y*boardY);
    }
  }
  ctx.drawImage(white,boardX*3,boardY*3);
  ctx.drawImage(white,boardX*4,boardY*4);
  ctx.drawImage(black,boardX*4,boardY*3);
  ctx.drawImage(black,boardX*3,boardY*4);
}
</script>



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
      socket.send('{"operation":"require","msg","hoge"}')
    }

    socket.onmessage=function(e){
      //console.log(e.data)
      var obj=JSON.parse(e.data);//送られてきたJSONを受け取る
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
