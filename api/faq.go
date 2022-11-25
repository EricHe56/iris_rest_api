package api

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	. "iris_rest_api/init"
	. "iris_rest_api/models"

	//"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type FaqApi struct {
	CodeDescription map[int]string
	W               http.ResponseWriter
	R               *http.Request
	RequestBody     []byte
	RndData         float64
}

func (x FaqApi) Insert(req Faq) (code int, data struct {
	ID primitive.ObjectID `json:"id" q:",返回id"`
}, e error) {
	ctxMongo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req.ID = primitive.NewObjectID()
	req.Ctime = time.Now().Unix()
	req.Mtime = time.Now().Unix()

	_, e = GlobalMgoClient.Database(DB).Collection(C_FAQ).InsertOne(ctxMongo, req)
	if e != nil {
		StdPrint.Error(e)
		code = -1
	} else {
		code = 0
		data.ID = req.ID
	}

	return
}

func (x FaqApi) Replace(req Faq) (code int, data struct {
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
	data.UpdateResult, e = GlobalMgoClient.Database(DB).Collection(C_FAQ).UpdateOne(ctxMongo, selector,
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

func (x FaqApi) Page(req PageDataKeyWord) (code int, data struct {
	List  []Faq `json:"list" q:", List"`
	Total int64 `json:"total" q:",total"`
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
	data.Total, _ = GlobalMgoClient.Database(DB).Collection(C_FAQ).CountDocuments(ctxMongo, selector)
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
		cursor, e = GlobalMgoClient.Database(DB).Collection(C_FAQ).Find(ctxMongo, selector, &findOptions)
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
		cursor, e = GlobalMgoClient.Database(DB).Collection(C_FAQ).Aggregate(ctxMongo, pipeLine)
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

func (x FaqApi) Delete(req struct {
	ID []primitive.ObjectID `json:"id" q:",要删除的ID数组"`
}) (code int, data struct{}, e error) {
	code = -1

	ctxMongo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if req.ID != nil {
		_, e = GlobalMgoClient.Database(DB).Collection(C_FAQ).DeleteMany(ctxMongo, bson.M{"_id": bson.M{"$in": req.ID}})
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

func (x FaqApi) Get(req struct {
	ID []primitive.ObjectID `json:"id" q:",要查询的ID数组"`
}) (code int, data struct {
	List []Faq `json:"list" q:",列表"`
	Size int   `json:"size" q:",个数"`
}, e error) {
	code = -1

	ctxMongo, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cursor *mongo.Cursor
	cursor, e = GlobalMgoClient.Database(DB).Collection(C_FAQ).Find(ctxMongo, bson.M{"_id": bson.M{"$in": req.ID}})
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
