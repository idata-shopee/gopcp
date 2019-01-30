package gopcp

import "encoding/json"

//PcpClient pcp client
type PcpClient struct{}

//Call call function in  pcp server
func (c *PcpClient) Call(funName string, params []interface{}) (interface{}, error) {
	return nil, nil
}

func (c *PcpClient) ToJSON(res interface{}) (str string, err error) {
	bytes, err := json.Marshal(res)
	if err != nil {
		return
	}
	str = string(bytes[:])
	return
}
