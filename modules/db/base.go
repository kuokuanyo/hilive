package db

import (
	"database/sql"
	"sync"
)

// Base struct
type Base struct {
	DB *sql.DB
	// sync.Once為唯一鎖，在代碼需要被執行時，只會被執行一次
	Once sync.Once
}
