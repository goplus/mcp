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

package stdio

import (
	"log"

	"github.com/goplus/mcp/server/svx"
	"github.com/mark3labs/mcp-go/server"
)

const (
	Scheme = "stdio"
)

func init() {
	svx.Register(Scheme, ListenAndServe)
}

// -----------------------------------------------------------------------------

// addr = "stdio:"
func ListenAndServe(addr string, svr *server.MCPServer) error {
	log.Println("Serving MCP server with stdio ...")
	return server.ServeStdio(svr)
}

// -----------------------------------------------------------------------------
