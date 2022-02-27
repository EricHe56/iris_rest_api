// +build create_router

package routes_code_creator

import (
	"encoding/json"
	"fmt"
	"iris_rest_api/models"
	"iris_rest_api/utils"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type ApiRouteInfo struct {
	ApiType                string // api类型 如： api.UserApi
	ApiFunction            string // api函数名称
	ApiFunctionDescription string // api函数说明
	ReqType                string // api函数body请求参数结构类型，对应函数的req声明
	RtnDataType            string // api函数返回data参数结构类型，对应函数返回的data声明
	RoutePath              string // api函数路径
}

type ApiInfoGroup struct {
	ApiType string
	ApiList []ApiRouteInfo
}

type JsonFieldInfo struct {
	Name        string
	Type        string
	Required    bool
	Description string
}

type StructInfo struct {
	StructName string
	Fields     []JsonFieldInfo
}

type ApiDocInfo struct {
	ApiDocTitle   string
	Time          string
	ApiInfoGroups []ApiInfoGroup
	StructInfos   []StructInfo
}

const TEMPLATE_FILE_NAME = "./routes_code_creator/apiRouteFuncTemplate.code"
const API_ROUTES_NAME = "./apiRoutes.go"
const API_DOC_NAME = "./doc/_apiDoc.html"
const API_DOC_TITLE = "Iris_Rest_Api接口文档"

var apiRoutesHeader = `
// +build !create_router

package main
import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"gopkg.in/mgo.v2/bson"
	"iris_rest_api/api"
	. "iris_rest_api/init"
	"iris_rest_api/models"
)

func loadRouteHandlers() {
{{routeHandlers}}
	return
}
`
var ApiFuncTemplateCode = ""
var ApiFuncRoutesCode = apiRoutesHeader
var RouterHandlerList = ""
var apiDoc_All = API_DOC_TITLE + "\n"
var apiDoc_Categories = "目录：\n"
var apiDoc_Content = ""
var apiInfoGroups []ApiInfoGroup = make([]ApiInfoGroup, 0)
var structInfos []StructInfo = make([]StructInfo, 0)
var apiDocInfo = ApiDocInfo{
	ApiDocTitle:   API_DOC_TITLE,
	Time:          time.Now().Format("2006-01-02T15:04:05.000Z07:00"),
	ApiInfoGroups: apiInfoGroups,
	StructInfos:   structInfos,
}
var apiDocStructAll = "X.公共数据结构\n"

var apiCodeAll = ""
var modelsCodeAll = ""

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

func getRtnDataType(methodType string) (rtnDataType string) {
	regexp1, _ := regexp.Compile("func\\([\\w\\W]*\\) \\(int, ")
	var str1 = regexp1.ReplaceAllString(methodType, "")
	regexp2, _ := regexp.Compile(", error\\)$")
	rtnDataType = regexp2.ReplaceAllString(str1, "")
	return
}

func customizeFunction(apiFuncTemplateCode string, apiRouteInfo ApiRouteInfo) {
	var funcCode = apiFuncTemplateCode
	var routeFunctionName = strings.ReplaceAll(apiRouteInfo.ApiType, "api.", "") + "_" + apiRouteInfo.ApiFunction
	funcCode = strings.ReplaceAll(funcCode, "{{apiType}}", apiRouteInfo.ApiType)
	var ReqType = apiRouteInfo.ReqType
	if apiRouteInfo.ReqType == "" {
		ReqType = "interface{}"
	}
	funcCode = strings.ReplaceAll(funcCode, "{{reqType}}", ReqType)
	funcCode = strings.ReplaceAll(funcCode, "{{apiFunction}}", apiRouteInfo.ApiFunction)
	funcCode = strings.ReplaceAll(funcCode, "{{routeFunctionName}}", routeFunctionName)
	var funcInput = "(apiX, req)"
	if apiRouteInfo.ReqType == "" {
		funcInput = "(apiX)"
	}
	funcCode = strings.ReplaceAll(funcCode, "{{apiFunctionInput}}", funcInput)
	ApiFuncRoutesCode += funcCode

	var routeHandeler = "	App.Handle(\"ANY\", \"/" + apiRouteInfo.RoutePath + "\", " + routeFunctionName + ")\n"
	RouterHandlerList += routeHandeler
	return
}

func registerApi(x interface{}) (apiInfoGroup ApiInfoGroup) {
	apiInfoGroup = ApiInfoGroup{
		ApiType: "",
		ApiList: make([]ApiRouteInfo, 0),
	}

	v := reflect.ValueOf(x)
	t := v.Type()
	apiInfoGroup.ApiType = t.String()

	fmt.Printf("apiType: %s\n", t)

	for i := 0; i < v.NumMethod(); i++ {
		var apiRouteInfo = ApiRouteInfo{
			ApiType:                "",
			ApiFunction:            "",
			ApiFunctionDescription: "",
			ReqType:                "",
			RtnDataType:            "",
			RoutePath:              "",
		}
		methType := v.Method(i).Type()
		apiRouteInfo.ApiType = t.String()
		apiRouteInfo.ApiFunction = t.Method(i).Name
		apiRouteInfo.ApiFunctionDescription = getCodeComment(strings.Split(apiRouteInfo.ApiType, ".")[1] + "." + apiRouteInfo.ApiFunction)
		apiRouteInfo.ReqType = getReqType(methType.String())
		apiRouteInfo.RtnDataType = getRtnDataType(methType.String())
		var pathGroupName = strings.ReplaceAll(apiRouteInfo.ApiType, "api.", "")
		pathGroupName = strings.ReplaceAll(pathGroupName, "Api", "")
		apiRouteInfo.RoutePath = getRoutePath(pathGroupName + "_" + apiRouteInfo.ApiFunction)
		//apiRouteInfo.RoutePath = getRoutePath(apiRouteInfo.ApiFunction)
		customizeFunction(ApiFuncTemplateCode, apiRouteInfo)
		apiInfoGroup.ApiList = append(apiInfoGroup.ApiList, apiRouteInfo)
		fmt.Printf("func (%s) %s%s\n", t.Name(), t.Method(i).Name,
			strings.TrimPrefix(methType.String(), "func"))
	}
	return
}

func CreateApiRoutesCode(apiInterfaces ...interface{}) bool {
	apiCodeAll, _ = getAllCodeText("./api")
	//modelsCodeAll, _ = getAllCodeText("./models")
	apiFuncTemplateCode, err := utils.ReadFileInString(TEMPLATE_FILE_NAME)
	ApiFuncTemplateCode = apiFuncTemplateCode
	if err == nil {
		for i, _ := range apiInterfaces {
			var apiInfoGroup = registerApi(apiInterfaces[i])
			apiInfoGroups = append(apiInfoGroups, apiInfoGroup)
		}
		err = os.Remove(API_ROUTES_NAME)
		ApiFuncRoutesCode = strings.ReplaceAll(ApiFuncRoutesCode, "{{routeHandlers}}", RouterHandlerList)
		err = utils.WriteFile(API_ROUTES_NAME, ApiFuncRoutesCode)
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
	}
	create_doc_txt()
	create_doc_json()
	return true
}

func getModlesType() {
	var j = 0
	sections, offsets := models.Typelinks()
	for i, base := range sections {
		for _, offset := range offsets[i] {
			typeAddr := models.Add(base, uintptr(offset), "")
			typ := reflect.TypeOf(*(*interface{})(unsafe.Pointer(&typeAddr)))
			//fmt.Println(typ.String())
			typeName := typ.String()
			if strings.Contains(typeName, "*models.") || strings.Contains(typeName, "*api.") {
				j += 1
				var structInfo = StructInfo{
					StructName: typ.Elem().String(),
					Fields:     make([]JsonFieldInfo, 0),
				}
				apiDoc_Categories += "\tX." + strconv.Itoa(j) + ". " + typ.Elem().String() + "\n"
				apiDocStructAll += "\tX." + strconv.Itoa(j) + ". " + typ.Elem().String() + "\n"
				fmt.Println(typ.Elem().String())
				//fmt.Printf("%v", typ.Elem().Kind())
				s := typ.Elem()
				if s.Kind() != reflect.Struct {
					continue
				}
				//v := reflect.New(typ)
				for i := 0; i < s.NumField(); i++ { // s must struct
					fieldType := s.Field(i)
					var required = false
					var description = fieldType.Tag.Get("q")
					if strings.Contains(description, "required") {
						required = true
					}
					description = strings.Replace(description, "required,", "", 1)
					if strings.HasPrefix(description, ",") {
						description = strings.Replace(description, ",", "", 1)
					}
					var jsonFieldInfo = JsonFieldInfo{
						Name:        fieldType.Tag.Get("json"),
						Type:        fieldType.Type.String(),
						Required:    required,
						Description: description,
					}
					structInfo.Fields = append(structInfo.Fields, jsonFieldInfo)
					apiDocStructAll += "\t\t" + fieldType.Tag.Get("json") + "\t\t" + fieldType.Type.String() + "\t\t" + fieldType.Tag.Get("q") + "\n"
				}
				structInfos = append(structInfos, structInfo)
			}
		}
	}
}

func create_doc_txt() {
	for i, _ := range apiInfoGroups {
		apiDoc_Categories += strconv.Itoa(i+1) + ". " + apiInfoGroups[i].ApiType + "\n"
		apiDoc_Content += strconv.Itoa(i+1) + ". " + apiInfoGroups[i].ApiType + "\n"
		for j, _ := range apiInfoGroups[i].ApiList {
			apiDoc_Categories += "\t" + strconv.Itoa(i+1) + "." + strconv.Itoa(j+1) + ". " + apiInfoGroups[i].ApiList[j].ApiFunction + "\n"
			apiDoc_Content += "\t\t" + strconv.Itoa(i+1) + "." + strconv.Itoa(j+1) + ". " + apiInfoGroups[i].ApiList[j].ApiFunction + "\n"
			apiDoc_Content += "\t\t\t" + "路径： \n"
			apiDoc_Content += "\t\t\t\t/" + apiInfoGroups[i].ApiList[j].RoutePath + "\n"
			apiDoc_Content += "\t\t\t" + "Post请求body： \n"
			apiDoc_Content += "\t\t\t\t" + apiInfoGroups[i].ApiList[j].ReqType + "\n"
			apiDoc_Content += "\t\t\t" + "返回信息： \n"
			apiDoc_Content += "\t\t\t\t" + "code： int 返回错误代码，0为操作成功\n"
			apiDoc_Content += "\t\t\t\t" + "message： string 返回代码对应的提示信息\n"
			apiDoc_Content += "\t\t\t\t" + "data： \n"
			apiDoc_Content += "\t\t\t\t\t" + apiInfoGroups[i].ApiList[j].RtnDataType + "\n"
		}
	}
	apiDoc_Categories += "X. 附录：\n"
	getModlesType()
	apiDoc_All += apiDoc_Categories + "\n\n" + apiDoc_Content + "\n\n" + apiDocStructAll
	_ = os.Remove(API_DOC_NAME + ".txt")
	err := utils.WriteFile(API_DOC_NAME+".txt", apiDoc_All)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	return
}

func create_doc_json() {
	apiDocInfo = ApiDocInfo{
		ApiDocTitle:   API_DOC_TITLE,
		Time:          time.Now().Format("2006-01-02T15:04:05.000Z07:00"),
		ApiInfoGroups: apiInfoGroups,
		StructInfos:   structInfos,
	}
	bytesApiDocInfo, err := json.Marshal(apiDocInfo)
	if err == nil {
		_ = os.Remove(API_DOC_NAME + ".js")
		err := utils.WriteFile(API_DOC_NAME+".js", "var apiDoc = "+string(bytesApiDocInfo)+";")
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
	}

	return
}

func getAllCodeText(dirPath string) (codeText string, e error) {
	files, _, err := utils.GetFilesAndDirs(dirPath)
	if err != nil {
		e = err
		return
	}
	var fileText = ""
	for _, v := range files {
		fileText, err = utils.ReadFileInString(v)
		if err != nil {
			e = err
			return
		}
		codeText += "\n//" + v + "\n" + fileText
	}
	return
}

func getCodeComment(methodName string) (comments string) {
	var regexpStr = `(?ism)^//##` + methodName + `:[!-~ \pP\p{Han}\t]*$`
	// var regexpStr = `(?im)^//##` + methodName + `:.*$`
	regexp1, _ := regexp.Compile(regexpStr)
	strMatched := regexp1.FindAllString(apiCodeAll, -1)
	var tmpStr = ""
	for _, v := range strMatched {
		tmpStr = strings.ReplaceAll(v, "//##"+methodName+":", "")
		if strings.HasPrefix(strings.TrimSpace(tmpStr), "<") && strings.HasSuffix(tmpStr, ">") {
			comments += tmpStr
		} else {
			comments += tmpStr + "\n"
		}
	}
	fmt.Println(comments)
	return
}
