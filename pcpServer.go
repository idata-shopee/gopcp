package gopcp

// function node
type FunNode struct {
	funName string
	params  []interface{}
}

// simpe calling protocol
// grammer based on json
// ["fun1", 1, 2, ["fun2", 3]] => fun1(1, 2, fun2(3))

type PcpServer struct {
	sandbox Sandbox
}

func (pcpServer PcpServer) parseJson(source string) {}

//func getPcpServer(sandbox Sandbox) PcpServer {
// merge sandbox with default sandbox
//}
