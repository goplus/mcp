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
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/server"
)

const (
	GopPackage   = "github.com/goplus/yap/test"
	GopTestClass = true
)

// -----------------------------------------------------------------------------

type MainApp struct {
}

// Gopt_MainApp_TestMain is required by Go+ compiler as the TestMain entry of a YAP testing project.
func Gopt_MainApp_TestMain(app any, m *testing.M) {
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	}
	os.Exit(m.Run())
}

// -----------------------------------------------------------------------------

// App is the application for MCP Server testing.
type App struct {
	client  *client.SSEMCPClient
	baseURL string
}

func (p *App) initApp() *App {
	return p
}

// MCPAppType represents the interface of a MCP Server application.
type MCPAppType interface {
	SetLAS(las func(addr string, svr *server.MCPServer) error)
	Main()
}

// TestServer runs a MCP server by httptest.Server.
func (p *App) TestServer(app MCPAppType) {
	app.SetLAS(func(addr string, svr *server.MCPServer) (err error) {
		ts := server.NewTestServer(svr)
		p.baseURL = ts.URL
		p.client, err = client.NewSSEMCPClient(ts.URL)
		if err != nil {
			return
		}
		err = p.client.Start(context.Background())
		return
	})
	app.Main()
}

// -----------------------------------------------------------------------------
