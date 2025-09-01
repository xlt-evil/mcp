# MCP配置示例

## Claude Desktop 配置

在Claude Desktop中，你可以在设置中添加MCP服务器配置：

### 1. Hello MCP服务器配置
```json
{
  "mcpServers": {
    "hello-server": {
      "command": "C:\\path\\to\\your\\project\\sayhi-server.exe",
      "args": [],
      "env": {},
      "description": "Hello MCP学习服务器 - 提供问候工具"
    }
  }
}
```

### 2. 数据库MCP服务器配置
```json
{
  "mcpServers": {
    "database-server": {
      "command": "C:\\path\\to\\your\\project\\database-server.exe",
      "args": ["--config", "config/database.yaml"],
      "env": {},
      "description": "数据库MCP服务器 - 提供数据库查询工具"
    }
  }
}
```

### 3. Redis MCP服务器配置
```json
{
  "mcpServers": {
    "redis-server": {
      "command": "C:\\path\\to\\your\\project\\redis-server.exe",
      "args": ["--config", "config/redis.yaml"],
      "env": {},
      "description": "Redis MCP服务器 - 提供Redis数据库操作工具"
    }
  }
}
```

### 4. 同时配置三个服务器
```json
{
  "mcpServers": {
    "hello-server": {
      "command": "C:\\path\\to\\your\\project\\sayhi-server.exe",
      "args": [],
      "env": {},
      "description": "Hello MCP学习服务器 - 提供问候工具"
    },
    "database-server": {
      "command": "C:\\path\\to\\your\\project\\database-server.exe",
      "args": ["--config", "config/database.yaml"],
      "env": {},
      "description": "数据库MCP服务器 - 提供数据库查询工具"
    },
    "redis-server": {
      "command": "C:\\path\\to\\your\\project\\redis-server.exe",
      "args": ["--config", "config/redis.yaml"],
      "env": {},
      "description": "Redis MCP服务器 - 提供Redis数据库操作工具"
    }
  }
}
```

## 其他MCP客户端配置

### 1. 通用配置格式
```json
{
  "mcpServers": {
    "server-name": {
      "command": "path/to/executable",
      "args": ["arg1", "arg2"],
      "env": {
        "ENV_VAR": "value"
      },
      "description": "服务器描述"
    }
  }
}
```

### 2. 环境变量配置
```json
{
  "mcpServers": {
    "database-server": {
      "command": "./database-server",
      "args": [],
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "3306",
        "DB_USER": "root",
        "DB_PASSWORD": "password",
        "DB_NAME": "database"
      },
      "description": "数据库MCP服务器"
    }
  }
}
```

## 配置说明

### 参数说明
- **command**: 服务器可执行文件的路径
- **args**: 命令行参数数组
- **env**: 环境变量键值对
- **description**: 服务器描述信息

### 路径说明
- **相对路径**: `./sayhi-server` (相对于工作目录)
- **绝对路径**: `C:\\path\\to\\sayhi-server.exe` (Windows)
- **绝对路径**: `/path/to/sayhi-server` (Linux/Mac)

### 环境变量
- 可以通过环境变量覆盖配置文件中的设置
- 支持数据库连接参数的环境变量配置
- 可以设置日志级别、调试模式等

## 故障排除

### 常见问题
1. **路径错误**: 确保可执行文件路径正确
2. **权限问题**: 确保有执行权限
3. **依赖缺失**: 确保所有依赖都已安装
4. **端口冲突**: 检查是否有端口冲突

### 调试技巧
1. 在命令行中手动运行服务器
2. 检查服务器日志输出
3. 验证配置文件格式
4. 测试网络连接

## 最佳实践

### 1. 路径管理
- 使用绝对路径避免路径问题
- 将可执行文件放在固定目录
- 使用环境变量管理路径

### 2. 配置管理
- 为不同环境创建不同配置
- 使用版本控制管理配置
- 定期备份配置文件

### 3. 安全考虑
- 不要在配置文件中硬编码密码
- 使用环境变量管理敏感信息
- 限制配置文件的访问权限 