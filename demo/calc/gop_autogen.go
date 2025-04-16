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
//line demo/calc/calc_mtest.gox:27:1
	this.ToolApp.Main(_gop_arg0, _gop_arg1, _gop_arg2)
//line demo/calc/calc_tool.gox:1:1
	this.Tool("calculate", func() {
//line demo/calc/calc_tool.gox:2:1
		this.Description("Perform basic arithmetic operations")
//line demo/calc/calc_tool.gox:3:1
		this.String("operation", func() {
//line demo/calc/calc_tool.gox:4:1
			this.Required()
//line demo/calc/calc_tool.gox:5:1
			this.Description("The operation to perform (add, subtract, multiply, divide)")
//line demo/calc/calc_tool.gox:6:1
			this.Enum("add", "subtract", "multiply", "divide")
		})
//line demo/calc/calc_tool.gox:8:1
		this.Float("x", func() {
//line demo/calc/calc_tool.gox:9:1
			this.Required()
//line demo/calc/calc_tool.gox:10:1
			this.Description("First number")
		})
//line demo/calc/calc_tool.gox:12:1
		this.Float("y", func() {
//line demo/calc/calc_tool.gox:13:1
			this.Required()
//line demo/calc/calc_tool.gox:14:1
			this.Description("Second number")
		})
	})
//line demo/calc/calc_tool.gox:18:1
	op := this.Gop_Env("operation").(string)
//line demo/calc/calc_tool.gox:19:1
	x := this.Gop_Env("x").(float64)
//line demo/calc/calc_tool.gox:20:1
	y := this.Gop_Env("y").(float64)
//line demo/calc/calc_tool.gox:22:1
	switch op {
//line demo/calc/calc_tool.gox:23:1
	case "add":
//line demo/calc/calc_tool.gox:24:1
		return server.Number__0(x + y)
//line demo/calc/calc_tool.gox:25:1
	case "subtract":
//line demo/calc/calc_tool.gox:26:1
		return server.Number__0(x - y)
//line demo/calc/calc_tool.gox:27:1
	case "multiply":
//line demo/calc/calc_tool.gox:28:1
		return server.Number__0(x * y)
//line demo/calc/calc_tool.gox:29:1
	case "divide":
//line demo/calc/calc_tool.gox:30:1
		if y == 0 {
//line demo/calc/calc_tool.gox:31:1
			return server.WithError__0("cannot divide by zero")
		}
//line demo/calc/calc_tool.gox:33:1
		return server.Number__0(x / y)
//line demo/calc/calc_tool.gox:34:1
	default:
//line demo/calc/calc_tool.gox:35:1
		return server.WithError__0("invalid operation")
	}
}
func (this *calc) Classclone() server.ToolProto {
	_gop_ret := *this
	return &_gop_ret
}
func main() {
//line demo/calc/calc_tool.gox:22:1
	new(MCPApp).Main()
}
