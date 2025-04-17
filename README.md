MCP Go+ 🚀
=====

[![Build Status](https://github.com/goplus/mcp/actions/workflows/go.yml/badge.svg)](https://github.com/goplus/mcp/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goplus/mcp)](https://goreportcard.com/report/github.com/goplus/mcp)
[![GitHub release](https://img.shields.io/github/v/tag/goplus/mcp.svg?label=release)](https://github.com/goplus/mcp/releases)
[![GoDoc](https://pkg.go.dev/badge/github.com/goplus/mcp.svg)](https://pkg.go.dev/github.com/goplus/mcp)
[![Language](https://img.shields.io/badge/language-Go+-blue.svg)](https://github.com/goplus/gop)
<!--
[![Coverage Status](https://codecov.io/gh/goplus/mcp/branch/main/graph/badge.svg)](https://codecov.io/gh/goplus/mcp)
-->

A Go+ implementation of the Model Context Protocol (MCP), enabling seamless integration between LLM applications and external data sources and tools.

This repo contains two [Go+ classfiles](https://github.com/goplus/gop/blob/main/doc/classfile.md). They are [mcp](#mcp-mcp-server-framework) (MCP Server Framework) and [mcptest](#mcptest-mcp-server-test-framework) (MCP Server Test Framework).

The classfile [mcp](#mcp-mcp-server-framework) has the file suffix `_mcp.gox` (the MCP Server), `_res.gox` (a MCP Resource or ResourceTemplate), `_tool.gox` (a MCP Tool) and `_prompt.gox` (a MCP Prompt). The classfile [mcptest](#mcptest-mcp-server-test-framework) has the file suffix `_mtest.gox`.

## mcp: MCP Server Framework

Here is a MCP Server example ([source code](demo/hello)). It has two files: `main_mcp.gox` (the MCP Server) and `hello_tool.gox` (a MCP Tool).

### A MCP Server Example: hello

First let us initialize a hello project:

```
gop mod init hello
```

Then we have it reference the [mcp](https://pkg.go.dev/github.com/goplus/mcp) classfile:

```
gop get github.com/goplus/mcp@latest
```

The content of `main_mcp.gox` (the MCP Server) is as follows

```go
server "Tool Demo 🚀", "1.0.0"
```

The content of `hello_tool.gox` (a MCP Tool) is as follows

```go
tool "helloWorld", => {
	description "Say hello to someone"
	string "name", => {
		required
		description "Name of the person to greet"
	}
}

name, ok := ${name}.(string)
if !ok {
	return newError("name must be a string")
}

return text("Hello, ${name}!")
```

Execute the following commands:

```
gop mod tidy
gop run .
```

A simplest MCP Server is running now.

## mcptest: MCP Server Test Framework

To test the above [hello MCP Server](#a-mcp-server-example-hello) ([source code](demo/hello)), you only need to implement a `hello_mtest.gox` file:

```go
mock new(MCPApp)

initialize nil
ret {}

call "helloWorld", {"name": "Ken"}
ret {
	"content": [{"type": "text", "text": "Hello, Ken!"}],
}
```

Then run `gop test` or `gop test -v` to execute the unit test.
