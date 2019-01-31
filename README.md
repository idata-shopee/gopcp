# gopcp

A simple Lisp-Style protocol for communication (golang version)

## Features

- lisp style composation

- function sandbox. User can control every function and detail by providing the sandbox.

- based on json grammer

## Quick Example

```go
import (
  "github.com/idata-shopee/gopcp"
  "errors"
)

// create sandbox
sandBox := gopcp.GetSandbox(map[string]*gopcp.BoxFunc{
	"add": gopcp.ToSandboxFun(func(args []interface{}, pcpServer *gopcpc.PcpServer) (interface{}, error) {
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
})

pcpServer := NewPcpServer(sandbox)

// we can just use json array string
// output: 3
r1, e1 := pcpServer.Execute("[\"add\", 1, 2]")

// we can also use PcpClient, instead of raw string
var pcpClient = PcpClient{}
callText, jerr := pcpClient.ToJSON(pcpClient.Call("add", 1, 2))

r2, e2 := pcpServer.Execute(callText)
```
