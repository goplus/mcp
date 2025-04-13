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
	"github.com/goplus/mcp/mtest/rtx"
	"github.com/qiniu/x/test"
)

// -----------------------------------------------------------------------------

// Request represents a request to a MCP server.
type Request struct {
	method string
	params rtx.M
	ctx    *CaseApp
	resp   rtx.M
	rerr   error
}

func (p *Request) t() CaseT {
	return p.ctx.CaseT
}

// Params sets the request parameters.
func (p *Request) Params(params rtx.M) *Request {
	p.params = params
	return p
}

// Resp returns the response.
func (p *Request) Resp() rtx.M {
	return p.resp
}

// LastErr returns the last error.
func (p *Request) LastErr() error {
	return p.rerr
}

const (
	Gopo_Request_Ret = ".Send,.RetWith"
)

// RetWith checks the response with the given value.
func (p *Request) RetWith(resp any) *Request {
	t := p.t()
	t.Helper()
	p.Send()
	if p.rerr != nil {
		t.Fatal(p.rerr)
	}
	test.Gopt_Case_MatchAny(t, resp, p.resp, "resp")
	return p
}

// Send s the request to the server and returns the response.
func (p *Request) Send() *Request {
	app := p.ctx
	p.resp, p.rerr = app.rt.RoundTrip(app.ctx, p.method, p.params)
	return p
}

// -----------------------------------------------------------------------------
