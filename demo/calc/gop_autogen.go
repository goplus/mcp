// Code generated by gop (Go+); DO NOT EDIT.

package main

import (
	"context"
	"github.com/goplus/mcp/server"
	_ "github.com/goplus/mcp/server/sse"
	"github.com/mark3labs/mcp-go/mcp"
)

const _ = true

type calc struct {
	server.ToolApp
	*MCPApp
}
type MCPApp struct {
	server.MCPApp
}
//line demo/calc/main_mcp.gox:5
func (this *MCPApp) MainEntry() {
//line demo/calc/main_mcp.gox:5:1
	this.Server("Calculator", "1.0.0")
//line demo/calc/main_mcp.gox:6:1
	this.Run("sse://:8080")
}
func (this *MCPApp) Main() {
	server.Gopt_MCPApp_Main(this, nil, []server.ToolProto{new(calc)}, nil)
}
//line demo/calc/calc_tool.gox:1
func (this *calc) Main(_gop_arg0 context.Context, _gop_arg1 mcp.CallToolRequest, _gop_arg2 *server.ToolAppProto) mcp.Content {
//line demo/calc/calc_mtest.gox:7:1
	this.ToolApp.Main(_gop_arg0, _gop_arg1, _gop_arg2)
//line demo/calc/calc_tool.gox:1:1
	this.Tool("mul", func() {
//line demo/calc/calc_tool.gox:2:1
		this.Description("Multiply two numbers")
//line demo/calc/calc_tool.gox:3:1
		this.Float("x", func() {
//line demo/calc/calc_tool.gox:4:1
			this.Required()
//line demo/calc/calc_tool.gox:5:1
			this.Description("One of multiplicative parameters")
		})
//line demo/calc/calc_tool.gox:7:1
		this.Float("y", func() {
//line demo/calc/calc_tool.gox:8:1
			this.Required()
//line demo/calc/calc_tool.gox:9:1
			this.Description("One of multiplicative parameters")
		})
	})
//line demo/calc/calc_tool.gox:13:1
	x, ok1 := this.Gop_Env("x").(float64)
//line demo/calc/calc_tool.gox:14:1
	y, ok2 := this.Gop_Env("y").(float64)
//line demo/calc/calc_tool.gox:15:1
	if !ok1 || !ok2 {
//line demo/calc/calc_tool.gox:16:1
		panic("multiplicative parameters x and y must be numbers")
	}
//line demo/calc/calc_tool.gox:19:1
	return server.Multiple(server.Text("multiply result:"), server.Number__0(x*y))
}
func (this *calc) Classclone() server.ToolProto {
	_gop_ret := *this
	return &_gop_ret
}
func main() {
//line demo/calc/calc_tool.gox:19:1
	new(MCPApp).Main()
}
