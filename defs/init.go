package defs

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
)

func init() {
	Log()
}

//格式化日志
func Log() {
	//配置
	log := logs.NewLogger(10000)
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	logConfig := fmt.Sprintf(`{"filename":"log/%d/%d/%d-%d-%d.log"}`, year, month,  year, month, day)

	//测试
	log.SetLogger("file", logConfig)
	log.EnableFuncCallDepth(true)
	log.Async()
}
