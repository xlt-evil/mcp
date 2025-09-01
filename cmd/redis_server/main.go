package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"hello-mcp-server/config"
	"hello-mcp-server/redis"
	"hello-mcp-server/types"
)

// RedisMCPServer Redis MCP服务器
type RedisMCPServer struct {
	serverInfo   types.ServerInfo
	redisManager *redis.RedisManager
	redisConfig  *config.RedisConfig
}

func NewRedisMCPServer(configPath string) *RedisMCPServer {
	// 加载配置文件
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("Warning: Failed to load config file %s, using default config: %v", configPath, err)
		cfg = config.LoadDefaultConfig()
	}

	// 获取Redis配置
	redisConfig := cfg.GetRedisConfig()

	// 创建Redis管理器
	redisManager := redis.NewRedisManager(redisConfig)

	return &RedisMCPServer{
		serverInfo: types.ServerInfo{
			Name:    "redis-mcp-server",
			Version: "1.0.0",
		},
		redisManager: redisManager,
		redisConfig:  redisConfig,
	}
}

func (s *RedisMCPServer) handleInitialize(params *types.InitializeParams) *types.InitializeResult {
	log.Printf("Initialize request: protocolVersion=%s, client=%s %s",
		params.ProtocolVersion, params.ClientInfo.Name, params.ClientInfo.Version)

	// 尝试连接Redis
	if err := s.redisManager.Connect(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		log.Printf("Successfully connected to Redis: %s", s.redisConfig.GetAddr())
	}

	return &types.InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: types.ServerCapabilities{
			Tools: &types.ToolsCapability{
				ListChanged: false,
			},
		},
		ServerInfo: s.serverInfo,
	}
}

func (s *RedisMCPServer) handleListTools() *types.ListToolsResult {
	return &types.ListToolsResult{
		Tools: []types.Tool{
			{
				Name:        "redis_get",
				Description: "获取Redis键的值",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]types.Property{
						"key": {
							Type:        "string",
							Description: "要获取的键名",
						},
					},
					Required: []string{"key"},
				},
			},
			{
				Name:        "redis_set",
				Description: "设置Redis键值对",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]types.Property{
						"key": {
							Type:        "string",
							Description: "键名",
						},
						"value": {
							Type:        "string",
							Description: "值",
						},
						"expiration": {
							Type:        "string",
							Description: "过期时间（如：1h, 30m, 24h）",
						},
					},
					Required: []string{"key", "value"},
				},
			},
			{
				Name:        "redis_del",
				Description: "删除Redis键",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]types.Property{
						"keys": {
							Type:        "array",
							Description: "要删除的键名列表",
							Items: &types.Property{
								Type: "string",
							},
						},
					},
					Required: []string{"keys"},
				},
			},
			{
				Name:        "redis_keys",
				Description: "获取匹配模式的键列表",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]types.Property{
						"pattern": {
							Type:        "string",
							Description: "键模式（如：user:*）",
						},
					},
					Required: []string{"pattern"},
				},
			},
			{
				Name:        "redis_type",
				Description: "获取键的数据类型",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]types.Property{
						"key": {
							Type:        "string",
							Description: "键名",
						},
					},
					Required: []string{"key"},
				},
			},
			{
				Name:        "redis_ttl",
				Description: "获取键的TTL（生存时间）",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]types.Property{
						"key": {
							Type:        "string",
							Description: "键名",
						},
					},
					Required: []string{"key"},
				},
			},
			{
				Name:        "redis_info",
				Description: "获取Redis服务器信息",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]types.Property{
						"section": {
							Type:        "string",
							Description: "信息部分（如：server, clients, memory）",
						},
					},
				},
			},
			{
				Name:        "redis_dbsize",
				Description: "获取当前数据库的键数量",
				InputSchema: types.InputSchema{
					Type:       "object",
					Properties: map[string]types.Property{},
				},
			},
			{
				Name:        "redis_flushdb",
				Description: "清空当前数据库",
				InputSchema: types.InputSchema{
					Type:       "object",
					Properties: map[string]types.Property{},
				},
			},
			{
				Name:        "redis_execute",
				Description: "执行自定义Redis命令",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]types.Property{
						"command": {
							Type:        "string",
							Description: "Redis命令",
						},
						"args": {
							Type:        "array",
							Description: "命令参数",
							Items: &types.Property{
								Type: "string",
							},
						},
					},
					Required: []string{"command"},
				},
			},
			{
				Name:        "redis_status",
				Description: "检查Redis连接状态",
				InputSchema: types.InputSchema{
					Type:       "object",
					Properties: map[string]types.Property{},
				},
			},
		},
	}
}

