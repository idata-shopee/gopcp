package gopcp

import (
	"encoding/json"
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
		"stringify": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
			bytes, err := json.Marshal(args[0])
			if err != nil {
				return nil, err
			}
			return string(bytes), nil
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
	runPcpCall(t, pcpServer, "[\"+\", 1, 2, 3]", expect)
}

func TestPcPClient(t *testing.T) {
	var pcpClient = PcpClient{}
	callText, _ := pcpClient.ToJSON(pcpClient.Call("+", 1, 2))
	assertEqual(t, callText, "[\"+\",1,2]", "")
}

func TestMarshal(t *testing.T) {
	var pcpClient = PcpClient{}

	var value = make(chan int)

	_, err := pcpClient.ToJSON(pcpClient.Call("+", value))
	if err == nil {
		t.Errorf("expect error, but no error")
	}
}

func TestUnmarshal(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	_, err := pcpServer.Execute("{{", nil)
	if err == nil {
		t.Errorf("expect error, but no error")
	}
}

func TestPcpServerMissingArgs(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	ret, err := pcpServer.Execute("[]", nil)
	if err != nil {
		t.Errorf("errored")
	}
	ret2 := ret.([]interface{})
	assertEqual(t, len(ret2), 0, "")
}

func TestCallWithArray(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	var pcpClient = PcpClient{}
	callText, _ := pcpClient.ToJSON(pcpClient.Call("sum", []interface{}{1, 2}))
	res, err := pcpServer.Execute(callText, nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	if res != 3.0 {
		t.Errorf("expect %v !=  actual %v", 3.0, res)
	}
}

func TestPcPClientNest(t *testing.T) {
	var pcpClient = PcpClient{}
	callText, _ := pcpClient.ToJSON(pcpClient.Call("+", 1, pcpClient.Call("succ", 2)))
	assertEqual(t, callText, "[\"+\",1,[\"succ\",2]]", "")
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
	runPcpCall(t, pcpServer, "[\"sum\", [\"List\", [\"+\", 6, 4], 1, 2]]", 13.0)
}

func TestMapFunction(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, `["stringify", ["Map", "age", 3, "weight", 45]]`, `{"age":3,"weight":45}`)
	runPcpCall(t, pcpServer, `["stringify", ["Map"]]`, `{}`)
}

func TestPureData(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, `["stringify", [1,2,3]]`, `[1,2,3]`)
}

func TestMapException(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCallExpectError(t, pcpServer, `["stringify", ["Map", "age"]]`)
	runPcpCallExpectError(t, pcpServer, `["stringify", ["Map", "age", 1, 2]]`)
	runPcpCallExpectError(t, pcpServer, `["stringify", ["Map", 1, 2]]`)
}

func TestIfException(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCallExpectError(t, pcpServer, "[\"if\", 2]")
}

func TestIfException2(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCallExpectError(t, pcpServer, "[\"if\", [\"error\", \"exception!!!\"], 1, 2]")
}

func TestErrorException(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCallExpectError(t, pcpServer, "[\"error\", 123]")
	runPcpCallExpectError(t, pcpServer, "[\"error\"]")
}

func TestRawData(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, "[\"sum\", [\"'\", 1, 2, 3]]", 6.0)
}

func TestMissingFunName(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCallExpectError(t, pcpServer, "[\"fakkkkke\"]")
}

func TestAddType(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCallExpectError(t, pcpServer, "[\"+\", null]")
}

func TestSubstractType(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, `["-", 3, 2]`, 1.0)
	runPcpCallExpectError(t, pcpServer, "[\"-\", 1]")
	runPcpCallExpectError(t, pcpServer, "[\"-\", 1, null]")
}

func TestMultiplyType(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, `["*", 3, 2, 4]`, 24.0)
	runPcpCallExpectError(t, pcpServer, "[\"*\", null]")
}

func TestDivideType(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, `["/", 3, 2]`, 1.5)
	runPcpCallExpectError(t, pcpServer, "[\"/\", 1]")
	runPcpCallExpectError(t, pcpServer, "[\"/\", 1, 0]")
	runPcpCallExpectError(t, pcpServer, "[\"/\", 1, null]")
}

func TestPropFunction(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, `["prop", {"a": 1.2}, "a"]`, 1.2)
	runPcpCallExpectError(t, pcpServer, `["prop", 1, null]`)
	runPcpCallExpectError(t, pcpServer, `["prop", 1]`)
}

func TestEqualFunction(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, `["==", 1, 1]`, true)
	runPcpCall(t, pcpServer, `["==", 1, 1.0]`, true)
	runPcpCall(t, pcpServer, `["==", 1, 2]`, false)
	runPcpCall(t, pcpServer, `["==", "a", "a"]`, true)
	runPcpCall(t, pcpServer, `["==", "a", "A"]`, false)
	runPcpCallExpectError(t, pcpServer, `["==", 1]`)
}

func TestNotEqualFunction(t *testing.T) {
	pcpServer := NewPcpServer(simpleSandbox())
	runPcpCall(t, pcpServer, `["!=", 1, 1]`, false)
	runPcpCall(t, pcpServer, `["!=", 1, 1.0]`, false)
	runPcpCall(t, pcpServer, `["!=", 1, 2]`, true)
	runPcpCall(t, pcpServer, `["!=", "a", "a"]`, false)
	runPcpCall(t, pcpServer, `["!=", "a", "A"]`, true)
	runPcpCallExpectError(t, pcpServer, `["!=", 1]`)
}

func runParseAstJson(t *testing.T, source string) {
	var arr interface{}
	err := json.Unmarshal([]byte(source), &arr)
	if err != nil {
		t.Errorf(err.Error())
	}
	ast := ParseJsonObjectToAst(arr)

	bs, err := json.Marshal(ParseAstToJsonObject(ast))
	if err != nil {
		t.Errorf(err.Error())
	}

	assertEqual(t, string(bs), source, "")
}

func TestParseAstToJsonObject(t *testing.T) {
	runParseAstJson(t, "1")
	runParseAstJson(t, `"hello"`)
	runParseAstJson(t, `true`)
	runParseAstJson(t, `null`)
	runParseAstJson(t, `{}`)
	runParseAstJson(t, `["'",1,2,3]`) // [1,2,3] <=> [',1,2,3]
	runParseAstJson(t, `["if",true,["get",3]]`)
	runParseAstJson(t, `["'",true,["'",3]]`)
	runParseAstJson(t, `["if",true,["'","ok"]]`)
	runParseAstJson(t, `{"a":[1,2,3]}`)
}
