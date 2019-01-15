package app

import (
	"dappswin/conf"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// News 轮播图信息
type News struct {
	ID          int    `json:"id"`
	Language    string `json:"language"`
	Subject     string `json:"subject"`
	ImageSource string `json:"image_source"`
	ImageLink   string `json:"image_link"`
	Date        string `json:"date"`
}

func getNews(c *gin.Context) {

	data, err := ioutil.ReadFile(conf.C.GetString("general.newsSource"))
	if err != nil {
		c.JSON(NewMsg(http.StatusInternalServerError, "内部读取新闻错误"))
		return
	}
	news := &[]News{}
	if err := json.Unmarshal(data, news); err != nil {
		c.JSON(NewMsg(http.StatusInternalServerError, "内部读取新闻错误"))
		return
	}
	c.JSON(NewMsg(http.StatusOK, news))
}
