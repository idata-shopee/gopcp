package test

import (
	"errors"
	"github.com/idata-shopee/gopcp"
	"testing"
)

func simpleSandbox() *gopcp.Sandbox {
	funcMap := map[string]*gopcp.BoxFunc{
		"add": gopcp.ToSandboxFun(func(args []interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
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
	}
	sandBox := gopcp.GetSandbox(funcMap)
	return sandBox
}

func runPcpCall(t *testing.T, pcpServer *gopcp.PcpServer, callText string, expect interface{}) {
	res, err := pcpServer.Execute(callText)
	if err != nil {
		t.Errorf(err.Error())
	}
	if res != expect {
		t.Errorf("expect %v !=  actual %v", expect, res)
	}
}

func TestBase(t *testing.T) {
	pcpServer := gopcp.NewPcpServer(simpleSandbox())
	var expect float64 = 6 // golang convert number to float64 in json
	runPcpCall(t, pcpServer, "[\"add\", 1, 2, 3]", expect)
}
