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

//func NewApp() *iris.Application {
//
//	App = iris.New() // iris.Default()
//
//	customStdPrint := logger.New(logger.Config{
//		// Status displays status code
//		Status: true,
//		// IP displays request's remote address
//		IP: true,
//		// Method displays the http method
//		Method: true,
//		// Path displays the request path
//		Path: true,
//		// Query appends the url query to the Path.
//		Query: true,
//
//		// Columns: true,
//
//		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
//		// will be added to the logs.
//		MessageContextKeys: []string{"rnd"},
//
//		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
//		//MessageHeaderKeys: []string{"User-Agent"},
//	})
//
//	App.Use(recover.New())
//	App.Use(customStdPrint)
//	// 2006-01-02T15:04:05.000Z07:00
//	// "2006-02-01 15:04:05.000"
//	StdPrint = App.Logger().SetTimeFormat("2006-01-02T15:04:05.000Z07:00")
//	StdPrint.Info("App Server Starting ...")
//
//	iniFile := "conf/dev.conf"
//	if len(os.Args) > 1 && os.Args[1] == "dev" || DevMode { // 给test留入口
//		DevMode = true
//	} else {
//		iniFile = "conf/prod.conf"
//	}
//
//	if DevMode && BuildApiDoc { // 仅在dev模式生成文档
//		yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
//			On:       true,
//			DocTitle: "Iris_Template_Api",
//			DocPath:  "apiDoc.html",
//			BaseUrls: map[string]string{"Production": "My Api", "Staging": "abc"},
//		})
//		App.Use(irisyaag.New())
//	}
//
//	IniConfiger, _ = config.NewConfig("ini", iniFile)
//
//	// 日志
//	maxSize, _ := IniConfiger.Int("log::maxsize")
//	maxBackups, _ := IniConfiger.Int("log::maxfiles")
//	maxAge, _ := IniConfiger.Int("log::maxdays")
//
//	logFile := &lumberjack.Logger{
//		Filename:   IniConfiger.String("log::filename"),
//		MaxSize:    maxSize,
//		MaxBackups: maxBackups,
//		MaxAge:     maxAge,
//	}
//	Info = log.New(logFile, "[debug]", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
//	Info.Println("guide service starting...")
//
//	// Redis 连接池
//	RedisClient = &redis.Pool{
//		MaxIdle:     100,
//		MaxActive:   4000,
//		IdleTimeout: 180 * time.Second,
//		Dial: func() (redis.Conn, error) {
//			c, err := redis.Dial("tcp", IniConfiger.String("redis::host")+":"+IniConfiger.String("redis::port"), redis.DialPassword(IniConfiger.String("redis::pwd")))
//			if nil != err {
//				StdPrint.Error("redis connect error: ", err)
//				return nil, err
//			}
//			StdPrint.Info("redis connected")
//			return c, nil
//		},
//	}
//
//	// mongodb 连接池
//	dialInfo := &mgo.DialInfo{
//		Addrs:    []string{IniConfiger.String("mongodb::host") + ":" + IniConfiger.String("mongodb::port")},
//		Username: IniConfiger.String("mongodb::username"),
//		Password: IniConfiger.String("mongodb::password"),
//	}
//
//	globalMgoSession, err := mgo.DialWithInfo(dialInfo)
//	if err != nil {
//		StdPrint.Error("mongo connect error: ", err)
//		Info.Println(err)
//	}
//	GlobalMgoSession = globalMgoSession
//	GlobalMgoSession.SetMode(mgo.Monotonic, true)
//	GlobalMgoSession.SetPoolLimit(300)
//	StdPrint.Info("mongo connected")
//
//	// 	CORS 跨域
//	cors := IniConfiger.String("cors")
//	StdPrint.Info("CORS is " + cors)
//	if cors == "true" {
//		corsMiddleware := func(ctx iris.Context) {
//			ctx.Header("Access-Control-Allow-Origin", "*")
//			ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
//			ctx.Header("Access-Control-Allow-Credentials", "true")
//			ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Content-Type, ACCEPT, Authorization")
//			ctx.Next()
//		} // or	"github.com/iris-contrib/middleware/cors"
//		App.Use(corsMiddleware)
//		App.AllowMethods(iris.MethodOptions)
//	}
//
//	App.Use(Pre_Handler)
//
//	App.Handle("ANY", "/ping", func(ctx iris.Context) {
//		_, _ = ctx.JSON(iris.Map{"message": "pong"})
//		//bytes, _ := ctx.GetBody()
//		body := ctx.Values().GetString("body")
//		rnd, _ := ctx.Values().GetFloat64("rnd")
//		StdPrint.Info(rnd, " MainHandler: ", ctx.Method(), "|", ctx.Request().RemoteAddr, "=>", ctx.Request().Host+ctx.Request().RequestURI, "\t", body)
//	})
//
//	return App
//}

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
