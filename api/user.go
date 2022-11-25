package api

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"iris_rest_api/models"
	"net/http"
)

type UserApi struct {
	CodeDescription map[int]string
	W               http.ResponseWriter
	R               *http.Request
	RequestBody     []byte
	RndData         float64
}

func (x UserApi) Get(req struct {
	Name string `json:"name" bson:"name" q:",名称"`
	Age  int    `json:"age" bson:"age" q:",年龄"`
}) (code int, data models.D_User, e error) {
	data.Name = req.Name
	data.Age = req.Age
	if req.Age > 100 {
		code = -1
	} else {
		code = 0
	}
	//data = req
	return
}

func (x UserApi) Put(req models.D_User) (code int, data models.D_User, e error) {
	data.Name = req.Name
	data.Age = req.Age
	if req.Age > 100 {
		code = -1
	} else {
		code = 0
	}
	data = req
	return
}

func (x UserApi) ReqDocBuildTest(req struct {
	List []struct {
		ID       string `json:"id"  bson:"_id" q:",线路编号"`
		Cover    string `json:"cover" bson:"cover" q:",封面图片"`
		Name     string `json:"name" bson:"name" q:",线路名称"`
		Abbr     string `json:"abbr" bson:"abbr" q:",线路简介"`
		PlanSets []struct {
			Index    int    `json:"index" bson:"index" q:",顺序号"`
			Name     string `json:"name" bson:"name" q:",日行程名称"`
			Duration int    `json:"duration" bson:"duration" q:",游览时间，单位分钟"`
			Length   int    `json:"length" bson:"length" q:",长度，单位米"`
			Plans    []struct {
				Index  int       `json:"index" bson:"index" q:",顺序号"`
				Name   string    `json:"name" bson:"name" q:",景区名称"`
				Images []string  `json:"images" bson:"images" q:",图片数组"`
				Abbr   string    `json:"abbr" bson:"abbr" q:",景区简介"`
				Loc    []float32 `json:"loc" bson:"loc" q:",经纬度信息 [经度-纬度]"`
			} `json:"plans" bson:"plans" q:",线路内容"`
		} `json:"plansets" bson:"plansets" q:",线路内容，每天一条记录"`
	} `json:"list"  q:",当前页列表"`
	Total int `json:"total"  q:",记录总数"`
}) (code int, data models.D_User, e error) {
	data = models.D_User{
		ID:    primitive.NewObjectID(),
		Name:  "",
		Age:   0,
		Ctime: 0,
		Mtime: 0,
	}
	return
}
