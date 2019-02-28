package gopcp

import (
	"errors"
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
}}
