package app

import (
	"dappswin/models"
	"fmt"
	"sort"

	"github.com/golang/glog"
)

const (
	handled = iota
	pending
	unknow
)

// ICORules ico等级， 100EOS以下 多奖励0.5%
var ICORules = map[int]float64{
	100 * 1e4:   0,
	500 * 1e4:   50,
	1000 * 1e4:  100,
	1500 * 1e4:  150,
	2000 * 1e4:  200,
	20000 * 1e4: 250,
}

func getRefund(time int64, value float64) float64 {

	if time-eosConf.ICOStartTime > oneDayMills {
		return value * eosConf.EOS_CGG
	}

	sortedKeys := []int{}
	for k, _ := range ICORules {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	for _, k := range sortedKeys {
		if value < float64(k) {
			return value*eosConf.EOS_CGG*ICORules[k]/1e4 + value*eosConf.EOS_CGG
		}
	}

	// here return the max reward
	return value*eosConf.EOS_CGG*250/1e4 + value*eosConf.EOS_CGG
}

var icochan = make(chan *models.ICO, 4096)
var oneDayMills = int64(24 * 60 * 60 * 1000)

func checkICORoutine() {
	for {
		select {
		case ico := <-icochan:
			glog.Infof("有人投ICO募资了，who=>%s 额度是%s EOS", ico.Account, fmt.Sprintf("%.4f", float64(ico.Amount)/1e4))

			amount := getRefund(ico.TimeMills, float64(ico.Amount)) / 1e4
			quantity := fmt.Sprintf("%.4f ", amount) + eosConf.TokenSymbol
			glog.Infof("奖励购买的代币===========> to %s, quantity: %s", ico.Account, quantity)
			if hash, err := sendTokens(ico.Account, quantity); err == nil {
				glog.Infof("to %s, %s, hash is %s", ico.Account, quantity, hash)
				if err := db.Model(ico).Update("status", handled).Error; err != nil {
					glog.Error(err)
				}
			}
		}
	}
}
