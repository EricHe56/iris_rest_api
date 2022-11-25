package init

import (
	"context"
	_ "fmt"
	"github.com/astaxie/beego/config"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	_ "iris_rest_api/utils"
	"log"
	"math/rand"
	_ "reflect"
	_ "regexp"
	_ "strings"
	"time"
)

const TEST_DOC_NAME = "./doc/_testDoc.html"

const TEST_DOC_TITLE = "Iris Rest Api"

const (
	//DB      = "test"
	C_ADMIN = "admin"
	C_FAQ   = "faq"
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

// var GlobalMgoSession *mgo.Session
var GlobalMgoClient *mongo.Client
var DevMode bool = false
var BuildApiDoc bool = false
var IniConfiger config.Configer
var JsonOptions = iris.JSON{
	StreamingJSON: false,
	UnescapeHTML:  false,
	Indent:        "",
	Prefix:        "",
}

var DB = "test"

func Pre_Handler(ctx iris.Context) {
	//ctx.Application().Logger().Println("Before Handler--Method: ", ctx.Method(), "Runs before %s", ctx.Path())
	path := ctx.Path()
	if path != "/admin/login" {
		xToken := ctx.GetHeader("X-Token")
		if !RedisIsExist(xToken) {
			_, _ = ctx.JSON(ResponseBody{
				Code:    50008, //50008: Illegal token, re-login
				Message: "X-Token Error",
				Data:    "",
			}, JsonOptions)
			return
		} else {
			value, _ := RedisGetString(xToken)
			ttl, _ := RedisTTL(xToken)
			StdPrint.Info(value, ttl)
		}
	}

	bytes, _ := ctx.GetBody()
	ctx.Values().Set("body", string(bytes))
	rnd := rand.Float64()
	ctx.Values().Set("rnd", rnd)
	// Log can be Warn, Error, Debug, Fatal
	StdPrint.Info(rnd, " Pre-Handler: ", ctx.Method(), "|", ctx.Request().RemoteAddr, "=>", ctx.Request().Host+ctx.Request().RequestURI, "\t", string(bytes))
	ctx.Next()
}

//func CloneSession() *mgo.Session {
//	return GlobalMgoSession.Clone()
//}

func CreateIndex(collection string, key string) (idxName string, err error) {
	ctxMongo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	idxModel := mongo.IndexModel{
		Keys: bson.D{
			{key, 1},
		},
	}
	idxName, err = GlobalMgoClient.Database(DB).Collection(collection).Indexes().CreateOne(ctxMongo, idxModel)
	return
}
