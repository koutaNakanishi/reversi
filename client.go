package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gorilla/websocket"
)

func (c *client) read() {
	messageInfo := MessageInfo{Operation: "tmp", Msg: "名無し"}
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {

			fmt.Println(string(msg))
			if err := json.Unmarshal(msg, &messageInfo); err != nil {
				fmt.Println("JSON UNMARSHAL ERROR:", err)
			} else if messageInfo.Operation == "require" {
				c.require()
			} else if messageInfo.Operation == "put" {
				c.put(messageInfo)
			}
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {

	for msg := range c.send { //部屋から受け取ったメッセージをクライアントに送信
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}

func (c *client) WriteMessageInfo(operation, msg string) {
	sendMessageInfo := MessageInfo{Operation: operation, Msg: msg}
	sendJSON, err := json.Marshal(sendMessageInfo)
	if err != nil {
		fmt.Println("JSON MARCHAL ERR:", err)
	}
	if err := c.socket.WriteMessage(websocket.TextMessage, sendJSON); err != nil {
		fmt.Println("ERROR IN writeRequire")
	}
}

func (c *client) require() {
	c.WriteMessageInfo("require", c.room.game.GetBoardStr()) //要はクライアントにrequireの要求結果を送信する
}

func (c *client) put(messageInfo MessageInfo) {
	x, _ := strconv.Atoi(string(messageInfo.Msg[0]))
	y, _ := strconv.Atoi(string(messageInfo.Msg[1]))

	canPut := c.room.game.PutStone(c, x, y)
	fmt.Println("canput:", canPut)

}
