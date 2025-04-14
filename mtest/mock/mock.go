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

type Transport struct {
	svr *server.MCPServer
}

func New(svr *server.MCPServer) Transport {
	return Transport{
		svr: svr,
	}
}

func (p Transport) Close() error {
	return nil
}

func (p Transport) OnNotify(notify func(method string, params rtx.M)) {
	panic("todo")
}

var (
	requestID atomic.Int64
)

func (p Transport) RoundTrip(ctx context.Context, method string, params rtx.M) (ret rtx.M, err error) {
	id := requestID.Add(1)
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
	resp := p.svr.HandleMessage(ctx, requestBytes)
	switch resp := resp.(type) {
	case mcp.JSONRPCResponse:
		b, e := json.Marshal(resp.Result)
		if e != nil {
			return nil, e
		}
		err = json.Unmarshal(b, &ret)
	case mcp.JSONRPCError:
		err = &jsonrpcError{
			Code:    resp.Error.Code,
			Message: resp.Error.Message,
		}
	default:
		panic("unexpected response type")
	}
	return
}

type jsonrpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (p *jsonrpcError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", p.Code, p.Message)
}

// -----------------------------------------------------------------------------
