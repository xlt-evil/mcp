# MCP资源 (Resources) 详解

## 什么是MCP资源？

MCP资源是Model Context Protocol中的另一个核心概念，它允许AI模型访问和管理外部数据源，如文件、数据库、网络服务等。资源提供了结构化的数据访问接口。

## 资源的核心概念

### 1. 资源定义
资源是具有以下特征的数据实体：
- **唯一标识**：URI格式的资源地址
- **类型信息**：MIME类型和元数据
- **访问权限**：读取、写入、删除等操作
- **生命周期**：创建、更新、删除等状态

### 2. 资源类型
- **文件资源**：本地文件系统中的文件
- **网络资源**：HTTP/HTTPS可访问的资源
- **数据库资源**：数据库中的记录或查询结果
- **内存资源**：程序运行时创建的数据
- **流式资源**：实时数据流

## MCP资源协议

### 1. 资源列表 (resources/list)
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "resources/list",
  "params": {
    "uri": "file:///path/to/directory"
  }
}
```

### 2. 资源读取 (resources/read)
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "resources/read",
  "params": {
    "uri": "file:///path/to/file.txt"
  }
}
```

### 3. 资源写入 (resources/write)
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "resources/write",
  "params": {
    "uri": "file:///path/to/file.txt",
    "data": "Hello, World!"
  }
}
```

### 4. 资源删除 (resources/delete)
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "resources/delete",
  "params": {
    "uri": "file:///path/to/file.txt"
  }
}
```

## 资源结构定义

### 1. 资源信息
```go
type Resource struct {
    URI         string            `json:"uri"`
    Name        string            `json:"name"`
    Description string            `json:"description,omitempty"`
    MimeType    string            `json:"mimeType"`
    Size        *int64            `json:"size,omitempty"`
    Created     *time.Time        `json:"created,omitempty"`
    Modified    *time.Time        `json:"modified,omitempty"`
    Expires     *time.Time        `json:"expires,omitempty"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}

type ListResourcesResult struct {
    Resources []Resource `json:"resources"`
}
```

### 2. 资源操作参数
```go
type ReadResourceParams struct {
    URI string `json:"uri"`
}

type ReadResourceResult struct {
    Contents []ContentItem `json:"contents"`
}

type WriteResourceParams struct {
    URI  string        `json:"uri"`
    Data []ContentItem `json:"data"`
}

type DeleteResourceParams struct {
    URI string `json:"uri"`
}
```

## 实现资源服务器

### 1. 资源管理器
```go
type ResourceManager struct {
    basePath string
    resources map[string]*Resource
}

func NewResourceManager(basePath string) *ResourceManager {
    return &ResourceManager{
        basePath: basePath,
        resources: make(map[string]*Resource),
    }
}

func (rm *ResourceManager) ListResources(uri string) ([]Resource, error) {
    // 实现资源列表逻辑
    var resources []Resource
    
    // 扫描目录或数据库
    // 构建资源列表
    
    return resources, nil
}

func (rm *ResourceManager) ReadResource(uri string) ([]ContentItem, error) {
    // 实现资源读取逻辑
    var contents []ContentItem
    
    // 根据URI类型读取资源
    // 返回内容项
    
    return contents, nil
}

func (rm *ResourceManager) WriteResource(uri string, data []ContentItem) error {
    // 实现资源写入逻辑
    
    // 验证权限
    // 写入数据
    // 更新元数据
    
    return nil
}

func (rm *ResourceManager) DeleteResource(uri string) error {
    // 实现资源删除逻辑
    
    // 验证权限
    // 删除资源
    // 清理元数据
    
    return nil
}
```

### 2. 资源处理函数
```go
func (s *HelloMCPServer) handleListResources(params *ListResourcesParams) (*ListResourcesResult, *JSONRPCError) {
    resources, err := s.resourceManager.ListResources(params.URI)
    if err != nil {
        return nil, &JSONRPCError{
            Code:    -32603,
            Message: fmt.Sprintf("Failed to list resources: %v", err),
        }
    }
    
    return &ListResourcesResult{Resources: resources}, nil
}

func (s *HelloMCPServer) handleReadResource(params *ReadResourceParams) (*ReadResourceResult, *JSONRPCError) {
    contents, err := s.resourceManager.ReadResource(params.URI)
    if err != nil {
        return nil, &JSONRPCError{
            Code:    -32603,
            Message: fmt.Sprintf("Failed to read resource: %v", err),
        }
    }
    
    return &ReadResourceResult{Contents: contents}, nil
}

