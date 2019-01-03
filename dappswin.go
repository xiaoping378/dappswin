package main

import (
	"dappswin/app"
	"dappswin/common"
	"dappswin/conf"
	"dappswin/database"
	"dappswin/logs"
	"dappswin/models"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

// TODO: wss

func main() {
	defer database.Close()
	defer glog.Flush()

	logs.Init()
	conf.Init()
	common.Init()
	database.Init()
	models.Init()
	app.Init()

	go app.ResolveRoutine()

	r := gin.Default()
	api := r.Group("/api")
	app.WSRegister(r.Group("/"))
	app.UserRegister(api)
	app.EosRegister(api)

	r.Use(static.Serve("/", static.LocalFile("./views", true)))

	// server := &http.Server{
	// 	Addr:    ":" + conf.C.GetString("gin.port"),
	// 	Handler: r,
	// }
	// gracehttp.Serve(server)
	// TODO: enable grace + autotls
	glog.Error(autotls.Run(r, "dappswin.io", "www.dappswin.io"))

}
