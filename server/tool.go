/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"context"
	"errors"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Text creates a new CallToolResult with a text content
func Text(text string) *mcp.CallToolResult {
	return mcp.NewToolResultText(text)
}

// Number creates a new CallToolResult with a number content
func Number__0(val float64) *mcp.CallToolResult {
	return mcp.NewToolResultText(strconv.FormatFloat(val, 'f', -1, 64))
}

// Number creates a new CallToolResult with a number content
func Number__1(val float64, prec int) *mcp.CallToolResult {
	return mcp.NewToolResultText(strconv.FormatFloat(val, 'f', prec, 64))
}

// Number creates a new CallToolResult with a number content
func Number__2(val float64, fmt byte, prec int) *mcp.CallToolResult {
	return mcp.NewToolResultText(strconv.FormatFloat(val, fmt, prec, 64))
}

/*
// Image creates a new CallToolResult with an image content
func Image(text, imageData, mimeType string) *mcp.CallToolResult {
	return mcp.NewToolResultImage(text, imageData, mimeType)
}

// Resource creates a new CallToolResult with a resource content
func Resource(text string, resource mcp.ResourceContents) *mcp.CallToolResult {
	return mcp.NewToolResultResource(text, resource)
}
*/

// -----------------------------------------------------------------------------

type stop struct{}

type ToolAppProto struct {
	tool mcp.Tool
	opts []mcp.PropertyOption
}

// ToolApp is a worker class of a MCPServer classfile.
type ToolApp struct {
	*ToolAppProto
	ctx     context.Context
	request mcp.CallToolRequest
	isClone bool
}

// Gop_Env returns the value of the specified parameter.
func (p *ToolApp) Gop_Env(name string) any {
	return p.request.Params.Arguments[name]
}

// Main is required by Go+ compiler as the entry of a MCPServer tool.
func (p *ToolApp) Main(ctx context.Context, request mcp.CallToolRequest, t *ToolAppProto) *mcp.CallToolResult {
	if t == nil {
		p.ctx = ctx
		p.request = request
		p.isClone = true
	} else {
		p.ToolAppProto = t
	}
	return nil
}

func initToolApp(self ToolProto) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(stop); !ok {
				panic(e)
			}
		}
	}()
	self.Main(context.TODO(), mcp.CallToolRequest{}, &ToolAppProto{
		tool: mcp.NewTool(""),
	})
}

// Tool calls fn to initialize the tool.
func (p *ToolApp) Tool(name string, fn func()) {
	if !p.isClone {
		p.tool.Name = name
		fn()
		panic(stop{})
	}
}

// For a tool:
// Description sets a description to the Tool.
// The description should provide a clear, human-readable explanation of
// what the tool does.
//
// For a tool property:
// Description adds a description to a property in the JSON Schema.
// The description should explain the purpose and expected values of the property.
func (p *ToolApp) Description(description string) {
	if p.opts != nil {
		p.opts = append(p.opts, mcp.Description(description))
	} else {
		mcp.WithDescription(description)(&p.tool)
	}
}

// String adds a string property to the tool schema.
// It accepts property options to configure the string property's behavior
// and constraints.
func (p *ToolApp) String(name string, fn ...func()) {
	if len(fn) > 0 {
		p.opts = make([]mcp.PropertyOption, 0, 2)
		fn[0]()
	}
	mcp.WithString(name, p.opts...)(&p.tool)
	p.opts = nil
}

// Float adds a number property to the tool schema.
// It accepts property options to configure the number property's behavior
// and constraints.
func (p *ToolApp) Float(name string, fn ...func()) {
	if len(fn) > 0 {
		p.opts = make([]mcp.PropertyOption, 0, 2)
		fn[0]()
	}
	mcp.WithNumber(name, p.opts...)(&p.tool)
	p.opts = nil
}

// Bool adds a boolean property to the tool schema.
// It accepts property options to configure the boolean property's behavior
// and constraints.
func (p *ToolApp) Bool(name string, fn ...func()) {
	if len(fn) > 0 {
		p.opts = make([]mcp.PropertyOption, 0, 2)
		fn[0]()
	}
	mcp.WithBoolean(name, p.opts...)(&p.tool)
	p.opts = nil
}

// List adds an array property to the tool schema.
// It accepts property options to configure the array property's behavior
// and constraints.
func (p *ToolApp) List(name string, fn ...func()) {
	if len(fn) > 0 {
		p.opts = make([]mcp.PropertyOption, 0, 2)
		fn[0]()
	}
	mcp.WithArray(name, p.opts...)(&p.tool)
	p.opts = nil
}

// Object adds an object property to the tool schema.
// It accepts property options to configure the object property's behavior
// and constraints.
func (p *ToolApp) Object(name string, fn ...func()) {
	if len(fn) > 0 {
		p.opts = make([]mcp.PropertyOption, 0, 2)
		fn[0]()
	}
	mcp.WithObject(name, p.opts...)(&p.tool)
	p.opts = nil
}

// Required marks a property as required in the tool's input schema.
// Required properties must be provided when using the tool.
func (p *ToolApp) Required() {
	p.opts = append(p.opts, mcp.Required())
}

func (p *ToolApp) addTo(self ToolProto, svr *server.MCPServer) {
	clone := self.Classclone
	initToolApp(self)
	svr.AddTool(p.tool, func(ctx context.Context, request mcp.CallToolRequest) (ret *mcp.CallToolResult, err error) {
		defer func() {
			if e := recover(); e != nil {
				switch e := e.(type) {
				case string:
					err = errors.New(e)
				case error:
					err = e
				default:
					err = errors.New("unknown error")
				}
			}
		}()
		ret = clone().Main(ctx, request, nil)
		return
	})
}

// -----------------------------------------------------------------------------
