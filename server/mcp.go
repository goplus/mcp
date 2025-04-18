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

package server

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"strconv"

	"github.com/goplus/mcp/server/svx"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	_ "github.com/goplus/mcp/server/stdio"
)

const (
	GopPackage = true
)

// -----------------------------------------------------------------------------

type multiContents struct {
	data  []mcp.Content
	isErr bool
	mcp.Content
}

// JsonContent represents a JSON content.
type JsonContent struct {
	// JSON represents data to be serialized as JSON.
	JSON any
}

// Text creates a new TextContent.
func Text__0(text string) mcp.Content {
	return mcp.NewTextContent(text)
}

// Text creates a new TextContent with a JSON content.
func Text__1(v JsonContent) mcp.Content {
	b, err := json.Marshal(v.JSON)
	if err != nil {
		panic(err)
	}
	return mcp.NewTextContent(string(b))
}

// Number creates a new TextContent with a number content.
func Number__0(val float64) mcp.Content {
	return mcp.NewTextContent(strconv.FormatFloat(val, 'f', -1, 64))
}

// Number creates a new TextContent with a number content.
func Number__1(val float64, prec int) mcp.Content {
	return mcp.NewTextContent(strconv.FormatFloat(val, 'f', prec, 64))
}

// Number creates a new TextContent with a number content.
func Number__2(val float64, fmt byte, prec int) mcp.Content {
	return mcp.NewTextContent(strconv.FormatFloat(val, fmt, prec, 64))
}

// Image creates a new ImageContent with an image.
func Image__0(mimeType, imageData string) mcp.Content {
	return &mcp.ImageContent{
		Type:     "image",
		Data:     imageData,
		MIMEType: mimeType,
	}
}

// Image creates a new Content with both text and image content.
func Image__1(text, mimeType, imageData string) mcp.Content {
	return Multiple(
		Text__0(text),
		Image__0(mimeType, imageData),
	)
}

// Embedded creates a new EmbeddedResource.
func Embedded__0(text *mcp.TextResourceContents) mcp.Content {
	return mcp.NewEmbeddedResource(text)
}

// Embedded creates a new EmbeddedResource.
func Embedded__1(blob *mcp.BlobResourceContents) mcp.Content {
	return mcp.NewEmbeddedResource(blob)
}

// Embedded creates a new EmbeddedResource.
func Embedded__2(v *JsonResourceContents) mcp.Content {
	return mcp.NewEmbeddedResource(Content__2(v))
}

// Embedded creates a new EmbeddedResource.
func Embedded__3(text *TextResourceByteContents) mcp.Content {
	return mcp.NewEmbeddedResource(Content__3(text))
}

// Multiple creates a new Content with multiple contents.
func Multiple(contents ...mcp.Content) mcp.Content {
	return &multiContents{
		data: contents,
	}
}

// NewError creates a new content with an error message.
// Any errors that originate from the tool SHOULD be reported inside the result object.
func NewError__0(text string) mcp.Content {
	return &multiContents{
		data: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: text,
			},
		},
		isErr: true,
	}
}

// NewError creates a new content with an error message.
// If an error is provided, its details will be appended to the text message.
// Any errors that originate from the tool SHOULD be reported inside the result object.
func NewError__1(text string, err error) mcp.Content {
	if err != nil {
		text = text + ": " + err.Error()
	}
	return NewError__0(text)
}

// -----------------------------------------------------------------------------

type JsonResourceContents struct {
	// The URI of this resource.
	URI string
	// JSON represents data to be serialized as JSON.
	JSON any
}

// TextResourceByteContents represents a text resource with byte contents.
type TextResourceByteContents struct {
	// The URI of this resource.
	URI string
	// The MIME type of this resource, if known.
	MIMEType string
	// The text of the item. This must only be set if the item can actually be
	// represented as text (not binary data).
	Text []byte
}

// Content returns a TextResourceContents.
func Content__0(text *mcp.TextResourceContents) *mcp.TextResourceContents {
	return text
}

