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

func (game *Game) run() { //ゲームが走る=対戦中

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
	if len(clients) < MAX_PLAYER { //人が抜けてしまった時
		fmt.Printf("Can't continue the game clients %v\n", len(clients))
		game.setState(STATE_FINISHED) //TODO FINISHじゃなくてconitnueに移動するように
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

	if len(clients) == MAX_PLAYER {

		game.setState(STATE_RUNNING)
		game.stones[clients[0]] = BOARD_WHITE
		game.stones[clients[1]] = BOARD_BLACK

		firstPlayer := clients[0]
		game.nowPlayer = clients[0]
		game.nextPlayer = clients[1]
		game.nowPlayer.WriteMessageInfo("board", game.GetBoardStr()) //TODO 送信する部分はPutStoneの外に出すべき
		game.nextPlayer.WriteMessageInfo("board", game.GetBoardStr())
		firstPlayer.WriteMessageInfo("notice", "you") //初めのプレイヤー
		fmt.Println("start the game")
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
	stone := game.stones[client]
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
		game.exePut(client, x, y)

	}
	return canPut
}

func (game *Game) exePut(client *client, x, y int) {
	fmt.Println(game.nowPlayer)
	fmt.Println(game.nextPlayer)
	game.board.ban[y][x] = game.stones[client]
	//////TODO ロジックと通信は分けたい
	game.nowPlayer.WriteMessageInfo("board", game.GetBoardStr()) //TODO 送信する部分はPutStoneの外に出すべき
	game.nextPlayer.WriteMessageInfo("board", game.GetBoardStr())
	game.handCount++
	game.nowPlayer.WriteMessageInfo("notice", "enemy")
	game.nextPlayer.WriteMessageInfo("notice", "you")
	game.nextPlayer, game.nowPlayer = game.nowPlayer, game.nextPlayer
}

func (game *Game) GetBoardStr() string { //boardをJSONとして送るための文字列に変換
	ret := ""
	for y := 0; y < game.board.y; y++ {
		for x := 0; x < game.board.x; x++ {
			ret += strconv.Itoa(game.board.ban[y][x])
		}
	}
	return ret
}
