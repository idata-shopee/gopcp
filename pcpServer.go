package gopcp

import (
	"encoding/json"
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
	sandbox Sandbox
}

func parseJSON(source string) (*FunNode, error) {
	arr := []interface{}{}
	if err := json.Unmarshal([]byte{source}, &arr); err != nil {
		return err
	}
	return parseAst(arr)
}

// Execute ....
func (pcpServer *PcpServer) Execute(source string) (interface{}, error) {
	node, err := parseJSON(source)
	if err != nil {
		return err
	}
	return pcpServer.executeAst(node)
}

func (pcpServer *PcpServer) executeAst(node *FunNode) (interface{}, error) {
	if node != nil {
		fn := pcpServer.sandbox.get(node.funName)
		funcType := fn.funType
		fun := fn.fun
		params := []interface{}{}
		if funcType == SandboxTypeNormal {
			for field := range node.params {
				var res interface{}
				switch field.(type) {
				case FunNode:
					funcNode := field.(FunNode)
					res = executeAst(funcNode)
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
}

func parseAst(arr []interface{}) (node *FunNode, err error) {
	if len(arr) == 0 {
		return
	}
	funName := arr[0]
	if reflect.TypeOf(funcName).Kind() != reflect.String {
		err = error.New("first element in array must be string")
		return
	}
	curNode := &FunNode{funName: funName, params: []interface{}{}}
	for i, v := range arr[1:] {
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Array {
			var newArr []interface{}
			copy(newArr[:], val)
			newNode, err := parseAst(newArr)
			curNode.params = append(curNode.params, newNode)
		} else {
			curNode.params = append(curNode.params, val)
		}
	}
	return
}

// NewPcpServer merge sandbox with default sandbox
func NewPcpServer(sandbox *Sandbox) *PcpServer {
	box := NewSandbox(defBox)
	defSandBox.Extend(sandbox)
	return &PcpServer{box}
}
