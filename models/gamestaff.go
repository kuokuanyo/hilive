package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// GameStaffModel 參加遊戲人員資料表欄位
type GameStaffModel struct {
	Base `json:"-"`

	ID          int64  `json:"id"`
	UserID      string `json:"user_id"`
	ActivityID  string `json:"activity_id"`
	GameID      string `json:"game_id"`
	ApplyStatus string `json:"apply_status"`
	ApplyTime   string `json:"apply_time"`
}

// DefaultGameStaffModel 預設GameStaffModel
func DefaultGameStaffModel() GameStaffModel {
	return GameStaffModel{Base: Base{TableName: "activity_apply_game_staff"}}
}

// GetGameStaffModelAndID 設置GameStaffModel與ID
func GetGameStaffModelAndID(tablename, id string) GameStaffModel {
	idInt, _ := strconv.Atoi(id)
	return GameStaffModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m GameStaffModel) SetConn(conn db.Connection) GameStaffModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m GameStaffModel) SetTx(tx *dbsql.Tx) GameStaffModel {
	m.Base.Tx = tx
	return m
}

// AddStaff 增加遊戲人員
func (m GameStaffModel) AddStaff(userid, activityid, gameid, status string) (GameStaffModel, error) {
	// 檢查是否有該活動遊戲
	_, err := m.SetTx(m.Base.Tx).Table("activity_all_game").
		Where("activity_id", "=", activityid).
		Where("game_id", "=", gameid).First()
	if err != nil {
		return m, errors.New("查詢不到活動ID及遊戲ID，請重新設置")
	}

	// 檢查是否有該用戶
	_, err = m.SetTx(m.Base.Tx).Table("users").Where("userid", "=", userid).First()
	if err != nil {
		return m, errors.New("查詢不到該用戶ID，請輸入正確用戶ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"user_id":      userid,
		"activity_id":  activityid,
		"game_id":      gameid,
		"apply_status": status,
	})

	m.ID = id
	m.UserID = userid
	m.ActivityID = activityid
	m.GameID = gameid
	m.ApplyStatus = status
	return m, err
}

// UpdateStaff 更新遊戲人員資料
func (m GameStaffModel) UpdateStaff(userid, activityid, gameid, status string) (int64, error) {
	// 檢查是否有該活動遊戲
	_, err := m.SetTx(m.Base.Tx).Table("activity_all_game").
		Where("activity_id", "=", activityid).
		Where("game_id", "=", gameid).First()
	if err != nil {
		return 0, errors.New("查詢不到活動ID及遊戲ID，請重新設置")
	}

	// 檢查是否有該用戶
	_, err = m.SetTx(m.Base.Tx).Table("users").Where("userid", "=", userid).First()
	if err != nil {
		return 0, errors.New("查詢不到該用戶ID，請輸入正確用戶ID")
	}

	fieldValues := sql.Value{
		"user_id":      userid,
		"activity_id":  activityid,
		"game_id":      gameid,
		"apply_status": status,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsStaffExist 檢查人員是否已經參加遊戲
func (m GameStaffModel) IsStaffExist(userid, gameid, id string) bool {
	if id == "" {
		check, _ := m.Table(m.TableName).Where("user_id", "=", userid).
			Where("game_id", "=", gameid).First()
		return check != nil
	}
	check, _ := m.Table(m.TableName).
		Where("user_id", "=", userid).
		Where("game_id", "=", gameid).
		Where("id", "!=", id).
		First()
	return check != nil
}
