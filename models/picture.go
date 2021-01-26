package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
	"strings"
)

// PictureModel 圖片牆的資料表欄位
type PictureModel struct {
	Base `json:"-"`

	ID           int64    `json:"id"`
	ActivityID   string   `json:"activity_id"`
	StartTime    string   `json:"start_time"`
	EndTime      string   `json:"end_time"`
	SwitchSecond int      `json:"switch_second"`
	PlayOrder    string   `json:"play_order"`
	PicturePath  []string `json:"picture_path"`
}

// DefaultPictureModel 預設PictureModel
func DefaultPictureModel() PictureModel {
	return PictureModel{Base: Base{TableName: "activity_set_picture"}}
}

// GetPictureModelAndID 設置PictureModel與ID
func GetPictureModelAndID(tablename, id string) PictureModel {
	idInt, _ := strconv.Atoi(id)
	return PictureModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m PictureModel) SetConn(conn db.Connection) PictureModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m PictureModel) SetTx(tx *dbsql.Tx) PictureModel {
	m.Base.Tx = tx
	return m
}

// AddPicture 增加圖片牆資料
func (m PictureModel) AddPicture(activityid, start, end string, second int, order, path, background string) (PictureModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":   activityid,
		"start_time":    start,
		"end_time":      end,
		"switch_second": second,
		"play_order":    order,
		"picture_path":  path,
		"background":    background,
	})

	m.ID = id
	m.ActivityID = activityid
	m.StartTime = start
	m.EndTime = end
	m.SwitchSecond = second
	m.PlayOrder = order
	m.PicturePath = strings.Split(path, "\n")
	return m, err
}

// UpdatePicture 更新圖片牆資料
func (m PictureModel) UpdatePicture(activityid, start, end string, second int, order, path, background string) (int64, error) {
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	fieldValues := sql.Value{
		"activity_id":   activityid,
		"start_time":    start,
		"end_time":      end,
		"switch_second": second,
		"play_order":    order,
		"picture_path":  path,
		"background":    background,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsActivityExist 檢查活動是否已經設置圖片牆
func (m PictureModel) IsActivityExist(activityid, id string) bool {
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
