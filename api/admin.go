package api

import (
	"crypto/md5"
	"encoding/hex"
	jwt "github.com/gbrlsnchs/jwt/v3"
	"gopkg.in/mgo.v2/bson"
	. "iris_rest_api/init"
	. "iris_rest_api/models"
	"net/http"
	"time"
)

type AdminApi struct {
	CodeDescription map[int]string
	W               http.ResponseWriter
	R               *http.Request
}

type Admin struct {
	ID           bson.ObjectId `json:"id"  bson:"_id" q:",编号"`
	Name         string        `json:"name" bson:"name" q:",名称"`
	Introduction string        `json:"introduction" bson:"introduction" q:",介绍"`
	Avatar       string        `json:"avatar" bson:"avatar" q:",头像"`
	Role         []string      `json:"role" bson:"role" q:",角色"`
	Password     string        `json:"password" bson:"password" q:",密码"`
	Ctime        int64         `json:"ctime" bson:"ctime" q:",创建时间"`
}

type CustomPayload struct {
	jwt.Payload
	AdminUser Admin `json:"admin,omitempty"`
}

var hs = jwt.NewHS256([]byte("WuLong"))

func (x AdminApi) Login(req struct {
	UserName string `json:"username" bson:"username" q:",用户名"`
	Password string `json:"password" bson:"password" q:",密码"`
}) (code int, data struct {
	Token string `json:"token" bson:"token" q:",令牌"`
}, e error) {
	session := CloneSession()
	defer session.Close()

	// copy req
	req1 := &req
	req2 := *req1
	reqNew := &req2

	md5Ctx := md5.New()
	md5Ctx.Write([]byte("xg6-2]<PgU9@?s{z" + reqNew.Password))
	cipherStr := md5Ctx.Sum(nil)
	reqNew.Password = hex.EncodeToString(cipherStr)

	var curAdmin Admin
	selector := bson.M{"$and": []bson.M{bson.M{"name": reqNew.UserName}, bson.M{"password": reqNew.Password}}}
	e = session.DB(DB).C(C_ADMIN).Find(selector).One(&curAdmin)

	if e == nil {
		now := time.Now()
		curAdmin.Password = ""
		pl := CustomPayload{
			Payload: jwt.Payload{
				Issuer:         "Go",
				Subject:        "someone",
				Audience:       jwt.Audience{"https://golang.org", "https://jwt.io"},
				ExpirationTime: jwt.NumericDate(now.Add(24 * time.Hour)),
				NotBefore:      jwt.NumericDate(now.Add(30 * time.Minute)),
				IssuedAt:       jwt.NumericDate(now),
				JWTID:          "Admin",
			},
			AdminUser: curAdmin,
		}

		newToken, e := jwt.Sign(pl, hs)
		if e != nil {
			// ...
			code = 50300
			Info.Println(x.R.RequestURI, req, code, e)
		} else {
			code = 0
			data.Token = string(newToken)
			_ = RedisSetString(data.Token, "1", 8*3600)
		}
	} else {
		code = 50300
		Info.Println(x.R.RequestURI, req, code, e)
	}

	return
}

func (x AdminApi) Info(req struct {
	Token string `json:"token" q:", 令牌"`
}) (code int, data Admin, e error) {

	var pl CustomPayload
	hd, err := jwt.Verify([]byte(req.Token), hs, &pl)
	if err == nil {
		Info.Println(hd)

		session := CloneSession()
		defer session.Close()

		var curAdmin Admin
		selector := bson.M{"$and": []bson.M{bson.M{"_id": pl.AdminUser.ID}, bson.M{"name": pl.AdminUser.Name}}}
		e = session.DB(DB).C(C_ADMIN).Find(selector).One(&curAdmin)
		if e == nil {
			code = 0
			curAdmin.ID = ""
			curAdmin.Password = ""
			data = curAdmin
		}
	} else {
		code = 50300
	}

	return
}

