package api

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ShowMeJson(data interface{}) string {
	dataBytes, _ := json.Marshal(data)
	return string(dataBytes)
}

func toSelector(filter map[string]interface{}) bson.M {
	selector := []bson.M{}
	for k, v := range filter {
		if k == "id" {
			strValue, _ := primitive.ObjectIDFromHex(v.(string))
			selector = append(selector, bson.M{"_id": strValue})
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
		return bson.M{}
	}
}

//func toSortor(order []string) string {
//	return strings.Join(order, ", ")
//}

func toSortor(order []string) (sorter bson.D) {
	sorter = bson.D{}
	if order != nil {
		for _, v := range order {
			if v[0] == '-' {
				sorter = append(sorter, bson.E{
					Key:   v[1:],
					Value: -1,
				})
			} else if v[0] == '+' {
				sorter = append(sorter, bson.E{
					Key:   v[1:],
					Value: 1,
				})
			} else {
				sorter = append(sorter, bson.E{
					Key:   v,
					Value: 1,
				})
			}
		}
	}
	return sorter
}

func toRegexSelector(filter map[string]interface{}) bson.M {
	selector := []bson.M{}
	for k, v := range filter {
		if k == "id" {
			strValue, _ := primitive.ObjectIDFromHex(v.(string))
			selector = append(selector, bson.M{"_id": strValue})
		} else {
			switch v.(type) {
			case float64:
				// json传递的大于6位的数字需要强制按照float64转换
				selector = append(selector, bson.M{k: int(v.(float64))})
			case string:
				selector = append(selector, bson.M{k: bson.M{"$regex": v}})
			default:
				selector = append(selector, bson.M{k: v})
			}
		}
	}
	if len(selector) > 0 {
		return bson.M{"$or": selector}
	} else {
		return nil
	}
}
