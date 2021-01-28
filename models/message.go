package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
	"strings"
)

// MessageModel 訊息牆的資料表欄位
type MessageModel struct {
	Base `json:"-"`

	ID                  int64    `json:"id"`
	ActivityID          string   `json:"activity_id"`
	PictureMessage      string   `json:"picture_message"`
	PictureAuto         string   `json:"picture_auto"`
	RefreshSecond       int      `json:"refresh_second"`
	PreventStatusUpdate string   `json:"prevent_status_update"`
	Message             []string `json:"message"`
}

// DefaultMessageModel 預設MessageModel
func DefaultMessageModel() MessageModel {
	return MessageModel{Base: Base{TableName: "activity_set_message"}}
}

// GetMessageModelAndID 設置MessageModel與ID
func GetMessageModelAndID(tablename, id string) MessageModel {
	idInt, _ := strconv.Atoi(id)
	return MessageModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m MessageModel) SetConn(conn db.Connection) MessageModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m MessageModel) SetTx(tx *dbsql.Tx) MessageModel {
	m.Base.Tx = tx
	return m
}

// AddMessage 增加訊息牆資料
func (m MessageModel) AddMessage(activityid, pictureMessage, auto, prevent, message string, second int) (MessageModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":           activityid,
		"picture_message":       pictureMessage,
		"picture_auto":          auto,
		"refresh_second":        second,
		"prevent_status_update": prevent,
		"message":               message,
	})

	m.ID = id
	m.ActivityID = activityid
	m.PictureMessage = pictureMessage
	m.PictureAuto = auto
	m.RefreshSecond = second
	m.PreventStatusUpdate = prevent
	m.Message = strings.Split(message, "\n")

	return m, err
}

// UpdateMessage 更新訊息設置資料
func (m MessageModel) UpdateMessage(activityid, pictureMessage, auto, prevent, message string, second int) (int64, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	fieldValues := sql.Value{
		"activity_id":           activityid,
		"picture_message":       pictureMessage,
		"picture_auto":          auto,
		"refresh_second":        second,
		"prevent_status_update": prevent,
		"message":               message,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsActivityExist 檢查活動是否已經設置訊息
func (m MessageModel) IsActivityExist(activityid, id string) bool {
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
