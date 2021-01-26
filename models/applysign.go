package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// ApplysignModel 活動資料的資料表欄位
type ApplysignModel struct {
	Base `json:"-"`

	ID         int64  `json:"id"`
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
	UserName   string `json:"user_name"`
	UserAvater string `json:"user_avater"`
	Status     int    `json:"status"`
	SignTime   int    `json:"sign_time"`
}

// DefaultApplysignModel 預設ApplysignModel
func DefaultApplysignModel() ApplysignModel {
	return ApplysignModel{Base: Base{TableName: "activity_applysign"}}
}

// GetApplysignModelAndID 設置ApplysignModel與ID
func GetApplysignModelAndID(tablename, id string) ApplysignModel {
	idInt, _ := strconv.Atoi(id)
	return ApplysignModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (s ApplysignModel) SetConn(conn db.Connection) ApplysignModel {
	s.Conn = conn
	return s
}

// SetTx 設置Tx
func (s ApplysignModel) SetTx(tx *dbsql.Tx) ApplysignModel {
	s.Base.Tx = tx
	return s
}

// AddApplysign 增加報名簽到資料
func (s ApplysignModel) AddApplysign(userid, activityid, username, avater string, status int) (ApplysignModel, error) {
	// 檢查是否有該活動
	_, err := s.SetTx(s.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return s, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}
	// 檢查是否有該用戶
	_, err = s.SetTx(s.Base.Tx).Table("users").Select("id").Where("userid", "=", userid).First()
	if err != nil {
		return s, errors.New("查詢不到該用戶ID，請輸入正確用戶ID")
	}

	id, err := s.SetTx(s.Base.Tx).Table(s.TableName).Insert(sql.Value{
		"user_id":     userid,
		"activity_id": activityid,
		"user_name":   username,
		"user_avater": avater,
		"status":      status,
	})

	s.ID = id
	s.ActivityID = activityid
	s.UserID = userid
	s.UserName = username
	s.UserAvater = avater
	s.Status = status
	return s, err
}

// UpdateActivityApplysign 更新報名簽到資料
func (s ApplysignModel) UpdateActivityApplysign(userid, activityid, username, avater string, status int) (int64, error) {
	model, err := s.SetTx(s.Base.Tx).Table(s.Base.TableName).Where("id", "=", s.ID).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動")
	}
	if model["activity_id"] != activityid {
		return 0, errors.New("資料中的活動ID不符合，無法更新資料")
	}
	if model["user_id"] != userid {
		return 0, errors.New("資料中的使用者ID不符合，無法更新資料")
	}

	fieldValues := sql.Value{
		"user_name":   username,
		"user_avater": avater,
		"status":      status,
	}

	return s.SetTx(s.Tx).Table(s.Base.TableName).
		Where("id", "=", s.ID).Update(fieldValues)
}

// IsSignExist 檢查是否已報到簽名
func (s ApplysignModel) IsSignExist(activityid, userid, id string) bool {
	if id == "" {
		check, _ := s.Table(s.TableName).Where("activity_id", "=", activityid).
			Where("user_id", "=", userid).First()
		return check != nil
	}
	check, _ := s.Table(s.TableName).
		Where("activity_id", "=", activityid).
		Where("user_id", "=", userid).Where("id", "!=", id).
		First()
	return check != nil
}
