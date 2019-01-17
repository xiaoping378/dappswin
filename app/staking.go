package app

import (
	"dappswin/common"
	"dappswin/conf"
	"dappswin/models"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/shopspring/decimal"
)

const (
	staked = iota
	unstaking
	released
)

// TotalLockedAmount 质押的总数量
var TotalLockedAmount float64
var totalSync sync.RWMutex

func setTotalLockedAmount(amount float64) {
	totalSync.Lock()
	defer totalSync.Unlock()
	TotalLockedAmount = amount
}

func getTotalLockedAmount() float64 {
	totalSync.RLock()
	defer totalSync.RUnlock()
	cached := TotalLockedAmount
	return cached
}

// CirculatingAmount 市场流通的额度
var CirculatingAmount float64
var circulatSync sync.RWMutex

func setCirculatingAmount(amount float64) {
	circulatSync.Lock()
	defer circulatSync.Unlock()
	CirculatingAmount = amount
}

func getCirculatingAmount() float64 {
	circulatSync.RLock()
	defer circulatSync.RUnlock()
	cached := CirculatingAmount
	return cached
}

func updateTotalLocked() {
	result := getBalance(eosConf.LockAccount, eosConf.TokenAccount, eosConf.TokenSymbol)
	if result == "" {
		return
	}
	amount, _ := strconv.ParseFloat(result, 64)
	glog.Info("Total locked amount is ", amount)
	setTotalLockedAmount(amount)
}

func updateCirculat() {
	result := getBalance(eosConf.OfficialLockAccount, eosConf.TokenAccount, eosConf.TokenSymbol)
	if result == "" {
		return
	}
	amount, _ := strconv.ParseFloat(result, 64)
	glog.Info("official locked amount is ", amount)
	setCirculatingAmount(eosConf.TotalCGGAmoount - amount)
}

func locktotalStatus(c *gin.Context) {
	c.JSON(NewMsg(200, map[string]interface{}{
		"total_locked": TotalLockedAmount,
		"circulating":  CirculatingAmount,
		"percent":      fmt.Sprintf("%0.4f", TotalLockedAmount/CirculatingAmount*100),
	}))
}

type Account struct {
	Name string `json:"name" binding:"required,max=12"`
}

func stakedStatus(c *gin.Context) {
	body := Account{}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(NewMsg(400, "输入参数错误"))
		return
	}

	stake := models.Stake{}
	db.Where("name = ? AND status = ?", body.Name, staked).First(&stake)

	c.JSON(NewMsg(200, map[string]interface{}{
		"staked":  stake.Amount,
		"percent": stake.Amount.Div(decimal.NewFromFloat(TotalLockedAmount)).Mul(decimal.NewFromFloat(100)).StringFixed(2),
	}))

}

type unstakePost struct {
	Name   string  `json:"name" binding:"required,max=12"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

func unstake(c *gin.Context) {
	body := unstakePost{}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(NewMsg(400, "输入参数错误"))
		return
	}

	stake := models.Stake{}
	db.Where("name = ? AND status = ?", body.Name, staked).First(&stake)
	if stake.Amount.LessThan(decimal.NewFromFloat(body.Amount)) {
		c.JSON(NewMsg(400, "赎回的数额大于实际质押的了"))
		return
	}
	amount := stake.Amount.Sub(decimal.NewFromFloat(body.Amount))
	// TODO 事务
	db.Model(&models.Stake{}).Where("name = ? AND status = ?", body.Name, staked).Update("amount", amount)
	unstake := models.Stake{}
	if unfound := db.Where("name = ? AND status = ?", body.Name, unstaking).First(&unstake).RecordNotFound(); unfound {
		db.Save(&models.Stake{Name: body.Name, Amount: decimal.NewFromFloat(body.Amount), Date: common.JSONTime{Time: time.Now().Add(conf.C.GetDuration("eos.UnstakePeriod"))}, Status: unstaking})
		c.JSON(NewMsg(200, "赎回成功，在排队中"))
		return
	}
	unAmount := unstake.Amount.Add(decimal.NewFromFloat(body.Amount))
	db.Model(&models.Stake{}).Where("name = ? AND status = ?", body.Name, unstaking).Update(&models.Stake{Amount: unAmount, Date: common.JSONTime{Time: time.Now().Add(conf.C.GetDuration("eos.UnstakePeriod"))}})

	c.JSON(NewMsg(200, "赎回成功，在排队中"))

}

func unstakeStatus(c *gin.Context) {
	body := Account{}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(NewMsg(400, "输入参数错误"))
		return
	}

	stake := models.Stake{}
	db.Where("name = ? AND status = ?", body.Name, unstaking).First(&stake)

	c.JSON(NewMsg(200, stake))

}

func checkNotifyRoutine() {
	ticker := time.NewTicker(1 * time.Minute)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			glog.Info("=====检查是有需要赎回到期的。。。。。。。")

			now := time.Now()

			index, limit := 0, 100
			for index = 0; ; index += limit {
				stakes := []models.Stake{}
				if nok := db.Model(&models.Stake{}).
					Where("status = ? AND date <= ?", unstaking, now).
					Offset(index).Limit(limit).Find(&stakes).RecordNotFound(); nok {
					break
				}

				for _, stake := range stakes {
					// key := fmt.Sprintf("%d@%s@%d", indivi.ID, indivi.TakeUpdateTime, indivi.TakeCount)
					// if sentMessages[key] {
					// 	continue
					// }
					db.Delete(&stake)
					sendTokens(eosConf.LockAccount, stake.Name, stake.Amount.StringFixed(4)+" CGG", "unstaked=>来自赎回的质押CGG")
					glog.Infof("=======赎回到期了 %#v", stake)

					// sentMessages[key] = true
				}

				if len(stakes) < limit {
					break
				}

			}

		}
	}

}
