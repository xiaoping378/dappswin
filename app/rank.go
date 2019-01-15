package app

import "github.com/gin-gonic/gin"

type RankPerDay struct {
	ID     int     `json:"id,omitempty"`
	User   string  `json:"user,omitempty"`
	Amount float64 `json:"amount,omitempty"`
	Reward float64 `json:"reward,omitempty"`
}

func rankPerDay(c *gin.Context) {
	c.JSON(NewMsg(200, &[]RankPerDay{
		{1, "wudixiaoping1", 200.1, 2},
		{2, "xiaopingeos1", 100.1, 1},
		{3, "aiaopingeos1", 90.1, 0.9},
		{4, "biaopingeos1", 80.1, 0.8},
		{5, "ciaopingeos1", 70.1, 0.7},
		{6, "diaopingeos1", 60.1, 0.6},
		{7, "eiaopingeos1", 50.1, 0.5},
		{8, "fiaopingeos1", 40.1, 0.4},
		{9, "giaopingeos1", 30.1, 0.3},
		{10, "hiaopingeos1", 20.1, 0.2},
	}))
}
