package models

// 查询参数
type PageData struct {
	Size   int                    `json:"size" q:",返回记录数量（对应每页记录数）"`
	Offset int                    `json:"offset" q:",返回记录偏移量数量，如： 0"`
	Order  []string               `json:"order" q:",排序条件：+/-后面跟排序字段，+代表正序，-代表逆序"`
	Filter map[string]interface{} `json:"filter" q:",过滤条件，{字段名:值}"`
}
