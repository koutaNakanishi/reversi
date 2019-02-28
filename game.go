package main

import "strconv"

const BOARD_EMPTY = 0
const BOARD_BLACK = 1
const BOARD_WHITE = 2

type Game struct {
	board *Board
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
	return game
}

func (board *Board) PutStone(stone, x, y int) bool { //おけたらtrue、置けなかったfalse
	canPut := false
	dirX := []int{1, 1, 0, -1, -1, -1, 0, 1}
	dirY := []int{0, 1, 1, 1, 0, -1, -1, -1}

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

func (board *Board) GetBoardStr() string {
	ret := ""
	for y := 0; y < board.y; y++ {
		for x := 0; x < board.x; x++ {
			ret += strconv.Itoa(board.ban[y][x])
		}
	}
	return ret
}
