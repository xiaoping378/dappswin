package models

import (
	"github.com/jinzhu/gorm"
)

// User 用户返佣系统
type User struct {
	gorm.Model
	Name          string  `json:"name" gorm:"size:12;unique"`
	PName         string  `json:"-" gorm:"size:12;index"`
	PNames        string  `json:"-"`
	Level         uint8   `json:"-"`
	ChildrenCount int     `json:"children_count"`
	TotalBet      float64 `json:"total_bet"`
	TotalRebate   float64 `json:"total_rebate"`
	Bet           float64 `json:"-"`
	Rebate        float64 `json:"-"`
	// ChildrenBet    uint64 `json:"-"`
	// ChildrenRebate uint32 `json:"-"`
}

// func updateParent(name string, bet uint64, isNew bool) {
// 	user := &User{}
// 	if db.Where("name = ?", name).First(user).RecordNotFound() {
// 		user = &User{Name: name}
// 	}
// 	if isNew {
// 		user.ChildrenCount++
// 	}

// 	user.TotalBet += bet
// 	// user.ChildrenBet += bet
// 	user.Level = getLevel(uint64(user.TotalBet / 1e4))

// 	if user.ID == 0 {
// 		db.Create(user)
// 	} else {
// 		db.Model(user).Update(user)
// 	}
// }

// func UpdateUserInfo(txmsg *Message) {
// 	var isNew bool
// 	user := &User{}
// 	tx, _ := txmsg.Data.(TX)
// 	if db.Where("name = ?", tx.From).First(user).RecordNotFound() {
// 		isNew = true
// 	}
// 	amout := tx.Amount()
// 	user.Bet += amout
// 	_, _, pAccount := ResolveMemo(tx.Memo)

// 	user.Name = tx.From
// 	if pAccount != "" {
// 		user.PName = pAccount
// 		updateParent(user.PName, user.Bet, isNew)
// 	}

// 	if !strings.Contains(user.PNames, user.PName) {
// 		user.PNames += ("," + user.PName)
// 	}
// 	user.TotalBet += amout
// 	user.Level = getLevel(uint64(user.TotalBet / 1e4))
// 	glog.V(7).Infof("User is %#v", user)

// 	if user.ID == 0 {
// 		db.Create(user)
// 	} else {
// 		db.Model(user).Update(user)
// 	}

// }
