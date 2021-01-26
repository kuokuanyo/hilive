package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// Sign3DModel 3D簽到牆的資料表欄位
type Sign3DModel struct {
	Base `json:"-"`

	ID          int64  `json:"id"`
	ActivityID  string `json:"activity_id"`
	AvatarShape string `json:"avatar_shape"`
	Display     string `json:"Display"`
	Background  string `json:"background"`
}

// DefaultSign3DModel 預設Sign3DModel
func DefaultSign3DModel() Sign3DModel {
	return Sign3DModel{Base: Base{TableName: "activity_set_3d"}}
}

// GetSign3DModelAndID 設置Sign3DModel與ID
func GetSign3DModelAndID(tablename, id string) Sign3DModel {
	idInt, _ := strconv.Atoi(id)
	return Sign3DModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m Sign3DModel) SetConn(conn db.Connection) Sign3DModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m Sign3DModel) SetTx(tx *dbsql.Tx) Sign3DModel {
	m.Base.Tx = tx
	return m
}

// Add3DSign 增加3D簽到牆資料
func (m Sign3DModel) Add3DSign(activityid, shape, display, background string) (Sign3DModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":  activityid,
		"avatar_shape": shape,
		"display":      display,
		"background":   background,
	})

	m.ID = id
	m.ActivityID = activityid
	m.AvatarShape = shape
	m.Display = display
	m.Background = background
	return m, err
}

// Update3DSign 更新3D簽到牆資料
func (m Sign3DModel) Update3DSign(activityid, shape, display, background string) (int64, error) {
	model, err := m.SetTx(m.Base.Tx).Table(m.Base.TableName).Where("id", "=", m.ID).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動")
	}
	if model["activity_id"] != activityid {
		return 0, errors.New("資料中的活動ID不符合，無法更新資料")
	}

	fieldValues := sql.Value{
		"avatar_shape": shape,
		"display":      display,
		"background":   background,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsActivityExist 檢查活動是否已經設置3D簽到牆
func (m Sign3DModel) IsActivityExist(activityid, id string) bool {
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
