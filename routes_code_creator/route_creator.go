// +build create_router

package routes_code_creator

import (
	"fmt"
	"iris_template/utils"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type ApiRouteInfo struct {
	ApiType     string // api类型 如： api.UserApi
	ApiFunction string // api函数名称
	ReqType     string // api函数body请求参数结构类型，对应函数的req声明
}

const TEMPLATE_FILE_NAME = "./routes_code_creator/apiRouteFuncTemplate.code"
const API_ROUTES_NAME = "./apiRoutes.go"

var apiRoutesHeader = `package main
import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"gopkg.in/mgo.v2/bson"
	"iris_template/api"
	. "iris_template/init"
	"iris_template/models"
)

func loadRouteHandlers() {
{{routeHandlers}}
	return
}
`
var ApiFuncTemplateCode = ""
var ApiFuncRoutesCode = apiRoutesHeader
var RouterHandlerList = ""

func getRoutePath(routeFunctionName string) (routePath string) {
	strA := strings.Split(routeFunctionName, "")
	str1 := ""
	for i, _ := range strA {
		var strT = strA[i]
		if strA[i] >= "A" && strA[i] <= "Z" && i > 0 {
			strT = "_" + strA[i]
		}
		str1 += strT
	}
	str2 := strings.ReplaceAll(str1, "__", "/")
	routePath = strings.ToLower(str2)
	fmt.Printf("routePath: /%s\n", routePath)
	return
}

func getReqType(methodType string) (reqType string) {
	regexp1, _ := regexp.Compile("\\) \\(int, [\\w\\W]+, error\\)$")
	var str1 = regexp1.ReplaceAllString(methodType, "")
	reqType = strings.Replace(str1, "func(", "", 1)
	return
}

func customizeFunction(apiFuncTemplateCode string, apiRouteInfo ApiRouteInfo) {
	var funcCode = apiFuncTemplateCode
	var routeFunctionName = strings.ReplaceAll(apiRouteInfo.ApiType, "api.", "") + "_" + apiRouteInfo.ApiFunction
	funcCode = strings.ReplaceAll(funcCode, "{{apiType}}", apiRouteInfo.ApiType)
	funcCode = strings.ReplaceAll(funcCode, "{{reqType}}", apiRouteInfo.ReqType)
	funcCode = strings.ReplaceAll(funcCode, "{{apiFunction}}", apiRouteInfo.ApiFunction)
	funcCode = strings.ReplaceAll(funcCode, "{{routeFunctionName}}", routeFunctionName)
	ApiFuncRoutesCode += funcCode

	var routePath = getRoutePath(routeFunctionName)
	var routeHandeler = "	App.Handle(\"ANY\", \"/" + routePath + "\", " + routeFunctionName + ")\n"
	RouterHandlerList += routeHandeler
	return
}

func registerApi(x interface{}) {
	var apiRouteInfo = ApiRouteInfo{
		ApiType:     "",
		ApiFunction: "",
		ReqType:     "",
	}

	v := reflect.ValueOf(x)
	t := v.Type()
	apiRouteInfo.ApiType = t.String()
	fmt.Printf("apiType: %s\n", t)

	for i := 0; i < v.NumMethod(); i++ {
		methType := v.Method(i).Type()
		apiRouteInfo.ApiFunction = t.Method(i).Name
		apiRouteInfo.ReqType = getReqType(methType.String())
		customizeFunction(ApiFuncTemplateCode, apiRouteInfo)

		fmt.Printf("func (%s) %s%s\n", t.Name(), t.Method(i).Name,
			strings.TrimPrefix(methType.String(), "func"))
	}
}

func CreateApiRoutesCode(apiInterfaces ...interface{}) bool {
	apiFuncTemplateCode, err := utils.ReadFileInString(TEMPLATE_FILE_NAME)
	ApiFuncTemplateCode = apiFuncTemplateCode
	if err == nil {
		for i, _ := range apiInterfaces {
			registerApi(apiInterfaces[i])
		}
		err = os.Remove(API_ROUTES_NAME)
		ApiFuncRoutesCode = strings.ReplaceAll(ApiFuncRoutesCode, "{{routeHandlers}}", RouterHandlerList)
		err = utils.WriteFile(API_ROUTES_NAME, ApiFuncRoutesCode)
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
	}
	return true
}