func (s *RedisMCPServer) handleCallTool(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	switch params.Name {
	case "redis_get":
		return s.handleRedisGet(params)
	case "redis_set":
		return s.handleRedisSet(params)
	case "redis_del":
		return s.handleRedisDel(params)
	case "redis_keys":
		return s.handleRedisKeys(params)
	case "redis_type":
		return s.handleRedisType(params)
	case "redis_ttl":
		return s.handleRedisTTL(params)
	case "redis_info":
		return s.handleRedisInfo(params)
	case "redis_dbsize":
		return s.handleRedisDBSize(params)
	case "redis_flushdb":
		return s.handleRedisFlushDB(params)
	case "redis_execute":
		return s.handleRedisExecute(params)
	case "redis_status":
		return s.handleRedisStatus(params)
	default:
		return nil, &types.JSONRPCError{
			Code:    -32601,
			Message: fmt.Sprintf("Unknown tool: %s", params.Name),
		}
	}
}

func (s *RedisMCPServer) handleRedisGet(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	key, ok := params.Arguments["key"].(string)
	if !ok || key == "" {
		return nil, &types.JSONRPCError{
			Code:    -32602,
			Message: "Key is required",
		}
	}

	if !s.redisManager.IsConnected() {
		if err := s.redisManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Redis connection failed: %v", err),
			}
		}
	}

	result := s.redisManager.Get(key)
	if !result.Success {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: result.Error,
		}
	}

	resultText := fmt.Sprintf("✅ 获取键值成功！\n\n")
	resultText += fmt.Sprintf("🔑 键名：%s\n", key)
	resultText += fmt.Sprintf("📄 值：%v\n", result.Data)

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *RedisMCPServer) handleRedisSet(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	key, ok := params.Arguments["key"].(string)
	if !ok || key == "" {
		return nil, &types.JSONRPCError{
			Code:    -32602,
			Message: "Key is required",
		}
	}

	value, ok := params.Arguments["value"].(string)
	if !ok {
		return nil, &types.JSONRPCError{
			Code:    -32602,
			Message: "Value is required",
		}
	}

	var expiration time.Duration
	if expStr, ok := params.Arguments["expiration"].(string); ok && expStr != "" {
		var err error
		expiration, err = time.ParseDuration(expStr)
		if err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32602,
				Message: fmt.Sprintf("Invalid expiration format: %v", err),
			}
		}
	}

	if !s.redisManager.IsConnected() {
		if err := s.redisManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Redis connection failed: %v", err),
			}
		}
	}

	result := s.redisManager.Set(key, value, expiration)
	if !result.Success {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: result.Error,
		}
	}

	resultText := fmt.Sprintf("✅ 设置键值成功！\n\n")
	resultText += fmt.Sprintf("🔑 键名：%s\n", key)
	resultText += fmt.Sprintf("📄 值：%s\n", value)
	if expiration > 0 {
		resultText += fmt.Sprintf("⏰ 过期时间：%s\n", expiration)
	}

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *RedisMCPServer) handleRedisDel(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	keysInterface, ok := params.Arguments["keys"].([]interface{})
	if !ok {
		return nil, &types.JSONRPCError{
			Code:    -32602,
			Message: "Keys array is required",
		}
	}

	var keys []string
	for _, k := range keysInterface {
		if key, ok := k.(string); ok {
			keys = append(keys, key)
		}
	}

	if len(keys) == 0 {
		return nil, &types.JSONRPCError{
			Code:    -32602,
			Message: "At least one key is required",
		}
	}

	if !s.redisManager.IsConnected() {
		if err := s.redisManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Redis connection failed: %v", err),
			}
		}
	}

	result := s.redisManager.Del(keys...)
	if !result.Success {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: result.Error,
		}
	}

	resultText := fmt.Sprintf("✅ 删除键成功！\n\n")
	resultText += fmt.Sprintf("🗑️  删除的键：%v\n", keys)
	resultText += fmt.Sprintf("📊 删除数量：%v\n", result.Data)

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *RedisMCPServer) handleRedisKeys(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	pattern, ok := params.Arguments["pattern"].(string)
	if !ok || pattern == "" {
		return nil, &types.JSONRPCError{
			Code:    -32602,
			Message: "Pattern is required",
		}
	}

	if !s.redisManager.IsConnected() {
		if err := s.redisManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Redis connection failed: %v", err),
			}
		}
	}

	result := s.redisManager.Keys(pattern)
	if !result.Success {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: result.Error,
		}
	}

	keys, ok := result.Data.([]string)
	if !ok {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: "Invalid keys data type",
		}
	}

	resultText := fmt.Sprintf("🔍 键匹配结果\n\n")
	resultText += fmt.Sprintf("🎯 匹配模式：%s\n", pattern)
	resultText += fmt.Sprintf("📊 匹配数量：%d\n\n", len(keys))

	if len(keys) > 0 {
		resultText += "📝 匹配的键：\n"
		for i, key := range keys {
			resultText += fmt.Sprintf("  %d. %s\n", i+1, key)
		}
	} else {
		resultText += "❌ 没有找到匹配的键"
	}

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *RedisMCPServer) handleRedisType(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	key, ok := params.Arguments["key"].(string)
	if !ok || key == "" {
		return nil, &types.JSONRPCError{
			Code:    -32602,
			Message: "Key is required",
		}
	}

	if !s.redisManager.IsConnected() {
		if err := s.redisManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Redis connection failed: %v", err),
			}
		}
	}

	result := s.redisManager.Type(key)
	if !result.Success {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: result.Error,
		}
	}

	resultText := fmt.Sprintf("🔍 键类型信息\n\n")
	resultText += fmt.Sprintf("🔑 键名：%s\n", key)
	resultText += fmt.Sprintf("📊 类型：%v\n", result.Data)

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *RedisMCPServer) handleRedisTTL(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	key, ok := params.Arguments["key"].(string)
	if !ok || key == "" {
		return nil, &types.JSONRPCError{
			Code:    -32602,
			Message: "Key is required",
		}
	}

	if !s.redisManager.IsConnected() {
		if err := s.redisManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Redis connection failed: %v", err),
			}
		}
	}

	result := s.redisManager.TTL(key)
	if !result.Success {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: result.Error,
		}
	}

	ttl, ok := result.Data.(float64)
	if !ok {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: "Invalid TTL data type",
		}
	}

	resultText := fmt.Sprintf("⏰ 键TTL信息\n\n")
	resultText += fmt.Sprintf("🔑 键名：%s\n", key)
	if ttl == -1 {
		resultText += "📊 TTL：永不过期\n"
	} else if ttl == -2 {
		resultText += "📊 TTL：键不存在\n"
	} else {
		resultText += fmt.Sprintf("📊 TTL：%.0f秒\n", ttl)
	}

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *RedisMCPServer) handleRedisInfo(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	section := ""
	if sec, ok := params.Arguments["section"].(string); ok {
		section = sec
	}

	if !s.redisManager.IsConnected() {
		if err := s.redisManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Redis connection failed: %v", err),
			}
		}
	}

	result := s.redisManager.Info(section)
	if !result.Success {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: result.Error,
		}
	}

	resultText := fmt.Sprintf("📊 Redis服务器信息\n\n")
	if section != "" {
		resultText += fmt.Sprintf("📋 信息部分：%s\n\n", section)
	}
	resultText += fmt.Sprintf("📄 详细信息：\n%s", result.Data)

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *RedisMCPServer) handleRedisDBSize(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	if !s.redisManager.IsConnected() {
		if err := s.redisManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Redis connection failed: %v", err),
			}
		}
	}

	result := s.redisManager.DBSize()
	if !result.Success {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: result.Error,
		}
	}

	resultText := fmt.Sprintf("📊 数据库大小信息\n\n")
	resultText += fmt.Sprintf("🗄️  数据库：%d\n", s.redisConfig.GetDB())
	resultText += fmt.Sprintf("📈 键数量：%v\n", result.Data)

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *RedisMCPServer) handleRedisFlushDB(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	if !s.redisManager.IsConnected() {
		if err := s.redisManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Redis connection failed: %v", err),
			}
		}
	}

	result := s.redisManager.FlushDB()
	if !result.Success {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: result.Error,
		}
	}

	resultText := fmt.Sprintf("🗑️  数据库清空成功！\n\n")
	resultText += fmt.Sprintf("🗄️  数据库：%d\n", s.redisConfig.GetDB())
	resultText += fmt.Sprintf("✅ 状态：%v\n", result.Data)

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *RedisMCPServer) handleRedisExecute(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	command, ok := params.Arguments["command"].(string)
	if !ok || command == "" {
		return nil, &types.JSONRPCError{
			Code:    -32602,
			Message: "Command is required",
		}
	}

	var args []interface{}
	if argsInterface, ok := params.Arguments["args"].([]interface{}); ok {
		args = argsInterface
	}

	if !s.redisManager.IsConnected() {
		if err := s.redisManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Redis connection failed: %v", err),
			}
		}
	}

	result := s.redisManager.ExecuteCommand(command, args...)
	if !result.Success {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: result.Error,
		}
	}

	resultText := fmt.Sprintf("⚡ 命令执行成功！\n\n")
	resultText += fmt.Sprintf("🔧 命令：%s\n", command)
	if len(args) > 0 {
		resultText += fmt.Sprintf("📝 参数：%v\n", args)
	}
	resultText += fmt.Sprintf("📊 结果：%v\n", result.Data)

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *RedisMCPServer) handleRedisStatus(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	isConnected := s.redisManager.IsConnected()

	resultText := fmt.Sprintf("🔍 Redis连接状态\n\n")
	resultText += fmt.Sprintf("🌐 地址：%s\n", s.redisConfig.GetAddr())
	resultText += fmt.Sprintf("🗄️  数据库：%d\n", s.redisConfig.GetDB())
	resultText += fmt.Sprintf("📊 状态：")

	if isConnected {
		resultText += "✅ 已连接\n"
	} else {
		resultText += "❌ 未连接\n"

		// 尝试重新连接
		if err := s.redisManager.Connect(); err != nil {
			resultText += fmt.Sprintf("🔄 重连失败：%v\n", err)
		} else {
			resultText += "🔄 重连成功！\n"
		}
	}

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *RedisMCPServer) processMessage(msg *types.JSONRPCMessage) *types.JSONRPCMessage {
	response := &types.JSONRPCMessage{
		JSONRPC: "2.0",
		ID:      msg.ID,
	}

	switch msg.Method {
	case "initialize":
		// 解析初始化参数
		paramsBytes, err := json.Marshal(msg.Params)
		if err != nil {
			response.Error = &types.JSONRPCError{
				Code:    -32602,
				Message: "Invalid params",
			}
			return response
		}

		var initParams types.InitializeParams
		if err := json.Unmarshal(paramsBytes, &initParams); err != nil {
			response.Error = &types.JSONRPCError{
				Code:    -32602,
				Message: "Invalid initialize params",
			}
			return response
		}

		response.Result = s.handleInitialize(&initParams)

	case "initialized":
		// initialized 通知不需要响应
		return nil

	case "tools/list":
		response.Result = s.handleListTools()

	case "tools/call":
		// 解析工具调用参数
		paramsBytes, err := json.Marshal(msg.Params)
		if err != nil {
			response.Error = &types.JSONRPCError{
				Code:    -32602,
				Message: "Invalid params",
			}
			return response
		}

		var callParams types.CallToolParams
		if err := json.Unmarshal(paramsBytes, &callParams); err != nil {
			response.Error = &types.JSONRPCError{
				Code:    -32602,
				Message: "Invalid call tool params",
			}
			return response
		}

		result, rpcErr := s.handleCallTool(&callParams)
		if rpcErr != nil {
			response.Error = rpcErr
		} else {
			response.Result = result
		}

	default:
		response.Error = &types.JSONRPCError{
			Code:    -32601,
			Message: fmt.Sprintf("Method not found: %s", msg.Method),
		}
	}

	return response
}

