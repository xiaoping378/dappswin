package app

import "github.com/gin-gonic/gin"

// Register 注册user相关路由
func Register(router *gin.RouterGroup) {
	router.GET("/ws", serveWs)
	router.POST("/chain/get_currency_balance", getCurrencyBalance)

	// for platfom
	router.GET("/news", getNews)
	router.POST("/user/date", dateUser)
	router.POST("/user/page", pageUser)
	router.POST("/user/bind", bindUser)
	router.POST("/bonus/pool", bonusPool)
	router.GET("/bonus/stats", bonusStats)
	router.GET("/arena", arenaStatus)
	router.GET("/rank/stats_per_day", rankPerDay)
	router.POST("/lucky/submit", openLucky)
	router.POST("/lucky/status", getLucky)

	// for game
	router.POST("/tx/page_tx", pageTxes)
	router.POST("/game/page_lottery", pageLottery)
}
