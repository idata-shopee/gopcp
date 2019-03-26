package gopcp

import (
	"bytes"
	"encoding/json"
	"strings"
)

//PcpClient pcp client
type PcpClient struct{}

type CallResult struct {
	result interface{}
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
			args = append(args, item.result)
		case []interface{}:
			constParam := append([]interface{}{"'"}, item...)
			args = append(args, constParam)
		default:
			args = append(args, item)
		}
	}

	return CallResult{append([]interface{}{funName}, args...)}
}

func (c *PcpClient) ToJSON(callResult CallResult) (str string, err error) {
	bytes, err := JSONMarshal(callResult.result)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(bytes[:]), " \n"), nil
}
