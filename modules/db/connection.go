package db

import (
	"database/sql"
	"fmt"
	"hilive/modules/config"
	"hilive/modules/service"
)

// Connection 資料庫連接的處理程序，Connection也屬於Service(interface)
type Connection interface {
	// Service(interface) method
	Name() string

	// 初始化資料庫
	InitDB(cfg config.Database) Connection

	// 查詢
	Query(query string, args ...interface{}) ([]map[string]interface{}, error)

	// 執行
	Exec(query string, args ...interface{}) (sql.Result, error)

	// 取得Tx
	GetTx() *sql.Tx

	// QueryWithTx 利用tx查詢資料
	QueryWithTx(tx *sql.Tx, query string, args ...interface{}) ([]map[string]interface{}, error)

	// ExecWithTx 利用tx執行命令
	ExecWithTx(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error)
}

// GetConnectionFromService 透過資料庫引擎從Service取得Connection(interface)
func GetConnectionFromService(srvs service.List) Connection {
	// Get 取得匹配參數的Service
	if conn, ok := srvs.Get(config.GetDatabaseDriver()).(Connection); ok {
		return conn
	}
	panic("取得Connection(interface)發生錯誤")
}

// GetConnectionByDriver 取得資料庫引擎的struct(同時也是Connection(interface))
func GetConnectionByDriver(driver string) Connection {
	switch driver {
	case "mysql":
		return GetDefaultMysql()
	default:
		panic("找不到資料庫引擎!")
	}
}

// ConvertServiceToConnection 將Service轉換Connection
func ConvertServiceToConnection(s interface{}) Connection {
	if c, ok := s.(Connection); ok {
		return c
	}
	panic("Service轉換Connection失敗")
}

// GetAggregationExpression 判斷資料庫引擎取得聚合表達式
func GetAggregationExpression(driver, field, headField, delimiter string) string {
	switch driver {
	case "mysql":
		return fmt.Sprintf("group_concat(%s separator '%s') as %s", field, delimiter, headField)
	default:
		panic("取得聚合表達式發生錯誤")
	}
}
