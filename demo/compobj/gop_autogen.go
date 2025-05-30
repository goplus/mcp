// Code generated by xgo (XGo); DO NOT EDIT.

package main

import (
	"context"
	"github.com/goplus/mcp/server"
	"github.com/mark3labs/mcp-go/mcp"
)

const _ = true

type compobj struct {
	server.ToolApp
	*MCPApp
}
type MCPApp struct {
	server.MCPApp
}
//line demo/compobj/main_mcp.gox:1
func (this *MCPApp) MainEntry() {
//line demo/compobj/main_mcp.gox:1:1
	this.Server("Composite Object Demo 🚀", "1.0.0")
}
func (this *MCPApp) Main() {
	_gop_obj0 := &compobj{MCPApp: this}
	_gop_lst1 := []server.ToolProto{_gop_obj0}
	server.Gopt_MCPApp_Main(this, nil, _gop_lst1, nil)
}
//line demo/compobj/compobj_tool.gox:1
func (this *compobj) Main(_gop_arg0 context.Context, _gop_arg1 mcp.CallToolRequest, _gop_arg2 *server.ToolAppProto) mcp.Content {
	this.ToolApp.Main(_gop_arg0, _gop_arg1, _gop_arg2)
//line demo/compobj/compobj_tool.gox:1:1
	this.Tool("compositeObject", func() {
//line demo/compobj/compobj_tool.gox:2:1
		this.Description("A composite object tool")
//line demo/compobj/compobj_tool.gox:3:1
		this.String("name", func() {
//line demo/compobj/compobj_tool.gox:4:1
			this.Required()
//line demo/compobj/compobj_tool.gox:5:1
			this.Description("User name")
		})
//line demo/compobj/compobj_tool.gox:7:1
		this.Object("profile", func() {
//line demo/compobj/compobj_tool.gox:8:1
			this.Required()
//line demo/compobj/compobj_tool.gox:9:1
			this.Description("User profile")
//line demo/compobj/compobj_tool.gox:10:1
			this.Float("age")
//line demo/compobj/compobj_tool.gox:11:1
			this.Array("works", func() {
//line demo/compobj/compobj_tool.gox:12:1
				this.Description("Work experience")
//line demo/compobj/compobj_tool.gox:13:1
				this.Object("items", func() {
//line demo/compobj/compobj_tool.gox:14:1
					this.String("company", func() {
//line demo/compobj/compobj_tool.gox:15:1
						this.Required()
//line demo/compobj/compobj_tool.gox:16:1
						this.Description("Company name")
					})
//line demo/compobj/compobj_tool.gox:18:1
					this.String("start", func() {
//line demo/compobj/compobj_tool.gox:19:1
						this.Required()
//line demo/compobj/compobj_tool.gox:20:1
						this.Description("Start time")
					})
//line demo/compobj/compobj_tool.gox:22:1
					this.String("end", func() {
//line demo/compobj/compobj_tool.gox:23:1
						this.Description("End time")
					})
				})
			})
		})
	})
//line demo/compobj/compobj_tool.gox:30:1
	name, ok := this.Gop_Env("name").(string)
//line demo/compobj/compobj_tool.gox:31:1
	if !ok {
//line demo/compobj/compobj_tool.gox:32:1
		panic("name must be a string")
	}
//line demo/compobj/compobj_tool.gox:35:1
	profile, ok := this.Gop_Env("profile").(map[string]interface{})
//line demo/compobj/compobj_tool.gox:36:1
	if !ok {
//line demo/compobj/compobj_tool.gox:37:1
		panic("profile must be an object")
	}
//line demo/compobj/compobj_tool.gox:40:1
	return server.Text__1(server.JsonContent{JSON: map[string]interface{}{"name": name, "profile": profile}})
}
func (this *compobj) Classclone() server.ToolProto {
	_gop_ret := *this
	return &_gop_ret
}
func main() {
	new(MCPApp).Main()
}
