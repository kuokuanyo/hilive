package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// PrizeStaffModel 中獎人員資料表欄位
type PrizeStaffModel struct {
	Base `json:"-"`

	ID             int64  `json:"id"`
	UserID         string `json:"user_id"`
	ActivityID     string `json:"activity_id"`
	GameID         string `json:"game_id"`
	WinTime        string `json:"win_time"`
	PrizeName      string `json:"prize_name"`
	RedeemMethod   string `json:"redeem_method"`
	RedeemStatus   string `json:"redeem_status"`
	RedeemPassword string `json:"redeem_password"`
}

// DefaultPrizeStaffModel 預設PrizeStaffModel
func DefaultPrizeStaffModel() PrizeStaffModel {
	return PrizeStaffModel{Base: Base{TableName: "activity_prize_staff"}}
}

// GetPrizeStaffModelAndID 設置PrizeStaffModel與ID
func GetPrizeStaffModelAndID(tablename, id string) PrizeStaffModel {
	idInt, _ := strconv.Atoi(id)
	return PrizeStaffModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m PrizeStaffModel) SetConn(conn db.Connection) PrizeStaffModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m PrizeStaffModel) SetTx(tx *dbsql.Tx) PrizeStaffModel {
	m.Base.Tx = tx
	return m
}

// AddStaff 增加遊戲人員
func (m PrizeStaffModel) AddStaff(userid, activityid, gameid, time,
	name, method, status, password string) (PrizeStaffModel, error) {
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
		"user_id":         userid,
		"activity_id":     activityid,
		"game_id":         gameid,
		"win_time":        time,
		"prize_name":      name,
		"redeem_method":   method,
		"redeem_status":   status,
		"redeem_password": password,
	})

	m.ID = id
	m.UserID = userid
	m.ActivityID = activityid
	m.GameID = gameid
	m.WinTime = time
	m.PrizeName = name
	m.RedeemMethod = method
	m.RedeemStatus = status
	m.RedeemPassword = password
	return m, err
}

// UpdateStaff 更新中獎人員資料
func (m PrizeStaffModel) UpdateStaff(userid, activityid, gameid, time,
	name, method, status, password string) (int64, error) {
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
		"user_id":         userid,
		"activity_id":     activityid,
		"game_id":         gameid,
		"win_time":        time,
		"prize_name":      name,
		"redeem_method":   method,
		"redeem_status":   status,
		"redeem_password": password,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

