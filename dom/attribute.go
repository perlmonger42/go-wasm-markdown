package dom

import "syscall/js"

type Callback func(this js.Value, args []js.Value) interface{}

type GoAttribute interface {
	SetAttributeOn(domElement js.Value)
	GetName() string
	DebugValue() string
}

// StringAttribute represents a DOM attribute with a string value.
type stringAttribute struct {
	Name  string
	Value string
}

var _ GoAttribute = &stringAttribute{}

func NewStringAttribute(name, value string) GoAttribute {
	return &stringAttribute{Name: name, Value: value}
}

func (attr *stringAttribute) GetName() string {
	return attr.Name
}

func (attr *stringAttribute) DebugValue() string {
	return attr.Value
}

func (attr *stringAttribute) SetAttributeOn(domElement js.Value) {
	domElement.Set(attr.Name, js.ValueOf(attr.Value))
}

// functionAttribute represents a DOM attribute with a function value.
type functionAttribute struct {
	Name   string
	Value  Callback
	FnName string
}

var _ GoAttribute = &functionAttribute{}

func NewFunctionAttribute(name string, functionName string, value Callback) GoAttribute {
	return &functionAttribute{Name: name, FnName: functionName, Value: value}
}

func (attr *functionAttribute) GetName() string {
	return attr.Name
}

func (attr *functionAttribute) DebugValue() string {
	return attr.FnName
}

func (attr *functionAttribute) SetAttributeOn(domElement js.Value) {
	domElement.Set(attr.Name, js.FuncOf(attr.Value))
}
