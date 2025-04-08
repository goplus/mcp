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

package svx

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/mark3labs/mcp-go/server"
)

var (
	// ErrUnknownScheme is returned when the scheme of the address is unknown.
	ErrUnknownScheme = errors.New("unknown scheme")
)

// -----------------------------------------------------------------------------

type LAS = func(addr string, svr *server.MCPServer) error

var (
	svxs = make(map[string]LAS, 4)
)

// Register registers a LAS with specified scheme.
func Register(scheme string, las LAS) {
	svxs[scheme] = las
}

func ListenAndServe(addr string, svr *server.MCPServer) error {
	scheme := schemeOf(addr)
	if las, ok := svxs[scheme]; ok {
		return las(addr, svr)
	}
	return &fs.PathError{Op: "svx.ListenAndServe", Err: ErrUnknownScheme, Path: addr}
}

func schemeOf(url string) (scheme string) {
	pos := strings.IndexAny(url, ":/")
	if pos > 0 {
		if url[pos] == ':' {
			return url[:pos]
		}
	}
	return ""
}

// -----------------------------------------------------------------------------
