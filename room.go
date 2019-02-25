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
	forward chan MessageInfo //誰かが送信したメッセージ
	join    chan *client     //入室してきたクライアント
	leave   chan *client     //体質していくクライアント
	clients map[*client]bool //入室しているクライアント一覧
}

func (r *room) GetRoomNum() int {
	return len(r.clients)
}
func newRoom() *room {
	return &room{
		forward: make(chan MessageInfo),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join: //クライアントが入室してきた時
			r.clients[client] = true //mapにクライアントを追加

		case client := <-r.leave: //クライアントが体質した時
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- ([]byte)(msg.CreateMessage()):
				default:
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
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
	fmt.Println("len(r.clients)=", len(r.clients))
	if len(r.clients) >= MAX_CONNECTION_PER_ROOM {
		fmt.Println("Can't connect this room")
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
