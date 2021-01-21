package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// QuestionModel 提問牆的資料表欄位
type QuestionModel struct {
	Base `json:"-"`

	ID           int64  `json:"id"`
	ActivityID   string `json:"activity_id"`
	MessageCheck string `json:"message_check"`
	Anonymous    string `json:"anonymous"`
	HideAnswered string `json:"hide_answered"`
	Qrcode       string `json:"qrcode"`
}

// DefaultQuestionModel 預設MessageModel
func DefaultQuestionModel() QuestionModel {
	return QuestionModel{Base: Base{TableName: "activity_set_question"}}
}

// GetQuestionModelAndID 設置QuestionModel與ID
func GetQuestionModelAndID(tablename, id string) QuestionModel {
	idInt, _ := strconv.Atoi(id)
	return QuestionModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (q QuestionModel) SetConn(conn db.Connection) QuestionModel {
	q.Conn = conn
	return q
}

// SetTx 設置Tx
func (q QuestionModel) SetTx(tx *dbsql.Tx) QuestionModel {
	q.Base.Tx = tx
	return q
}

// AddQuestion 增加提問牆資料
func (q QuestionModel) AddQuestion(activityid, check, anonymous, answer, qrcode string) (QuestionModel, error) {
	// 檢查是否有該活動
	_, err := q.SetTx(q.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return q, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := q.SetTx(q.Base.Tx).Table(q.TableName).Insert(sql.Value{
		"activity_id":   activityid,
		"message_check": check,
		"anonymous":     anonymous,
		"hide_answered": answer,
		"qrcode":        qrcode,
	})

	q.ID = id
	q.ActivityID = activityid
	q.MessageCheck = check
	q.Anonymous = anonymous
	q.HideAnswered = answer
	q.Qrcode = qrcode

	return q, err
}

// UpdateQuestion 更新提問牆資料
func (q QuestionModel) UpdateQuestion(activityid, check, anonymous, answer, qrcode string) (int64, error) {
	_, err := q.SetTx(q.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	fieldValues := sql.Value{
		"activity_id":   activityid,
		"message_check": check,
		"anonymous":     anonymous,
		"hide_answered": answer,
		"qrcode":        qrcode,
	}

	return q.SetTx(q.Tx).Table(q.Base.TableName).
		Where("id", "=", q.ID).Update(fieldValues)
}

// IsActivityExist 檢查活動是否已經設置提問牆
func (q QuestionModel) IsActivityExist(activityid, id string) bool {
	if id == "" {
		check, _ := q.Table(q.TableName).Where("activity_id", "=", activityid).First()
		return check != nil
	}
	check, _ := q.Table(q.TableName).
		Where("activity_id", "=", activityid).
		Where("id", "!=", id).
		First()
	return check != nil
}
