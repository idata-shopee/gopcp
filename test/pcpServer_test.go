package test

import (
	"errors"
	"github.com/idata-shopee/gopcp"
	"testing"
)

func TestPcpServer(t *testing.T) {
	addFunc := &gopcp.BoxFunc{
		FunType: gopcp.SandboxTypeNormal,
		Fun: func(args []interface{}, pcpServer *gopcp.PcpServer) (interface{}, error) {
			var res float64
			for _, arg := range args {
				if val, ok := arg.(float64); !ok {
					return nil, errors.New("args should be int")
				} else {
					res += val
				}
			}
			return res, nil
		},
	}

	funcMap := map[string]*gopcp.BoxFunc{
		"add": addFunc,
	}
	sandBox := gopcp.NewSandbox(funcMap)
	pcpServer := gopcp.NewPcpServer(sandBox)
	res, err := pcpServer.Execute("[\"add\", 1, 2]")
	if err != nil {
		t.Errorf(err.Error())
	}
	if res != float64(3) {
		t.Errorf("expect %f, actual %f", float64(3), res)
	}
	//
}
