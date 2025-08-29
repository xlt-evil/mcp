# Hello MCP Server - MCP协议学习项目

这是一个用于学习Model Context Protocol (MCP)的Go语言示例项目。MCP是一个用于AI模型与外部工具和服务进行安全通信的协议。

## 项目简介

这个项目实现了一个完整的MCP服务器，提供以下功能：
- 实现MCP协议的基本消息处理
- 提供工具调用服务 (Tools)
- 支持提示词管理 (Prompts)
- 实现资源管理 (Resources)
- 记录问候日志到文件
- 支持JSON-RPC 2.0通信

## 功能特性

### 工具列表
- **say_hello**: 向指定的人发送问候消息
  - 参数: `person_name` (必需), `greeting_message` (可选)
  - 功能: 记录问候日志并返回友好的回应

### 协议支持
- ✅ JSON-RPC 2.0
- ✅ MCP 2024-11-05 协议版本
- ✅ 工具列表和调用 (Tools)
- ✅ 提示词管理 (Prompts)
- ✅ 资源管理 (Resources)
- ✅ 错误处理

## 快速开始

### 环境要求
- Go 1.19+
- 支持MCP的客户端（如Claude Desktop）

### 编译运行
```bash
# 使用构建脚本（推荐）
build.bat

# 手动构建
go build -o sayhi-server cmd/sayhi_server/main.go
go build -o database-server cmd/database_server/main.go

# 运行Hello服务器
./sayhi-server

# 运行数据库服务器
./database-server
```

### 客户端配置
在支持MCP的客户端中添加服务器配置：

**Hello MCP服务器配置：**
```json
{
  "mcpServers": {
    "hello-server": {
      "command": "path/to/sayhi-server",
      "args": [],
      "env": {},
      "description": "Hello MCP学习服务器 - 提供问候工具"
    }
  }
}
```

**数据库MCP服务器配置：**
```json
{
  "mcpServers": {
    "database-server": {
      "command": "path/to/database-server",
      "args": ["--config", "config/database.yaml"],
      "env": {},
      "description": "数据库MCP服务器 - 提供数据库查询工具"
    }
  }
}
```

**同时配置两个服务器：**
```json
{
  "mcpServers": {
    "hello-server": {
      "command": "path/to/sayhi-server",
      "args": [],
      "env": {},
      "description": "Hello MCP学习服务器 - 提供问候工具"
    },
    "database-server": {
      "command": "path/to/database-server",
      "args": ["--config", "config/database.yaml"],
      "env": {},
      "description": "数据库MCP服务器 - 提供数据库查询工具"
    }
  }
}
```

> 📝 **注意**: 请将 `path/to/` 替换为实际的可执行文件路径。详细配置说明请参考 `mcp-config-examples.md` 文件。

## MCP协议学习要点

### 1. 消息流程
1. **initialize**: 客户端初始化连接
2. **initialized**: 初始化完成通知
3. **tools/list**: 获取可用工具列表
4. **tools/call**: 调用指定工具

### 2. 消息结构
- 所有消息都遵循JSON-RPC 2.0规范
- 包含method、params、result、error等字段
- 支持请求-响应模式

### 3. 错误处理
- 使用标准错误码
- 提供详细的错误信息
- 支持错误数据传递

## 项目结构

```
.
├── cmd/                     # 可执行程序
│   ├── sayhi_server/        # Hello MCP服务器
│   │   ├── main.go         # 主程序
│   │   └── README.md       # 说明文档
│   └── database_server/     # 数据库MCP服务器
│       ├── main.go         # 主程序
│       └── README.md       # 说明文档
├── config/                  # 配置管理
│   ├── database.go         # 数据库配置结构
│   └── database.yaml       # 数据库配置文件
├── database/                # 数据库管理
│   └── manager.go          # 数据库管理器
├── types/                   # 共享类型
│   └── mcp_types.go        # MCP类型定义
├── examples/                # 使用示例
├── docs/                    # 详细文档
│   ├── MCP_PROTOCOL.md     # MCP协议详解
│   ├── MCP_PROMPTS.md      # 提示词管理详解
│   ├── MCP_RESOURCES.md    # 资源管理详解
│   ├── MCP_ARCHITECTURE.md # 架构设计详解
│   └── LEARNING_GUIDE.md   # 学习指南
├── README.md                # 项目说明文档
├── mcp-config.json          # MCP配置文件示例
├── go.mod                   # Go模块文件
├── build.bat                # 构建脚本
└── .gitignore               # Git忽略文件
```

## 学习资源

- [MCP官方文档](https://modelcontextprotocol.io/)
- [MCP GitHub仓库](https://github.com/modelcontextprotocol)
- [JSON-RPC 2.0规范](https://www.jsonrpc.org/specification)

## 贡献

欢迎提交Issue和Pull Request来改进这个学习项目！

## 许可证

MIT License 