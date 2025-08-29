package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"hello-mcp-server/config"

	_ "github.com/go-sql-driver/mysql"
)

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	config *config.DatabaseConfig
	db     *sql.DB
}

// QueryResult 查询结果结构
type QueryResult struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
	Count   int             `json:"count"`
	Error   string          `json:"error,omitempty"`
}

// NewDatabaseManager 创建数据库管理器
func NewDatabaseManager(cfg *config.DatabaseConfig) *DatabaseManager {
	return &DatabaseManager{
		config: cfg,
	}
}

// Connect 连接数据库
func (dm *DatabaseManager) Connect() error {
	if !dm.config.IsValid() {
		return fmt.Errorf("invalid database configuration")
	}

	dsn := dm.config.GetDSN()
	log.Printf("Connecting to database: %s@%s:%d/%s",
		dm.config.User, dm.config.Host, dm.config.Port, dm.config.Name)

	var err error
	dm.db, err = sql.Open(dm.config.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}

	// 设置连接池参数
	dm.db.SetMaxOpenConns(25)
	dm.db.SetMaxIdleConns(5)
	dm.db.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	if err := dm.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	log.Printf("Successfully connected to database: %s", dm.config.Name)
	return nil
}

// Close 关闭数据库连接
func (dm *DatabaseManager) Close() error {
	if dm.db != nil {
		return dm.db.Close()
	}
	return nil
}

// ExecuteQuery 执行查询（修复方法接收器）
func (dm *DatabaseManager) ExecuteQuery(query string) *QueryResult {
	if dm.db == nil {
		return &QueryResult{
			Error: "Database not connected",
		}
	}

	// 执行查询
	rows, err := dm.db.Query(query)
	if err != nil {
		return &QueryResult{
			Error: fmt.Sprintf("Query execution failed: %v", err),
		}
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return &QueryResult{
			Error: fmt.Sprintf("Failed to get columns: %v", err),
		}
	}

	// 准备结果容器
	var resultRows [][]interface{}
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// 读取数据
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return &QueryResult{
				Error: fmt.Sprintf("Failed to scan row: %v", err),
			}
		}

		// 复制值到新切片
		row := make([]interface{}, len(columns))
		for i, v := range values {
			row[i] = v
		}
		resultRows = append(resultRows, row)
	}

	if err := rows.Err(); err != nil {
		return &QueryResult{
			Error: fmt.Sprintf("Error during rows iteration: %v", err),
		}
	}

	return &QueryResult{
		Columns: columns,
		Rows:    resultRows,
		Count:   len(resultRows),
	}
}

// GetTableInfo 获取表信息
func (dm *DatabaseManager) GetTableInfo() ([]string, error) {
	if dm.db == nil {
		return nil, fmt.Errorf("database not connected")
	}

	var tables []string
	query := "SHOW TABLES"
	rows, err := dm.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %v", err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// GetTableSchema 获取表结构
func (dm *DatabaseManager) GetTableSchema(tableName string) (*QueryResult, error) {
	if dm.db == nil {
		return nil, fmt.Errorf("database not connected")
	}

	query := fmt.Sprintf("DESCRIBE %s", tableName)
	result := dm.ExecuteQuery(query)
	if result.Error != "" {
		return nil, fmt.Errorf("failed to get table schema: %v", result.Error)
	}

	return result, nil
}

// IsConnected 检查是否已连接
func (dm *DatabaseManager) IsConnected() bool {
	if dm.db == nil {
		return false
	}
	return dm.db.Ping() == nil
}
