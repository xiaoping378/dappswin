package app

import (
	"dappswin/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// Generated by https://quicktype.io

type PageUser struct {
	PageIndex int    `json:"page_index" binding:"required,gt=0,lt=100"`
	OrderBy   string `json:"order_by" binding:"required,len=9"`
	PageSize  int    `json:"page_size" binding:"required,gt=0,lt=100"`
	PID       string `json:"pid" binding:"required,len=12"`
}

// UserRegister 注册user相关路由
func UserRegister(router *gin.RouterGroup) {
	router.POST("/page_date_summary", dateUser)
	router.POST("/page_user", pageUser)
}

func pageUser(c *gin.Context) {
	pu := &PageUser{}
	if err := c.ShouldBind(pu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"count": 0, "data": nil, "error": err.Error()})
		return
	}
	users := []*models.User{}
	user := &models.User{}
	var count int
	index := (pu.PageIndex - 1) * pu.PageSize

	db.Where(models.User{PName: pu.PID}).Offset(index).Limit(pu.PageSize - 1).Order(pu.OrderBy + " desc").Find(&users).Count(&count)
	if db.Error != nil {
		glog.Error(db.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"count": 0, "data": nil, "error": "内部服务错误"})
		return
	}
	// 这个业务需求需要加上自己 最前边 送到前端展示
	db.Where(models.User{Name: pu.PID}).First(user)
	users = append([]*models.User{user}, users...)
	// end

	c.JSON(http.StatusOK, gin.H{"count": count, "data": users, "error": nil})
}

func dateUser(c *gin.Context) {
	c.JSON(200, "NULL")
}

type UserMsg struct {
	Count int64         `json:"count"`
	Data  []models.User `json:"data"`
	Error interface{}   `json:"error"`
}
