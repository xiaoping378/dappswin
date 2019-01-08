package conf

import (
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/viper"
)

// C read from etc < $home < $pwd < env
var C = viper.New()

func Init() {

	C.AddConfigPath("/etc/dappswin/")
	C.AddConfigPath("$HOME/.dappswin")
	C.AddConfigPath(".")
	C.SetConfigName("dappswin")
	C.SetConfigType("toml")

	C.SetEnvPrefix("DAPPSWIN")
	C.AutomaticEnv()
	// support read nested key from env: eth.host = DAPPSWIN_ETH_HOST
	C.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := C.ReadInConfig(); err != nil {
		glog.Error(err.Error())
	}
	glog.Info("Using config file:", C.ConfigFileUsed())

	loadDefaultConfig()

}

func loadDefaultConfig() {
	C.SetDefault("eth.host", "127.0.0.1:8545")
	C.SetDefault("eth.fromBLock", 0)
	C.SetDefault("redis.host", "127.0.0.1:6379")
	C.SetDefault("redis.password", "")
	C.SetDefault("gin.host", ":8378")
	// Cfg.SetDefault("mongo.host", "127.0.0.1:27017")
}
