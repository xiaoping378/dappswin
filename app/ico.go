package app

import (
	"dappswin/models"
	"fmt"

	"github.com/golang/glog"
)

const (
	handled = iota
	pending
	unknow
)

// ICORules ico等级， 100EOS以下 多奖励0.5%
var ICORules = map[uint64]float64{
	100 * 1e4:  0.005,
	500 * 1e4:  0.01,
	1000 * 1e4: 0.015,
	1500 * 1e4: 0.02,
	2000 * 1e4: 0.025,
}

func getRefund(value uint64) float64 {

	for amount, reward := range ICORules {
		if value < amount {
			return float64(value*eosConf.EOS_CGG) * (1 + reward)
		}
	}
	// here return the max reward
	return float64(value*eosConf.EOS_CGG) * (1 + 0.025)
}

var icochan = make(chan *models.ICO, 4096)
var oneDayMills = int64(24 * 60 * 60 * 1000)

func checkICORoutine() {
	for {
		select {
		case ico := <-icochan:
			glog.V(7).Infof("有人投ICO募资了，who=>%s 额度是%s EOS", ico.Account, fmt.Sprintf("%.4f", float64(ico.Amount)/1e4))
			var amount float64
			if ico.TimeMills-eosConf.ICOStartTime <= oneDayMills {
				amount = getRefund(ico.Amount) / 1e4
			}
			amount = float64(ico.Amount*eosConf.EOS_CGG) / 1e4
			quantity := fmt.Sprintf("%.4f", amount) + " " + eosConf.TokenSymbol
			glog.Infof("奖励购买的代币===========> to %s, quantity: %s", ico.Account, quantity)
			if err := sendTokens(ico.Account, quantity); err == nil {
				if err := db.Model(ico).Update("status", handled).Error; err != nil {
					glog.Error(err)
				}
			}
		}
	}
}
