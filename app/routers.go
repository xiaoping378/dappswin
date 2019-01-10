package app

import "github.com/gin-gonic/gin"

// Register 注册user相关路由
func Register(router *gin.RouterGroup) {
	router.GET("/ws", serveWs)
	router.POST("/chain/get_currency_balance", getCurrencyBalance)

	router.POST("/page_date_summary", dateUser)
	router.POST("/page_user", pageUser)

	router.POST("/tx/page_tx", pageTxes)
	router.POST("/game/page_lottery", pageLottery)
}
