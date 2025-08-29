package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"hello-mcp-server/config"
	"hello-mcp-server/database"
	"hello-mcp-server/types"
)

// 数据库查询工具参数
type DatabaseQueryParams struct {
	SQL string `json:"sql"`
}

type DatabaseQueryResult struct {
	Success bool                  `json:"success"`
	Data    *database.QueryResult `json:"data,omitempty"`
	Error   string                `json:"error,omitempty"`
}

// DatabaseMCPServer 数据库MCP服务器
type DatabaseMCPServer struct {
	serverInfo types.ServerInfo
	dbManager  *database.DatabaseManager
	dbConfig   *config.DatabaseConfig
}

func NewDatabaseMCPServer(configPath string) *DatabaseMCPServer {
	// 加载配置文件
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("Warning: Failed to load config file %s, using default config: %v", configPath, err)
		cfg = config.LoadDefaultConfig()
	}

	// 获取数据库配置
	dbConfig := cfg.GetDatabaseConfig()

	// 创建数据库管理器
	dbManager := database.NewDatabaseManager(dbConfig)

	return &DatabaseMCPServer{
		serverInfo: types.ServerInfo{
			Name:    "database-mcp-server",
			Version: "1.0.0",
		},
		dbManager: dbManager,
		dbConfig:  dbConfig,
	}
}

func (s *DatabaseMCPServer) handleInitialize(params *types.InitializeParams) *types.InitializeResult {
	log.Printf("Initialize request: protocolVersion=%s, client=%s %s",
		params.ProtocolVersion, params.ClientInfo.Name, params.ClientInfo.Version)

	// 尝试连接数据库
	if err := s.dbManager.Connect(); err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
	} else {
		log.Printf("Successfully connected to database: %s", s.dbConfig.Name)
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

func (s *DatabaseMCPServer) handleListTools() *types.ListToolsResult {
	return &types.ListToolsResult{
		Tools: []types.Tool{
			{
				Name:        "database_query",
				Description: "执行SQL查询并返回结果",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]types.Property{
						"sql": {
							Type:        "string",
							Description: "要执行的SQL查询语句",
						},
					},
					Required: []string{"sql"},
				},
			},
			{
				Name:        "database_tables",
				Description: "获取数据库中的所有表名",
				InputSchema: types.InputSchema{
					Type:       "object",
					Properties: map[string]types.Property{},
				},
			},
			{
				Name:        "database_schema",
				Description: "获取指定表的结构信息",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]types.Property{
						"table_name": {
							Type:        "string",
							Description: "要查看结构的表名",
						},
					},
					Required: []string{"table_name"},
				},
			},
			{
				Name:        "database_status",
				Description: "检查数据库连接状态",
				InputSchema: types.InputSchema{
					Type:       "object",
					Properties: map[string]types.Property{},
				},
			},
		},
	}
}

func (s *DatabaseMCPServer) handleCallTool(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	switch params.Name {
	case "database_query":
		return s.handleDatabaseQuery(params)
	case "database_tables":
		return s.handleDatabaseTables(params)
	case "database_schema":
		return s.handleDatabaseSchema(params)
	case "database_status":
		return s.handleDatabaseStatus(params)
	default:
		return nil, &types.JSONRPCError{
			Code:    -32601,
			Message: fmt.Sprintf("Unknown tool: %s", params.Name),
		}
	}
}

