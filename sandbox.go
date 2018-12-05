package gopcp

import "fmt"

type GeneralFun = func([]interface{}, PcpServer) interface{}

// SandBoxType
const (
	SandboxTypeNormal = "normal_sandbox_type"
	SandboxTypeLazy   = "lazy_sandbox_type"
)

// BoxFun
type BoxFun struct {
	funType string // SandBoxType
	fun     GeneralFun
}

// Sandbox
type Sandbox struct {
	funcMap map[string]BoxFunc
}

func NewSandbox(val *map[string]BoxFunc) {
	sandbox := &SandBox{}
	if val != nil {
		sandbox.funcMap = *val
	} else {
		sandbox.funcMap = map[string]BoxFunc{}
	}
	return sandbox
}

// Get get sandbox method
func (sandBox *SandBox) Get(name string) (*BoxFunc, error) {
	if val, ok := sandBox.funcMap[name]; ok {
		return &val, nil
	}
	return nil, error.New(fmt.Sprintf("function [%s] doesn't exist in sandBox", name))
}

// Set set sandbox method
func (sandBox *SandBox) Set(name string, val BoxFunc) {
	sandBox.funcMap[name] = val
}

// Extend merge newSandBox's value to origin sandBox
func (sandBox *SandBox) Extend(newSandBox *SandBox) {
	if newSandBox == nil {
		return
	}
	for k, v := range newSandBox.funcMap {
		sandBox.Set(k, v)
	}
}
