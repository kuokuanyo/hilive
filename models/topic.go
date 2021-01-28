package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// TopicModel 主題牆的資料表欄位
type TopicModel struct {
	Base `json:"-"`

	ID         int64  `json:"id"`
	ActivityID string `json:"activity_id"`
	Background string `json:"background"`
}

// DefaultTopicModel 預設TopicModel
func DefaultTopicModel() TopicModel {
	return TopicModel{Base: Base{TableName: "activity_set_topic"}}
}

// GetTopicModelAndID 設置TopicModel與ID
func GetTopicModelAndID(tablename, id string) TopicModel {
	idInt, _ := strconv.Atoi(id)
	return TopicModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (t TopicModel) SetConn(conn db.Connection) TopicModel {
	t.Conn = conn
	return t
}

// SetTx 設置Tx
func (t TopicModel) SetTx(tx *dbsql.Tx) TopicModel {
	t.Base.Tx = tx
	return t
}

// AddTopic 增加主題牆資料
func (t TopicModel) AddTopic(activityid, background string) (TopicModel, error) {
	// 檢查是否有該活動
	_, err := t.SetTx(t.Base.Tx).Table("activity").Where("activity_id", "=", activityid).First()
	if err != nil {
		return t, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := t.SetTx(t.Base.Tx).Table(t.TableName).Insert(sql.Value{
		"activity_id": activityid,
		"background":  background,
	})

	t.ID = id
	t.ActivityID = activityid
	t.Background = background

	return t, err
}

// UpdateTopic 更新主題牆資料
func (t TopicModel) UpdateTopic(activityid, background string) (int64, error) {
	// 檢查是否有該活動
	_, err := t.SetTx(t.Base.Tx).Table("activity").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	fieldValues := sql.Value{
		"activity_id": activityid,
		"background":  background,
	}

	return t.SetTx(t.Tx).Table(t.Base.TableName).
		Where("id", "=", t.ID).Update(fieldValues)
}

// IsActivityExist 檢查活動是否已經設置主題牆
func (t TopicModel) IsActivityExist(activityid, id string) bool {
	if id == "" {
		check, _ := t.Table(t.TableName).Where("activity_id", "=", activityid).First()
		return check != nil
	}
	check, _ := t.Table(t.TableName).
		Where("activity_id", "=", activityid).
		Where("id", "!=", id).
		First()
	return check != nil
}
