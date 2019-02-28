package gopcp

import "encoding/json"

//PcpClient pcp client
type PcpClient struct{}

type CallResult struct {
	result interface{}
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
	bytes, err := json.Marshal(callResult.result)
	if err != nil {
		return "", err
	}
	return string(bytes[:]), nil
}
