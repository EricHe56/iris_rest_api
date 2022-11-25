package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type D_User struct {
	ID    primitive.ObjectID `json:"id"  bson:"_id" q:",编号"`
	Name  string             `json:"name" bson:"name" q:",名称"`
	Age   int                `json:"age" bson:"age" q:",年龄"`
	Ctime int64              `json:"ctime" bson:"_ctime" q:",创建时间"`
	Mtime int64              `json:"mtime" bson:"_mtime" q:",修改时间"`
}

type Faq struct {
	ID             primitive.ObjectID `json:"id"  bson:"_id" q:",编号"`
	Question       string             `json:"question" bson:"question" q:",问题"`
	Answer         string             `json:"answer" bson:"answer" q:",答案"`
	Status         int                `json:"status" bson:"status" q:",状态 1=上线 10=下线"`
	Lang           string             `json:"lang" bson:"lang" q:",语言版本，cn: 中文;en: 英文;jp: 日文;kr: 韩文;"`
	Ctime          int64              `json:"ctime" bson:"_ctime" q:",创建时间"`
	Mtime          int64              `json:"mtime" bson:"_mtime" q:",变更时间"`
	MultiLanguages []FaqMultiLanguage `json:"multi_languages" bson:"multi_languages" q:",多语言信息包数组，包含多语言支持信息。"`
}

type FaqMultiLanguage struct {
	Question string `json:"question" bson:"question" q:",问题"`
	Answer   string `json:"answer" bson:"answer" q:",答案"`
	Lang     string `json:"lang" bson:"lang" q:",语言版本，cn: 中文;en: 英文;jp: 日文;kr: 韩文;"`
}
