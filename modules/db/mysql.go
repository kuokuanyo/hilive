package db

import (
	"context"
	"database/sql"
	"errors"
	"hilive/modules/config"
	"regexp"
	"strings"
)

// Mysql struct
type Mysql struct {
	Base
}

// GetDefaultMysql 設置Mysql(struct)
func GetDefaultMysql() *Mysql {
	return &Mysql{}
}

// -----Connection的所有方法-----start

// Name Service的方法
func (db *Mysql) Name() string {
	return "mysql"
}

// InitDB 初始化資料庫引擎
func (db *Mysql) InitDB(cfg config.Database) Connection {
	// Once為唯一鎖，在需要被執行時只會被執行一次
	db.Once.Do(func() {
		if cfg.Dsn == "" {
			cfg.Dsn = cfg.User + ":" + cfg.Pwd + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" +
				cfg.Name
		}
		// 連接資料庫引擎
		sqlDB, err := sql.Open("mysql", cfg.Dsn)
		if err != nil {
			if sqlDB != nil {
				sqlDB.Close()
			}
			panic("連接資料庫引擎發生錯誤")
		} else {
			sqlDB.SetMaxIdleConns(cfg.MaxIdleCon)
			sqlDB.SetMaxOpenConns(cfg.MaxOpenCon)
			db.DB = sqlDB
		}
		if err := sqlDB.Ping(); err != nil {
			panic("初始化資料庫發生錯誤")
		}
	})
	return db
}

// Query 查詢sql資料
func (db *Mysql) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	return CommonQuery(db.DB, query, args...)
}

// Exec 執行sql命令
func (db *Mysql) Exec(query string, args ...interface{}) (sql.Result, error) {
	return CommonExec(db.DB, query, args...)
}

// GetTx 取得Tx
func (db *Mysql) GetTx() *sql.Tx {
	return CommonBeginTxWithLevel(db.Base.DB, sql.LevelDefault)
}

// -----Connection的所有方法-----end

// CommonQuery 查詢資料
func CommonQuery(db *sql.DB, query string, args ...interface{}) ([]map[string]interface{}, error) {
	// 查詢
	rs, err := db.Query(query, args...)
	if err != nil {
		panic("查詢資料發生錯誤")
	}

	// 最後要關閉*sql.rows
	defer func() {
		if rs != nil {
			_ = rs.Close()
		}
	}()

	// 取得欄位名稱
	col, colErr := rs.Columns()
	if colErr != nil {
		return nil, errors.New("取得欄位名稱發生錯誤")
	}
	// 取得欄位類別資訊
	typeVal, err := rs.ColumnTypes()
	if err != nil {
		return nil, errors.New("取得欄位類別時發生錯誤")
	}

	results := make([]map[string]interface{}, 0)

	for rs.Next() {
		// 欄位數值
		var colVar = make([]interface{}, len(col))

		r, _ := regexp.Compile(`\\((.*)\\)`)
		for i := 0; i < len(col); i++ {
			// typeName 大寫類別名稱(ex:INT)
			typeName := strings.ToUpper(r.ReplaceAllString(typeVal[i].DatabaseTypeName(), ""))
			// SetColVarType 轉換type
			SetColVarType(&colVar, i, typeName)
		}

		result := make(map[string]interface{})

		if scanErr := rs.Scan(colVar...); scanErr != nil {
			return nil, errors.New("印出資料時發生錯誤")
		}

		for j := 0; j < len(col); j++ {
			// typeName 大寫類別名稱(ex:INT)
			typeName := strings.ToUpper(r.ReplaceAllString(typeVal[j].DatabaseTypeName(), ""))
			// SetResultValue 轉後type並設置該type值(依照資料表type)
			SetResultValue(&result, col[j], colVar[j], typeName)
		}
		results = append(results, result)
	}
	if err := rs.Err(); err != nil {
		return nil, errors.New("查詢資料時發生錯誤")
	}
	return results, nil
}

// CommonExec 執行一般資料庫命令
func CommonExec(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	rs, err := db.Exec(query, args...)
	if err != nil {
		return nil, errors.New("執行資料庫命令發生錯誤")
	}
	return rs, nil
}

// CommonBeginTxWithLevel 透過LevelDefault and db取得Tx(struct)
func CommonBeginTxWithLevel(db *sql.DB, level sql.IsolationLevel) *sql.Tx {
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: level})
	if err != nil {
		panic(err)
	}
	return tx
}
