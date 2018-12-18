package database

import (
	"log"
	"os"

	"dappswin/conf"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

// Db 全局连接 TODO: ensure 启动了连接池
var Db *gorm.DB

// Init 初始化myql
func Init() {
	glog.Info("Connecting mysql ...")
	db, err := gorm.Open("mysql", getDSN())
	if err != nil {
		glog.Exitln("failed to connect database", err)
	}
	if err := db.DB().Ping(); err != nil {
		glog.Exitln("mysql 不可通")
	}
	// TODO output to glog.
	db.SetLogger(log.New(os.Stdout, "\r\n", 0))

	Db = db
	// Client.FlushDB()
}

// Close 关闭连接
func Close() {
	Db.Close()
}

func getDSN() string {
	user := conf.C.GetString("mysql.user") + ":" + conf.C.GetString("mysql.password")
	host := conf.C.GetString("mysql.host") + ":" + conf.C.GetString("mysql.port")
	dsn := user + "@tcp(" + host + ")/dappswin?charset=utf8&parseTime=True&loc=Local"
	return dsn
}
