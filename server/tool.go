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

// Text creates a new CallToolResult with a text content
func Text(text string) *mcp.CallToolResult {
	return mcp.NewToolResultText(text)
}

// -----------------------------------------------------------------------------

type stop struct{}

type ToolType struct {
	tool  mcp.Tool
	clone func() any
}

// ToolApp is a worker class of a MCPServer classfile.
type ToolApp struct {
	*ToolType
	ctx     context.Context
	request mcp.CallToolRequest
	isClone bool
}

func (p *ToolApp) Gop_Env(name string) any {
	panic("todo")
}

// Main is required by Go+ compiler as the entry of a MCPServer tool.
func (p *ToolApp) Main(ctx context.Context, request mcp.CallToolRequest, t *ToolType) *mcp.CallToolResult {
	if t == nil {
		p.ctx = ctx
		p.request = request
		p.isClone = true
	} else {
		p.ToolType = t
	}
	return nil
}

func (p *ToolApp) initTool(name string, clone func() any) {
	defer func() {
		e := recover()
		if e, ok := e.(stop); !ok {
			panic(e)
		}
	}()
	p.Main(context.TODO(), mcp.CallToolRequest{}, &ToolType{
		tool: mcp.Tool{
			Name: name,
		},
		clone: clone,
	})
}

func (p *ToolApp) Tool(fn func()) {
	if !p.isClone {
		fn()
		panic(stop{})
	}
}

// Description sets a description to the Tool.
// The description should provide a clear, human-readable explanation of
// what the tool does.
func (p *ToolApp) Description(description string) {
	p.tool.Description = description
}

func (p *ToolApp) String(name string, fn func()) {
	panic("todo")
}

func (p *ToolApp) Required() {
	panic("todo")
}

func (p *ToolApp) addTo(self iHandlerProto, svr *server.MCPServer) {
	clone := self.Classclone
	p.initTool(self.Classfname(), clone)
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
		ret = clone().(*ToolApp).Main(ctx, request, nil)
		return
	})
}

// -----------------------------------------------------------------------------
