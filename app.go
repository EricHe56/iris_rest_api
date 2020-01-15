package main

import (
	_ "fmt"
	"github.com/kataras/iris/v12"
	"iris_template/api"
	. "iris_template/init"
	"iris_template/routes_code_creator"
	_ "iris_template/utils"
	_ "reflect"
	_ "regexp"
	_ "strings"
)

func main() {

	isRuningForCreateRouter := routes_code_creator.CreateApiRoutesCode(
		// 这里添加新的接口类
		api.UserApi{},
		api.FaqApi{},
	)
	if !isRuningForCreateRouter {
		App = NewApp()
		loadRouteHandlers()
		port := IniConfiger.String("httpport")
		_ = App.Run(iris.Addr(":"+port), iris.WithConfiguration(iris.Configuration{ //默认配置:
			DisableStartupLog:                 false,
			DisableInterruptHandler:           false,
			DisablePathCorrection:             false,
			EnablePathEscape:                  false,
			FireMethodNotAllowed:              false,
			DisableBodyConsumptionOnUnmarshal: false,
			DisableAutoFireStatusCode:         false,
			TimeFormat:                        "2006-01-02T15:04:05.000Z07:00", // "Mon, 02 Jan 2006 15:04:05.000 GMT",
			Charset:                           "UTF-8",
		}))
	}
}
