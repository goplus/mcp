// Code generated by gop (Go+); DO NOT EDIT.

package main

import (
	"github.com/goplus/mcp/mtest"
	"testing"
)

type case_hello struct {
	mtest.CaseApp
}
//line demo/hello/hello_mtest.gox:1
func (this *case_hello) Main() {
//line demo/hello/hello_mtest.gox:1:1
	this.TestServer__0(new(MCPApp))
//line demo/hello/hello_mtest.gox:3:1
	this.Initialize(nil)
//line demo/hello/hello_mtest.gox:4:1
	this.RetWith(map[string]interface{}{})
//line demo/hello/hello_mtest.gox:6:1
	this.Call("helloWorld", map[string]any{"name": "Ken"})
//line demo/hello/hello_mtest.gox:7:1
	this.RetWith(map[string][]map[string]string{"content": []map[string]string{map[string]string{"type": "text", "text": "Hello, Ken!"}}})
}
func Test_hello(t *testing.T) {
	mtest.Gopt_CaseApp_TestMain(new(case_hello), t)
}
