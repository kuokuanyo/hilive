package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// DanmuModel 彈幕的資料表欄位
type DanmuModel struct {
	Base `json:"-"`

	ID           int64  `json:"id"`
	ActivityID   string `json:"activity_id"`
	DanmuLoop    string `json:"danmu_loop"`
	Postion      string `json:"position"`
	DisplayUser  string `json:"display_user"`
	DanmuSize    string `json:"danmu_size"`
	DanmuSpeed   string `json:"danmu_speed"`
	DanmuDensity string `json:"danmu_density"`
	DanmuOpacity string `json:"danmu_opacity"`
}

// DefaultDanmuModel 預設DanmuModel
func DefaultDanmuModel() DanmuModel {
	return DanmuModel{Base: Base{TableName: "activity_set_danmu"}}
}

// GetDanmuModelAndID 設置MessageModel與ID
func GetDanmuModelAndID(tablename, id string) DanmuModel {
	idInt, _ := strconv.Atoi(id)
	return DanmuModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m DanmuModel) SetConn(conn db.Connection) DanmuModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m DanmuModel) SetTx(tx *dbsql.Tx) DanmuModel {
	m.Base.Tx = tx
	return m
}

// AddDanmu 增加彈幕資料
func (m DanmuModel) AddDanmu(activityid, loop, position, display, size, speed, density, opacity string) (DanmuModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":   activityid,
		"danmu_loop":    loop,
		"position":      position,
		"display_user":  display,
		"danmu_size":    size,
		"danmu_speed":   speed,
		"danmu_density": density,
		"danmu_opacity": opacity,
	})

	m.ID = id
	m.ActivityID = activityid
	m.DanmuLoop = loop
	m.Postion = position
	m.DisplayUser = display
	m.DanmuSize = size
	m.DanmuSpeed = speed
	m.DanmuDensity = density
	m.DanmuOpacity = opacity
	return m, err
}

// UpdateDanmu 更新彈幕資料
func (m DanmuModel) UpdateDanmu(activityid, loop, position, display, size, speed, density, opacity string) (int64, error) {
	model, err := m.SetTx(m.Base.Tx).Table(m.Base.TableName).Where("id", "=", m.ID).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動")
	}
	if model["activity_id"] != activityid {
		return 0, errors.New("資料中的活動ID不符合，無法更新資料")
	}

	fieldValues := sql.Value{
		"danmu_loop":    loop,
		"position":      position,
		"display_user":  display,
		"danmu_size":    size,
		"danmu_speed":   speed,
		"danmu_density": density,
		"danmu_opacity": opacity,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsActivityExist 檢查活動是否已經設置彈幕
func (m DanmuModel) IsActivityExist(activityid, id string) bool {
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
