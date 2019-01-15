package app

import "github.com/gin-gonic/gin"

type Arena struct {
	ArenaAmount    float64 `json:"arena_amount,omitempty"`
	TotalArenaSell float64 `json:"total_arena_sell,omitempty"`
	TodayArenaSell float64 `json:"today_arena_sell,omitempty"`
	TotalArenaBuy  float64 `json:"total_arena_buy,omitempty"`
	TodayArenaBuy  float64 `json:"today_arena_buy,omitempty"`
	Round          int     `json:"round,omitempty"`
	EndTime        int     `json:"end_time,omitempty"`
	LatestPrice    float64 `json:"latest_price,omitempty"`
	LatestUser     string  `json:"latest_user,omitempty"`
	LatestTime     string  `json:"latest_time,omitempty"`
}

func arenaStatus(c *gin.Context) {
	c.JSON(NewMsg(200, &Arena{
		0.002,
		231.2,
		1.2,
		200234,
		123.2,
		17,
		1003234,
		234.1,
		"wudipingeos2",
		"2019/01/12 09:34:34",
	}))
}
