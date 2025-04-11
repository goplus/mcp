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

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Role represents the sender or recipient of messages and data in a
// conversation.
type Role = mcp.Role

const (
	RoleUser      = mcp.RoleUser
	RoleAssistant = mcp.RoleAssistant
)

// -----------------------------------------------------------------------------

type PromptAppProto struct {
	prompt mcp.Prompt
	opts   []mcp.ArgumentOption
	parent *server.MCPServer
}

type PromptApp struct {
	*PromptAppProto
	withContext
	request mcp.GetPromptRequest
	isClone bool
}

// Gop_Env returns the value of the specified parameter.
func (p *PromptApp) Gop_Env(name string) string {
	return p.request.Params.Arguments[name]
}

// Main is required by Go+ compiler as the entry of a MCPServer prompt.
func (p *PromptApp) Main(ctx context.Context, request mcp.GetPromptRequest, t *PromptAppProto) (string, []mcp.PromptMessage) {
	if t == nil {
		p.ctx = ctx
		p.request = request
		p.isClone = true
	} else {
		p.PromptAppProto = t
		p.svr = t.parent
	}
	return "", nil
}

func initPromptApp(self PromptProto, svr *server.MCPServer) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(stop); !ok {
				panic(e)
			}
		}
	}()
	self.Main(context.TODO(), mcp.GetPromptRequest{}, &PromptAppProto{
		prompt: mcp.NewPrompt(""),
		parent: svr,
	})
}

// Prompt calls fn to initialize the prompt.
func (p *PromptApp) Prompt__0(name string, fn func()) {
	if !p.isClone {
		p.prompt.Name = name
		fn()
		panic(stop{})
	}
}

// Prompt creates a new PromptMessage
func (p *PromptApp) Prompt__1(role Role, content mcp.Content) mcp.PromptMessage {
	return mcp.NewPromptMessage(role, content)
}

// For a prompt:
// Description sets a description to the Prompt.
// The description should provide a clear, human-readable explanation of what the prompt does.
//
// For a prompt property:
// Description adds a description to a prompt argument.
// The description should explain the purpose and expected values of the argument.
func (p *PromptApp) Description(description string) {
	if p.opts != nil {
		p.opts = append(p.opts, mcp.ArgumentDescription(description))
	} else {
		mcp.WithPromptDescription(description)(&p.prompt)
	}
}

// Arg adds an argument to the prompt's argument list.
// The argument will be configured based on the provided options.
func (p *PromptApp) Arg(name string, fn ...func()) {
	if len(fn) > 0 {
		p.opts = make([]mcp.ArgumentOption, 0, 2)
		fn[0]()
	}
	mcp.WithArgument(name, p.opts...)(&p.prompt)
	p.opts = nil
}

// Required marks an argument as required in the prompt.
// Required arguments must be provided when getting the prompt.
func (p *PromptApp) Required() {
	p.opts = append(p.opts, mcp.RequiredArgument())
}

func (p *PromptApp) addTo(self PromptProto, svr *server.MCPServer) {
	clone := self.Classclone
	initPromptApp(self, svr)
	svr.AddPrompt(p.prompt, func(ctx context.Context, request mcp.GetPromptRequest) (ret *mcp.GetPromptResult, err error) {
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
		description, messages := clone().Main(ctx, request, nil)
		ret = mcp.NewGetPromptResult(description, messages)
		return
	})
}

// -----------------------------------------------------------------------------
