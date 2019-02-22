package utils

import (
	"fmt"
	"github.com/go-ini/ini"
	"os"
)

var Config *ini.File
var CsrfExcept *ini.File
var RootPath string

func init()  {
	//RootPath项目绝对路径
	RootPath="/Users/baidu/data/golang/gopath/src/go-gin-mvc"
	var err error
	Config, err = ini.Load(RootPath+"/conf/config.ini");
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	CsrfExcept, err = ini.Load(RootPath+"/conf/csrf_except.ini");
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
}

