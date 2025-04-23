package helper

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modern-go/parse/model"
)

// MarshalMap 把model.Map对象转成json
func MarshalMap(m model.Map) ([]byte, error) {
	return json.Marshal(convertModelMap2GeneralMap(m))
}

func convertModelMap2GeneralMap(m model.Map) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		var keyValue string
		switch kk := k.(type) {
		case string:
			keyValue = kk
		case int:
			keyValue = strconv.Itoa(kk)
		case int64:
			keyValue = strconv.FormatInt(kk, 10)
		case float64:
			keyValue = strconv.FormatFloat(kk, 'b', -1, 64)
		default:
			keyValue = fmt.Sprintf("%v", kk)
		}
		vv := v
		if mMap, ok := v.(model.Map); ok {
			vv = convertModelMap2GeneralMap(mMap)
		} else if mList, ok := v.(model.List); ok {
			for idx, item := range mList {
				if itemMap, ok := item.(model.Map); ok {
					mList[idx] = convertModelMap2GeneralMap(itemMap)
				}
			}
		}
		result[keyValue] = vv
	}
	return result
}

func IntToBytes(n int,b byte) ([]byte,error) {
	switch b {
	case 1:
		tmp := int8(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(),nil
	case 2:
		tmp := int16(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(),nil
	case 3,4:
		tmp := int32(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(),nil
	}
	return nil,fmt.Errorf("IntToBytes b param is invaild")
}
