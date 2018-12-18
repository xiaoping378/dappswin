package app

import (
	"dappswin/models"

	"github.com/gin-gonic/gin"
)

// UserRegister 注册user相关路由
func UserRegister(router *gin.RouterGroup) {
	router.POST("/page_date_summary", getUser)
	router.POST("/page_user", pageUser)
}

func getUser(c *gin.Context) {
	c.JSON(200, "haha")
}

func pageUser(c *gin.Context) {
	c.JSON(200, "haha")
}

type UserMsg struct {
	Count int64         `json:"count"`
	Data  []models.User `json:"data"`
	Error interface{}   `json:"error"`
}
