package app

import (
	"dappswin/conf"
	"dappswin/database"

	"github.com/jinzhu/gorm"
)

var eosConf *EosConf
var db *gorm.DB

// InitEos 启动Eos Resolver
func Init() {
	eosConf = newEosConf()
	db = database.Db
	go gameRoutine()
	go resloveTXRoutine()
	go checkWinRoutine()
	go Huber.run()
}

// EosConf :
type EosConf struct {
	RPCURL         string
	ChainID        string
	WalletPassword string
	FetchIdleDur   int    // 查询blk时间间隔
	FromBlkNum     uint32 // 从哪个blocknum开始查询
	GameAccount    string
}

func newEosConf() *EosConf {
	dur := conf.C.GetInt("eos.FetchIdleDur")
	num := conf.C.GetInt64("eos.FromBlkNum")
	return &EosConf{
		RPCURL:         conf.C.GetString("eos.RPCURL"),
		ChainID:        conf.C.GetString("eos.ChainID"),
		WalletPassword: conf.C.GetString("eos.WalletPassword"),
		FetchIdleDur:   dur,
		FromBlkNum:     uint32(num),
		GameAccount:    conf.C.GetString("eos.GameAccount"),
	}
}
