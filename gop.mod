xgo 1.5

project *_mcp.gox MCPApp github.com/goplus/mcp/server

class *_res.gox ResourceApp ResourceProto
class *_tool.gox ToolApp ToolProto
class *_prompt.gox PromptApp PromptProto

import log

project main_mtest.gox MainApp github.com/goplus/mcp/mtest github.com/qiniu/x/test

class *_mtest.gox CaseApp
