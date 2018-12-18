package app

import (
	"encoding/json"

	"dappswin/models"

	"github.com/golang/glog"
)

// TODO creat message hub, switch message.Data.(type)
var txschan = make(chan *models.Message, 4096)

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
			if len(txs) == 0 {
				glog.V(7).Infof("空块 %d", txsMsg.BlockNum)
				break
			}
			txmsg := []*models.Message{}

			for _, tx := range txs {
				if tx.Status != "executed" {
					continue
				}
				msg := &models.Message{}
				// tx.Trx is string
				if _, ok := tx.Trx.(string); ok {
					continue
				}
				// tx.Trx is map[string]interface {}
				trxBuf, _ := json.Marshal(tx.Trx)
				trxS := &models.TRX{}
				if err := json.Unmarshal(trxBuf, trxS); err != nil {
					glog.V(7).Infof("tx.Trx is not supported, ID %v, blockNum is %d", err, txsMsg.BlockNum)
					continue
				}

				for _, action := range trxS.TX.Actions {
					if action.Account == "eosio.token" && action.Name == "transfer" {
						if eosConf.EnableICO && action.Data.To == eosConf.ICOAccount {
							t := models.TX{Quantity: action.Data.Quantity}
							ico := &models.ICO{Hash: trxS.ID, Account: action.Data.From, Amount: t.Amount(), Status: pending, TimeMills: txsMsg.Time}
							models.AddIcoRecord(ico)
							icochan <- ico
						}
						if action.Data.To != eosConf.GameAccount {
							continue
						}
						if game, _, _ := models.ResolveMemo(action.Data.Memo); game != "lottery" {
							continue
						}
						msg.BlockNum = txsMsg.BlockNum
						msg.Time = txsMsg.Time
						msg.Type = txsMsg.Type
						msg.Hash = trxS.ID
						msg.Data = action.Data
						msg.HandleTimeStamp()

						txdb := &models.Tx{
							TxID: trxS.ID, BlockNum: txsMsg.BlockNum,
							From: action.Data.From, To: action.Data.To, Quantity: action.Data.Quantity, Memo: action.Data.Memo,
							Status: pending, Time: txsMsg.Time, TimeMintue: txsMsg.Time / 1000 / 60}

						go models.AddTx(txdb)

						// 更新用户投注信息
						models.UpdateUserInfo(msg)

						txmsg = append(txmsg, msg)
					}
				}
			}
			if len(txmsg) != 0 {
				buf, _ := json.Marshal(txmsg)
				Huber.broadcast <- buf
			}
		}
	}
}
