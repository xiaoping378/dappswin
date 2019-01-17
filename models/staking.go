package models

import (
	"dappswin/common"

	"github.com/shopspring/decimal"
)

// Stake 质押的结构体
type Stake struct {
	ID   uint   `grom:"PRIMARY_KEY" json:"-"`
	Name string `json:"name"`
	// 20,8, 整数部分取20位，小数部分支持8位
	Amount decimal.Decimal `json:"amount" sql:"type:decimal(20,8)"`
	// 赎回到账时间
	Date   common.JSONTime `json:"date"`
	Status int             `grom:"index" json:"status"`
}