func (x AdminApi) Logout(req Admin) (code int, data string, e error) {
	xToken := x.R.Header.Get("X-Token")
	_ = RedisDelString(xToken)
	code = 0
	data = "success"
	return
}

func (x AdminApi) Insert(req Admin) (code int, data struct {
	ID bson.ObjectId `json:"id" q:",返回id"`
}, e error) {
	session := CloneSession()
	defer session.Close()
	req.ID = bson.NewObjectId()
	//req.Stime = time.Now().Unix()
	//req.Etime = time.Now().Unix()

	e = session.DB(DB).C(C_ADMIN).Insert(req)
	if e != nil {
		Info.Println(e)
		code = 50300
	} else {
		code = 0
		data.ID = req.ID
	}

	Info.Println(x.R.RequestURI, req, code, data)
	return
}

func (x AdminApi) Replace(req Admin) (code int, data struct{}, e error) {
	session := CloneSession()
	defer session.Close()

	selector := bson.M{"_id": req.ID}
	e = session.DB(DB).C(C_ADMIN).Update(selector, req)
	if e != nil {
		Info.Println(e)
		code = 50300
		return
	} else {
		code = 0
		return
	}
	return
}

func (x AdminApi) Page(req PageData) (code int, data struct {
	List  []Admin `json:"list" q:",列表"`
	Total int     `json:"total" q:",总数"`
}, e error) {
	session := CloneSession()
	defer session.Close()
	sorter := toSortor(req.Order)
	if len(sorter) > 0 {
		data.Total, e = session.DB(DB).C(C_ADMIN).Find(toSelector(req.Filter)).Count()
		e = session.DB(DB).C(C_ADMIN).Find(toSelector(req.Filter)).Sort(sorter).Skip(req.Offset).Limit(req.Size).All(&data.List)
	} else {
		data.Total, e = session.DB(DB).C(C_ADMIN).Find(toSelector(req.Filter)).Count()
		e = session.DB(DB).C(C_ADMIN).Find(toSelector(req.Filter)).Skip(req.Offset).Limit(req.Size).All(&data.List)
	}
	if e != nil {
		Info.Println(e)
		code = 50300
		return
	} else {
		code = 0
		return
	}
	return
}

//##AdminApi.Delete: 用于XXXXX
//##AdminApi.Delete: author: <a href="mailto:myemail@abc.com">myemail@abc.com</a>
//##AdminApi.Delete: 可以添加多行说明，并且可以使用少量html
func (x AdminApi) Delete(req struct {
	ID []bson.ObjectId `json:"id" q:",要删除的ID数组"`
}) (code int, data struct{}, e error) {
	session := CloneSession()
	defer session.Close()
	session.DB(DB).C(C_ADMIN).RemoveAll(bson.M{"_id": bson.M{"$in": req.ID}})
	code = 0
	return
}

func (x AdminApi) Get(req struct {
	ID []bson.ObjectId `json:"id" q:",要查询的ID数组"`
}) (code int, data struct {
	List []Admin `json:"list" q:",列表"`
	Size int     `json:"size" q:",个数"`
}, e error) {
	session := CloneSession()
	defer session.Close()
	session.DB(DB).C(C_ADMIN).Find(bson.M{"_id": bson.M{"$in": req.ID}}).All(&data.List)
	data.Size = len(data.List)
	code = 0
	return
}

func (x AdminApi) GetCur(req struct {
	Curtime int64  `json:"curtime" q:",当前时间"`
	Lang    string `json:"lang" q:",语言"`
}) (code int, data struct {
	List []Admin `json:"list" q:",列表"`
	Size int     `json:"size" q:",个数"`
}, e error) {
	session := CloneSession()
	defer session.Close()
	session.DB(DB).C(C_ADMIN).Find(bson.M{"_stime": bson.M{"$lt": req.Curtime}, "_etime": bson.M{"$gt": req.Curtime}, "lang": req.Lang}).All(&data.List)
	data.Size = len(data.List)
	code = 0
	return
}
