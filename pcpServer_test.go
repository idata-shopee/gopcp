package gopcp

import (
	"errors"
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, expect interface{}, actual interface{}, message string) {
	if expect == actual {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("expect %v !=  actual %v", expect, actual)
	}
	t.Fatal(message)
}

func simpleSandbox() *Sandbox {
	funcMap := map[string]*BoxFunc{
		"add": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
			var res float64
			for _, arg := range args {
				if val, ok := arg.(float64); !ok {
					return nil, errors.New("args should be int")
				} else {
					res += val
				}
			}
			return res, nil
		}),
		"concat": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
			res := ""
			for _, arg := range args {
				if val, ok := arg.(string); !ok {
					return nil, errors.New("args should be string")
				} else {
					res += val
				}
			}
			return res, nil
		}),
		">": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
			a, _ := args[0].(float64)
			b, _ := args[1].(float64)
			return a > b, nil
		}),
		"sum": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
			list, ok := args[0].([]interface{})
			if !ok {
				return nil, errors.New("args should be int list")
			}
			v := 0.0
			for _, item := range list {
				itemValue, iok := item.(float64)
				if !iok {
					return nil, errors.New("args should be int list")
				}
				v += itemValue
			}
			return v, nil
		}),
	}
	sandBox := GetSandbox(funcMap)
	return sandBox
}

func runPcpCall(t *testing.T, pcpServer *PcpServer, callText string, expect interface{}) {
	res, err := pcpServer.Execute(callText, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	if res != expect {
		t.Errorf("expect %v !=  actual %v", expect, res)
	}
}

func runPcpCallExpectError(t *testing.T, pcpServer *PcpServer, callText string) {
	_, err := pcpServer.Execute(callText, nil)
	if err == nil {
		t.Errorf("expect error, but no error")
	}
}

func TestBase(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	var expect float64 = 6 // golang convert number to float64 in json
	runPcpCall(t, pcpServer, "[\"add\", 1, 2, 3]", expect)
}

func TestPcPClient(t *testing.T) {
	var pcpClient = PcpClient{}
	callText, _ := pcpClient.ToJSON(pcpClient.Call("add", 1, 2))
	assertEqual(t, callText, "[\"add\",1,2]", "")
}

func TestPcPClientNest(t *testing.T) {
	var pcpClient = PcpClient{}
	callText, _ := pcpClient.ToJSON(pcpClient.Call("add", 1, pcpClient.Call("succ", 2)))
	assertEqual(t, callText, "[\"add\",1,[\"succ\",2]]", "")
}

func TestConcat(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, "[\"concat\", \"hello\", \" \", \"world!\"]", "hello world!")
}

func TestIfFail(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, "[\"if\", [\">\", 3, 4], 1, 2]", 2.0)
	runPcpCall(t, pcpServer, "[\"if\", 0, 1, 2]", 2.0)
	runPcpCall(t, pcpServer, "[\"if\", null, 1, 2]", 2.0)
}

func TestIfSuccess(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, "[\"if\", [\">\", 6, 4], 1, 2]", 1.0)
	runPcpCall(t, pcpServer, "[\"if\", true, 1, 2]", 1.0)
	runPcpCall(t, pcpServer, "[\"if\", 1, 1, 2]", 1.0)
}

func TestListFunction(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, "[\"sum\", [\"list\", [\"add\", 6, 4], 1, 2]]", 13.0)
}

func TestIfException(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCallExpectError(t, pcpServer, "[\"if\", 2]")
}

func TestRawData(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, "[\"sum\", [\"'\", 1, 2, 3]]", 6.0)
}

func TestMissingFunName(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCallExpectError(t, pcpServer, "[\"fakkkkke\"]")
}
