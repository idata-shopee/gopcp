package gopcp

import (
	"encoding/json"
	"sync"
)

// FunNode function node
type FunNode struct {
	funName string
	params  []interface{}
}

// PcpServer simpe calling protocol
// grammer based on json
// ["fun1", 1, 2, ["fun2", 3]] => fun1(1, 2, fun2(3))
type PcpServer struct {
	sandbox *Sandbox
}

// Execute ....
func (pcpServer *PcpServer) Execute(source string, attachment interface{}) (interface{}, error) {
	var arr interface{}
	if err := json.Unmarshal([]byte(source), &arr); err != nil {
		return nil, err
	}

	return pcpServer.ExecuteJsonObj(arr, attachment)
}

// @param arr. Pure Json Object.
func (pcpServer *PcpServer) ExecuteJsonObj(arr interface{}, attachment interface{}) (interface{}, error) {
	ast := parseAst(arr)
	return pcpServer.ExecuteAst(ast, attachment)
}

func (p *PcpServer) ExecuteAst(ast interface{}, attachment interface{}) (interface{}, error) {
	switch funNode := ast.(type) {
	case FunNode:
		if sandboxFun, err := p.sandbox.Get(funNode.funName); err != nil {
			return nil, err
		} else {
			if sandboxFun.FunType == SandboxTypeNormal {
				// for normal mode, resolve params first
				paramRets := make([]interface{}, len(funNode.params))
				var wg sync.WaitGroup

				var err error = nil

				for i, param := range funNode.params {
					wg.Add(1)
					go func(i int, param interface{}) {
						defer wg.Done()
						paramRet, paramErr := p.ExecuteAst(param, attachment)
						if paramErr != nil {
							// error happened, set it
							err = paramErr
						} else {
							paramRets[i] = paramRet
						}
					}(i, param)
				}

				wg.Wait()

				if err != nil {
					return nil, err
				} else {
					return sandboxFun.Fun(paramRets, attachment, p)
				}
			} else if sandboxFun.FunType == SandboxTypeLazy {
				// execute lazy sandbox function
				return sandboxFun.Fun(funNode.params, attachment, p)
			}
		}
	}

	return ast, nil
}

func parseAst(source interface{}) interface{} {
	switch arr := source.(type) {
	case []interface{}:
		if len(arr) == 0 {
			return arr
		}

		switch head := arr[0].(type) {
		case string:
			if head == "'" {
				return arr[1:]
			} else {
				var params []interface{}

				for i := 1; i < len(arr); i++ {
					params = append(params, parseAst(arr[i]))
				}

				return FunNode{head, params}
			}
		default:
			return arr
		}
	default:
		return arr
	}
}

// NewPcpServer merge sandbox with default sandbox
func NewPcpServer(sandbox *Sandbox) *PcpServer {
	return &PcpServer{sandbox}
}
