package gopcp

import "encoding/json"

//PcpClient pcp client
type PcpClient struct{}

//Call call function in  pcp server
func (c *PcpClient) Call(funName string, params []interface{}) (interface{}, error) {
}

func (c *PcpClient) ToJSON(res interface{}) (string, error) {
	return json.Marshal(res)
}
