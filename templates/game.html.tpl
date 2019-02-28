<html>
  <head>
    <script>
      //preload
      var socket;
      var board = new Image();    //画像オブジェクト作成
      var white = new Image();
      var black = new Image();
      board.src = "http://localhost:8081/templates/board.png";  //写真のパスを指定する
      black.src = "http://localhost:8081/templates/black.png";
      white.src = "http://localhost:8081/templates/white.png";
      var maxX=8;
      var maxY=8;
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
var ctx = document.getElementById("cv1").getContext("2d");
var canvas=document.getElementById("cv1");
var boardX=100,boardY=100;
var stoneX=100,stoneY=100;
var myStone=0;
var myTurn=false;
function draw1(){

  for(var y=0;y<maxY;y++){
    for(var x=0;x<maxX;x++){
      ctx.drawImage(board,x*boardX,y*boardY);
    }
  }
}
</script>

<script>

function sendMessageInfo(operation,message){
  sendJSON=JSON.stringify({"operation":operation,"message":message});
  console.log(socket);
  socket.send(sendJSON);
}
function fetchBoard(board){
  for(var y=0;y<maxY;y++){
    for(var x=0;x<maxX;x++){
      var now=board[y*maxY+x];
      //console.log(now);
      if(now=="1"){
        ctx.drawImage(black,boardX*x,boardY*y);
      }else if(now=="2"){
        ctx.drawImage(white,boardX*x,boardY*y);
      }
    }
  }
}

function fetchInfo(msg){
  if(msg=="you"){
    myTurn=true;
    console.log(myTurn);
  }
}

</script>

<script>
//canvasがクリックされた時
function onDown(e){
  var x = e.clientX - canvas.offsetLeft;
  var y = e.clientY - canvas.offsetTop;
  console.log("x:", x, "y:", y,"myturn",myTurn);
  if(myTurn){
    sendMessageInfo("put",String(Math.floor(x/boardX))+String(Math.floor(y/boardY)));//クライアントは打つ位置だけ知らせる
  }
}

canvas.addEventListener("mousedown",onDown,false);
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

    }

    socket.onmessage=function(e){
      //console.log(e.data)
      var obj=JSON.parse(e.data);//送られてきたJSONを受け取る
      console.log("operation="+obj.operation+" , message="+obj.message)

      if(obj.operation=="board"){
        fetchBoard(obj.message)
      }
      if(obj.operation=="notice"){
        fetchInfo(obj.message)
      }
    }
    socket.onopen=function(e){
      in_room=true;
      var obj;

      socket.send('{"operation":"require","message":"hoge"}')//wsが繋がった時まず盤面をとってくる
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
