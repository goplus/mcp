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

// Req create a new request given method.
func (p *CaseApp) Req__0(method string) *Request {
	p.Request = &Request{method: method, ctx: p}
	return p.Request
}

// Req returns current request object.
func (p *CaseApp) Req__1() *Request {
	return p.Request
}

// OnNotify registers a notification handler.
func (p *CaseApp) OnNotify(notify func(method string, params map[string]any)) {
	p.client.OnNotification(func(in mcp.JSONRPCNotification) {
		notify(in.Method, makeParams(in.Params))
	})
}

func makeParams(in mcp.NotificationParams) map[string]any {
	if in.Meta == nil {
		return in.AdditionalFields
	}
	m := make(map[string]any, len(in.AdditionalFields)+1)
	maps.Copy(m, in.AdditionalFields)
	m["_meta"] = in.Meta
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
}

func (p *CaseApp) initCaseApp(app *App, t CaseT) {
	p.App = app
	p.CaseT = t
	p.ctx = context.TODO()
}

// -----------------------------------------------------------------------------
