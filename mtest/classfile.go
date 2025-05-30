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

package mtest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/goplus/mcp/mtest/mock"
	"github.com/goplus/mcp/mtest/rtx"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/server"
)

const (
	GopPackage   = "github.com/qiniu/x/test"
	GopTestClass = true
)

// Dump prints the arguments in a formatted JSON style.
func Dump(args ...any) {
	in := make([]any, len(args))
	for i, arg := range args {
		if _, ok := arg.(error); ok {
			in[i] = arg
		} else if b, e := json.MarshalIndent(arg, "", "  "); e == nil {
			in[i] = string(b)
		} else {
			in[i] = arg
		}
	}
	fmt.Println(in...)
}

// -----------------------------------------------------------------------------

type MainApp struct {
}

// Gopt_MainApp_TestMain is required by XGo compiler as the TestMain entry of a YAP testing project.
func Gopt_MainApp_TestMain(app any, m *testing.M) {
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	}
	os.Exit(m.Run())
}

// -----------------------------------------------------------------------------

// App is the application for MCP Server testing.
type App struct {
	rt rtx.RoundTripper
}

func (p *App) initApp() *App {
	return p
}

func (p *App) shutdown() {
	if rt := p.rt; rt != nil {
		p.rt = nil
		rt.Close()
	}
}

// MCPAppType represents the interface of a MCP Server application.
type MCPAppType interface {
	SetLAS(las func(addr string, svr *server.MCPServer) error)
	Main()
}

// Mock runs a MCP server by a mock transport.
func (p *App) Mock(app MCPAppType) {
	app.SetLAS(func(addr string, svr *server.MCPServer) (err error) {
		p.rt = mock.New(svr)
		log.Println("Mocking MCP server")
		return nil
	})
	app.Main()
}

// TestServer runs a MCP server by httptest.Server.
func (p *App) TestServer__0(app MCPAppType) {
	p.TestServer__1("/sse", app)
}

// TestServer runs a MCP server by httptest.Server.
func (p *App) TestServer__1(path string, app MCPAppType) {
	app.SetLAS(func(addr string, svr *server.MCPServer) (err error) {
		ts := server.NewTestServer(svr)
		log.Println("Serving MCP server at", ts.URL)
		baseURL := ts.URL + path
		client, err := transport.NewSSE(baseURL)
		if err != nil {
			log.Println("NewSSEMCPClient:", err)
			return
		}
		err = client.Start(context.Background())
		if err != nil {
			log.Println("SSEMCPClient.Start:", err)
		}
		p.rt = rtx.New(client)
		return
	})
	app.Main()
}

// -----------------------------------------------------------------------------
