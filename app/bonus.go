package app

import "github.com/gin-gonic/gin"

type bonusPoolRsp struct {
	StartTime      int     `json:"start_time,omitempty"`
	ExpectEarnings float64 `json:"expect_earnings,omitempty"`
	EarningsPerCGG float64 `json:"earnings_per_cgg,omitempty"`
	BonusAmount    float64 `json:"bonus_amount,omitempty"`
	CoinName       string  `json:"coin_name"`
}

type bonusPoolPost struct {
	Name string `json:"name" binding:"required,max=12"`
}

func bonusPool(c *gin.Context) {
	body := &bonusPoolPost{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(NewMsg(400, "输入参数有误"))
		return
	}

	c.JSON(NewMsg(200, &[]bonusPoolRsp{{10 * 60 * 60, 0.0, 0.000024, 10.2, "EOS"},
		{10 * 60 * 60, 0.0, 0.000014, 100.2, "CGG"}}))
}

type bonusStatusRsp struct {
	TotalAmountEos float64 `json:"total_amount_eos,omitempty"`
	TotalAmountCGG float64 `json:"total_amount_cgg,omitempty"`
}

func bonusStats(c *gin.Context) {

	c.JSON(NewMsg(200, &bonusStatusRsp{1.2, 3002.1}))
}
