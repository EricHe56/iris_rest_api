package models

// 查询参数
type PageData struct {
	Size   int                    `json:"size" q:",返回记录数量（对应每页记录数）"`
	Offset int                    `json:"offset" q:",返回记录偏移量数量，如： 0"`
	Order  []string               `json:"order" q:",排序条件：+/-后面跟排序字段，+代表正序，-代表逆序"`
	Filter map[string]interface{} `json:"filter" q:",过滤条件，{字段名:值}"`
}

// 查询参数
type PageDataKeyWord struct {
	Keyword       string                 `json:"keyword" q:",Keyword for searching"`
	KeywordFields []string               `json:"keyword_fields" q:",searching fields list"`
	Size          int64                  `json:"size" q:", page size"`
	Offset        int64                  `json:"offset" q:", offset: data return from default 0"`
	Sort          []string               `json:"sort" q:",sort order：[+score, -ctime](asc on score, desc on ctime)"`
	Filter        map[string]interface{} `json:"filter" q:",filter: {field_name:value}"`
	Location      []float64              `json:"location" bson:"location" q:",Location[latitude, longitude]"`
	LocationField string                 `json:"location_field" bson:"location_field" q:", location field name in DB. default is the field which has geo index if only one geo index. otherwise the parameter will be required if geo index more than one."`
	MinDistance   int                    `json:"min_distance" bson:"min_distance" q:", min distance in meters. default 0m"`
	MaxDistance   int                    `json:"max_distance" bson:"max_distance" q:", max distance in meters. default 10000m"`
}
