package main

import (
	"encoding/json"
	"fmt"

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
				writeRequire(c) //要はクライアントにrequireの要求結果を送信する
			}
		} else {
			break
		}
	}
	c.socket.Close()
}

func writeRequire(c *client) {
	sendMessageInfo := MessageInfo{Operation: "board", Msg: c.room.game.GetBoardStr()}
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
