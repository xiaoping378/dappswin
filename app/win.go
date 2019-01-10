package app

import (
	"encoding/json"
	"fmt"

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
				glog.Infof("betInfo is %s, %s, %d", betInfo, game.Result, tx.Amount)
				winTimes, betnum := HandleBetInfo(betInfo, []byte(game.Result))
				if winTimes > 0 {
					winValue := calcBenefit(winTimes, betnum, tx.Amount)
					// TODO call cleos unlock transfer EOS and lock, if OK update SQL.
					glog.Infof("开始发放奖励, to %s,  %s", tx.From, winValue)

					// {"bookid":0,"status":0,"to":"kunoichi3141","amount":"0.3920 EOS","memo":"win|25736640:50090:e"}
					memo := "win|" + fmt.Sprint(game.Id) + ":" + game.Result + ":" + betInfo
					reward := &models.Reward{Amount: winValue + " " + coinNames[tx.CoinID], ID: i, Status: 0, To: tx.From, Memo: memo}

					//构造赢家消息
					msg = &models.Message{gameType, game.BlockNum, tx.TxID, tx.TimeMills, reward}
					// msg.HandleTimeStamp()
					msgs = append(msgs, msg)
				}
				tx.Status = handled
				models.UpdateTxById(&tx)
			}
			if len(msgs) > 0 {
				buf, _ := json.Marshal(msgs)
				Huber.broadcast <- buf
			}
		}
	}
}

// "0.1000"
func calcBenefit(times int, betTimes int, amount uint64) string {
	s := fmt.Sprintf("%.4f", float64(amount)/1e4/float64(betTimes)*float64(times)*98/100)
	return s
}
