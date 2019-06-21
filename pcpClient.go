package gopcp

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
)

//PcpClient pcp client
type PcpClient struct{}

type CallResult struct {
	Result interface{}
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

//Call call function in  pcp server
//
func (c *PcpClient) Call(funName string, params ...interface{}) CallResult {
	var args []interface{}

	for _, param := range params {
		switch item := param.(type) {
		case CallResult:
			args = append(args, item.Result)
		default:
			if reflect.ValueOf(item).Kind() == reflect.Slice {
				args = append(args, append([]interface{}{"'"}, getItems(item)...))
			} else {
				args = append(args, item)
			}
		}
	}

	return CallResult{append([]interface{}{funName}, args...)}
}

func getItems(item interface{}) []interface{} {
	var ans []interface{}
	items := reflect.ValueOf(item)
	for i := 0; i < items.Len(); i++ {
		ans = append(ans, items.Index(i).Interface())
	}
	return ans
}

func (c *PcpClient) ToJSON(callResult CallResult) (str string, err error) {
	bytes, err := JSONMarshal(callResult.Result)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(bytes[:]), " \n"), nil
}
