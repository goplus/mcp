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

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// -----------------------------------------------------------------------------

type ToolAppProto struct {
	tool mcp.Tool
	opts []mcp.PropertyOption
}

// ToolApp is a work class of a MCPServer classfile.
type ToolApp struct {
	*ToolAppProto
	ctx     context.Context
	request mcp.CallToolRequest
	isClone bool
	kind    byte
}

const (
	kindString  = 1
	kindNumber  = 2
	kindBoolean = 3
	kindArray   = 4
	kindObject  = 5
)

// Gop_Env returns the value of the specified parameter.
func (p *ToolApp) Gop_Env(name string) any {
	return p.request.Params.Arguments[name]
}

// Main is required by Go+ compiler as the entry of a MCPServer tool.
func (p *ToolApp) Main(ctx context.Context, request mcp.CallToolRequest, t *ToolAppProto) mcp.Content {
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
	p.kind = kindString
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
	p.kind = kindNumber
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
	p.kind = kindBoolean
	if len(fn) > 0 {
		p.opts = make([]mcp.PropertyOption, 0, 2)
		fn[0]()
	}
	mcp.WithBoolean(name, p.opts...)(&p.tool)
	p.opts = nil
}

// Array adds an array property to the tool schema.
// It accepts property options to configure the array property's behavior
// and constraints.
func (p *ToolApp) Array(name string, fn ...func()) {
	p.kind = kindArray
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
	p.kind = kindObject
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

// Defval sets the default value for a string property.
// This value will be used if the property is not explicitly provided.
func (p *ToolApp) Defval__0(value string) {
	if p.kind != kindString {
		panic("defval: not a string property")
	}
	p.opts = append(p.opts, mcp.DefaultString(value))
}

// Defval sets the default value for a number property.
// This value will be used if the property is not explicitly provided.
func (p *ToolApp) Defval__1(value float64) {
	if p.kind != kindNumber {
		panic("defval: not a number property")
	}
	p.opts = append(p.opts, mcp.DefaultNumber(value))
}

// Defval sets the default value for a boolean property.
// This value will be used if the property is not explicitly provided.
func (p *ToolApp) Defval__2(value bool) {
	if p.kind != kindBoolean {
		panic("defval: not a boolean property")
	}
	p.opts = append(p.opts, mcp.DefaultBool(value))
}

// Maxval sets the maximum value for a number property. The number value
// must not exceed this maximum.
func (p *ToolApp) Maxval(value float64) {
	if p.kind != kindNumber {
		panic("maxval: not a number property")
	}
	p.opts = append(p.opts, mcp.Max(value))
}

// Min sets the minimum value for a number property. The number value must
// not be less than this minimum.
func (p *ToolApp) Minval(value float64) {
	if p.kind != kindNumber {
		panic("minval: not a number property")
	}
	p.opts = append(p.opts, mcp.Min(value))
}

// Maxlen sets the maximum length for a string or an array property. The
// string or array length must not exceed this length.
func (p *ToolApp) Maxlen(n int) {
	switch p.kind {
	case kindString:
		p.opts = append(p.opts, mcp.MaxLength(n))
	case kindArray:
		p.opts = append(p.opts, mcp.MaxItems(n))
	default:
		panic("maxlen: not a string or array property")
	}
}

// Minlen sets the minimum length for a string or an array property. The
// string or array length must not be less than this length.
func (p *ToolApp) Minlen(n int) {
	switch p.kind {
	case kindString:
		p.opts = append(p.opts, mcp.MinLength(n))
	case kindArray:
		p.opts = append(p.opts, mcp.MinItems(n))
	default:
		panic("minlen: not a string or array property")
	}
}

// Enum specifies a list of allowed values for a string property. The property
// value must be one of the specified enum values.
func (p *ToolApp) Enum(values ...string) {
	if p.kind != kindString {
		panic("enum: not a string property")
	}
	p.opts = append(p.opts, mcp.Enum(values...))
}

// Unique specifies that the array property must contain unique items.
func (p *ToolApp) Unique() {
	if p.kind != kindArray {
		panic("unique: not an array property")
	}
	p.opts = append(p.opts, mcp.UniqueItems(true))
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
		content := clone().Main(ctx, request, nil)
		ret = new(mcp.CallToolResult)
		if multi, ok := content.(*multiContents); ok {
			ret.Content = multi.data
		} else {
			ret.Content = []mcp.Content{content}
		}
		return
	})
}

// -----------------------------------------------------------------------------
