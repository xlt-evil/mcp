# 数据库MCP服务器

这是一个支持数据库查询的MCP服务器，可以连接到MySQL数据库并执行SQL查询操作。

## 功能特性

### 🔍 数据库查询工具
- **database_query**: 执行SQL查询并返回结果
- **database_tables**: 获取数据库中的所有表名
- **database_schema**: 获取指定表的结构信息
- **database_status**: 检查数据库连接状态

### 🗄️ 支持的数据库
- MySQL (主要支持)
- PostgreSQL (计划支持)

### 🛡️ 安全特性
- 连接池管理
- 自动重连机制
- 查询结果限制
- 错误处理和日志记录

## 快速开始

### 1. 环境要求
- Go 1.19+
- MySQL 5.7+ 或 8.0+
- 网络访问权限到数据库服务器

### 2. 配置数据库连接
编辑 `config/database.yaml` 文件：

```yaml
database:
  enabled: true
  driver: "mysql"
  host: "host"
  port: 3306
  user: "="
  password: ""
  name: ""
```

### 3. 构建和运行
```bash
# 构建
go build -o database-mcp-server cmd/database_server/main.go

# 运行
./database-mcp-server
```

## 工具使用说明

### 1. database_query
执行SQL查询并返回结果。

**参数:**
- `sql` (必需): SQL查询语句

**示例:**
```json
{
  "name": "database_query",
  "arguments": {
    "sql": "SELECT * FROM users LIMIT 10"
  }
}
```

**返回结果:**
- 查询列信息
- 数据行（限制显示前10行）
- 总行数统计

### 2. database_tables
获取数据库中的所有表名。

**参数:** 无

**示例:**
```json
{
  "name": "database_tables",
  "arguments": {}
}
```

**返回结果:**
- 数据库名称
- 表总数
- 表名列表

### 3. database_schema
获取指定表的结构信息。

**参数:**
- `table_name` (必需): 要查看结构的表名

**示例:**
```json
{
  "name": "database_schema",
  "arguments": {
    "table_name": "users"
  }
}
```

**返回结果:**
- 表名
- 字段数量
- 每个字段的详细信息（名称、类型、可空性、键类型、默认值、额外信息）

### 4. database_status
检查数据库连接状态。

**参数:** 无

**示例:**
```json
{
  "name": "database_status",
  "arguments": {}
}
```

**返回结果:**
- 数据库配置信息
- 连接状态
- 重连尝试结果

## 配置选项

### 环境变量
可以通过环境变量覆盖配置文件：

```bash
export DB_HOST="localhost"
export DB_PORT="3306"
export DB_USER="myuser"
export DB_PASSWORD="mypassword"
export DB_NAME="mydatabase"
```

### 连接池配置
```yaml
pool:
  max_open_conns: 25      # 最大打开连接数
  max_idle_conns: 5       # 最大空闲连接数
  conn_max_lifetime: "5m" # 连接最大生命周期
```

### 安全配置
```yaml
security:
  ssl_mode: "preferred"    # SSL模式
  timeout: "30s"           # 连接超时
  read_timeout: "10s"      # 读取超时
  write_timeout: "10s"     # 写入超时
```

## 错误处理

### 常见错误码
- `-32601`: 未知工具
- `-32602`: 无效参数
- `-32603`: 内部错误（数据库连接失败、查询执行失败等）
- `-32700`: JSON解析错误

### 错误恢复
- 自动重连机制
- 连接池管理
- 详细的错误日志

## 性能优化

### 查询限制
- 结果行数限制（默认显示前10行）
- 连接池复用
- 超时控制

### 监控指标
- 连接状态
- 查询执行时间
- 错误率统计

## 安全考虑

### 访问控制
- 数据库用户权限限制
- 网络访问控制
- 查询结果过滤

### 数据保护
- 敏感信息不记录到日志
- 查询结果大小限制
- 连接超时保护

## 故障排除

### 连接问题
1. 检查网络连通性
2. 验证数据库凭据
3. 确认防火墙设置
4. 检查数据库服务状态

### 查询问题
1. 验证SQL语法
2. 检查表权限
3. 确认数据库存在
4. 查看错误日志

## 扩展功能

### 计划中的功能
- 支持更多数据库类型
- 查询缓存机制
- 批量操作支持
- 事务管理
- 存储过程调用

### 自定义扩展
- 插件系统
- 自定义查询模板
- 结果格式化器
- 查询优化建议

## 贡献指南

欢迎提交Issue和Pull Request来改进这个项目！

### 开发环境设置
1. Fork项目
2. 创建功能分支
3. 提交更改
4. 创建Pull Request

### 代码规范
- 遵循Go语言规范
- 添加适当的注释
- 编写测试用例
- 更新相关文档 