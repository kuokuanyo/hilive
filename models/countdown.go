package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// CountdownModel 倒數計時的資料表欄位
type CountdownModel struct {
	Base `json:"-"`

	ID          int64  `json:"id"`
	ActivityID  string `json:"activity_id"`
	Second      int    `json:"second"`
	IndexURL    string `json:"index_url"`
	AvatarShape string `json:"avatar_shape"`
}

// DefaultCountdownModel 預設CountdownModel
func DefaultCountdownModel() CountdownModel {
	return CountdownModel{Base: Base{TableName: "activity_set_countdown"}}
}

// GetCountdownModelAndID 設置CountdownModel與ID
func GetCountdownModelAndID(tablename, id string) CountdownModel {
	idInt, _ := strconv.Atoi(id)
	return CountdownModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m CountdownModel) SetConn(conn db.Connection) CountdownModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m CountdownModel) SetTx(tx *dbsql.Tx) CountdownModel {
	m.Base.Tx = tx
	return m
}

// AddCountdown 增加倒數計時資料
func (m CountdownModel) AddCountdown(activityid, url, shape string, second int) (CountdownModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":  activityid,
		"second":       second,
		"index_url":    url,
		"avatar_shape": shape,
	})

	m.ID = id
	m.ActivityID = activityid
	m.Second = second
	m.IndexURL = url
	m.AvatarShape = shape
	return m, err
}

// UpdateCountdown 更新倒數計時資料
func (m CountdownModel) UpdateCountdown(activityid, url, shape string, second int) (int64, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}
	fieldValues := sql.Value{
		"activity_id":  activityid,
		"second":       second,
		"index_url":    url,
		"avatar_shape": shape,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsActivityExist 檢查活動是否已經設置倒數計時
func (m CountdownModel) IsActivityExist(activityid, id string) bool {
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
