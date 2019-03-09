package main

import (
	"fmt"
	"strconv"
)

const BOARD_EMPTY = 0
const BOARD_BLACK = 1
const BOARD_WHITE = 2
const STATE_MATCHING = 1000
const STATE_PAUSING = 1003
const STATE_RUNNING = 1001
const STATE_FINISHED = 1002
const MAX_PLAYER = 2

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

func (game *Game) run() { //ゲームが走る=対戦中

	clients := *(game.clients)
	fmt.Println("len(game.clients:" + strconv.Itoa(len(clients)))
	for {
		if game.state == STATE_MATCHING {
			game.runMatching()
		} else if game.state == STATE_RUNNING {
			game.runRunning()
		} else if game.state == STATE_PAUSING {
			game.runPausing()
		} else if game.state == STATE_FINISHED {
			fmt.Println("FINISH THE GAME")
			return
		}
	}
}

func (game *Game) runRunning() { //ゲームが走る=対戦中
	clients := *(game.clients)
	if len(clients) < MAX_PLAYER {
		fmt.Println("ゲームを続行できません")
		panic(len(clients))
	}

	if game.handCount == 3 { //TODO しっかり終了時の処理をする

		game.setState(STATE_FINISHED)
		//////TODO ロジックと通信は分けたい
		//clients[0].WriteMessageInfo("notice", "finish") //TODO 送信する部分はPutStoneの外に出すべき
		//clients[1].WriteMessageInfo("notice", "finish")
	}
}

func (game *Game) GetState() int {
	return game.state
}

func (game *Game) runMatching() { //待機中(マッチング中ともいえる)
	clients := *(game.clients)
	//if rand.Intn(10000) == 1 {
	//	fmt.Println("runMathing")
	//}
	if len(clients) == MAX_PLAYER {

		//game.state = STATE_RUNNING
		game.setState(STATE_RUNNING)
		game.stones[clients[0]] = BOARD_WHITE
		game.stones[clients[1]] = BOARD_BLACK

		firstPlayer := clients[0]
		game.nowPlayer = clients[0]
		game.nextPlayer = clients[1]
		clients[0].WriteMessageInfo("board", game.GetBoardStr()) //TODO 送信する部分はPutStoneの外に出すべき
		clients[1].WriteMessageInfo("board", game.GetBoardStr())
		firstPlayer.WriteMessageInfo("notice", "you") //初めのプレイヤー
		fmt.Println("ゲーム開始")
	}
}

func (game *Game) setState(state int) {
	*game.gameState <- state
	game.state = state
}

func (game *Game) runPausing() {

}

func (game *Game) PutStone(client *client, x, y int) bool { //おけたらtrue、置けなかったfalse
	canPut := false
	dirX := []int{1, 1, 0, -1, -1, -1, 0, 1}
	dirY := []int{0, 1, 1, 1, 0, -1, -1, -1}
	board := game.board
	stone := game.stones[client] //TODO これを毎回書かなくてもよい用に
	clients := *(game.clients)
	for dir := 0; dir < 8; dir++ {
		canPutNowDir := false
		for cnt, nowX, nowY := 0, x+dirX[dir], y+dirY[dir]; nowX >= 0 && nowY >= 0 && nowX < board.x && nowY < board.y && board.ban[nowY][nowX] != BOARD_EMPTY; cnt, nowX, nowY = cnt+1, nowX+dirX[dir], nowY+dirY[dir] {
			if board.ban[nowY][nowX] == stone {
				if cnt > 0 {
					canPutNowDir = true
					canPut = true
				}
				break
			}

		}
		if !canPutNowDir {
			continue
		}
		//fmt.Println("okeru")
		for cnt, nowX, nowY := 0, x+dirX[dir], y+dirY[dir]; nowX >= 0 && nowY >= 0 && nowX < board.x && nowY < board.y && board.ban[nowY][nowX] != stone; cnt, nowX, nowY = cnt+1, nowX+dirX[dir], nowY+dirY[dir] {
			board.ban[nowY][nowX] = stone
		}
	}
	if canPut {
		fmt.Println(game.nowPlayer)
		fmt.Println(game.nextPlayer)
		board.ban[y][x] = stone
		//////TODO ロジックと通信は分けたい
		clients[0].WriteMessageInfo("board", game.GetBoardStr()) //TODO 送信する部分はPutStoneの外に出すべき
		clients[1].WriteMessageInfo("board", game.GetBoardStr())
		game.handCount++
		game.nowPlayer.WriteMessageInfo("notice", "enemy")
		game.nextPlayer.WriteMessageInfo("notice", "you")
		game.nextPlayer, game.nowPlayer = game.nowPlayer, game.nextPlayer
	}
	return canPut
}

func (game *Game) GetBoardStr() string {
	ret := ""
	for y := 0; y < game.board.y; y++ {
		for x := 0; x < game.board.x; x++ {
			ret += strconv.Itoa(game.board.ban[y][x])
		}
	}
	return ret
}
