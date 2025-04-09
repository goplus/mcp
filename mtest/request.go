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
	"errors"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/qiniu/x/test"
)

var (
	// ErrUnknownMethod is returned when the method is not recognized.
	ErrUnknownMethod = errors.New("unknown method")
)

// -----------------------------------------------------------------------------

// Request represents a request to a MCP server.
type Request struct {
	method string
	params map[string]any
	ctx    *CaseApp
	resp   map[string]any
	rerr   error
}

func (p *Request) t() CaseT {
	return p.ctx.CaseT
}

// Params sets the request parameters.
func (p *Request) Params(params map[string]any) *Request {
	p.params = params
	return p
}

// Resp returns the response.
func (p *Request) Resp() map[string]any {
	return p.resp
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
	p.resp, p.rerr = p.doSend()
	return p
}

func (p *Request) doSend() (ret map[string]any, err error) {
	app := p.ctx
	c := app.client
	if fn, ok := routes[p.method]; ok {
		req, e := json.Marshal(map[string]any{
			"method": p.method,
			"params": p.params,
		})
		if e != nil {
			return nil, e
		}
		resp, e := fn(c, app.ctx, req)
		if e != nil {
			return nil, e
		}
		b, e := json.Marshal(resp)
		if e != nil {
			return nil, e
		}
		err = json.Unmarshal(b, &ret)
		return
	}
	return nil, ErrUnknownMethod
}

var routes = map[string]func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error){
	"initialize": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.InitializeRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return c.Initialize(ctx, in)
	},
}

// -----------------------------------------------------------------------------
