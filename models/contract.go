package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// ContractModel 簽約牆的資料表欄位
type ContractModel struct {
	Base `json:"-"`

	ID                     int64  `json:"id"`
	ActivityID             string `json:"activity_id"`
	ContractTitle          string `json:"contract_title"`
	ContractBackground     string `json:"contract_background"`
	SignatureAnimationSize string `json:"signature_animation_size"`
	SignatureAreaSize      string `json:"signature_area_size"`
	MobileDevice           string `json:"MobileDevice"`
	DeviceDirection        string `json:"device_direction"`
	MobileBackground       string `json:"MobileBackground"`
	CreateTime             string `json:"create_time"`
}

// DefaultContractModel 預設ContractModel
func DefaultContractModel() ContractModel {
	return ContractModel{Base: Base{TableName: "activity_set_contract"}}
}

// GetContractModelAndID 設置ContractModel與ID
func GetContractModelAndID(tablename, id string) ContractModel {
	idInt, _ := strconv.Atoi(id)
	return ContractModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m ContractModel) SetConn(conn db.Connection) ContractModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m ContractModel) SetTx(tx *dbsql.Tx) ContractModel {
	m.Base.Tx = tx
	return m
}

// AddContract 增加合約牆資料
func (m ContractModel) AddContract(activityid, title, contractBack, animation, area,
	mobile, direction, mobileback string) (ContractModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":              activityid,
		"contract_title":           title,
		"contract_background":      contractBack,
		"signature_animation_size": animation,
		"signature_area_size":      area,
		"mobile_device":            mobile,
		"device_direction":         direction,
		"mobile_background":        mobileback,
	})

	m.ID = id
	m.ActivityID = activityid
	m.ContractTitle = title
	m.ContractBackground = contractBack
	m.SignatureAnimationSize = animation
	m.SignatureAreaSize = area
	m.MobileDevice = mobile
	m.DeviceDirection = direction
	m.MobileBackground = mobileback
	return m, err
}

// UpdateContract 更新簽約牆資料
func (m ContractModel) UpdateContract(activityid, title, contractBack, animation, area,
	mobile, direction, mobileback string) (int64, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	fieldValues := sql.Value{
		"activity_id":              activityid,
		"contract_title":           title,
		"contract_background":      contractBack,
		"signature_animation_size": animation,
		"signature_area_size":      area,
		"mobile_device":            mobile,
		"device_direction":         direction,
		"mobile_background":        mobileback,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}
