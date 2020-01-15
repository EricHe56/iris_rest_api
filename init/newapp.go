// +build !create_router

package init

import (
	"github.com/astaxie/beego/config"
	"github.com/betacraft/yaag/yaag"
	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"gopkg.in/mgo.v2"
	"gopkg.in/natefinch/lumberjack.v2"
	"iris_template/irisyaag"
	"log"
	"os"
	"time"
)

func NewApp() *iris.Application {

	App = iris.New() // iris.Default()

	customStdPrint := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		// Query appends the url query to the Path.
		Query: true,

		// Columns: true,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: []string{"rnd"},

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		//MessageHeaderKeys: []string{"User-Agent"},
	})

	App.Use(recover.New())
	App.Use(customStdPrint)
	// 2006-01-02T15:04:05.000Z07:00
	// "2006-02-01 15:04:05.000"
	StdPrint = App.Logger().SetTimeFormat("2006-01-02T15:04:05.000Z07:00")
	StdPrint.Info("App Server Starting ...")

	iniFile := "conf/dev.conf"
	if len(os.Args) > 1 && os.Args[1] == "dev" || DevMode { // 给test留入口
		DevMode = true
	} else {
		iniFile = "conf/prod.conf"
	}

	if DevMode && BuildApiDoc { // 仅在dev模式生成文档
		yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
			On:       true,
			DocTitle: "Iris_Template_Api",
			DocPath:  "apiDoc.html",
			BaseUrls: map[string]string{"Production": "My Api", "Staging": "abc"},
		})
		App.Use(irisyaag.New())
	}

	IniConfiger, _ = config.NewConfig("ini", iniFile)

	// 日志
	maxSize, _ := IniConfiger.Int("log::maxsize")
	maxBackups, _ := IniConfiger.Int("log::maxfiles")
	maxAge, _ := IniConfiger.Int("log::maxdays")

	logFile := &lumberjack.Logger{
		Filename:   IniConfiger.String("log::filename"),
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	}
	Info = log.New(logFile, "[debug]", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)
	Info.Println("guide service starting...")

	// Redis 连接池
	RedisClient = &redis.Pool{
		MaxIdle:     100,
		MaxActive:   4000,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", IniConfiger.String("redis::host")+":"+IniConfiger.String("redis::port"), redis.DialPassword(IniConfiger.String("redis::pwd")))
			if nil != err {
				StdPrint.Error("redis connect error: ", err)
				return nil, err
			}
			StdPrint.Info("redis connected")
			return c, nil
		},
	}

	// mongodb 连接池
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{IniConfiger.String("mongodb::host") + ":" + IniConfiger.String("mongodb::port")},
		Username: IniConfiger.String("mongodb::username"),
		Password: IniConfiger.String("mongodb::password"),
	}

	globalMgoSession, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		StdPrint.Error("mongo connect error: ", err)
		Info.Println(err)
	}
	GlobalMgoSession = globalMgoSession
	GlobalMgoSession.SetMode(mgo.Monotonic, true)
	GlobalMgoSession.SetPoolLimit(300)
	StdPrint.Info("mongo connected")

	// 	CORS 跨域
	cors := IniConfiger.String("cors")
	StdPrint.Info("CORS is " + cors)
	if cors == "true" {
		corsMiddleware := func(ctx iris.Context) {
			ctx.Header("Access-Control-Allow-Origin", "*")
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Content-Type, ACCEPT, Authorization")
			ctx.Next()
		} // or	"github.com/iris-contrib/middleware/cors"
		App.Use(corsMiddleware)
		App.AllowMethods(iris.MethodOptions)
	}

	App.Use(Pre_Handler)

	App.Handle("ANY", "/ping", func(ctx iris.Context) {
		_, _ = ctx.JSON(iris.Map{"message": "pong"})
		//bytes, _ := ctx.GetBody()
		body := ctx.Values().GetString("body")
		rnd, _ := ctx.Values().GetFloat64("rnd")
		StdPrint.Info(rnd, " MainHandler: ", ctx.Method(), "|", ctx.Request().RemoteAddr, "=>", ctx.Request().Host+ctx.Request().RequestURI, "\t", body)
	})

	return App
}
