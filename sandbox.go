package gopcp

type GeneralFun = func([]interface{}, PcpServer) interface{}

// funType: 0 -> normal sandbox function
// funType: 1 -> lazy sandbox function
type BoxFun struct {
	funType int
	fun     GeneralFun
}

type Sandbox struct{}
