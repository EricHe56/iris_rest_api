package main

import (
	"github.com/gavv/httpexpect"
	"github.com/kataras/iris/v12/httptest"
	. "iris_rest_api/init"
	"iris_rest_api/models"
	"testing"
)

// 用于外部访问测试生成文档，需要运行dev参数
func TestApiDocByRequest(t *testing.T) {
	//url := "http://127.0.0.1:8080"

	//code, body := utils.HttpRequest("POST", url+"/json5?v=100", "{\"name\":\"Tom西\", \"age\": 123}")
	//if code != 200 || body != "{\"message\":\"pong\"}" {
	//	t.Error("Failed: ", code, body)
	//}
	return
}

// 内部测试直接生成文档，不需要运行参数支持
func TestApi2Doc(t *testing.T) {
	DevMode = true
	BuildApiDoc = true
	app := NewApp()
	var e *httpexpect.Expect = httptest.New(t, app)

	//e.GET("/json2").Expect().Status(httptest.StatusOK).
	//	JSON().Equal(map[string]interface{}{"message": "pong"})

	//e.POST("/hello3").WithFormField("username", "kataras").Expect().Status(httptest.StatusOK).
	//	Body().Equal("{\"message\":\"pong\"}")

	e.POST("/reqbody").WithJSON(models.D_User{Name: "kataras"}).Expect().Status(httptest.StatusOK).
		JSON().Equal(map[string]interface{}{"message": "pong"})

	// give time to "yaag" to generate the doc, 5 seconds are more than enough
	//time.Sleep(5 * time.Second)
}
