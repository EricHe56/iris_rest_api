package api

import (
	//"encoding/json"
	//"math/rand"
	"gopkg.in/mgo.v2/bson"
	. "iris_rest_api/init"
	. "iris_rest_api/models"
	"net/http"
	"time"
)

type FaqApi struct {
	CodeDescription map[int]string
	W               http.ResponseWriter
	R               *http.Request
}

func (x FaqApi) Insert(req Faq) (code int, data struct {
	ID bson.ObjectId `json:"id" q:",返回id"`
}, e error) {
	session := CloneSession()
	defer session.Close()
	req.ID = bson.NewObjectId()
	req.Ctime = time.Now().Unix()
	req.Mtime = time.Now().Unix()

	e = session.DB(DB).C(C_FAQ).Insert(req)
	if e != nil {
		Info.Println(e)
		code = 50300
	} else {
		code = 0
		data.ID = req.ID
	}
	return
}

func (x FaqApi) Replace(req Faq) (code int, data struct{}, e error) {
	session := CloneSession()
	defer session.Close()

	selector := bson.M{"$and": []bson.M{bson.M{"_id": req.ID}, bson.M{"_mtime": req.Mtime}}}
	e = session.DB(DB).C(C_FAQ).Update(selector, req)
	if e != nil {
		Info.Println(e)
		code = 50300
	} else {
		code = 0
	}
	return
}

func (x FaqApi) Page(req PageData) (code int, data struct {
	List  []Faq `json:"list" q:",列表"`
	Total int   `json:"total" q:",总数"`
}, e error) {
	session := CloneSession()
	defer session.Close()
	sorter := toSortor(req.Order)
	if len(sorter) > 0 {
		data.Total, e = session.DB(DB).C(C_FAQ).Find(toSelector(req.Filter)).Count()
		e = session.DB(DB).C(C_FAQ).Find(toSelector(req.Filter)).Sort(sorter).Skip(req.Offset).Limit(req.Size).All(&data.List)
	} else {
		data.Total, e = session.DB(DB).C(C_FAQ).Find(toSelector(req.Filter)).Count()
		e = session.DB(DB).C(C_FAQ).Find(toSelector(req.Filter)).Skip(req.Offset).Limit(req.Size).All(&data.List)
	}
	if e != nil {
		Info.Println(e)
		code = 50300
	} else {
		code = 0
	}
	return
}

func (x FaqApi) Delete(req struct {
	ID []bson.ObjectId `json:"id" q:",要删除的ID数组"`
}) (code int, data struct{}, e error) {
	session := CloneSession()
	defer session.Close()
	_, e = session.DB(DB).C(C_FAQ).RemoveAll(bson.M{"_id": bson.M{"$in": req.ID}})
	code = 0
	return
}

func (x FaqApi) Get(req struct {
	ID []bson.ObjectId `json:"id" q:",要查询的ID数组"`
}) (code int, data struct {
	List []Faq `json:"list" q:",列表"`
	Size int   `json:"size" q:",个数"`
}, e error) {
	session := CloneSession()
	defer session.Close()
	e = session.DB(DB).C(C_FAQ).Find(bson.M{"_id": bson.M{"$in": req.ID}}).All(&data.List)
	data.Size = len(data.List)
	code = 0
	return
}
