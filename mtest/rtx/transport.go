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

package rtx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"sync/atomic"

	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

var (
	// ErrUnknownMethod is returned when the method is not recognized.
	ErrUnknownMethod = errors.New("unknown method")
)

// -----------------------------------------------------------------------------

type M = map[string]any

type RoundTripper interface {
	RoundTrip(ctx context.Context, method string, params M) (resp M, err error)
	OnNotify(notify func(method string, params M))
	Close() error
}

// -----------------------------------------------------------------------------

type Transport struct {
	t         transport.Interface
	requestID atomic.Int64
}

func New(t transport.Interface) *Transport {
	return &Transport{
		t: t,
	}
}

// Close the connection.
func (p *Transport) Close() error {
	return p.t.Close()
}

func (p *Transport) OnNotify(notify func(method string, params M)) {
	p.t.SetNotificationHandler(func(in mcp.JSONRPCNotification) {
		notify(in.Method, Params(in.Params.Meta, in.Params.AdditionalFields))
	})
}

func Params(meta, addition M) M {
	if meta == nil {
		return addition
	}
	m := make(M, len(addition)+1)
	maps.Copy(m, addition)
	m["_meta"] = meta
	return m
}

func (p *Transport) RoundTrip(ctx context.Context, method string, params M) (ret M, err error) {
	id := p.requestID.Add(1)
	request := transport.JSONRPCRequest{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      id,
		Method:  method,
		Params:  params,
	}
	resp, err := p.t.SendRequest(ctx, request)
	if err != nil {
		return
	}
	if e := resp.Error; e != nil {
		err = &Error{
			Code:    e.Code,
			Message: e.Message,
		}
	} else {
		err = json.Unmarshal(resp.Result, &ret)
	}
	return
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (p *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s", p.Code, p.Message)
}

// -----------------------------------------------------------------------------
