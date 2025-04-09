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
	"maps"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/qiniu/x/test"
)

// -----------------------------------------------------------------------------

type CaseT = test.CaseT

type CaseApp struct {
	*Request
	*App
	test.Case
	ctx context.Context
}

// Req create a new request with given method.
func (p *CaseApp) Req__0(method string) *Request {
	p.Request = &Request{method: method, ctx: p}
	return p.Request
}

// Req returns current request object.
func (p *CaseApp) Req__1() *Request {
	return p.Request
}

// Initialize creates a new request to initialize the MCP server.
func (p *CaseApp) Initialize(params map[string]any) *Request {
	return p.Req__0("initialize").Params(mapJoin(map[string]any{
		"protocolVersion": mcp.LATEST_PROTOCOL_VERSION,
		"capabilities": map[string]any{
			"roots":    map[string]any{"listChanged": true},
			"sampling": map[string]any{},
		},
		"clientInfo": map[string]any{
			"name":    "github.com/goplus/mcp/mtest",
			"version": "0.3.0",
		},
	}, params))
}

// OnNotify registers a notification handler.
func (p *CaseApp) OnNotify(notify func(method string, params map[string]any)) {
	p.client.OnNotification(func(in mcp.JSONRPCNotification) {
		notify(in.Method, makeParams(in.Params.Meta, in.Params.AdditionalFields))
	})
}

func makeParams(meta, addition map[string]any) map[string]any {
	if meta == nil {
		return addition
	}
	m := make(map[string]any, len(addition)+1)
	maps.Copy(m, addition)
	m["_meta"] = meta
	return m
}

func mapJoin(a, b map[string]any) map[string]any {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	m := make(map[string]any, len(a)+len(b))
	maps.Copy(m, a)
	maps.Copy(m, b)
	return m
}

// -----------------------------------------------------------------------------

var _ = (*CaseApp).initCaseApp

type iCaseProto interface {
	initCaseApp(*App, CaseT)
	Main()
}

// Gopt_CaseApp_TestMain is required by Go+ compiler as the entry of a YAP test case.
func Gopt_CaseApp_TestMain(c iCaseProto, t *testing.T) {
	app := new(App).initApp()
	c.initCaseApp(app, test.NewT(t))
	c.Main()
	app.shutdown()
}

func (p *CaseApp) initCaseApp(app *App, t CaseT) {
	p.App = app
	p.CaseT = t
	p.ctx = context.TODO()
}

// -----------------------------------------------------------------------------
