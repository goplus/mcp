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
	"testing"

	"github.com/goplus/mcp/server/svx"
	"github.com/mark3labs/mcp-go/server"
)

func TestServerWithBasePath(t *testing.T) {
	ret, err := NewServer("sse://:8080/base", nil)
	if err != nil {
		t.Fatal("failed to create server:", err)
	}
	sse := ret.Handler.(*server.SSEServer)
	if path := sse.CompleteSsePath(); path != `/base/sse` {
		t.Fatal("expected `/base/sse`, got", path)
	}
	if ret.Addr != ":8080" {
		t.Fatal("expected address :8080, got", ret.Addr)
	}
}

func TestServerWithEndpoint(t *testing.T) {
	ret, err := NewServer("sse://localhost:8080?sse=/newsse&msg=/msg", nil)
	if err != nil {
		t.Fatal("failed to create server:", err)
	}
	sse := ret.Handler.(*server.SSEServer)
	if path := sse.CompleteSsePath(); path != `/newsse` {
		t.Fatal("expected `/newsse`, got", path)
	}
	if path := sse.CompleteMessagePath(); path != `/msg` {
		t.Fatal("expected `/msg`, got", path)
	}
	if ret.Addr != "localhost:8080" {
		t.Fatal("expected address localhost:8080, got", ret.Addr)
	}
}

func TestErrServer(t *testing.T) {
	_, err := NewServer("newsse://localhost:8080", nil)
	if err != svx.ErrUnknownScheme {
		t.Fatal("NewServer:", err)
	}
	_, err = NewServer("sse://localhost:8080?;", nil)
	if err == nil || err.Error() != "invalid semicolon separator in query" {
		t.Fatal("NewServer:", err)
	}
	err = ListenAndServe("newsse://localhost:8080", nil)
	if err != svx.ErrUnknownScheme {
		t.Fatal("ListenAndServe:", err)
	}
	err = ListenAndServe("sse://1.1.1.1:8080", nil)
	if err == nil {
		t.Fatal("ListenAndServe: no error?")
	}
	err = ListenAndServeTLS("newsse://localhost:8080", "", "", nil)
	if err != svx.ErrUnknownScheme {
		t.Fatal("ListenAndServeTLS:", err)
	}
	err = ListenAndServeTLS("sse://1.1.1.1:8080", "", "", nil)
	if err == nil {
		t.Fatal("ListenAndServeTLS: no error?")
	}
}
