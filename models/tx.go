package models

import (
	"strconv"
	"strings"
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

type TRX struct {
	ID string      `json:"id"`
	TX Transaction `json:"transaction"`
}

type Action struct {
	Account string `json:"account"`
	Name    string `json:"name"`
	Data    TX     `json:"data"`
}

type Transaction struct {
	Actions []Action `json:"actions"`
}

// Tx dave to DB
type Tx struct {
	Id       int64  `gorm:"PRIMARY_KEY"`
	TxID     string `gorm:"size:64"`
	BlockNum uint32
	From     string `json:"from"`
	To       string `json:"to"`
	Quantity string `json:"quantity"`
	Memo     string `json:"memo"`
	// 判断是否待处理
	Status int8 `gorm:"index:status"`
	Time   int64
	// 提取属于哪一期游戏
	TimeMintue int64 `gorm:"index:timemintue"`
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
