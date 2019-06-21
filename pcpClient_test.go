package gopcp

import (
	"testing"
)

func TestClientBase(t *testing.T) {
	pcpClient := PcpClient{}
	t1, _ := pcpClient.ToJSON(pcpClient.Call("test", []int{1, 2, 3}))
	t2, _ := pcpClient.ToJSON(pcpClient.Call("test", []interface{}{1, 2, 3}))
	t3, _ := pcpClient.ToJSON(pcpClient.Call("test", map[string]int{"a": 1}))
	t4, _ := pcpClient.ToJSON(pcpClient.Call("test", 120))
	assertEqual(t, `["test",["'",1,2,3]]`, t1, "")
	assertEqual(t, `["test",["'",1,2,3]]`, t2, "")
	assertEqual(t, `["test",{"a":1}]`, t3, "")
	assertEqual(t, `["test",120]`, t4, "")
}
