package db

import (
	"database/sql"
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
