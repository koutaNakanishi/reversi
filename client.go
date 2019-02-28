package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

type MessageInfo struct { //クライアントから送られてくるメッセージ
	Operation string `json:"operation"`
	Msg       string `json:"message"`
}

//func (messageInfo *MessageInfo) CreateMessage() string { //各クライアントのブラウザに表示されるテキスト
//	ret := messageInfo.Name + ":" + messageInfo.Msg
//	return ret
//}

func (c *client) read() {
	messageInfo := MessageInfo{Operation: "tmp", Msg: "名無し"}
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			if string(msg) == "" {
				continue
			}
			fmt.Println(string(msg))
			if err := json.Unmarshal(msg, &messageInfo); err != nil {
				fmt.Println("JSON UNMARSHAL ERROR:", err)
			}
			//c.room.forward <- messageInfo //クライアントから受け取ったメッセージを送信
			if messageInfo.Operation == "require" {
				c.WriteRequire() //要はクライアントにrequireの要求結果を送信する
			}
			if messageInfo.Operation == "put" {

				x, _ := strconv.Atoi(string(messageInfo.Msg[0]))
				y, _ := strconv.Atoi(string(messageInfo.Msg[1]))
				fmt.Println(x + y)
				canPut := c.room.game.PutStone(c, x, y)
				fmt.Println("canput:", canPut)
				if canPut == false {
					//c.WriteNotice("you")
				} else {
					//c.WriteRequire() //盤面を教えてあげる
				}
			}
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) WriteRequire() {
	sendMessageInfo := MessageInfo{Operation: "board", Msg: c.room.game.GetBoardStr()}
	sendJSON, err := json.Marshal(sendMessageInfo)
	if err != nil {
		fmt.Println("JSON MARCHAL ERR:", err)
	}
	if err := c.socket.WriteMessage(websocket.TextMessage, sendJSON); err != nil {
		fmt.Println("ERROR IN writeRequire")
	}
}

func (c *client) WriteNotice(msg string) {
	sendMessageInfo := MessageInfo{Operation: "notice", Msg: msg}
	sendJSON, err := json.Marshal(sendMessageInfo)
	if err != nil {
		fmt.Println("JSON MARCHAL ERR:", err)
	}
	if err := c.socket.WriteMessage(websocket.TextMessage, sendJSON); err != nil {
		fmt.Println("ERROR IN writeRequire")
	}
}

func (c *client) write() {

	for msg := range c.send { //部屋から受け取ったメッセージをクライアントに送信
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
