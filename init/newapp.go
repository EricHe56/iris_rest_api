//go:build !create_router
// +build !create_router

package init

import (
	"context"
	"github.com/astaxie/beego/config"
	"github.com/betacraft/yaag/yaag"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/natefinch/lumberjack.v2"
	"iris_rest_api/irisyaag"
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
		App.HandleDir("/_doc", "./doc")
	} else {
		iniFile = "conf/prod.conf"
	}

	if len(os.Args) > 2 && os.Args[2] == "doc" || DevMode { // 给test留入口
		BuildApiDoc = true
	}

	if DevMode && BuildApiDoc { // 仅在dev模式生成文档
		yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
			On:       true,
			DocTitle: TEST_DOC_TITLE,
			DocPath:  TEST_DOC_NAME,
			BaseUrls: map[string]string{"Production": TEST_DOC_TITLE, "Staging": time.Now().Format("2006-01-02T15:04:05.000Z07:00"), "Your Info": "Here"},
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
	Info.Println("app service starting...")

	// Redis 连接池
	var redisHost = IniConfiger.String("redis::host")
	var redisPort = IniConfiger.String("redis::port")
	var redisPwd = IniConfiger.String("redis::pwd")
	RedisClient = RedisCnnPool(redisHost, redisPort, redisPwd)

	// mongodb 连接池
	//var mongoHost = IniConfiger.String("mongodb::host")
	//var mongoPort = IniConfiger.String("mongodb::port")
	//var mongoUser = IniConfiger.String("mongodb::username")
	//var mongoPwd = IniConfiger.String("mongodb::password")
	var mongoUri = IniConfiger.String("mongodb::uri")
	DB = IniConfiger.String("mongodb::database")
	//dialInfo := &mgo.DialInfo{
	//	Addrs:    []string{mongoHost + ":" + mongoPort},
	//	Username: mongoUser,
	//	Password: mongoPwd,
	//}
	//
	//globalMgoSession, err := mgo.DialWithInfo(dialInfo)
	//if err != nil {
	//	StdPrint.Error("mongo connect error: ", err)
	//	Info.Println(err)
	//}
	//GlobalMgoSession = globalMgoSession
	//GlobalMgoSession.SetMode(mgo.Monotonic, true)
	//GlobalMgoSession.SetPoolLimit(300)
	//StdPrint.Info("mongo connected")

	// official driver
	ctxMongo, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()
	var optionsMongo = options.Client().ApplyURI(mongoUri)
	cltMongo, err := mongo.NewClient(optionsMongo)
	if err != nil {
		StdPrint.Info("mongo NewClient Error: ", err)
	}
	err = cltMongo.Connect(ctxMongo)
	if err != nil {
		StdPrint.Info("mongo Connect Error: ", err)
	}
	err = cltMongo.Ping(ctxMongo, readpref.Primary())
	if err != nil {
		StdPrint.Info("mongo Ping Error: ", err)
	}
	GlobalMgoClient = cltMongo
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
