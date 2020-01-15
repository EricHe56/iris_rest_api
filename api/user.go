package api

import (
	"iris_rest_api/models"
	"net/http"
)

type UserApi struct {
	CodeDescription map[int]string
	W               http.ResponseWriter
	R               *http.Request
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
