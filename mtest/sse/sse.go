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

package sse

import (
	"context"
	"encoding/json"
	"maps"

	"github.com/goplus/mcp/mtest/rtx"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// -----------------------------------------------------------------------------

type Transport struct {
	client *client.SSEMCPClient
}

func New(client *client.SSEMCPClient) *Transport {
	return &Transport{
		client: client,
	}
}

func (p *Transport) Close() error {
	return p.client.Close()
}

func (p *Transport) OnNotify(notify func(method string, params rtx.M)) {
	p.client.OnNotification(func(in mcp.JSONRPCNotification) {
		notify(in.Method, makeParams(in.Params.Meta, in.Params.AdditionalFields))
	})
}

func makeParams(meta, addition rtx.M) rtx.M {
	if meta == nil {
		return addition
	}
	m := make(rtx.M, len(addition)+1)
	maps.Copy(m, addition)
	m["_meta"] = meta
	return m
}

type jsonrpcRequest struct {
	Method string `json:"method"`
	Params rtx.M  `json:"params,omitempty"`
}

func (p *Transport) RoundTrip(ctx context.Context, method string, params rtx.M) (ret rtx.M, err error) {
	if fn, ok := routes[method]; ok {
		req, e := json.Marshal(jsonrpcRequest{method, params})
		if e != nil {
			return nil, e
		}
		resp, e := fn(p.client, ctx, req)
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
	return nil, rtx.ErrUnknownMethod
}

var routes = map[string]func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error){
	"initialize": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.InitializeRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return c.Initialize(ctx, in)
	},
	"ping": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		return nil, c.Ping(ctx)
	},
	"resources/list": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.ListResourcesRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return c.ListResources(ctx, in)
	},
	"resources/templates/list": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.ListResourceTemplatesRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return c.ListResourceTemplates(ctx, in)
	},
	"resources/read": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.ReadResourceRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return c.ReadResource(ctx, in)
	},
	"resources/subscribe": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.SubscribeRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return nil, c.Subscribe(ctx, in)
	},
	"resources/unsubscribe": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.UnsubscribeRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return nil, c.Unsubscribe(ctx, in)
	},
	"prompts/list": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.ListPromptsRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return c.ListPrompts(ctx, in)
	},
	"prompts/get": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.GetPromptRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return c.GetPrompt(ctx, in)
	},
	"tools/list": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.ListToolsRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return c.ListTools(ctx, in)
	},
	"tools/call": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.CallToolRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return c.CallTool(ctx, in)
	},
	"logging/setLevel": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.SetLevelRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return nil, c.SetLevel(ctx, in)
	},
	"completion/complete": func(c *client.SSEMCPClient, ctx context.Context, req []byte) (any, error) {
		var in mcp.CompleteRequest
		if err := json.Unmarshal(req, &in); err != nil {
			return nil, err
		}
		return c.Complete(ctx, in)
	},
}

// -----------------------------------------------------------------------------
