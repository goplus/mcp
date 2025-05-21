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
	"log"
	"net/http"
	"net/url"

	"github.com/goplus/mcp/server/svx"
	"github.com/mark3labs/mcp-go/server"
)

const (
	Scheme = "sse"
)

func init() {
	svx.Register(Scheme, ListenAndServe)
}

// -----------------------------------------------------------------------------

// ListenAndServe starts a new SSE server with the given address and server.
// addr = "sse://<Host>[<Path>?sse=<Endpoint>&msg=<MsgEndpoint>]"
func ListenAndServe(addr string, svr *server.MCPServer) (err error) {
	server, err := NewServer(addr, svr)
	if err != nil {
		return
	}
	log.Println("Serving MCP server at", addr)
	return server.ListenAndServe()
}

// ListenAndServeTLS starts a new SSE server with TLS support.
func ListenAndServeTLS(addr, certFile, keyFile string, svr *server.MCPServer) (err error) {
	server, err := NewServer(addr, svr)
	if err != nil {
		return
	}
	log.Println("Serving TLS MCP server at", addr)
	return server.ListenAndServeTLS(certFile, keyFile)
}

// NewServer creates a new SSE server with the given address and server.
func NewServer(addr string, svr *server.MCPServer) (ret *http.Server, err error) {
	opts, err := ParseAddr(addr)
	if err != nil {
		return
	}
	options := make([]server.SSEOption, 0, 3)
	if opts.Path != "" {
		options = append(options, server.WithStaticBasePath(opts.Path))
	}
	if opts.Endpoint != "" {
		options = append(options, server.WithSSEEndpoint(opts.Endpoint))
	}
	if opts.MsgEndpoint != "" {
		options = append(options, server.WithMessageEndpoint(opts.MsgEndpoint))
	}
	sse := server.NewSSEServer(svr, options...)
	return &http.Server{Addr: opts.Host, Handler: sse}, nil
}

// -----------------------------------------------------------------------------

type Options struct {
	Host        string
	Path        string
	Endpoint    string
	MsgEndpoint string
}

// ParseAddr parses the SSE address and returns all the options.
func ParseAddr(addr string) (ret Options, err error) {
	u, err := url.Parse(addr)
	if err != nil {
		return
	}
	if u.Scheme != Scheme {
		err = svx.ErrUnknownScheme
		return
	}
	ret.Host = u.Host
	ret.Path = u.Path
	if u.RawQuery == "" {
		return
	}
	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return
	}
	ret.Endpoint = params.Get("sse")
	ret.MsgEndpoint = params.Get("msg")
	return
}

// -----------------------------------------------------------------------------