func (s *HelloMCPServer) handleWriteResource(params *WriteResourceParams) (*JSONRPCMessage, *JSONRPCError) {
    err := s.resourceManager.WriteResource(params.URI, params.Data)
    if err != nil {
        return nil, &JSONRPCError{
            Code:    -32603,
            Message: fmt.Sprintf("Failed to write resource: %v", err),
        }
    }
    
    return &JSONRPCMessage{
        JSONRPC: "2.0",
        Result:  "Resource written successfully",
    }, nil
}

func (s *HelloMCPServer) handleDeleteResource(params *DeleteResourceParams) (*JSONRPCMessage, *JSONRPCError) {
    err := s.resourceManager.DeleteResource(params.URI)
    if err != nil {
        return nil, &JSONRPCError{
            Code:    -32603,
            Message: fmt.Sprintf("Failed to delete resource: %v", err),
        }
    }
    
    return &JSONRPCMessage{
        JSONRPC: "2.0",
        Result:  "Resource deleted successfully",
    }, nil
}
```

## 资源使用场景

### 1. 文件管理
```json
{
  "uri": "file:///home/user/documents/",
  "name": "Documents Directory",
  "description": "用户文档目录",
  "mimeType": "inode/directory",
  "metadata": {
    "owner": "user",
    "permissions": "755"
  }
}
```

### 2. 网络资源
```json
{
  "uri": "https://api.example.com/data",
  "name": "API Data",
  "description": "外部API数据",
  "mimeType": "application/json",
  "metadata": {
    "cache_control": "max-age=3600",
    "etag": "abc123"
  }
}
```

### 3. 数据库资源
```json
{
  "uri": "db://localhost/users/123",
  "name": "User Record",
  "description": "用户记录",
  "mimeType": "application/json",
  "metadata": {
    "table": "users",
    "primary_key": "123"
  }
}
```

## 资源安全考虑

### 1. 访问控制
- **身份验证**：验证用户身份
- **权限检查**：检查操作权限
- **路径验证**：防止路径遍历攻击
- **资源隔离**：隔离不同用户的资源

### 2. 数据保护
- **加密存储**：敏感数据加密
- **传输安全**：使用HTTPS等安全协议
- **审计日志**：记录所有操作
- **备份恢复**：数据备份和恢复

### 3. 限制措施
- **大小限制**：限制文件大小
- **类型限制**：限制文件类型
- **频率限制**：限制操作频率
- **配额管理**：用户存储配额

## 高级功能

### 1. 资源监控
```go
type ResourceMonitor struct {
    watchers map[string][]chan ResourceEvent
}

type ResourceEvent struct {
    Type     string   `json:"type"`     // created, modified, deleted
    URI      string   `json:"uri"`
    Resource Resource `json:"resource"`
}

func (rm *ResourceMonitor) WatchResource(uri string) <-chan ResourceEvent {
    // 实现资源监控
    ch := make(chan ResourceEvent, 10)
    rm.watchers[uri] = append(rm.watchers[uri], ch)
    return ch
}
```

### 2. 资源缓存
```go
type ResourceCache struct {
    cache map[string]CachedResource
    ttl   time.Duration
}

type CachedResource struct {
    Resource  Resource
    Data      []ContentItem
    CachedAt  time.Time
    ExpiresAt time.Time
}

func (rc *ResourceCache) Get(uri string) (*CachedResource, bool) {
    // 实现缓存逻辑
    if cached, exists := rc.cache[uri]; exists && time.Now().Before(cached.ExpiresAt) {
        return &cached, true
    }
    return nil, false
}
```

### 3. 资源同步
```go
type ResourceSync struct {
    localPath  string
    remotePath string
    interval   time.Duration
}

func (rs *ResourceSync) StartSync() {
    ticker := time.NewTicker(rs.interval)
    go func() {
        for range ticker.C {
            rs.syncResources()
        }
    }()
}

func (rs *ResourceSync) syncResources() {
    // 实现同步逻辑
    // 比较本地和远程资源
    // 执行必要的同步操作
}
```

## 最佳实践

### 1. 设计原则
- **统一接口**：提供一致的资源操作接口
- **错误处理**：完善的错误处理和恢复机制
- **性能优化**：实现缓存、分页等优化
- **可扩展性**：支持新的资源类型和操作

### 2. 实现建议
- **异步操作**：长时间操作使用异步处理
- **批量操作**：支持批量资源操作
- **增量更新**：实现增量同步和更新
- **状态管理**：维护资源状态和一致性

### 3. 测试策略
- **单元测试**：测试各个组件功能
- **集成测试**：测试完整的工作流程
- **性能测试**：测试资源操作性能
- **安全测试**：测试安全防护措施 