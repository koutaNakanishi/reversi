package main

import (
	"fmt"
	"strconv"
)

const BOARD_EMPTY = 0
const BOARD_BLACK = 1
const BOARD_WHITE = 2
const STATE_TERMINATING = 1000
const STATE_RUNNING = 1001
const STATE_STARTING = 1002
const MAX_PLAYER = 2

type Game struct {
	board     *Board
	state     int
	roomNum   int //部屋にいる人数
	handCount int //何手分ゲームが進んだか
	clients   []*client
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

func NewGame() *Game {

	_board := NewBoard()
	game := new(Game)
	game.board = _board
	game.roomNum = 0
	game.state = STATE_TERMINATING //始めはゲームが始まっていない
	game.handCount = 0
	return game
}

func (game *Game) run() { //ゲームが走る=対戦中

	for {
		if game.state == STATE_RUNNING {
			game.runRunning()
		} else if game.state == STATE_STARTING { //書記処理

		} else if game.state == STATE_TERMINATING {
			game.runTerminating()
		}
	}
}

func (game *Game) runRunning() { //ゲームが走る=対戦中
	if game.roomNum < MAX_PLAYER {
		fmt.Println("ゲームを続行できません")
		return
	}
}

func (game *Game) runTerminating() { //待機中
	if game.roomNum == MAX_PLAYER {

		game.state = STATE_RUNNING
	}
}

func (game Game) PutStone(stone, x, y int) bool { //おけたらtrue、置けなかったfalse
	canPut := false
	dirX := []int{1, 1, 0, -1, -1, -1, 0, 1}
	dirY := []int{0, 1, 1, 1, 0, -1, -1, -1}
	board := game.board
	for dir := 0; dir < 8; dir++ {
		canPutNowDir := false
		for cnt, nowX, nowY := 0, x+dirX[dir], y+dirY[dir]; board.ban[nowY][nowX] != BOARD_EMPTY && nowX >= 0 && nowY >= 0 && nowX < board.x && nowY < board.y; cnt, nowX, nowY = cnt+1, nowX+dirX[dir], nowY+dirY[dir] {
			if cnt > 0 && board.ban[nowY][nowX] == stone {
				canPutNowDir = true
				canPut = true
			}
		}
		if !canPutNowDir {
			continue
		}

		for cnt, nowX, nowY := 0, x+dirX[dir], y+dirY[dir]; board.ban[nowY][nowX] != stone && nowX >= 0 && nowY >= 0 && nowX < board.x && nowY < board.y; cnt, nowX, nowY = cnt+1, nowX+dirX[dir], nowY+dirY[dir] {
			board.ban[nowY][nowX] = stone
		}
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
