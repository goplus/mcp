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

package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/goplus/mcp/mtest/rtx"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// -----------------------------------------------------------------------------
// implement ClientSession

func (p *Transport) Initialize() {
}

func (p *Transport) Initialized() bool {
	return true
}

func (p *Transport) NotificationChannel() chan<- mcp.JSONRPCNotification {
	return p.ch
}

func (p *Transport) SessionID() string {
	panic("unreachable")
}

// -----------------------------------------------------------------------------

type notifyNode struct {
	notify func(method string, params rtx.M)
	prev   *notifyNode
}

type Transport struct {
	svr       *server.MCPServer
	ch        chan<- mcp.JSONRPCNotification
	notify    atomic.Pointer[notifyNode]
	requestID atomic.Int64
}

func New(svr *server.MCPServer) *Transport {
	ch := make(chan mcp.JSONRPCNotification, 8)
	ret := &Transport{
		svr: svr,
		ch:  ch,
	}
	go func(ch chan mcp.JSONRPCNotification, t *Transport) {
		for in := range ch {
			for lst := ret.notify.Load(); lst != nil; lst = lst.prev {
				lst.notify(in.Method, rtx.Params(in.Params.Meta, in.Params.AdditionalFields))
			}
		}
	}(ch, ret)
	return ret
}

func (p *Transport) Close() error {
	if ch := p.ch; ch != nil {
		p.ch = nil
		close(ch)
	}
	return nil
}

func (p *Transport) OnNotify(notify func(method string, params rtx.M)) {
	n := &notifyNode{
		notify: notify,
	}
	for {
		prev := p.notify.Load()
		n.prev = prev
		if p.notify.CompareAndSwap(prev, n) {
			break
		}
	}
}

func (p *Transport) RoundTrip(ctx context.Context, method string, params rtx.M) (ret rtx.M, err error) {
	id := p.requestID.Add(1)
	request := mcp.JSONRPCRequest{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      id,
		Request: mcp.Request{
			Method: method,
		},
		Params: params,
	}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	svr := p.svr
	ctx = svr.WithContext(ctx, p)
	resp := svr.HandleMessage(ctx, requestBytes)
	switch resp := resp.(type) {
	case mcp.JSONRPCResponse:
		b, e := json.Marshal(resp.Result)
		if e != nil {
			return nil, e
		}
		err = json.Unmarshal(b, &ret)
	case mcp.JSONRPCError:
		err = &rtx.Error{
			Code:    resp.Error.Code,
			Message: resp.Error.Message,
		}
	default:
		panic("unexpected response type")
	}
	return
}

// -----------------------------------------------------------------------------
