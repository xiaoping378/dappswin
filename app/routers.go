package app

import "github.com/gin-gonic/gin"

// AppRegister 注册user相关路由
func AppRegister(router *gin.RouterGroup) {
	router.GET("/ws", serveWs)

	router.POST("/page_date_summary", dateUser)
	router.POST("/page_user", pageUser)

	router.POST("/chain/get_currency_balance", getCurrencyBalance)

	router.POST("/tx/page_tx", pageTxes)
}
