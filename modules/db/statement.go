package db

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db/sql"
	"regexp"
	"strings"
	"sync"
)

// TxFn is the transaction callback function.
type TxFn func(tx *dbsql.Tx) (error, map[string]interface{})

// SQLPool 降低內存負擔，用於需要被重複分配、回收內存的地方
var SQLPool = sync.Pool{
	New: func() interface{} {
		return &SQL{
			FilterCondition: sql.FilterCondition{
				Fields:     make([]string, 0),
				TableName:  "",
				Args:       make([]interface{}, 0),
				Wheres:     make([]sql.Where, 0),
				WhereRaws:  "",
				UpdateRaws: make([]sql.RawUpdate, 0),
				Leftjoins:  make([]sql.Join, 0),
				Order:      "",
				Group:      "",
				Limit:      "",
			},
			conn: nil,
			crud: nil,
		}
	},
}

// newSQL 取得新的SQL(struct)
func newSQL() *SQL {
	return SQLPool.Get().(*SQL)
}

// SQL 包含過濾條件、CRUD方法、Conn...等
type SQL struct {
	sql.FilterCondition // sql過濾條件
	conn                Connection
	crud                sql.CRUD // CRUD方法
	tx                  *dbsql.Tx
}

// Table 設置SQL(struct)
func Table(table string, conn Connection) *SQL {
	// Get()可以取得sync.Pool裡設置的New函式
	sqlpool := newSQL()
	sqlpool.FilterCondition.TableName = table
	sqlpool.conn = conn
	// conn.Name = mysql、mssql...
	sqlpool.crud = sql.GetCRUDByDriver(conn.Name())
	return sqlpool
}

// Table 清除SQL(struct)後設置tablenane
func (s *SQL) Table(table string) *SQL {
	s.clean()
	s.FilterCondition.TableName = table
	return s
}

// TableAndCleanData 設置SQL(struct)並且還有清除過濾條件的設置
func TableAndCleanData(table string, conn Connection) *SQL {
	// Get()可以取得sync.Pool裡設置的New函式
	sqlpool := newSQL()
	sqlpool.conn = conn
	sqlpool.crud = sql.GetCRUDByDriver(conn.Name())
	sqlpool.clean()
	sqlpool.FilterCondition.TableName = table
	return sqlpool
}

// SetConnectionAndCRUD 設置SQL(struct)的Connection、CRUD
func SetConnectionAndCRUD(conn Connection) *SQL {
	sqlpool := newSQL()
	sqlpool.conn = conn
	sqlpool.crud = sql.GetCRUDByDriver(conn.Name())
	return sqlpool
}

// ShowColumns 取得所有欄位資訊
func (s *SQL) ShowColumns() ([]map[string]interface{}, error) {
	defer RecycleSQL(s)
	return s.conn.Query(s.crud.ShowColumns(s.FilterCondition.TableName))
}

// Insert 插入資料後回傳插入值的id
func (s *SQL) Insert(values sql.Value) (int64, error) {
	defer RecycleSQL(s)

	s.FilterCondition.Values = values
	s.crud.Insert(&s.FilterCondition)

	res, err := s.conn.Exec(s.FilterCondition.Statement, s.FilterCondition.Args...)
	if err != nil {
		return 0, err
	}
	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return 0, errors.New("沒有影響任何資料")
	}
	return res.LastInsertId()
}

// Update 執行更新命令
func (s *SQL) Update(values sql.Value) (int64, error) {
	defer RecycleSQL(s)

	s.FilterCondition.Values = values
	s.crud.Update(&s.FilterCondition)

	res, err := s.conn.Exec(s.FilterCondition.Statement, s.FilterCondition.Args...)
	if err != nil {
		return 0, err
	}
	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return 0, errors.New("沒有影響任何資料")
	}
	return res.LastInsertId()
}

// Delete 刪除資料
func (s *SQL) Delete() error {
	defer RecycleSQL(s)
	s.crud.Delete(&s.FilterCondition)

	res, err := s.conn.Exec(s.FilterCondition.Statement, s.FilterCondition.Args...)
	if err != nil {
		return errors.New("刪除資料發生錯誤")
	}
	if affectRow, _ := res.RowsAffected(); affectRow < 1 {
		return errors.New("沒有影響任何資料")
	}
	return nil
}

// All 查詢所有資料
func (s *SQL) All() ([]map[string]interface{}, error) {
	defer RecycleSQL(s)
	// 取得所有資料
	s.crud.Select(&s.FilterCondition)

	return s.conn.Query(s.FilterCondition.Statement, s.FilterCondition.Args...)
}

// FindByID 藉由id尋找資料
func (s *SQL) FindByID(arg interface{}) (map[string]interface{}, error) {
	return s.Where("id", "=", arg).First()
}

