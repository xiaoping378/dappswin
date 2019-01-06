package models

type Game struct {
	Id        int64  `gorm:"PRIMARY_KEY"`
	Result    string `gorm:"size:6"`
	BlockNum  uint32
	TimeStamp int64
	// 游戏属于的哪个分钟段的
	GameMintue int64
}

// AddGame insert a new Game into database and returns
// last inserted Id on success.
func AddGame(g *Game) (err error) {
	d := db.Create(g)
	return d.Error
}

// GetGameByMintue retrieves Game by Id. Returns error if
// Id doesn't exist
func GetGameByMintue(mintue int64) (v *Game, err error) {
	v = &Game{GameMintue: mintue}
	db.Where("GameMintue = ?", "mintue").First(&v)
	return v, db.Error
}
