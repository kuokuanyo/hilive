package models

import "hilive/modules/db"

// Base is base model structure. 紀錄資料表名稱...等資訊
type Base struct {
	TableName string
	Conn      db.Connection
}

// Table 設置SQL(struct)
func (b Base) Table(table string) *db.SQL {
	return db.Table(table, b.Conn)
}
