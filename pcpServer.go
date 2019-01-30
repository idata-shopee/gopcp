package gopcp

import (
	"encoding/json"
	"errors"
	"reflect"
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

func parseJSON(source string) (*FunNode, error) {
	arr := []interface{}{}
	if err := json.Unmarshal([]byte(source), &arr); err != nil {
		return nil, err
	}
	return parseAst(arr)
}

// Execute ....
func (pcpServer *PcpServer) Execute(source string) (interface{}, error) {
	node, err := parseJSON(source)
	if err != nil {
		return nil, err
	}
	return pcpServer.executeAst(node)
}

func (p *PcpServer) executeAst(node *FunNode) (interface{}, error) {
	if node != nil {
		fn, err := p.sandbox.Get(node.funName)
		funcType := fn.FunType
		fun := fn.Fun
		params := []interface{}{}
		if funcType == SandboxTypeNormal {
			for _, field := range node.params {
				var res interface{}
				switch field.(type) {
				case FunNode:
					funcNode := field.(FunNode)
					res, err = p.executeAst(&funcNode)
					if err != nil {
						return nil, err
					}
				default:
					res = field
				}
				params = append(params, res)
			}
			res, err := fun(params...)
			return res, err
		} else if funcType == SandboxTypeLazy {
			return nil, nil
		}
	}
	return nil, nil
}

func parseAst(arr []interface{}) (node *FunNode, err error) {
	if len(arr) == 0 {
		return
	}
	funName := arr[0]
	if reflect.TypeOf(funName).Kind() != reflect.String {
		err = errors.New("first element in array must be string")
		return
	}
	curNode := &FunNode{funName: funName.(string), params: []interface{}{}}
	for _, v := range arr[1:] {
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Array {
			ret := make([]interface{}, val.Len())
			for i := 0; i < val.Len(); i++ {
				ret[i] = val.Index(i).Interface()
			}
			newNode, err := parseAst(ret)
			if err != nil {
				return nil, err
			}
			curNode.params = append(curNode.params, newNode)
		} else {
			curNode.params = append(curNode.params, v)
		}
	}
	return curNode, nil
}

// NewPcpServer merge sandbox with default sandbox
func NewPcpServer(sandbox *Sandbox) *PcpServer {
	box := defBox
	box.Extend(sandbox)
	return &PcpServer{box}
}
