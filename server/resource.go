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
	"errors"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yosida95/uritemplate/v3"
)

// -----------------------------------------------------------------------------

type ResourceAppProto struct {
	resource mcp.Resource
	template mcp.ResourceTemplate
}

type ResourceApp struct {
	*ResourceAppProto
	ctx      context.Context
	request  mcp.ReadResourceRequest
	values   uritemplate.Values
	isClone  bool
	hasTempl bool
}

// RequestURI returns the URI of the request.
func (p *ResourceApp) RequestURI() string {
	return p.request.Params.URI
}

// Gop_Env returns the value of the specified parameter.
func (p *ResourceApp) Gop_Env(name string) any {
	if p.hasTempl {
		if p.values == nil {
			p.values = p.template.URITemplate.Match(p.request.Params.URI)
		}
		if v, ok := p.values[name]; ok {
			if v.T == uritemplate.ValueTypeString {
				return v.V[0]
			} else {
				return v.V
			}
		}
	}
	if v, ok := p.request.Params.Arguments[name]; ok {
		return v
	}
	return nil
}

// Main is required by Go+ compiler as the entry of a MCPServer resource.
func (p *ResourceApp) Main(ctx context.Context, request mcp.ReadResourceRequest, t *ResourceAppProto) []mcp.ResourceContents {
	if t == nil {
		p.ctx = ctx
		p.request = request
		p.isClone = true
	} else {
		p.ResourceAppProto = t
	}
	return nil
}

func initResourceApp(self ResourceProto) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(stop); !ok {
				panic(e)
			}
		}
	}()
	self.Main(context.TODO(), mcp.ReadResourceRequest{}, &ResourceAppProto{})
}

func hasTemplate(uri string) bool {
	return strings.IndexByte(uri, '{') >= 0
}

// Resource calls fn to initialize the resource.
func (p *ResourceApp) Resource(uri string, name string, fn func()) {
	if !p.isClone {
		p.hasTempl = hasTemplate(uri)
		if p.hasTempl {
			p.template = mcp.NewResourceTemplate(uri, name)
		} else {
			p.resource = mcp.NewResource(uri, name)
		}
		fn()
		panic(stop{})
	}
}

// Description adds a description to the Resource.
// The description should provide a clear, human-readable explanation of what the resource represents.
func (p *ResourceApp) Description(description string) {
	if p.hasTempl {
		mcp.WithTemplateDescription(description)(&p.template)
	} else {
		mcp.WithResourceDescription(description)(&p.resource)
	}
}

// MimeType sets the MIME type for the Resource.
// This should indicate the format of the resource's contents.
func (p *ResourceApp) MimeType(mime string) {
	if p.hasTempl {
		mcp.WithTemplateMIMEType(mime)(&p.template)
	} else {
		mcp.WithMIMEType(mime)(&p.resource)
	}
}

func (p *ResourceApp) addTo(self ResourceProto, svr *server.MCPServer) {
	clone := self.Classclone
	handle := func(ctx context.Context, request mcp.ReadResourceRequest) (ret []mcp.ResourceContents, err error) {
		defer func() {
			if e := recover(); e != nil {
				switch e := e.(type) {
				case string:
					err = errors.New(e)
				case error:
					err = e
				default:
					err = errors.New("unknown error")
				}
			}
		}()
		ret = clone().Main(ctx, request, nil)
		return
	}
	initResourceApp(self)
	if p.hasTempl {
		svr.AddResourceTemplate(p.template, handle)
	} else {
		svr.AddResource(p.resource, handle)
	}
}

// -----------------------------------------------------------------------------
