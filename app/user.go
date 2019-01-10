package app

import (
	"dappswin/models"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

// Generated by https://quicktype.io

// UserRegister 注册user相关路由
func UserRegister(router *gin.RouterGroup) {
	router.POST("/page_date_summary", dateUser)
	router.POST("/page_user", pageUser)
}

type PageUserPost struct {
	PageIndex int    `json:"page_index" binding:"required,gt=0,lt=100"`
	OrderBy   string `json:"order_by" binding:"required,len=9"`
	PageSize  int    `json:"page_size" binding:"required,gt=0,lt=100"`
	PName     string `json:"pid" binding:"required,max=12"`
}

type PageUserRsp struct {
	Count int            `json:"count"`
	Data  []*models.User `json:"data"`
}

func pageUser(c *gin.Context) {
	body := &PageUserPost{}
	if err := c.ShouldBind(body); err != nil {
		c.JSON(NewMsg(400, "输入参数有误"))
		return
	}
	users := []*models.User{}
	user := &models.User{}
	var count int
	index := (body.PageIndex - 1) * body.PageSize

	db.Where(models.User{PName: body.PName}).Offset(index).Limit(body.PageSize - 1).Order(body.OrderBy + " desc").Find(&users).Count(&count)
	if db.Error != nil {
		glog.Error(db.Error)
		c.JSON(NewMsg(500, "系统内部错误"))
		return
	}
	// 这个业务需求需要加上自己 最前边送到前端展示
	db.Where(models.User{Name: body.PName}).First(user)
	users = append([]*models.User{user}, users...)
	// end

	c.JSON(NewMsg(200, &PageUserRsp{count, users}))
}

func dateUser(c *gin.Context) {
	c.JSON(200, "NULL")
}

type loginUserPost struct {
	Name  string
	PName string
}

func hasParent(name string) bool {
	return false
}

// Generated by https://quicktype.io

type DataUserPost struct {
	PageIndex int64    `json:"page_index" binding:"required,gt=0,lt=100"`
	PageSize  int64    `json:"page_size" binding:"required,gt=0,lt=100"`
	Date      []string `json:"date" binding:"required,len=2"`
	PName     string   `json:"pid" binding:"required,max=12"`
}
