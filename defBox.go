package gopcp

import (
	"errors"
	"fmt"
)

var DefBox = &Sandbox{map[string]*BoxFunc{
	"if": ToLazySandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, errors.New("if grammer error. if must have at least 2 params, at most 3 params. eg: [\"if\", true, 1, 0], [\"if\", true, 1]")
		}

		conditionExp := args[0]
		successExp := args[1]
		var failExp interface{} = nil
		if len(args) > 2 {
			failExp = args[2]
		}

		// condition
		conditionRet, cerr := pcpServer.ExecuteAst(conditionExp, attachment)
		if cerr != nil {
			return nil, cerr
		}

		if conditionRet == false || conditionRet == 0.0 || conditionRet == nil {
			return pcpServer.ExecuteAst(failExp, attachment)
		} else {
			return pcpServer.ExecuteAst(successExp, attachment)
		}
	}),

	// basic data structure: List
	"List": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
		return args, nil
	}),

	// basic data structure: Map[string]interface{}
	"Map": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
		l := len(args)
		if l%2 == 1 {
			return nil, errors.New("Map grammer error. eg: [\"Map\", \"age\", 3, \"weight\", 45]")
		}

		m := make(map[string]interface{})

		for i := 0; i < l; i += 2 {
			key, ok := args[i].(string)
			if !ok {
				return nil, errors.New("Map grammer error. eg: [\"Map\", \"age\", 3, \"weight\", 45]")
			}
			m[key] = args[i+1]
		}

		return m, nil
	}),

	// get property from object
	"prop": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
		if len(args) != 2 {
			return nil, errors.New("Prop grammer error. eg: [\"Prop\", {\"a\": 1}, \"a\"]")
		}

		obj, ok1 := args[0].(map[string]interface{})
		propName, ok2 := args[1].(string)

		if !ok1 || !ok2 {
			return nil, fmt.Errorf("error types of function prop")
		} else {
			return obj[propName], nil
		}
	}),

	"error": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
		l := len(args)
		if l%2 < 1 {
			return nil, errors.New("Must specify error message. eg: [\"error\", \"Exception!\"]")
		}

		if msg, ok := args[0].(string); !ok {
			return nil, errors.New("Must specify error message (string). eg: [\"error\", \"Exception!\"]")
		} else {
			return nil, errors.New(msg)
		}
	}),

	"+": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
		var res float64 = 0
		for _, arg := range args {
			if val, ok := arg.(float64); !ok {
				return nil, errors.New("args of \"+\" should be float64")
			} else {
				res += val
			}
		}
		return res, nil
	}),

	"*": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
		var res float64 = 1
		for _, arg := range args {
			if val, ok := arg.(float64); !ok {
				return nil, errors.New("args of \"+\" should be float64")
			} else {
				res *= val
			}
		}
		return res, nil
	}),

	"-": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
		if len(args) != 2 {
			return nil, errors.New("- must have two arguments")
		}

		v1, ok1 := args[0].(float64)
		v2, ok2 := args[1].(float64)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("args of \"-\" should be float64, but got v1=%v, v2=%v", v1, v2)
		} else {
			return v1 - v2, nil
		}
	}),

	"/": ToSandboxFun(func(args []interface{}, attachment interface{}, pcpServer *PcpServer) (interface{}, error) {
		if len(args) != 2 {
			return nil, errors.New("- must have two arguments")
		}

		v1, ok1 := args[0].(float64)
		v2, ok2 := args[1].(float64)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("args of \"/\" should be float64, but got v1=%v, v2=%v", v1, v2)
		} else {
			if v2 == 0 {
				return nil, errors.New("divisor can not be 0")
			} else {
				return v1 / v2, nil
			}
		}
	}),
}}
