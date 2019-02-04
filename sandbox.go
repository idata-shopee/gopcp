package gopcp

import (
	"errors"
	"fmt"
)

// (params, attachment, pcpServer) -> (result, error)
type GeneralFun = func([]interface{}, interface{}, *PcpServer) (interface{}, error)

// SandBoxType
const (
	SandboxTypeNormal = 1
	SandboxTypeLazy   = 2
)

// BoxFun
type BoxFunc struct {
	FunType int // SandBoxType
	Fun     GeneralFun
}

// Sandbox
type Sandbox struct {
	funcMap map[string]*BoxFunc // name -> boxFunc
}

func GetSandbox(box map[string]*BoxFunc) *Sandbox {
	return (&Sandbox{box}).Extend(DefBox)
}

// Get get sandbox method
func (s *Sandbox) Get(name string) (*BoxFunc, error) {
	if val, ok := s.funcMap[name]; ok {
		return val, nil
	}
	return nil, errors.New(fmt.Sprintf("function [%s] doesn't exist in sandBox", name))
}

// Set set sandbox method
func (s *Sandbox) Set(name string, val *BoxFunc) {
	s.funcMap[name] = val
	return
}

// Extend merge newSandBox's value to origin sandBox
func (s *Sandbox) Extend(newSandBox *Sandbox) *Sandbox {
	for k, v := range newSandBox.funcMap {
		s.Set(k, v)
	}
	return s
}

func ToSandboxFun(fun GeneralFun) *BoxFunc {
	return &BoxFunc{SandboxTypeNormal, fun}
}

func ToLazySandboxFun(fun GeneralFun) *BoxFunc {
	return &BoxFunc{SandboxTypeLazy, fun}
}
