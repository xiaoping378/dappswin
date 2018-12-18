package models

import "github.com/jinzhu/gorm"

type ICO struct {
	gorm.Model
	Hash      string `grom:"index"`
	Account   string
	Amount    uint64
	Status    int
	TimeMills int64
}

func AddIcoRecord(m *ICO) error {
	d := db.Create(m)
	return d.Error
}
