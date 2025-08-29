package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// MCP客户端结构
type MCPClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

// 创建新的MCP客户端
func NewMCPClient(serverPath string) (*MCPClient, error) {
	cmd := exec.Command(serverPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start server: %v", err)
	}

	return &MCPClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
	}, nil
}

// 发送消息到服务器
func (c *MCPClient) SendMessage(msg interface{}) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = c.stdin.Write(append(msgBytes, '\n'))
	return err
}

// 接收服务器响应
func (c *MCPClient) ReceiveResponse() ([]byte, error) {
	reader := bufio.NewReader(c.stdout)
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	return []byte(strings.TrimSpace(line)), nil
}

// 关闭客户端
func (c *MCPClient) Close() error {
	if c.stdin != nil {
		c.stdin.Close()
	}
	if c.cmd != nil && c.cmd.Process != nil {
		return c.cmd.Process.Kill()
	}
	return nil
}

// 测试初始化流程
func testInitialize(client *MCPClient) error {
	fmt.Println("🔧 测试初始化流程...")

	// 发送初始化请求
	initMsg := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"roots": map[string]interface{}{
					"listChanged": true,
				},
				"sampling": map[string]interface{}{},
			},
			"clientInfo": map[string]interface{}{
				"name":    "Test Client",
				"version": "1.0.0",
			},
		},
	}

	if err := client.SendMessage(initMsg); err != nil {
		return fmt.Errorf("failed to send initialize: %v", err)
	}

	// 接收响应
	response, err := client.ReceiveResponse()
	if err != nil {
		return fmt.Errorf("failed to receive initialize response: %v", err)
	}

	fmt.Printf("✅ 初始化响应: %s\n", string(response))

	// 发送initialized通知
	initializedMsg := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "initialized",
	}

	if err := client.SendMessage(initializedMsg); err != nil {
		return fmt.Errorf("failed to send initialized: %v", err)
	}

	fmt.Println("✅ 初始化完成")
	return nil
}

// 测试工具列表
func testListTools(client *MCPClient) error {
	fmt.Println("🔧 测试工具列表...")

	listMsg := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "tools/list",
	}

	if err := client.SendMessage(listMsg); err != nil {
		return fmt.Errorf("failed to send tools/list: %v", err)
	}

	response, err := client.ReceiveResponse()
	if err != nil {
		return fmt.Errorf("failed to receive tools/list response: %v", err)
	}

	fmt.Printf("✅ 工具列表响应: %s\n", string(response))
	return nil
}

// 测试工具调用
func testCallTool(client *MCPClient) error {
	fmt.Println("🔧 测试工具调用...")

	callMsg := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      3,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "say_hello",
			"arguments": map[string]interface{}{
				"person_name":      "测试用户",
				"greeting_message": "你好，这是测试消息！",
			},
		},
	}

	if err := client.SendMessage(callMsg); err != nil {
		return fmt.Errorf("failed to send tools/call: %v", err)
	}

	response, err := client.ReceiveResponse()
	if err != nil {
		return fmt.Errorf("failed to receive tools/call response: %v", err)
	}

	fmt.Printf("✅ 工具调用响应: %s\n", string(response))
	return nil
}

// 主测试函数
func main() {
	log.Println("🚀 开始MCP服务器测试...")

	// 检查服务器可执行文件
	serverPath := "./mcp-server"
	if _, err := os.Stat(serverPath); os.IsNotExist(err) {
		log.Fatalf("服务器可执行文件不存在: %s", serverPath)
	}

	// 创建客户端
	client, err := NewMCPClient(serverPath)
	if err != nil {
		log.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)

	// 运行测试
	tests := []func(*MCPClient) error{
		testInitialize,
		testListTools,
		testCallTool,
	}

	for i, test := range tests {
		fmt.Printf("\n📋 测试 %d/%d\n", i+1, len(tests))
		if err := test(client); err != nil {
			log.Printf("❌ 测试失败: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n🎉 所有测试完成！")
}
