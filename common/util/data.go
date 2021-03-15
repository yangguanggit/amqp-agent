package util

import (
	"encoding/json"
)

//解析消息数据
func ParseMessage(message interface{}) ([]string, error) {
	var data []string
	switch value := message.(type) {
	case string:
		data = append(data, value)
	case []string:
		data = value
	case []interface{}:
		for _, v := range value {
			switch a := v.(type) {
			case string:
				data = append(data, a)
			default:
				js, err := json.Marshal(v)
				if err != nil {
					return nil, err
				}
				data = append(data, string(js))
			}
		}
	default:
		js, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		data = append(data, string(js))
	}

	return data, nil
}
