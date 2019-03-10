<html>
  <head>
    <script>
      //preload
      var socket;
      var board = new Image();    //画像オブジェクト作成
      var white = new Image();
      var black = new Image();
      var loading=new Image();
      board.src = "http://localhost:8081/templates/board.png";  //写真のパスを指定する
      black.src = "http://localhost:8081/templates/black.png";
      white.src = "http://localhost:8081/templates/white.png";
      loading.src="http://localhost:8081/templates/loading.gif";
      var maxX=8;
      var maxY=8;
    </script>
 </head>

 <body>



   ゲーム画面
    <form id="backButton">
       <input type="button" value="タイトルに戻る"/>
  </form>



<script src="//ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js"> </script>
<canvas id="cv1" width="800" height="800"></canvas>
<div id="matchingField">
  <img src="http://localhost:8081/templates/loading.gif">
  マッチングチュウ....
</div>


<script>
//描画系
var ctx = document.getElementById("cv1").getContext("2d");
var canvas=document.getElementById("cv1");
var boardX=100,boardY=100;
var stoneX=100,stoneY=100;
var myStone=0;
var myTurn=false;
function drawDefault(){

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
  if(msg=="enemy"){
    myTurn=false;
    console.log(myTurn);
  }
  if(msg=="finish"){
    alert("試合終了です。")
    socket.close();
  }
}

</script>

<script>
$(function(){

  $("#backButton").on('click',function(){
    console.log("backButton")
    setTimeout(function(){
      window.location.href = '/';
    }, 1000);
  })});
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

<script type="text/javascript" src="http://localhost:8081/templates/contactToServer.js"></script>

 </body>
</html>
