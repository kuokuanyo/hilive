package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// ScheduleModel 活動行程資料表欄位
type ScheduleModel struct {
	Base `json:"-"`

	ID              int64  `json:"id"`
	ActivityID      string `json:"activity_id"`
	ScheduleName    string `json:"schedule_name"`
	ScheduleContent string `json:"schedule_content"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
}

// DefaultScheduleModel ScheduleModel
func DefaultScheduleModel() ScheduleModel {
	return ScheduleModel{Base: Base{TableName: "activity_schedule"}}
}

// GetScheduleModelAndID 設置ScheduleModel與ID
func GetScheduleModelAndID(tablename, id string) ScheduleModel {
	idInt, _ := strconv.Atoi(id)
	return ScheduleModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (s ScheduleModel) SetConn(conn db.Connection) ScheduleModel {
	s.Conn = conn
	return s
}

// SetTx 設置Tx
func (s ScheduleModel) SetTx(tx *dbsql.Tx) ScheduleModel {
	s.Base.Tx = tx
	return s
}

// AddActivitySchedule 增加活動行程資料
func (s ScheduleModel) AddActivitySchedule(activityid, name, content, start, end string) (ScheduleModel, error) {
	// 檢查是否有該活動
	_, err := s.SetTx(s.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return s, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := s.SetTx(s.Base.Tx).Table(s.TableName).Insert(sql.Value{
		"activity_id":      activityid,
		"schedule_name":    name,
		"schedule_content": content,
		"start_time":       start,
		"end_time":         end,
	})

	s.ID = id
	s.ActivityID = activityid
	s.ScheduleName = name
	s.ScheduleContent = content
	s.StartTime = start
	s.EndTime = end
	return s, err
}

// UpdateActivitySchedule 更新活動行程資料
func (s ScheduleModel) UpdateActivitySchedule(activityid, name, content, start, end string) (int64, error) {
	_, err := s.SetTx(s.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	fieldValues := sql.Value{
		"activity_id":      activityid,
		"schedule_name":    name,
		"schedule_content": content,
		"start_time":       start,
		"end_time":         end,
	}

	return s.SetTx(s.Tx).Table(s.Base.TableName).
		Where("id", "=", s.ID).Update(fieldValues)
}