// Content returns a BlobResourceContents.
func Content__1(blob *mcp.BlobResourceContents) *mcp.BlobResourceContents {
	return blob
}

// Content returns a TextResourceContents.
func Content__2(v *JsonResourceContents) *mcp.TextResourceContents {
	b, err := json.Marshal(v.JSON)
	if err != nil {
		panic(err)
	}
	return &mcp.TextResourceContents{
		URI:      v.URI,
		MIMEType: "application/json",
		Text:     string(b),
	}
}

// Content returns a TextResourceContents.
func Content__3(text *TextResourceByteContents) *mcp.TextResourceContents {
	return &mcp.TextResourceContents{
		URI:      text.URI,
		MIMEType: text.MIMEType,
		Text:     string(text.Text),
	}
}

// -----------------------------------------------------------------------------

type withContext struct {
	ctx context.Context
	svr *server.MCPServer
}

// Notify sends a notification to the current client
func (p *withContext) Notify(method string, params map[string]any) {
	p.svr.SendNotificationToClient(p.ctx, method, params)
}

// -----------------------------------------------------------------------------

// MCPApp is the project class of a MCPServer classfile.
type MCPApp struct {
	svr  *server.MCPServer
	las  func(addr string, svr *server.MCPServer) error
	addr string
}

// Sys returns the underlying MCPServer instance.
// Don't use this method except for testing purposes.
func (p *MCPApp) Sys() *server.MCPServer {
	return p.svr
}

// Server creates a new MCP server instance with the given name and version.
func (p *MCPApp) Server(name, version string) {
	p.svr = server.NewMCPServer(name, version)
	p.addr = "stdio:"
	if p.las == nil {
		p.las = svx.ListenAndServe
	}
}

// Run sets the MCP server address.
func (p *MCPApp) Run(addr string) {
	p.addr = addr
}

// SetLAS sets the ListenAndServe function for the MCP server.
func (p *MCPApp) SetLAS(las func(addr string, svr *server.MCPServer) error) {
	p.las = las
}

func (p *MCPApp) serve() error {
	return p.las(p.addr, p.svr)
}

// -----------------------------------------------------------------------------

var _ = (*ResourceApp).addTo
var _ = (*ToolApp).addTo
var _ = (*PromptApp).addTo
var _ = (*MCPApp).serve

type stop struct{}

type iAppProto interface {
	serve() error
	Sys() *server.MCPServer
	MainEntry()
}

type ResourceProto interface {
	addTo(self ResourceProto, svr *server.MCPServer)
	Main(ctx context.Context, request mcp.ReadResourceRequest, t *ResourceAppProto) []mcp.ResourceContents
	Classclone() ResourceProto
}

type ToolProto interface {
	addTo(self ToolProto, svr *server.MCPServer)
	Main(ctx context.Context, request mcp.CallToolRequest, t *ToolAppProto) mcp.Content
	Classclone() ToolProto
}

type PromptProto interface {
	addTo(self PromptProto, svr *server.MCPServer)
	Main(ctx context.Context, request mcp.GetPromptRequest, t *PromptAppProto) (string, []mcp.PromptMessage)
	Classclone() PromptProto
}

// Gopt_MCPApp_Main is required by Go+ compiler as the entry of a MCPServer project.
func Gopt_MCPApp_Main(app iAppProto, resources []ResourceProto, tools []ToolProto, prompts []PromptProto) {
	app.MainEntry()
	svr := app.Sys()
	for _, r := range resources {
		initProj(r, app)
		r.addTo(r, svr)
	}
	for _, h := range tools {
		initProj(h, app)
		h.addTo(h, svr)
	}
	for _, p := range prompts {
		initProj(p, app)
		p.addTo(p, svr)
	}
	err := app.serve()
	if err != nil {
		log.Panicln(err)
	}
}

func initProj(work, app any) {
	reflect.ValueOf(work).Elem().Field(1).Set(reflect.ValueOf(app))
}

// -----------------------------------------------------------------------------
