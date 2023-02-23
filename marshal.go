package amqpx

import (
	"encoding/json"
	"strconv"
)

// json反序列化实现
func _unmarshal(b []byte, value interface{}) error {
	if len(b) == 0 {
		return nil
	}
	switch value := value.(type) {
	case nil:
		return nil
	case *[]byte:
		clone := make([]byte, len(b))
		copy(clone, b)
		*value = clone
		return nil
	case *string:
		*value = string(b)
		return nil
	default:
		return json.Unmarshal(b, value)
	}
}

func _marshal(value interface{}) ([]byte, error) {
	var sValue string
	switch value := value.(type) {
	case nil:
		return nil, nil
	case []byte:
		return value, nil
	case string:
		return []byte(value), nil
	case bool:
		sValue = strconv.FormatBool(value)
		return []byte(sValue), nil
	case float64:
		sValue = strconv.FormatFloat(value, 'f', -1, 64)
		return []byte(sValue), nil
	case float32:
		sValue = strconv.FormatFloat(float64(value), 'f', -1, 64)
		return []byte(sValue), nil
	case int:
		sValue = strconv.Itoa(value)
		return []byte(sValue), nil
	case int64:
		sValue = strconv.FormatInt(value, 10)
		return []byte(sValue), nil
	case int32:
		sValue = strconv.Itoa(int(value))
		return []byte(sValue), nil
	case int16:
		sValue = strconv.FormatInt(int64(value), 10)
		return []byte(sValue), nil
	case int8:
		sValue = strconv.FormatInt(int64(value), 10)
		return []byte(sValue), nil
	case uint:
		sValue = strconv.FormatUint(uint64(value), 10)
		return []byte(sValue), nil
	case uint64:
		sValue = strconv.FormatUint(value, 10)
		return []byte(sValue), nil
	case uint32:
		sValue = strconv.FormatUint(uint64(value), 10)
		return []byte(sValue), nil
	case uint16:
		sValue = strconv.FormatUint(uint64(value), 10)
		return []byte(sValue), nil
	case uint8:
		sValue = strconv.FormatUint(uint64(value), 10)
		return []byte(sValue), nil
	case error:
		return []byte(value.Error()), nil
	}
	return json.Marshal(value)
}
