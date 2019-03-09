package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const MAX_CONNECTION_PER_ROOM = 2 //1部屋に繋げる最大の人数
const MAX_ROOM_NUM = 10           //最大のサーバ全体の部屋の数
type room struct {
	forward    chan MessageInfo //誰かが送信したメッセージ
	join       chan *client     //入室してきたクライアント
	leave      chan *client     //体質していくクライアント
	gameState  chan int
	clientsMap map[*client]bool //入室しているクライアント一覧
	clients    []*client        //部屋に居る人達のスライス
	game       *Game            //ゲームの今の状況
}

func (r *room) GetRoomNum() int {
	return len(r.clientsMap)
}
func newRoom() *room {
	fmt.Println("OK")
	ret := &room{
		forward:    make(chan MessageInfo),
		join:       make(chan *client),
		leave:      make(chan *client),
		gameState:  make(chan int, 100),
		clientsMap: make(map[*client]bool),
	}
	ret.game = NewGame(&ret.gameState, &ret.clients)
	return ret
}

func (r *room) run() {
	//go checkGameState(r)
	for {
		select {
		case client := <-r.join: //クライアントが入室してきた時
			r.clientsMap[client] = true //mapにクライアントを追加
			r.clients = append(r.clients, client)
			fmt.Println("room.r.client.len", len(r.clients))
		case client := <-r.leave: //クライアントが体質した時
			delete(r.clientsMap, client)
			close(client.send)
		case state := <-r.gameState:
			checkGameState(state, r)

		case msg := <-r.forward: //誰からのメッセージが来た時
			fmt.Println(msg)

		}

	}
}
func checkGameState(state int, r *room) {
	//fmt.Println(r.game.GetState())
	if state == STATE_FINISHED {
		for _, c := range r.clients { //TODO ゲーム終了時の実際の部屋やgameオブジェクトの処理はroom.goで
			c.WriteMessageInfo("notice", "finish")
			c.socket.Close()
			fmt.Println("RESET THE GAME")
		}
		fmt.Println(r.game, r.game.state)
		r.game = NewGame(&r.gameState, &r.clients)
		fmt.Println(r.game, r.game.state)
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

	fmt.Println("len(r.clientsMap)=", len(r.clientsMap))
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
		fmt.Println("leave room")
	}()

	go client.write()
	client.read()
}

func remove(clients []*client, search *client) []*client {
	result := []*client{}
	for _, v := range clients {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}
