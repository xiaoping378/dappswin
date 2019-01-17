package models

import (
	"dappswin/conf"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// TX lottery action.
type TX struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Quantity string `json:"quantity"`
	Memo     string `json:"memo"`
}

// ResolveMemo return game, bet, account
//
// "lottery:[357][357]@parentacccount"
func ResolveMemo(m string) (game, bet, account string) {
	var g, b, a string
	var str2 []string

	if !strings.Contains(m, ":") {
		return g, b, a
	}
	str := strings.Split(m, ":")
	g = str[0]

	if !strings.Contains(str[1], "@") {
		return g, str[1], a
	}
	str2 = strings.Split(str[1], "@")
	b = str2[0]
	a = str2[1]

	return g, b, a
}

// Amount return 0.0010 * 10000
//
// "0.0010 EOS"
func (t *TX) Amount() uint64 {
	s := strings.Split(t.Quantity, " ")
	f, _ := strconv.ParseFloat(s[0], 64)
	return uint64(f * 1e4)
}

type TransactionResp struct {
	Status string      `json:"status"`
	Trx    interface{} `json:"trx"`
}

func (tx TransactionResp) GetActions() ([]Action, string) {
	actions := []Action{}
	if tx.Status != "executed" {
		return nil, ""
	}
	// tx.Trx is string
	if _, ok := tx.Trx.(string); ok {
		return nil, ""
	}
	// tx.Trx is map[string]interface {}
	trxBuf, _ := json.Marshal(tx.Trx)
	trxS := &TRX{}
	if err := json.Unmarshal(trxBuf, trxS); err != nil {
		glog.V(7).Infof("tx.Trx is not supported %v", err)
		return nil, ""
	}
	actions = trxS.TX.Actions
	return actions, trxS.ID
}

type TRX struct {
	ID string      `json:"id"`
	TX Transaction `json:"transaction"`
}

type Action struct {
	Account string `json:"account"`
	Name    string `json:"name"`
	Data    TX     `json:"data"`
}

func (a Action) Coin() string {
	switch a.Account {
	case "eosio.token":
		return "EOS"
	case conf.C.GetString("eos.TokenAccount"):
		return conf.C.GetString("eos.TokenSymbol")
	default:
		// glog.Warningf("not supported coin %s", a.Account)
		return ""
	}
}

func (a Action) IsTransfer() bool {
	if a.Name == "transfer" {
		return true
	}
	return false
}

type Transaction struct {
	Actions []Action `json:"actions"`
}

// Tx dave to DB
type Tx struct {
	Id       int64   `gorm:"PRIMARY_KEY" json:"id"`
	TxID     string  `gorm:"size:64" json:"hash"`
	BlockNum uint32  `json:"blocknum"`
	From     string  `json:"from"`
	To       string  `json:"to"`
	Amount   float64 `json:"amount"`
	CoinID   int     `json:"coinID"`
	Memo     string  `json:"memo"`
	// 判断是否待处理
	Status    int8  `gorm:"index:status" json:"status"`
	TimeMills int64 `json:"timestamp"`
	// 提取属于哪一期游戏
	TimeMintue int64 `gorm:"index:timemintue" json:"time_mintue"`
}

// AddTx insert a new Tx into database and returns
// last inserted Id on success.
func AddTx(m *Tx) (err error) {
	d := db.Create(m)
	return d.Error
}

func GetTxsInGame(time int64, status int) ([]Tx, error) {
	var txes []Tx
	db.Where("time_mintue = ? and status = ?", time, status).Find(&txes)
	return txes, db.Error
}

// UpdateTx updates Tx by Id and returns error if
// the record to be updated doesn't exist
func UpdateTxById(m *Tx) (err error) {
	db.Save(m)
	return db.Error
}
