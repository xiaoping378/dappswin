package app

import (
	"dappswin/common"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"dappswin/models"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/shopspring/decimal"
)

// TODO creat message hub, switch message.Data.(type)
var txschan = make(chan *models.Message, 4096)
var votechan = make(chan *models.Tx, 4096)

func resloveTXRoutine() {
	for {
		select {
		case txsMsg := <-txschan:
			// TODO save to db.tx
			txs, ok := txsMsg.Data.([]models.TransactionResp)
			if !ok {
				glog.Error("txs asscert failed!")
				break
			}

			txsmsg := []*models.Message{}
			for _, tx := range txs {
				// glog.Infof("%#v", tx)
				actions, hash := tx.GetActions()
				if actions == nil {
					continue
				}

				glog.V(7).Infof("%#v", actions)
				for _, action := range actions {
					if !action.IsTransfer() {
						continue
					}
					coinName := action.Coin()
					if coinName == "" {
						continue
					}
					if msg := handleTX(coinName, hash, action, txsMsg); msg != nil {
						txsmsg = append(txsmsg, msg)
					}
				}
			}
			if len(txsmsg) != 0 {
				buf, _ := json.Marshal(txsmsg)
				Huber.broadcast <- buf
			}
		}
	}
}

func handleTX(coin string, hash string, action models.Action, txsMsg *models.Message) *models.Message {
	msg := models.Message{}
	glog.V(7).Infof("coin %s to %s", coin, action.Data.To)
	if eosConf.EnableICO && coin == "EOS" && action.Data.To == eosConf.ICOAccount {
		t := models.TX{Quantity: action.Data.Quantity}
		ico := &models.ICO{Hash: hash, Account: action.Data.From, Amount: t.Amount(), Status: pending, TimeMills: txsMsg.TimeMills}
		models.AddIcoRecord(ico)
		icochan <- ico
		return nil
	}

	// 操作质押的逻辑
	if coin == "CGG" && action.Data.To == eosConf.LockAccount {
		str := strings.Split(action.Data.Quantity, " ")
		amount, _ := strconv.ParseFloat(str[0], 64)
		// db.Save(&models.Stake{Name: action.Data.From, Amount: amount, Date: common.JSONTime{Time: time.Now().Add(24 * time.Hour)}, Status: staked})
		stake := models.Stake{Name: action.Data.From, Amount: decimal.NewFromFloat(amount), Date: common.JSONTime{Time: time.Now()}, Status: staked}
		if unfound := db.Where("name = ?", action.Data.From).First(&stake).RecordNotFound(); unfound {
			db.Save(&stake)
			return nil
		}
		decimalAmount := stake.Amount.Add(decimal.NewFromFloat(amount))
		db.Model(&models.Stake{}).Where("name = ?", action.Data.From).Update(&models.Stake{Amount: decimalAmount})
		return nil
	}

	if action.Data.To != eosConf.GameAccount {
		return nil
	}

	game, _, _ := models.ResolveMemo(action.Data.Memo)

	msg.BlockNum = txsMsg.BlockNum
	msg.TimeMills = txsMsg.TimeMills
	t, ok := wsTypes[game+coin+"Buy"]
	if !ok {
		return nil
	}
	msg.Type = t
	msg.Hash = hash
	msg.Data = action.Data

	str := strings.Split(action.Data.Quantity, " ")
	amount, _ := strconv.ParseFloat(str[0], 64)
	glog.Infof("Coming Bet is %s, %s,  %s, %f, timemills: %d , %d期游戏", action.Data.Quantity, str[0], action.Data.From, amount, txsMsg.TimeMills, txsMsg.TimeMills/1000/60)

	txdb := &models.Tx{
		TxID: hash, BlockNum: txsMsg.BlockNum,
		From: action.Data.From, To: action.Data.To, Amount: amount, CoinID: coinIDs[coin], Memo: action.Data.Memo,
		Status: pending, TimeMills: txsMsg.TimeMills, TimeMintue: txsMsg.TimeMills / 1000 / 60}
	// 计算累积投注
	votechan <- txdb

	go models.AddTx(txdb)

	// cancel paltform
	// go updateUsersFromTX(txdb)
	return &msg
}

func votedRoutine() {
	for {
		select {
		case vote := <-votechan:
			glog.Infof("计算累积投注 %#v, cachegameid is %d", vote, cachedgameid)
			// if cachedgameid != 0 && vote.TimeMills/1000/60 > cachedgameid {
			// 	if vote.CoinID == eos {
			// 		totalVotedEOS = decimal.NewFromFloat(vote.Amount)
			// 	} else if vote.CoinID == cgg {
			// 		totalVotedCGG = decimal.NewFromFloat(vote.Amount)
			// 	}
			// 	pushVoteMsg(vote)
			// 	break
			// }

			if vote.CoinID == eos {
				totalVotedEOS = totalVotedEOS.Add(decimal.NewFromFloat(vote.Amount))
			} else if vote.CoinID == cgg {
				totalVotedCGG = totalVotedCGG.Add(decimal.NewFromFloat(vote.Amount))
			}

			pushVoteMsg(vote)

		}
	}
}

func pushVoteMsg(vote *models.Tx) {
	votemsgs := []*models.Message{}
	votemsg := &models.Message{}
	votemsg.Type = wsTypes["lottery"+coinNames[vote.CoinID]+"TotalVoted"]
	votemsg.BlockNum = vote.BlockNum
	votemsg.Hash = vote.TxID
	votemsg.TimeMills = vote.TimeMills
	if vote.CoinID == eos {
		votemsg.Data = map[string]string{"total_voted": totalVotedEOS.StringFixed(4)}
	} else if vote.CoinID == cgg {
		votemsg.Data = map[string]string{"total_voted": totalVotedCGG.StringFixed(4)}
	}
	votemsgs = append(votemsgs, votemsg)
	buf, _ := json.Marshal(votemsgs)
	Huber.broadcast <- buf
}

type pageTXPost struct {
	PageIndex int    `json:"page_index" binding:"required,gt=0,lt=100"`
	PageSize  int    `json:"page_size" binding:"required,gt=0,lt=100"`
	Name      string `json:"name" binding:"required,max=12"`
}

type pageTXRsp struct {
	Count int          `json:"count"`
	Data  []*models.Tx `json:"data"`
}

func pageTxes(c *gin.Context) {
	body := &pageTXPost{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(NewMsg(400, "输入参数有误"))
		return
	}
	txes := []*models.Tx{}
	var count int
	index := (body.PageIndex - 1) * body.PageSize

	if err := db.Where(models.Tx{From: body.Name}).Offset(index).Limit(body.PageSize).Order("time_mintue desc").Find(&txes).Count(&count).Error; err != nil {
		c.JSON(NewMsg(500, "系统内部错误"))
		return
	}

	c.JSON(NewMsg(200, &pageTXRsp{count, txes}))
}
