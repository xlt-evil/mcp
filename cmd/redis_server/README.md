# Redis MCP Server

这是一个基于 MCP (Model Context Protocol) 的 Redis 服务器，提供 Redis 数据库操作的工具接口。

## 功能特性

### 工具列表

* **redis_get**: 获取Redis键的值
  * 参数: `key` (必需) - 要获取的键名
  * 功能: 返回指定键的值

* **redis_set**: 设置Redis键值对
  * 参数: `key` (必需), `value` (必需), `expiration` (可选)
  * 功能: 设置键值对，支持过期时间

* **redis_del**: 删除Redis键
  * 参数: `keys` (必需) - 要删除的键名列表
  * 功能: 删除指定的键

* **redis_keys**: 获取匹配模式的键列表
  * 参数: `pattern` (必需) - 键模式（如：user:*）
  * 功能: 返回匹配模式的键列表

* **redis_type**: 获取键的数据类型
  * 参数: `key` (必需) - 键名
  * 功能: 返回键的数据类型（string, hash, list, set, zset等）

* **redis_ttl**: 获取键的TTL（生存时间）
  * 参数: `key` (必需) - 键名
  * 功能: 返回键的剩余生存时间

* **redis_info**: 获取Redis服务器信息
  * 参数: `section` (可选) - 信息部分
  * 功能: 返回Redis服务器详细信息

* **redis_dbsize**: 获取当前数据库的键数量
  * 参数: 无
  * 功能: 返回当前数据库中的键数量

* **redis_flushdb**: 清空当前数据库
  * 参数: 无
  * 功能: 清空当前数据库中的所有键

* **redis_execute**: 执行自定义Redis命令
  * 参数: `command` (必需), `args` (可选)
  * 功能: 执行任意Redis命令

* **redis_status**: 检查Redis连接状态
  * 参数: 无
  * 功能: 检查并显示Redis连接状态

## 配置

服务器使用 `config/redis.yaml` 配置文件，支持以下配置项：

```yaml
redis:
  enabled: true
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  
  # 连接池配置
  pool:
    max_idle: 10
    max_active: 100
    idle_timeout: "5m"
  
  # 超时配置
  timeout:
    connect: "5s"
    read: "3s"
    write: "3s"
  
  # 日志配置
  logging:
    enabled: true
    level: "info"
    file: "redis.log"
```

## 编译运行

```bash
# 编译
go build -o redis-server cmd/redis_server/main.go

# 运行（使用默认配置）
./redis-server

# 运行（指定配置文件）
./redis-server --config /path/to/redis.yaml
```

## 客户端配置

在支持MCP的客户端中添加服务器配置：

```json
{
  "mcpServers": {
    "redis-server": {
      "command": "path/to/redis-server",
      "args": ["--config", "config/redis.yaml"],
      "env": {},
      "description": "Redis MCP服务器 - 提供Redis数据库操作工具"
    }
  }
}
```

## 使用示例

### 设置键值
```json
{
  "name": "redis_set",
  "arguments": {
    "key": "user:1001",
    "value": "{\"name\":\"张三\",\"age\":25}",
    "expiration": "1h"
  }
}
```

### 获取键值
```json
{
  "name": "redis_get",
  "arguments": {
    "key": "user:1001"
  }
}
```

### 查找键
```json
{
  "name": "redis_keys",
  "arguments": {
    "pattern": "user:*"
  }
}
```

### 执行自定义命令
```json
{
  "name": "redis_execute",
  "arguments": {
    "command": "HGETALL",
    "args": ["user:1001"]
  }
}
```

## 依赖

- Go 1.19+
- github.com/redis/go-redis/v9
- gopkg.in/yaml.v3

## 协议支持

* ✅ JSON-RPC 2.0
* ✅ MCP 2024-11-05 协议版本
* ✅ 工具列表和调用 (Tools)
* ✅ 错误处理
* ✅ 连接池管理
* ✅ 超时控制 