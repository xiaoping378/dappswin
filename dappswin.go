package main

import (
	"dappswin/app"
	"dappswin/common"
	"dappswin/conf"
	"dappswin/database"
	"dappswin/logs"
	"dappswin/models"
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
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

	r.Static("lottery", "./views")

	server := &http.Server{
		Addr:    ":" + conf.C.GetString("gin.port"),
		Handler: r,
	}
	gracehttp.Serve(server)
}
