package api

import (
	"gopkg.in/mgo.v2/bson"
	"strings"
)

func toSelector(filter map[string]interface{}) bson.M {
	selector := []bson.M{}
	for k, v := range filter {
		if k == "id" {
			selector = append(selector, bson.M{"_id": bson.ObjectIdHex(v.(string))})
		} else {
			switch v.(type) {
			case float64:
				// json传递的大于6位的数字需要强制按照float64转换
				selector = append(selector, bson.M{k: int(v.(float64))})
			default:
				selector = append(selector, bson.M{k: v})
			}
		}
	}
	if len(selector) > 0 {
		return bson.M{"$and": selector}
	} else {
		return nil
	}
}

func toSortor(order []string) string {
	return strings.Join(order, ", ")
}
