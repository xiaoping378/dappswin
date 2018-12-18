package common

import (
	"fmt"

	"github.com/golang/glog"
)

func Init() {
	glog.Info(fmt.Sprintf("Current version is %v (%v/%v)", CurrentVersion, BuildDate, BuildHash))
}

var versions = []string{
	"0.3",
}

var CurrentVersion string = versions[0]
var BuildDate string
var BuildHash string
