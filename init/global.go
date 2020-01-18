package init

import (
	_ "fmt"
	"github.com/astaxie/beego/config"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"gopkg.in/mgo.v2"
	_ "iris_rest_api/utils"
	"log"
	"math/rand"
	_ "reflect"
	_ "regexp"
	_ "strings"
)

const TEST_DOC_NAME = "./doc/_testDoc.html"

const TEST_DOC_TITLE = "Iris Rest Api"

const (
	DB    = "test"
	C_FAQ = "faq"
)

var GlobalCodeDescription = map[int]string{
	0:  "OK",
	-1: "错误",
}

type ResponseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var App *iris.Application
var StdPrint *golog.Logger
var Info *log.Logger
var RedisClient *redis.Pool
var GlobalMgoSession *mgo.Session
var DevMode bool = false
var BuildApiDoc bool = false
var IniConfiger config.Configer
var JsonOptions = iris.JSON{
	StreamingJSON: false,
	UnescapeHTML:  false,
	Indent:        "",
	Prefix:        "",
}

func Pre_Handler(ctx iris.Context) {
	//ctx.Application().Logger().Println("Before Handler--Method: ", ctx.Method(), "Runs before %s", ctx.Path())
	bytes, _ := ctx.GetBody()
	ctx.Values().Set("body", string(bytes))
	rnd := rand.Float64()
	ctx.Values().Set("rnd", rnd)
	// Log can be Warn, Error, Debug, Fatal
	StdPrint.Info(rnd, " Pre-Handler: ", ctx.Method(), "|", ctx.Request().RemoteAddr, "=>", ctx.Request().Host+ctx.Request().RequestURI, "\t", string(bytes))
	ctx.Next()
}

func CloneSession() *mgo.Session {
	return GlobalMgoSession.Clone()
}

func CreateIndex(collection string, key string) (err error) {
	session := CloneSession()
	defer session.Close()
	index := mgo.Index{
		Key: []string{key}, // 索引字段， 默认升序,若需降序在字段前加-
		//Unique:     true,		// 唯一索引 同mysql唯一索引
		//DropDups:   true,		// 索引重复替换旧文档,Unique为true时失效
		//Background: true,		// 后台创建索引
		//Sparse:     true,		// Only index documents containing the Key fields
	}
	err = session.DB(DB).C(collection).EnsureIndex(index)
	return
}
