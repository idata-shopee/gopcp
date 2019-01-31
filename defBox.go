package gopcp

import (
	"errors"
)

var DefBox = &Sandbox{map[string]*BoxFunc{
	"if": ToLazySandboxFun(func(args []interface{}, pcpServer *PcpServer) (interface{}, error) {
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
		conditionRet, cerr := pcpServer.ExecuteAst(conditionExp)
		if cerr != nil {
			return nil, cerr
		}

		if conditionRet == false || conditionRet == 0.0 || conditionRet == nil {
			return pcpServer.ExecuteAst(failExp)
		} else {
			return pcpServer.ExecuteAst(successExp)
		}

		return nil, nil
	})}}