func (s *DatabaseMCPServer) handleDatabaseQuery(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	// 获取SQL参数
	sqlQuery, ok := params.Arguments["sql"].(string)
	if !ok || sqlQuery == "" {
		return nil, &types.JSONRPCError{
			Code:    -32602,
			Message: "SQL query is required",
		}
	}

	// 检查数据库连接
	if !s.dbManager.IsConnected() {
		if err := s.dbManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Database connection failed: %v", err),
			}
		}
	}

	// 执行查询
	result := s.dbManager.ExecuteQuery(sqlQuery)
	if result.Error != "" {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: fmt.Sprintf("Query execution failed: %v", result.Error),
		}
	}

	// 格式化结果
	resultText := fmt.Sprintf("✅ 查询执行成功！\n\n📊 查询结果：\n")
	resultText += fmt.Sprintf("📋 列数：%d\n", len(result.Columns))
	resultText += fmt.Sprintf("📝 行数：%d\n\n", result.Count)

	// 添加列名
	if len(result.Columns) > 0 {
		resultText += "🏷️  列名：\n"
		for i, col := range result.Columns {
			resultText += fmt.Sprintf("  %d. %s\n", i+1, col)
		}
		resultText += "\n"
	}

	// 添加数据行（限制显示前10行）
	maxRows := 10
	if result.Count > maxRows {
		resultText += fmt.Sprintf("📊 数据行（显示前%d行）：\n", maxRows)
	} else {
		resultText += "📊 数据行：\n"
	}

	for i, row := range result.Rows {
		if i >= maxRows {
			break
		}
		resultText += fmt.Sprintf("  行 %d: ", i+1)
		for j, cell := range row {
			if j > 0 {
				resultText += " | "
			}
			resultText += fmt.Sprintf("%v", cell)
		}
		resultText += "\n"
	}

	if result.Count > maxRows {
		resultText += fmt.Sprintf("\n... 还有 %d 行数据未显示", result.Count-maxRows)
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

func (s *DatabaseMCPServer) handleDatabaseTables(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	// 检查数据库连接
	if !s.dbManager.IsConnected() {
		if err := s.dbManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Database connection failed: %v", err),
			}
		}
	}

	// 获取表列表
	tables, err := s.dbManager.GetTableInfo()
	if err != nil {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: fmt.Sprintf("Failed to get tables: %v", err),
		}
	}

	// 格式化结果
	resultText := fmt.Sprintf("📋 数据库表列表\n\n")
	resultText += fmt.Sprintf("🗄️  数据库：%s\n", s.dbConfig.Name)
	resultText += fmt.Sprintf("📊 表总数：%d\n\n", len(tables))

	if len(tables) > 0 {
		resultText += "📝 表名列表：\n"
		for i, table := range tables {
			resultText += fmt.Sprintf("  %d. %s\n", i+1, table)
		}
	} else {
		resultText += "❌ 没有找到任何表"
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

func (s *DatabaseMCPServer) handleDatabaseSchema(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	// 获取表名参数
	tableName, ok := params.Arguments["table_name"].(string)
	if !ok || tableName == "" {
		return nil, &types.JSONRPCError{
			Code:    -32602,
			Message: "Table name is required",
		}
	}

	// 检查数据库连接
	if !s.dbManager.IsConnected() {
		if err := s.dbManager.Connect(); err != nil {
			return nil, &types.JSONRPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Database connection failed: %v", err),
			}
		}
	}

	// 获取表结构
	schema, err := s.dbManager.GetTableSchema(tableName)
	if err != nil {
		return nil, &types.JSONRPCError{
			Code:    -32603,
			Message: fmt.Sprintf("Failed to get table schema: %v", err),
		}
	}

	// 格式化结果
	resultText := fmt.Sprintf("🏗️  表结构信息\n\n")
	resultText += fmt.Sprintf("📋 表名：%s\n", tableName)
	resultText += fmt.Sprintf("📊 字段数：%d\n\n", len(schema.Columns))

	if len(schema.Columns) > 0 {
		resultText += "📝 字段详情：\n"
		for i, row := range schema.Rows {
			if len(row) >= 6 {
				fieldName := fmt.Sprintf("%v", row[0])
				fieldType := fmt.Sprintf("%v", row[1])
				nullable := fmt.Sprintf("%v", row[2])
				key := fmt.Sprintf("%v", row[3])
				defaultValue := fmt.Sprintf("%v", row[4])
				extra := fmt.Sprintf("%v", row[5])

				resultText += fmt.Sprintf("  %d. %s\n", i+1, fieldName)
				resultText += fmt.Sprintf("     类型: %s\n", fieldType)
				resultText += fmt.Sprintf("     可空: %s\n", nullable)
				resultText += fmt.Sprintf("     键: %s\n", key)
				resultText += fmt.Sprintf("     默认值: %s\n", defaultValue)
				resultText += fmt.Sprintf("     额外: %s\n", extra)
				resultText += "\n"
			}
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

func (s *DatabaseMCPServer) handleDatabaseStatus(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	// 检查连接状态
	isConnected := s.dbManager.IsConnected()

	// 格式化结果
	resultText := fmt.Sprintf("🔍 数据库连接状态\n\n")
	resultText += fmt.Sprintf("🗄️  数据库：%s\n", s.dbConfig.Name)
	resultText += fmt.Sprintf("🌐 主机：%s:%d\n", s.dbConfig.Host, s.dbConfig.Port)
	resultText += fmt.Sprintf("👤 用户：%s\n", s.dbConfig.User)
	resultText += fmt.Sprintf("🔌 驱动：%s\n", s.dbConfig.Driver)
	resultText += fmt.Sprintf("📊 状态：")

	if isConnected {
		resultText += "✅ 已连接\n"
	} else {
		resultText += "❌ 未连接\n"

		// 尝试重新连接
		if err := s.dbManager.Connect(); err != nil {
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

func (s *DatabaseMCPServer) processMessage(msg *types.JSONRPCMessage) *types.JSONRPCMessage {
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

func (s *DatabaseMCPServer) sendMessage(msg *types.JSONRPCMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	fmt.Println(string(msgBytes))
	return nil
}

func (s *DatabaseMCPServer) run() {
	log.Println("Database MCP Server starting...")
	log.Printf("Database config: %s@%s:%d/%s",
		s.dbConfig.User, s.dbConfig.Host, s.dbConfig.Port, s.dbConfig.Name)

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

	// 关闭数据库连接
	if err := s.dbManager.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
	}
}

func main() {
	// 设置日志输出到stderr
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 默认配置文件路径
	configPath := "config/database.yaml"

	// 检查命令行参数
	if len(os.Args) > 2 && os.Args[1] == "--config" {
		configPath = os.Args[2]
	}

	log.Printf("Using config file: %s", configPath)

	server := NewDatabaseMCPServer(configPath)
	server.run()
}
