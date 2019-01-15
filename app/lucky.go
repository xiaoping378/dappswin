package app

import "github.com/gin-gonic/gin"

type luckyPost struct {
	Name string `json:"name" binding:"required,max=12"`
}

type luckyRsp struct {
	LuckyCode int  `json:"lucky_code"`
	Count     int  `json:"count"`
	IsAble    bool `json:"is_able"`
}

func openLucky(c *gin.Context) {
	body := &luckyPost{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(NewMsg(400, "输入参数有误"))
		return
	}

	c.JSON(NewMsg(200, &luckyRsp{2, 2, true}))
}

func getLucky(c *gin.Context) {
	body := &luckyPost{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(NewMsg(400, "输入参数有误"))
		return
	}

	c.JSON(NewMsg(200, &luckyRsp{2, 1, true}))
}
