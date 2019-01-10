package app

import (
	"encoding/json"
	"fmt"
	"strconv"

	"dappswin/models"

	"github.com/golang/glog"
)

var winchan = make(chan *models.Game, 4096)

func checkWinRoutine() {
	for {
		select {
		case game := <-winchan:
			glog.Infof("开奖啦，激动人心的时刻到来了。。。%#v", game)
			txs, err := models.GetTxsInGame(game.GameMintue, pending)
			if err != nil {
				glog.Error("getTxsInGame error :", err)
			}
			msgs := []*models.Message{}
			msg := &models.Message{}

			for i, tx := range txs {
				_, betInfo, _ := models.ResolveMemo(tx.Memo)
				glog.Infof("betInfo is %s, %s, %f", betInfo, game.Result, tx.Amount)
				winTimes, betnum := HandleBetInfo(betInfo, []byte(game.Result))
				if winTimes > 0 {
					winValue := calcBenefit(winTimes, betnum, tx.Amount)
					// TODO call cleos unlock transfer EOS and lock, if OK update SQL.
					// {"bookid":0,"status":0,"to":"kunoichi3141","amount":"0.3920 EOS","memo":"win|25736640:50090:e"}
					memo := "win|" + fmt.Sprint(game.GameMintue) + ":" + game.Result + ":" + betInfo
					reward := &models.Reward{Amount: winValue + " " + coinNames[tx.CoinID], ID: i, Status: handled, To: tx.From, Memo: memo}
					glog.Infof("开始发放奖励, %#v", reward)

					//构造赢家消息
					gameStr, _, _ := models.ResolveMemo(tx.Memo)
					msg = &models.Message{wsTypes[gameStr+coinNames[tx.CoinID]+"Win"], game.BlockNum, tx.TxID, tx.TimeMills, reward}
					// msg.HandleTimeStamp()
					msgs = append(msgs, msg)
					quan := winValue + " " + coinNames[tx.CoinID]
					go func(tx *models.Tx, quan string) {
						if hash, err := sendTokens(tx.From, quan); err == nil {
							amount, _ := strconv.ParseFloat(winValue, 64)
							models.AddTx(&models.Tx{TxID: hash, From: eosConf.GameAccount, To: tx.From,
								Amount: amount, CoinID: tx.CoinID, Memo: memo, Status: handled, TimeMills: tx.TimeMills, TimeMintue: tx.TimeMintue})
						}
					}(&tx, quan)

				}

				tx.Status = handled
				models.UpdateTxById(&tx)

				// TODO make winTX, save to DB
			}
			if len(msgs) > 0 {
				buf, _ := json.Marshal(msgs)
				Huber.broadcast <- buf
			}
		}
	}
}

// "0.1000"
func calcBenefit(times int, betTimes int, amount float64) string {
	s := fmt.Sprintf("%.4f", amount/float64(betTimes)*float64(times)*98/100)
	return s
}
