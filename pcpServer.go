package gopcp

import (
	"encoding/json"
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

func parseJSON(source string) (interface{}, error) {
	var arr interface{}
	if err := json.Unmarshal([]byte(source), &arr); err != nil {
		return nil, err
	}
	return parseAst(arr), nil
}

// Execute ....
func (pcpServer *PcpServer) Execute(source string, attachment interface{}) (interface{}, error) {
	ast, err := parseJSON(source)
	if err != nil {
		return nil, err
	}
	return pcpServer.ExecuteAst(ast, attachment)
}

func (p *PcpServer) ExecuteAst(ast interface{}, attachment interface{}) (interface{}, error) {
	switch funNode := ast.(type) {
	case FunNode:
		sandboxFun, err := p.sandbox.Get(funNode.funName)
		if err != nil {
			return nil, err
		}

		if sandboxFun.FunType == SandboxTypeNormal {
			// for normal mode, resolve params first
			var paramRets []interface{}
			for _, param := range funNode.params {
				paramRet, paramErr := p.ExecuteAst(param, attachment)
				if paramErr != nil {
					return nil, paramErr
				}
				paramRets = append(paramRets, paramRet)
			}

			return sandboxFun.Fun(paramRets, attachment, p)
		} else if sandboxFun.FunType == SandboxTypeLazy {
			// execute lazy sandbox function
			return sandboxFun.Fun(funNode.params, attachment, p)
		}

	default:
		return ast, nil
	}
	return nil, nil // impossible for this line
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