// First 查詢單筆資料
func (s *SQL) First() (map[string]interface{}, error) {
	// 執行結束後清空
	defer RecycleSQL(s)

	// 取得所有資料
	s.crud.Select(&s.FilterCondition)

	// 查詢
	results, err := s.conn.Query(s.FilterCondition.Statement, s.FilterCondition.Args...)
	if err != nil {
		return nil, errors.New("查詢資料發生錯誤")
	}
	if len(results) < 1 {
		return nil, errors.New("查無此資料")
	}

	return results[0], nil
}

// Select 處理欄位
func (s *SQL) Select(fields ...string) *SQL {
	s.FilterCondition.Fields = fields
	s.FilterCondition.Functions = make([]string, len(fields))

	reg, _ := regexp.Compile("(.*?)\\((.*?)\\)")
	for k, field := range fields {
		res := reg.FindAllStringSubmatch(field, -1)
		if len(res) > 0 && len(res[0]) > 2 {
			s.FilterCondition.Functions[k] = res[0][1]
			s.FilterCondition.Fields[k] = res[0][2]
		}
	}
	return s
}

// Where 設置sql where條件
func (s *SQL) Where(field, operation string, arg interface{}) *SQL {
	s.FilterCondition.Wheres = append(s.FilterCondition.Wheres, sql.Where{
		Field:     field,
		Operation: operation,
		Value:     "?",
	})
	s.FilterCondition.Args = append(s.FilterCondition.Args, arg)
	return s
}

// WhereRaw 設置sql whereraw條件
func (s *SQL) WhereRaw(raw string, args ...interface{}) *SQL {
	s.FilterCondition.WhereRaws = raw
	s.FilterCondition.Args = append(s.FilterCondition.Args, args...)
	return s
}

// WhereIn 設置where in 多個數值語法
func (s *SQL) WhereIn(field string, arg []interface{}) *SQL {
	if len(arg) == 0 {
		panic("wherein語法參數不能為空")
	}
	s.FilterCondition.Wheres = append(s.FilterCondition.Wheres,
		sql.Where{
			Field:     field,
			Operation: "in",
			Value:     "(" + strings.Repeat("?,", len(arg)-1) + "?)",
		})
	s.FilterCondition.Args = append(s.FilterCondition.Args, arg...)
	return s
}

// LeftJoin 添加Join語法
func (s *SQL) LeftJoin(table, fieldA, operation, fieldB string) *SQL {
	s.FilterCondition.Leftjoins = append(s.FilterCondition.Leftjoins,
		sql.Join{
			FieldA:    fieldA,
			FieldB:    fieldB,
			Table:     table,
			Operation: operation,
		})
	return s
}

// OrderBy 資料排列順序
func (s *SQL) OrderBy(fields ...string) *SQL {
	if len(fields) == 0 {
		panic("orderBy參數設置錯誤")
	}
	for i := 0; i < len(fields); i++ {
		if i == len(fields)-2 {
			s.FilterCondition.Order += " " + fields[i] + " " + fields[i+1]
			return s
		}
		s.FilterCondition.Order += " " + fields[i] + " and "
	}
	return s
}

// RecycleSQL 清空SQL(struct)所有資訊設置至SQLPool中
func RecycleSQL(s *SQL) {
	s.clean()
	s.conn = nil
	s.crud = nil
	// 將清空的SQL設置於SQLPool
	SQLPool.Put(s)
}

// clean 清空過濾條件
func (s *SQL) clean() {
	s.FilterCondition.Functions = make([]string, 0)
	s.FilterCondition.Group = ""
	s.FilterCondition.Values = make(map[string]interface{})
	s.FilterCondition.Fields = make([]string, 0)
	s.FilterCondition.TableName = ""
	s.FilterCondition.Wheres = make([]sql.Where, 0)
	s.FilterCondition.Leftjoins = make([]sql.Join, 0)
	s.FilterCondition.Args = make([]interface{}, 0)
	s.FilterCondition.Order = ""
	s.FilterCondition.Offset = ""
	s.FilterCondition.Limit = ""
	s.FilterCondition.WhereRaws = ""
	s.FilterCondition.UpdateRaws = make([]sql.RawUpdate, 0)
	s.FilterCondition.Statement = ""
}

// WithTx 設置至SQL.tx
func (s *SQL) WithTx(tx *dbsql.Tx) *SQL {
	s.tx = tx
	return s
}

// WithTransaction 取得Tx，持續並行commit、rollback
func (s *SQL) WithTransaction(fn TxFn) (res map[string]interface{}, err error) {
	// 取得Tx
	tx := s.conn.GetTx()

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()
	err, res = fn(tx)
	return
}
