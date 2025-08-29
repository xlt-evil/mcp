package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"hello-mcp-server/types"
)

// HelloMCPServer 问候MCP服务器
type HelloMCPServer struct {
	serverInfo types.ServerInfo
}

func NewHelloMCPServer() *HelloMCPServer {
	return &HelloMCPServer{
		serverInfo: types.ServerInfo{
			Name:    "hello-mcp-server",
			Version: "1.0.0",
		},
	}
}

func (s *HelloMCPServer) handleInitialize(params *types.InitializeParams) *types.InitializeResult {
	log.Printf("Initialize request: protocolVersion=%s, client=%s %s",
		params.ProtocolVersion, params.ClientInfo.Name, params.ClientInfo.Version)

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

func (s *HelloMCPServer) handleListTools() *types.ListToolsResult {
	return &types.ListToolsResult{
		Tools: []types.Tool{
			{
				Name:        "say_hello",
				Description: "向指定的人说你好，记录问候信息并返回友好的回应",
				InputSchema: types.InputSchema{
					Type: "object",
					Properties: map[string]types.Property{
						"person_name": {
							Type:        "string",
							Description: "要问候的人的姓名",
						},
						"greeting_message": {
							Type:        "string",
							Description: "可选的自定义问候消息，默认为'你好'",
						},
					},
					Required: []string{"person_name"},
				},
			},
		},
	}
}

func (s *HelloMCPServer) handleCallTool(params *types.CallToolParams) (*types.CallToolResult, *types.JSONRPCError) {
	if params.Name != "say_hello" {
		return nil, &types.JSONRPCError{
			Code:    -32601,
			Message: fmt.Sprintf("Unknown tool: %s", params.Name),
		}
	}

	// 获取参数
	personName := "朋友"
	if name, ok := params.Arguments["person_name"].(string); ok && name != "" {
		personName = name
	}

	greetingMessage := "你好"
	if msg, ok := params.Arguments["greeting_message"].(string); ok && msg != "" {
		greetingMessage = msg
	}

	// 记录日志
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] 向 %s 说: %s", timestamp, personName, greetingMessage)

	// 写入日志文件
	if err := s.writeLog(logEntry); err != nil {
		log.Printf("Failed to write log: %v", err)
	}

	// 生成回应
	responses := []string{
		fmt.Sprintf("你好 %s！很高兴见到你！", personName),
		fmt.Sprintf("嗨 %s！希望你今天过得愉快！", personName),
		fmt.Sprintf("%s，你好！有什么可以帮助你的吗？", personName),
		fmt.Sprintf("你好 %s！欢迎使用MCP服务！", personName),
	}

	responseIndex := len(personName) % len(responses)
	response := responses[responseIndex]

	resultText := fmt.Sprintf("✨ 问候已发送！\n\n📝 日志记录：%s\n🎉 回应：%s", logEntry, response)

	return &types.CallToolResult{
		Content: []types.ContentItem{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *HelloMCPServer) writeLog(message string) error {
	file, err := os.OpenFile("hello_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(message + "\n")
	return err
}

func (s *HelloMCPServer) processMessage(msg *types.JSONRPCMessage) *types.JSONRPCMessage {
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

func (s *HelloMCPServer) sendMessage(msg *types.JSONRPCMessage) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	fmt.Println(string(msgBytes))
	return nil
}

func (s *HelloMCPServer) run() {
	log.Println("Hello MCP Server starting...")

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
}

func main() {
	// 设置日志输出到stderr
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	server := NewHelloMCPServer()
	server.run()
}
