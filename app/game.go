package app

import (
	"encoding/json"
	"strconv"

	"dappswin/models"

	"github.com/golang/glog"
)

var isGameDone bool = true

var gameResult string

// 奖励是5位数字
const gameCodeLen int = 5

var gameChan = make(chan *models.Block, 4096)

func gameRoutine() {
	for {
		select {
		case block := <-gameChan:
			glog.V(7).Infof("Gamemsg Coming %v, %v, 游戏开奖结束%v", block.Hash, isNumber(block.LastLetter()), isGameDone)
			tm := block.Time / 1000
			if !isGameDone && isNumber(block.LastLetter()) {
				gameResult += block.LastLetter()
				glog.Infof("gameResult....: %s, len is %d, %d, BlockTime is %d", gameResult, len(gameResult), block.Num, block.Time)
				if len(gameResult) == gameCodeLen {
					gameorm := &models.Game{Result: gameResult, BlockNum: block.Num, TimeStamp: tm, GameMintue: int64(tm/60) - 1}
					// TODO GameID 应该由时间算出来
					if err := models.AddGame(gameorm); err != nil {
						glog.Errorf("insert game error %v", err)
					}
					gameID := tm / 60

					// 推送到待处理交易区，判定输赢
					winchan <- gameorm

					r, _ := strconv.ParseInt(gameResult, 10, 32)
					gamews := &GameWS{gameID, r}
					pushGameMessage(block, gamews)
					// TODO push to win.go
					isGameDone = true
					gameResult = ""
				}
			} else if isGameDone && tm%60 == 0 {
				isGameDone = false
				if isNumber(block.LastLetter()) {
					gameResult += block.LastLetter()
				}
			}
		}
	}
}

// 根据交易消息里的时间判定属于哪一期游戏，判定

func pushGameMessage(blk *models.Block, game *GameWS) {
	m := &models.Message{winType, blk.Num, blk.Hash, blk.Time, game}
	m.HandleTimeStamp()
	// needless do this. block have handled it.
	// m.HandleTimeStamp()
	ms := []*models.Message{}
	ms = append(ms, m)
	result, _ := json.Marshal(ms)
	Huber.broadcast <- result
}

func isNumber(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

type GameWS struct {
	ID     int64 `json:"gameid"`
	Result int64 `json:"result"`
}
