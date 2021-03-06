package models

type Game struct {
	Id        int64  `gorm:"PRIMARY_KEY" json:"id"`
	Result    string `gorm:"size:6" json:"result"`
	BlockNum  uint32 `json:"blocknum"`
	TimeStamp int64  `json:"timestamp"`
	// 游戏属于的哪个分钟段的
	GameMintue int64  `json:"game_mintue"`
	Content    string `json:"content"`
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
