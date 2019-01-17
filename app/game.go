package app

import (
	"encoding/json"
	"strconv"

	"dappswin/models"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/shopspring/decimal"
)

var isGameDone bool = true

var gameResult string

// 奖励是5位数字
const gameCodeLen int = 6

var cachedgameid int64

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
					// 记录成上一分钟的值，数据库好匹配
					content := getGameContent(gameResult)
					gameorm := &models.Game{Result: gameResult, BlockNum: block.Num, TimeStamp: tm, GameMintue: int64(tm/60) - 1, Content: content}
					if err := models.AddGame(gameorm); err != nil {
						glog.Errorf("insert game error %v", err)
					}
					gameID := tm / 60
					cachedgameid = gameID

					// 推送到待处理交易区，判定输赢
					winchan <- gameorm

					r, _ := strconv.ParseInt(gameResult, 10, 32)
					gamews := &GameWS{gameID, r, content}
					pushGameMessage(block, gamews)
					// TODO push to win.go
					isGameDone = true
					gameResult = ""
				}
			} else if isGameDone && tm%60 == 0 {
				// 揭晓上一分钟的号码
				totalVotedCGG = decimal.NewFromFloat(0)
				totalVotedEOS = decimal.NewFromFloat(0)
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
	ID      int64  `json:"gameid"`
	Result  int64  `json:"result"`
	Content string `json:"content"`
}

func getGameContent(result string) string {
	var content string
	if result[lastStar] >= '0' && result[lastStar] <= '4' {
		content = "s|"
	} else {
		content = "b|"
	}

	if result[lastStar]%2 == 0 {
		content += "e"
	} else {
		content += "o"
	}

	return content
}

type gamePagePost struct {
	PageIndex int `json:"page_index" binding:"required,gt=0,lt=100"`
	PageSize  int `json:"page_size" binding:"required,gt=0,lt=100"`
}

type pageGameRsp struct {
	Count int            `json:"total_items"`
	Pages int            `json:"total_pages"`
	Data  []*models.Game `json:"data"`
}

func pageLottery(c *gin.Context) {
	body := &gamePagePost{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(NewMsg(400, "输入参数有误"))
		return
	}
	games := []*models.Game{}
	var count int
	index := (body.PageIndex - 1) * body.PageSize

	if err := db.Where(models.Game{}).Offset(index).Limit(body.PageSize).Order("id desc").Find(&games).Error; err != nil {
		c.JSON(NewMsg(500, "系统内部错误"))
		return
	}
	if err := db.Model(models.Game{}).Count(&count).Error; err != nil {
		c.JSON(NewMsg(500, "系统内部错误"))
		return
	}
	c.JSON(NewMsg(200, &pageGameRsp{count, (count / body.PageSize) + 1, games}))
}
