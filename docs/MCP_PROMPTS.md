# MCP提示词 (Prompts) 详解

## 什么是MCP提示词？

MCP提示词是Model Context Protocol中的一个重要概念，它允许AI模型获取和使用预定义的提示模板，从而提供更一致、更专业的响应。

## 提示词的核心概念

### 1. 提示词定义
提示词是预定义的文本模板，包含：
- **系统提示**：定义AI的行为和角色
- **用户提示**：常见的用户问题或请求
- **上下文提示**：提供背景信息和指导

### 2. 提示词类型
- **静态提示**：固定的文本内容
- **动态提示**：包含变量的模板
- **条件提示**：根据上下文选择的提示

## MCP提示词协议

### 1. 提示词列表 (prompts/list)
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "prompts/list"
}
```

### 2. 提示词获取 (prompts/get)
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "prompts/get",
  "params": {
    "name": "greeting_prompt"
  }
}
```

### 3. 提示词结构
```json
{
  "name": "greeting_prompt",
  "description": "问候用户的提示词",
  "prompt": "你是一个友好的助手，请用温暖的语言问候用户。",
  "arguments": {
    "type": "object",
    "properties": {
      "user_name": {
        "type": "string",
        "description": "用户姓名"
      },
      "time_of_day": {
        "type": "string",
        "enum": ["morning", "afternoon", "evening"]
      }
    }
  }
}
```

## 实现提示词服务器

### 1. 提示词结构定义
```go
type Prompt struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Prompt      string                 `json:"prompt"`
    Arguments   *InputSchema          `json:"arguments,omitempty"`
}

type ListPromptsResult struct {
    Prompts []Prompt `json:"prompts"`
}

type GetPromptParams struct {
    Name string `json:"name"`
}

type GetPromptResult struct {
    Prompt Prompt `json:"prompt"`
}
```

### 2. 提示词处理函数
```go
func (s *HelloMCPServer) handleListPrompts() *ListPromptsResult {
    return &ListPromptsResult{
        Prompts: []Prompt{
            {
                Name:        "greeting_prompt",
                Description: "生成友好的问候语",
                Prompt:      "你是一个友好的助手，请用温暖的语言问候用户。",
                Arguments: &InputSchema{
                    Type: "object",
                    Properties: map[string]Property{
                        "user_name": {
                            Type:        "string",
                            Description: "用户姓名",
                        },
                        "time_of_day": {
                            Type:        "string",
                            Description: "一天中的时间",
                        },
                    },
                },
            },
            {
                Name:        "help_prompt",
                Description: "提供帮助信息",
                Prompt:      "你是一个有用的助手，请为用户提供清晰、有帮助的指导。",
            },
        },
    }
}

func (s *HelloMCPServer) handleGetPrompt(params *GetPromptParams) (*GetPromptResult, *JSONRPCError) {
    // 根据名称查找提示词
    for _, prompt := range s.prompts {
        if prompt.Name == params.Name {
            return &GetPromptResult{Prompt: prompt}, nil
        }
    }
    
    return nil, &JSONRPCError{
        Code:    -32601,
        Message: fmt.Sprintf("Prompt not found: %s", params.Name),
    }
}
```

## 提示词使用场景

### 1. 角色定义
```json
{
  "name": "expert_consultant",
  "description": "专业咨询顾问角色",
  "prompt": "你是一位经验丰富的专业顾问，具有深厚的行业知识和丰富的实践经验。请以专业、权威的语气回答用户问题，并提供实用的建议和解决方案。"
}
```

### 2. 任务指导
```json
{
  "name": "code_review",
  "description": "代码审查指导",
  "prompt": "请对以下代码进行全面的审查，重点关注：1) 代码质量和可读性 2) 性能优化 3) 安全性问题 4) 最佳实践遵循。请提供具体的改进建议。"
}
```

### 3. 上下文提供
```json
{
  "name": "project_context",
  "description": "项目背景信息",
  "prompt": "这是一个Go语言MCP服务器项目，目标是实现Model Context Protocol。项目使用JSON-RPC 2.0通信，支持工具调用和资源管理。"
}
```

## 提示词最佳实践

### 1. 设计原则
- **明确性**：提示词应该清晰明确
- **一致性**：保持风格和语调一致
- **可重用性**：设计可复用的模板
- **参数化**：使用变量增加灵活性

### 2. 组织方式
- 按功能分类
- 按复杂度分层
- 提供版本控制
- 支持标签和搜索

### 3. 质量控制
- 测试提示词效果
- 收集用户反馈
- 持续优化改进
- 维护提示词库

## 与工具的集成

### 1. 提示词增强工具
```go
func (s *HelloMCPServer) handleCallTool(params *CallToolParams) (*CallToolResult, *JSONRPCError) {
    if params.Name == "say_hello" {
        // 获取问候提示词
        prompt := s.getPrompt("greeting_prompt")
        if prompt != nil {
            // 使用提示词增强响应
            return s.generateEnhancedResponse(prompt, params)
        }
    }
    // ... 其他处理逻辑
}
```

### 2. 动态提示词生成
```go
func (s *HelloMCPServer) generateDynamicPrompt(template string, args map[string]interface{}) string {
    // 使用模板引擎生成动态提示词
    prompt := template
    for key, value := range args {
        placeholder := fmt.Sprintf("{{%s}}", key)
        prompt = strings.ReplaceAll(prompt, placeholder, fmt.Sprintf("%v", value))
    }
    return prompt
}
```

## 扩展功能

### 1. 提示词管理
- 提示词版本控制
- 提示词分类和标签
- 提示词使用统计
- 提示词效果评估

### 2. 高级特性
- 条件提示词
- 提示词链式调用
- 提示词模板继承
- 多语言支持

### 3. 集成选项
- 与LLM框架集成
- 提示词市场
- 社区贡献
- 自动化优化 