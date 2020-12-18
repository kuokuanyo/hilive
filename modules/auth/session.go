package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"hilive/modules/config"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// DBDriver 紀錄session
type DBDriver struct {
	conn      db.Connection
	tableName string
}

// Session contains info of session
type Session struct {
	Expires time.Duration
	Cookie  string
	Sid     string
	Values  map[string]interface{}
	Driver  PersistenceDriver
	Context *gin.Context
}

// PersistenceDriver 持久性驅動
type PersistenceDriver interface {
	// Load 取得session資料表的cookie_values欄位(ex:{"user_id":1})
	Load(string) (map[string]interface{}, error)
	// Update 建立或更新session資料表資料
	Update(string, map[string]interface{}) error
}

// table 設置SQL(struct)
func (driver *DBDriver) table() *db.SQL {
	return db.Table(driver.tableName, driver.conn)
}

// DefaultDBDriver 預設DBDriver
func DefaultDBDriver(conn db.Connection) *DBDriver {
	return &DBDriver{
		conn:      conn,
		tableName: "session",
	}
}

// InitSession 初始化Session並取得session資料表的cookie_values欄位(ex:{"user_id":1})
func InitSession(ctx *gin.Context, conn db.Connection) (*Session, error) {
	session := new(Session)
	session.Expires = time.Second * time.Duration(config.GetSessionLifeTime())
	session.Cookie = "session"
	session.Driver = DefaultDBDriver(conn)
	session.Values = make(map[string]interface{})

	if cookie, err := ctx.Request.Cookie(session.Cookie); err == nil && cookie.Value != "" {
		session.Sid = cookie.Value
		// Load 取得session資料表的cookie_values欄位(ex:{"user_id":1})
		valueFromDriver, err := session.Driver.Load(cookie.Value)
		if err != nil {
			return nil, err
		}
		if len(valueFromDriver) > 0 {
			session.Values = valueFromDriver
		}
	} else {
		uid, _ := uuid.NewV4()
		session.Sid = uid.String()
	}
	session.Context = ctx
	return session, nil
}

// deleteOverdueSession 刪除超過時間的ccokie
func (driver *DBDriver) deleteOverdueSession() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
	}()

	var (
		duration   = strconv.Itoa(config.GetSessionLifeTime() + 1000)
		driverName = config.GetDatabaseDriver()
		raw        = ``
	)
	if driverName == "mysql" {
		raw = `unix_timestamp(created_at) < unix_timestamp() - ` + duration
	}
	if raw != "" {
		driver.table().WhereRaw(raw).Delete()
	}
}

// -----PersistenceDriver的方法-----start

// Load 取得session資料表的cookie_values欄位(ex:{"user_id":1})
func (driver *DBDriver) Load(sid string) (map[string]interface{}, error) {
	sesModel, err := driver.table().Where("sid", "=", sid).First()
	if err != nil {
		if fmt.Sprintf("%s", err) != "查無此資料" {
			return nil, errors.New("查詢session資料表發生錯誤")
		}
	}
	if sesModel == nil {
		return map[string]interface{}{}, nil
	}

	var values map[string]interface{}
	err = json.Unmarshal([]byte(sesModel["cookie_values"].(string)), &values)
	return values, err
}

// Update 建立或更新session資料表資料
func (driver *DBDriver) Update(sid string, values map[string]interface{}) error {
	// 刪除超過時間的ccokie
	go driver.deleteOverdueSession()

	if sid != "" {
		// 假設該資料沒有values值則刪除
		if len(values) == 0 {
			err := driver.table().Where("sid", "=", sid).Delete()
			if fmt.Sprintf("%s", err) != "沒有影響任何資料" {
				return errors.New("刪除session資料表資料發生錯誤")
			}
		}

		// json編碼
		valueByte, err := json.Marshal(values)
		if err != nil {
			return errors.New("json編碼發生錯誤")
		}
		sesValue := string(valueByte)

		sesModel, _ := driver.table().Where("sid", "=", sid).First()
		if sesModel == nil {
			if !config.GetNoLimitLoginIP() {
				err = driver.table().Where("cookie_values", "=", sesValue).Delete()
				if err != nil {
					if fmt.Sprintf("%s", err) != "沒有影響任何資料" {
						return errors.New("刪除session資料表資料發生錯誤")
					}
				}
			}

			_, err = driver.table().Insert(sql.Value{
				"cookie_values": sesValue,
				"sid":           sid,
			})
			if err != nil {
				return errors.New("插入資料至session資料表發生錯誤")
			}
		} else {
			// 更新資料
			_, err := driver.table().Where("sid", "=", sid).
				Update(sql.Value{
					"cookie_values": sesValue,
				})
			if err != nil {
				if fmt.Sprintf("%s", err) != "沒有影響任何資料" {
					return errors.New("更新session資料表發生錯誤")
				}
			}
		}
	}
	return nil
}

// -----PersistenceDriver的方法-----end
