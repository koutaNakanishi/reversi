package main

import (
	"html/template"
	"sync"

	"github.com/gorilla/websocket"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//部屋の管理に使用
type room struct {
	room_id    int
	forward    chan MessageInfo //誰かが送信したメッセージ
	join       chan *client     //入室してきたクライアント
	leave      chan *client     //体質していくクライアント
	gameState  chan int
	clientsMap map[*client]bool //入室しているクライアント一覧
	clients    []*client        //部屋に居る人達のスライス
	game       *Game            //ゲームの今の状況
}

func newRoom(room_id_ int) *room {
	ret := &room{
		room_id:    room_id_,
		forward:    make(chan MessageInfo),
		join:       make(chan *client),
		leave:      make(chan *client),
		gameState:  make(chan int, 100),
		clientsMap: make(map[*client]bool),
	}
	ret.game = NewGame(&ret.gameState, &ret.clients)
	return ret
}

func removeRoom(rooms []*room, search *room) []*room {
	result := []*room{}
	for _, v := range rooms {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}

type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func removeClient(clients []*client, search *client) []*client {
	result := []*client{}
	for _, v := range clients {
		if v != search {
			result = append(result, v)
		}
	}
	return result
}

//チャットでのやりとりの管理に使用
type MessageInfo struct { //クライアントから送られてくるメッセージ
	Operation string `json:"operation"`
	Msg       string `json:"message"`
}

////////ゲームのロジック部分に利用
type Game struct {
	board      *Board
	state      int
	roomNum    int //部屋にいる人数
	handCount  int //何手分ゲームが進んだか
	nowPlayer  *client
	nextPlayer *client
	clients    *[]*client
	stones     map[*client]int
	gameState  *chan int
}

func NewGame(c *chan int, clients *[]*client) *Game {

	_board := NewBoard()
	game := new(Game)
	game.board = _board
	game.roomNum = 0
	game.handCount = 0
	game.stones = make(map[*client]int)
	game.clients = clients
	game.gameState = c
	game.setState(STATE_MATCHING)

	return game
}

type Board struct {
	x, y int
	ban  [8][8]int
}

func NewBoard() *Board {
	board := new(Board)
	board.x = 8
	board.y = 8
	board.ban[3][3] = BOARD_WHITE
	board.ban[4][4] = BOARD_WHITE
	board.ban[3][4] = BOARD_BLACK
	board.ban[4][3] = BOARD_BLACK
	return board
}
