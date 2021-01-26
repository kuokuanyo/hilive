package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// SignModel 簽到牆的資料表欄位
type SignModel struct {
	Base `json:"-"`

	ID         int64  `json:"id"`
	ActivityID string `json:"activity_id"`
	Display    string `json:"Display"`
	Background string `json:"background"`
}

// DefaultSignModel 預設SignModel
func DefaultSignModel() SignModel {
	return SignModel{Base: Base{TableName: "activity_set_sign"}}
}

// GetSignModelAndID 設置SignModel與ID
func GetSignModelAndID(tablename, id string) SignModel {
	idInt, _ := strconv.Atoi(id)
	return SignModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m SignModel) SetConn(conn db.Connection) SignModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m SignModel) SetTx(tx *dbsql.Tx) SignModel {
	m.Base.Tx = tx
	return m
}

// AddSign 增加簽到牆資料
func (m SignModel) AddSign(activityid, display, background string) (SignModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id": activityid,
		"display":     display,
		"background":  background,
	})

	m.ID = id
	m.ActivityID = activityid
	m.Display = display
	m.Background = background
	return m, err
}

// UpdateSign 更新簽到牆資料
func (m SignModel) UpdateSign(activityid, display, background string) (int64, error) {
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	fieldValues := sql.Value{
		"activity_id": activityid,
		"display":     display,
		"background":  background,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsActivityExist 檢查活動是否已經設置簽到牆
func (m SignModel) IsActivityExist(activityid, id string) bool {
	if id == "" {
		check, _ := m.Table(m.TableName).Where("activity_id", "=", activityid).First()
		return check != nil
	}
	check, _ := m.Table(m.TableName).
		Where("activity_id", "=", activityid).
		Where("id", "!=", id).
		First()
	return check != nil
}
