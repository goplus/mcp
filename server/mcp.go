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
	"log"
	"reflect"

	"github.com/goplus/mcp/server/svx"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	_ "github.com/goplus/mcp/server/stdio"
)

const (
	GopPackage = true
)

// -----------------------------------------------------------------------------

// MCPApp is the project class of a MCPServer classfile.
type MCPApp struct {
	svr  *server.MCPServer
	addr string
}

func (p *MCPApp) mcpServer() *server.MCPServer {
	return p.svr
}

// Server creates a new MCP server instance with the given name and version.
func (p *MCPApp) Server(name, version string) {
	p.svr = server.NewMCPServer(name, version)
	p.addr = "stdio:"
}

// Run sets the MCP server address.
func (p *MCPApp) Run(addr string) {
	p.addr = addr
}

func (p *MCPApp) serve() error {
	return svx.ListenAndServe(p.addr, p.svr)
}

// -----------------------------------------------------------------------------

var _ = (*ToolApp).addTo
var _ = (*MCPApp).mcpServer
var _ = (*MCPApp).serve

type iAppProto interface {
	mcpServer() *server.MCPServer
	serve() error
	MainEntry()
}

type iHandlerProto interface {
	addTo(self iHandlerProto, svr *server.MCPServer)
	Main(ctx context.Context, request mcp.CallToolRequest, t *ToolAppProto) *mcp.CallToolResult
	Classfname() string
	Classclone() any
}

// Gopt_MCPApp_Main is required by Go+ compiler as the entry of a MCPServer project.
func Gopt_MCPApp_Main(app iAppProto, handlers ...iHandlerProto) {
	app.MainEntry()
	svr := app.mcpServer()
	for _, h := range handlers {
		reflect.ValueOf(h).Elem().Field(1).Set(reflect.ValueOf(app)) // (*handler).App = app
		h.addTo(h, svr)
	}
	err := app.serve()
	if err != nil {
		log.Panicln(err)
	}
}

// -----------------------------------------------------------------------------
