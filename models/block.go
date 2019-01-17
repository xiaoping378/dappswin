package models

import (
	"dappswin/database"
	"encoding/json"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func Init() {
	db = database.Db
	db.AutoMigrate(&Block{}, &User{}, &Tx{}, &Game{}, &ICO{}, &Stake{})
}

// Block ws send to hub
type Block struct {
	Hash string `json:"id"`
	Num  uint32 `json:"num" gorm:"PRIMARY_KEY"`
	Time int64  `json:"time"`
}

// Message make block msg to front
func (b Block) Message() []byte {

	m := Message{}
	m.Type = 0
	// same as func (m *Message) HandleTimeStamp()
	// b.Time = b.Time*2/1000 - 946684800
	b.Time = b.Time
	m.Data = b

	ms := []Message{}
	ms = append(ms, m)
	result, _ := json.Marshal(ms)
	return result
}

func (b *Block) LastLetter() string {
	if len(b.Hash) != 64 {
		glog.Errorf("Error on block ！！！！ %#v", b)
		return ""
	}
	return b.Hash[len(b.Hash)-1:]
}

func AddBlock(b *Block) (err error) {
	d := db.Create(b)
	return d.Error

}

func GetLastestBlock() (*Block, error) {
	blk := &Block{}
	d := db.Last(blk)
	return blk, d.Error
}

type Reward struct {
	Amount string `json:"amount"`
	ID     int    `json:"bookid"`
	Memo   string `json:"memo"`
	Status int    `json:"status"`
	To     string `json:"to"`
}

// Message ws send this.
type Message struct {
	Type      int         `json:"type"`
	BlockNum  uint32      `json:"blocknum,omitempty"`
	Hash      string      `json:"id,omitempty"`
	TimeMills int64       `json:"time,omitempty"`
	Data      interface{} `json:"data"`
}

// HandleTimeStamp 前端为了显示0.5所做的特殊处理
//
// (2142155917+946684800)/2*1000 == // 1544420358500
func (m *Message) HandleTimeStamp() {
	// m.Time = m.Time*2/1000 - 946684800

}
