var in_room=false;
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
