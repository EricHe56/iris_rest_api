package api

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	jwt "github.com/gbrlsnchs/jwt/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	. "iris_rest_api/init"
	. "iris_rest_api/models"
	"net/http"
	"time"
)

type AdminApi struct {
	CodeDescription map[int]string
	W               http.ResponseWriter
	R               *http.Request
	RequestBody     []byte
	RndData         float64
}

type Admin struct {
	ID           primitive.ObjectID `json:"id"  bson:"_id" q:",编号"`
	Name         string             `json:"name" bson:"name" q:",名称"`
	Introduction string             `json:"introduction" bson:"introduction" q:",介绍"`
	Avatar       string             `json:"avatar" bson:"avatar" q:",头像"`
	Role         []string           `json:"role" bson:"role" q:",角色"`
	Password     string             `json:"password" bson:"password" q:",密码"`
	Ctime        int64              `json:"ctime" bson:"ctime" q:",创建时间"`
	Mtime        int64              `json:"mtime" bson:"_mtime" q:",变更时间"`
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
	ctxMongo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
	//e = session.DB(DB).C(C_ADMIN).Find(selector).One(&curAdmin)
	e = GlobalMgoClient.Database(DB).Collection(C_ADMIN).FindOne(ctxMongo, selector).Decode(&curAdmin)
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

		ctxMongo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var curAdmin Admin
		selector := bson.M{"$and": []bson.M{bson.M{"_id": pl.AdminUser.ID}, bson.M{"name": pl.AdminUser.Name}}}
		//e = session.DB(DB).C(C_ADMIN).Find(selector).One(&curAdmin)
		e = GlobalMgoClient.Database(DB).Collection(C_ADMIN).FindOne(ctxMongo, selector).Decode(&curAdmin)
		if e == nil {
			code = 0
			//curAdmin.ID = ""
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
	_, _ = RedisDelString(xToken)
	code = 0
	data = "success"
	return
}

func (x AdminApi) Insert(req Admin) (code int, data struct {
	ID primitive.ObjectID `json:"id" q:",返回id"`
}, e error) {
	ctxMongo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req.ID = primitive.NewObjectID()
	req.Ctime = time.Now().Unix()
	req.Mtime = time.Now().Unix()

	_, e = GlobalMgoClient.Database(DB).Collection(C_ADMIN).InsertOne(ctxMongo, req)
	if e != nil {
		StdPrint.Error(e)
		code = -1
	} else {
		code = 0
		data.ID = req.ID
	}

	return
}

func (x AdminApi) Replace(req Admin) (code int, data struct {
	Mtime        int64               `json:"mtime" bson:"_mtime" q:",update time"`
	UpdateResult *mongo.UpdateResult `json:"update_result" bson:"update_result" q:", update result"`
}, e error) {
	code = -1
	ctxMongo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//req.ID = primitive.NewObjectID()
	//req.Ctime = time.Now().Unix()
	req.Mtime = time.Now().Unix()
	// allow overwriting, no mtime restrict
	//selector := bson.M{"$and": []bson.M{bson.M{"_id": req.ID}, bson.M{"_mtime": req.Mtime}}}
	selector := bson.M{"_id": req.ID}
	data.UpdateResult, e = GlobalMgoClient.Database(DB).Collection(C_ADMIN).UpdateOne(ctxMongo, selector,
		bson.D{
			{"$set", req},
		})
	if e != nil {
		StdPrint.Info(x.R.RequestURI+" rnd: ", x.RndData, "\t update:", "=> Error: ", e)
	} else {
		code = 0
		data.Mtime = req.Mtime
	}
	return
}

func (x AdminApi) Page(req PageDataKeyWord) (code int, data struct {
	List  []Admin `json:"list" q:",列表"`
	Total int64   `json:"total" q:",总数"`
}, e error) {
	code = -1
	ctxMongo, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sorter := toSortor(req.Sort)
	selector := toSelector(req.Filter)
	var cursor *mongo.Cursor
	if req.Keyword != "" && len(req.KeywordFields) > 0 {
		var keywordFilter = make(map[string]interface{}, 0)
		for _, v := range req.KeywordFields {
			keywordFilter[v] = req.Keyword
		}
		selector = bson.M{"$and": []bson.M{selector, toRegexSelector(keywordFilter)}}
	}
	data.Total, _ = GlobalMgoClient.Database(DB).Collection(C_ADMIN).CountDocuments(ctxMongo, selector)
	if req.Location == nil {
		var findOptions options.FindOptions
		if len(sorter) > 0 {
			findOptions = options.FindOptions{
				Limit: &req.Size,
				Skip:  &req.Offset,
				Sort:  sorter,
			}
		} else {
			findOptions = options.FindOptions{
				Limit: &req.Size,
				Skip:  &req.Offset,
			}
		}
		cursor, e = GlobalMgoClient.Database(DB).Collection(C_ADMIN).Find(ctxMongo, selector, &findOptions)
		if e != nil {
			StdPrint.Error(x.R.RequestURI, "\t GlobalMgoClient.Database(DB).Collection(C_OPERATOR_META).Find(ctxMongo, selector, &findOptions) Error: ", ShowMeJson(req), "\t", e)
			return
		}
	} else {
		// need adding Index for collection
		// db.collection.ensureIndex({"loc":"2dsphere"})
		if req.Size == 0 {
			req.Size = 10000
		}
		if req.LocationField == "" {
			req.LocationField = "location"
		}
		if req.MaxDistance == 0 {
			req.MaxDistance = 10000
		}
		pipeLine := []bson.M{
			{"$geoNear": bson.M{
				"near": bson.M{
					"type": "Point", "coordinates": req.Location,
				},
				"key":                req.LocationField,
				"minDistance":        req.MinDistance, // meters
				"maxDistance":        req.MaxDistance, // meters
				"distanceField":      "distance",
				"distanceMultiplier": 0.0006215, // 1/1609 convert meter to mile
				"spherical":          true,
				"query":              selector,
			}},
			{"$addFields": bson.M{
				"distance_difference": bson.M{
					"$add": []bson.M{
						{
							//"$multiply": []interface{}{1609, "$desired_distance"},
							"$multiply": bson.A{1, "$desired_distance"},
						},
						{
							"$multiply": bson.A{-1, "$distance"},
						},
					},
				},
			}},
			//{"$match": selector}, 	//同样的查询这里会缺少部分数据，必须使用$geoNear的query
			{"$skip": req.Offset},
			//{"$group": bson.M{"_id": "$pid", "count": bson.M{"$sum": 1}}},
			//{"$sort": bson.M{"count": -1}},
			{"$limit": req.Size},
		}
		cursor, e = GlobalMgoClient.Database(DB).Collection(C_ADMIN).Aggregate(ctxMongo, pipeLine)
		if e != nil {
			StdPrint.Error(x.R.RequestURI, "\t GlobalMgoClient.Database(DB).Collection(C_OPERATOR_META).Aggregate(ctxMongo, pipeLine) Error: ", ShowMeJson(req), "\t", e)
			return
		}
	}
	e = cursor.All(ctxMongo, &data.List)
	if e == nil {
		code = 0
	}

	return
}

// ##AdminApi.Delete: 用于XXXXX
// ##AdminApi.Delete: author: <a href="mailto:myemail@abc.com">myemail@abc.com</a>
// ##AdminApi.Delete: 可以添加多行说明，并且可以使用少量html
func (x AdminApi) Delete(req struct {
	ID []primitive.ObjectID `json:"id" q:",要删除的ID数组"`
}) (code int, data struct{}, e error) {
	code = -1

	ctxMongo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if req.ID != nil {
		_, e = GlobalMgoClient.Database(DB).Collection(C_ADMIN).DeleteMany(ctxMongo, bson.M{"_id": bson.M{"$in": req.ID}})
		if e != nil {
			StdPrint.Error(x.R.RequestURI+" rnd: ", x.RndData, "\t"+" Error: ", e, "\t", ShowMeJson(req))
			return
		}
		code = 0
		return
	}

	StdPrint.Info(x.R.RequestURI+" rnd: ", x.RndData, "\t"+" req: ", ShowMeJson(req))
	return
}

func (x AdminApi) Get(req struct {
	ID []primitive.ObjectID `json:"id" q:",要查询的ID数组"`
}) (code int, data struct {
	List []Admin `json:"list" q:",列表"`
	Size int     `json:"size" q:",个数"`
}, e error) {
	code = -1

	ctxMongo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cursor *mongo.Cursor
	cursor, e = GlobalMgoClient.Database(DB).Collection(C_ADMIN).Find(ctxMongo, bson.M{"_id": bson.M{"$in": req.ID}})
	if e != nil {
		StdPrint.Error(x.R.RequestURI+" rnd: ", x.RndData, "\t Find(ctxMongo, bson.M{\"_id\": bson.M{\"$in\": req.IdList}})"+" Error: ", e, "\t", ShowMeJson(req))
		return
	}
	e = cursor.All(ctxMongo, &data.List)
	if e != nil {
		StdPrint.Error(x.R.RequestURI+" rnd: ", x.RndData, "\t cursor.All(ctxMongo, &data.List)"+" Error: ", e, "\t", ShowMeJson(req))
		return
	}
	code = 0
	data.Size = len(data.List)
	StdPrint.Info(x.R.RequestURI+" rnd: ", x.RndData, "\t"+" req: ", ShowMeJson(req))
	return
}