func (s *RedisMCPServer) sendMessage(msg *types.JSONRPCMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	fmt.Println(string(msgBytes))
	return nil
}

func (s *RedisMCPServer) run() {
	log.Println("Redis MCP Server starting...")
	log.Printf("Redis config: %s", s.redisConfig.GetAddr())

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// 解析输入消息
		var msg types.JSONRPCMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			log.Printf("Failed to parse JSON: %v", err)
			errorMsg := &types.JSONRPCMessage{
				JSONRPC: "2.0",
				Error: &types.JSONRPCError{
					Code:    -32700,
					Message: "Parse error",
				},
			}
			s.sendMessage(errorMsg)
			continue
		}

		log.Printf("Received message: method=%s, id=%v", msg.Method, msg.ID)

		// 处理消息
		response := s.processMessage(&msg)
		if response != nil {
			if err := s.sendMessage(response); err != nil {
				log.Printf("Failed to send response: %v", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}

	// 关闭Redis连接
	if err := s.redisManager.Close(); err != nil {
		log.Printf("Failed to close Redis connection: %v", err)
	}
}

func main() {
	// 设置日志输出到stderr
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 默认配置文件路径
	configPath := "config/redis.yaml"

	// 检查命令行参数
	if len(os.Args) > 2 && os.Args[1] == "--config" {
		configPath = os.Args[2]
	}

	log.Printf("Using config file: %s", configPath)

	server := NewRedisMCPServer(configPath)
	server.run()
}
