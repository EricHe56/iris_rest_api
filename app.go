package main

import (
	_ "fmt"
	"github.com/kataras/iris/v12"
	"iris_rest_api/api"
	. "iris_rest_api/init"
	"iris_rest_api/routes_code_creator"
	_ "iris_rest_api/utils"
	_ "reflect"
	_ "regexp"
	_ "strings"
)

func main() {

	isRuningForCreateRouter := routes_code_creator.CreateApiRoutesCode(
		// TODO: add your api here
		api.AdminApi{},
		api.FaqApi{},
		api.UserApi{},
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
			DisableBodyConsumptionOnUnmarshal: true, // keep body
			DisableAutoFireStatusCode:         false,
			TimeFormat:                        "2006-01-02T15:04:05.000Z07:00", // "Mon, 02 Jan 2006 15:04:05.000 GMT",
			Charset:                           "UTF-8",
		}))
	}
}
