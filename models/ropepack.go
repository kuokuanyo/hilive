package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// RopepackModel 套紅包的資料表欄位
type RopepackModel struct {
	Base `json:"-"`

	ID          int64  `json:"id"`
	ActivityID  string `json:"activity_id"`
	GameID      string `json:"game_id"`
	Title       string `json:"title"`
	Percent     int    `json:"percent"`
	Second      int    `json:"second"`
	AllowRepeat string `json:"allow_repeat"`
	GameStatus  string `json:"game_status"`
}

// DefaultRopepackModel 預設RopepackModel
func DefaultRopepackModel() RopepackModel {
	return RopepackModel{Base: Base{TableName: "activity_set_ropepack"}}
}

// GetRopepackModelAndID 設置RopepackModel與ID
func GetRopepackModelAndID(tablename, id string) RopepackModel {
	idInt, _ := strconv.Atoi(id)
	return RopepackModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m RopepackModel) SetConn(conn db.Connection) RopepackModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m RopepackModel) SetTx(tx *dbsql.Tx) RopepackModel {
	m.Base.Tx = tx
	return m
}

// AddRopepack 增加套紅包獎品
func (m RopepackModel) AddRopepack(activityid, gameid, title, repeat,
	status string, second, percent int) (RopepackModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}
	if percent > 100 || percent < 0 {
		return m, errors.New("中獎機率設置區間為1~100，請重新設置")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":  activityid,
		"game_id":      gameid,
		"title":        title,
		"second":       second,
		"allow_repeat": repeat,
		"game_status":  status,
		"percent":      percent,
	})
	if err != nil {
		return m, err
	}

	// 將小遊戲資料加入activity_all_game
	_, err = m.SetTx(m.Base.Tx).Table("activity_all_game").Insert(sql.Value{
		"activity_id": activityid,
		"game_id":     gameid,
		"game":        "套紅包",
	})

	m.ID = id
	m.ActivityID = activityid
	m.GameID = gameid
	m.Title = title
	m.Second = second
	m.AllowRepeat = repeat
	m.GameStatus = status
	m.Percent = percent
	return m, err
}

// UpdateRopepack 更新套紅包資料
func (m RopepackModel) UpdateRopepack(activityid, gameid, title, repeat,
	status string, second, percent int) (int64, error) {
	// 檢查是否有該活動遊戲
	_, err := m.SetTx(m.Base.Tx).Table("activity_all_game").
		Where("activity_id", "=", activityid).
		Where("game_id", "=", gameid).
		First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID及遊戲ID，請輸重新設置")
	}

	if percent > 100 || percent < 0 {
		return 0, errors.New("中獎機率設置區間為1~100，請重新設置")
	}

	fieldValues := sql.Value{
		"activity_id":  activityid,
		"game_id":      gameid,
		"title":        title,
		"second":       second,
		"allow_repeat": repeat,
		"game_status":  status,
		"percent":      percent,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}
