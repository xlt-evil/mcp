# MCP学习指南

## 🎯 学习目标

通过这个项目，你将学会：
1. 理解MCP协议的基本概念和工作原理
2. 实现一个完整的MCP服务器
3. 掌握JSON-RPC 2.0的使用
4. 学习Go语言的最佳实践

## 📚 学习路径

### 第一阶段：理解MCP协议
1. **阅读文档**
   - 查看 `docs/MCP_PROTOCOL.md` 了解协议细节
   - 阅读 [MCP官方文档](https://modelcontextprotocol.io/)

2. **理解核心概念**
   - 服务器 vs 客户端
   - 工具定义和调用
   - 消息流程和状态管理

### 第二阶段：代码分析
1. **主程序结构** (`sayhi.go`)
   - 查看 `HelloMCPServer` 结构体
   - 理解消息处理流程
   - 分析工具实现

2. **协议实现**
   - JSON-RPC消息结构
   - 错误处理机制
   - 日志记录功能

### 第三阶段：实践操作
1. **构建和运行**
   ```bash
   # Windows
   build.bat
   
   # 手动构建
   go build -o mcp-server.exe sayhi.go
   go build -o test-client.exe examples/test_client.go
   ```

2. **测试功能**
   - 运行测试客户端
   - 观察消息交互
   - 查看日志文件

## 🔍 代码重点解析

### 1. 消息处理流程
```go
func (s *HelloMCPServer) processMessage(msg *JSONRPCMessage) *JSONRPCMessage {
    // 根据method分发到不同的处理函数
    switch msg.Method {
    case "initialize":
        // 处理初始化
    case "tools/list":
        // 返回工具列表
    case "tools/call":
        // 调用指定工具
    }
}
```

### 2. 工具定义
```go
type Tool struct {
    Name        string      `json:"name"`
    Description string      `json:"description"`
    InputSchema InputSchema `json:"inputSchema"`
}
```

### 3. 错误处理
```go
type JSONRPCError struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

## 🧪 动手实验

### 实验1：添加新工具
尝试在 `handleListTools()` 函数中添加一个新的工具：

```go
{
    Name:        "get_time",
    Description: "获取当前时间",
    InputSchema: InputSchema{
        Type: "object",
        Properties: map[string]Property{
            "format": {
                Type:        "string",
                Description: "时间格式 (可选)",
            },
        },
    },
},
```

然后在 `handleCallTool()` 中实现对应的处理逻辑。

### 实验2：修改问候逻辑
修改 `handleCallTool()` 函数，添加更多的问候选项或随机化逻辑。

### 实验3：添加配置支持
实现从配置文件读取服务器配置，而不是硬编码。

## 📝 调试技巧

### 1. 启用详细日志
```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### 2. 使用测试客户端
运行 `test-client.exe` 来测试服务器功能。

### 3. 查看日志文件
检查 `hello_log.txt` 文件了解工具调用记录。

## 🚀 进阶学习

### 1. 扩展功能
- 添加数据库支持
- 实现用户认证
- 支持WebSocket连接
- 添加插件系统

### 2. 性能优化
- 实现连接池
- 添加缓存机制
- 异步处理长时间任务
- 监控和指标收集

### 3. 安全性
- 输入验证和清理
- 访问控制
- 加密通信
- 审计日志

## 🔗 相关资源

- [Go语言官方文档](https://golang.org/doc/)
- [JSON-RPC 2.0规范](https://www.jsonrpc.org/specification)
- [MCP GitHub示例](https://github.com/modelcontextprotocol)
- [Go并发编程](https://golang.org/doc/effective_go.html#concurrency)

## 💡 常见问题

### Q: 为什么服务器启动后立即退出？
A: 检查是否有语法错误，确保Go环境正确安装。

### Q: 测试客户端无法连接？
A: 确保服务器正在运行，检查可执行文件路径。

### Q: 如何添加更多工具？
A: 在 `handleListTools()` 中添加工具定义，在 `handleCallTool()` 中实现处理逻辑。

## 🎉 学习成果

完成这个项目后，你将能够：
- 理解MCP协议的设计理念
- 实现自己的MCP服务器
- 掌握Go语言网络编程
- 为AI应用开发工具扩展

继续探索MCP生态系统的其他组件，如资源管理、文件系统访问等！ 