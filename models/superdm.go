package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// SuperdmModel 超級彈幕的資料表欄位
type SuperdmModel struct {
	Base `json:"-"`

	ID                int64  `json:"id"`
	ActivityID        string `json:"activity_id"`
	MessageCheck      string `json:"message_check"`
	EyeCatchingPrice  int    `json:"eye_catching_price"`
	LargeDanmuPrice   int    `json:"large_danmu_price"`
	StatusUpdatePrice int    `json:"status_update_price"`
	PictureDanmuPrice int    `json:"picture_danmu_price"`
}

// DefaultSuperdmModel 預設SuperdmModel
func DefaultSuperdmModel() SuperdmModel {
	return SuperdmModel{Base: Base{TableName: "activity_set_superdm"}}
}

// GetSuperdmModelAndID 設置SuperdmModel與ID
func GetSuperdmModelAndID(tablename, id string) SuperdmModel {
	idInt, _ := strconv.Atoi(id)
	return SuperdmModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m SuperdmModel) SetConn(conn db.Connection) SuperdmModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m SuperdmModel) SetTx(tx *dbsql.Tx) SuperdmModel {
	m.Base.Tx = tx
	return m
}

// AddSuperdm 增加超級彈幕資料
func (m SuperdmModel) AddSuperdm(activityid, check string, eye, large, statusUpdate, picture int) (SuperdmModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":         activityid,
		"message_check":       check,
		"eye_catching_price":  eye,
		"large_danmu_price":   large,
		"status_update_price": statusUpdate,
		"picture_danmu_price": picture,
	})

	m.ID = id
	m.ActivityID = activityid
	m.MessageCheck = check
	m.EyeCatchingPrice = eye
	m.LargeDanmuPrice = large
	m.StatusUpdatePrice = statusUpdate
	m.PictureDanmuPrice = picture
	return m, err
}

// UpdateSuperdm 更新超級彈幕資料
func (m SuperdmModel) UpdateSuperdm(activityid, check string, eye, large, statusUpdate, picture int) (int64, error) {
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	fieldValues := sql.Value{
		"activity_id":         activityid,
		"message_check":       check,
		"eye_catching_price":  eye,
		"large_danmu_price":   large,
		"status_update_price": statusUpdate,
		"picture_danmu_price": picture,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsActivityExist 檢查活動是否已經設置超級彈幕
func (m SuperdmModel) IsActivityExist(activityid, id string) bool {
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
