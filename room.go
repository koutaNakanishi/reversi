package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const MAX_CONNECTION_PER_ROOM = 2 //1部屋に繋げる最大の人数
const MAX_ROOM_NUM = 10           //最大のサーバ全体の部屋の数

func (r *room) GetRoomNum() int {
	return len(r.clients)
}

func (r *room) run() {

	for {
		select {
		case client := <-r.join: //クライアントが入室してきた時
			r.clientsMap[client] = true //mapにクライアントを追加
			r.clients = append(r.clients, client)
			fmt.Printf("join in room %v. now clients %v\n", r.room_id, len(r.clients))
		case client := <-r.leave: //クライアントが体質した時
			delete(r.clientsMap, client)
			r.clients = removeClient(r.clients, client)
			close(client.send)
			fmt.Printf("left in room %v. now clients %v\n", r.room_id, len(r.clients))
		case state := <-r.gameState:
			checkGameState(state, r)

		case msg := <-r.forward: //誰からのメッセージが来た時
			fmt.Println(msg)

		}

	}
}
func checkGameState(state int, r *room) {

	if state == STATE_FINISHED {
		//ここでゲーム終了処理
		fmt.Println("RESET THE GAME")
		for _, c := range r.clients { //TODO ゲーム終了時の実際の部屋やgameオブジェクトの処理はroom.goで
			c.WriteMessageInfo("notice", "finish")
			c.socket.Close()

		}
		rooms = removeRoom(rooms, r)

	}

}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	if len(r.clientsMap) >= MAX_CONNECTION_PER_ROOM {
		fmt.Println("Can't connect this room")
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()
}
