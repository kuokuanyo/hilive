package models

import (
	dbsql "database/sql"
	"errors"
	"fmt"
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
	ContractOrder          int    `json:"contract_order"`
}

// DefaultContractModel 預設HoldScreenModel
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
	mobile, direction, mobileback, create string, order int) (ContractModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	// 判斷order設置
	res, _ := m.SetTx(m.Base.Tx).Table("activity_set_contract").Where("activity_id", "=", activityid).All()
	count := len(res)
	if order > count+1 {
		return m, fmt.Errorf("該活動目前總共設置%d筆的活動資料，如要新增活動資料，活動排序欄位請設置%d以下(包含)的數值", count, count+1)
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
		"create_time":              create,
		"contract_order":           order,
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
	m.CreateTime = create
	m.ContractOrder = order
	return m, err
}

// UpdateContract 更新簽約牆資料
func (m ContractModel) UpdateContract(activityid, title, contractBack, animation, area,
	mobile, direction, mobileback, create string, order int) (int64, error) {
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	// 判斷order設置
	res, _ := m.SetTx(m.Base.Tx).Table("activity_set_contract").Where("activity_id", "=", activityid).All()
	count := len(res)
	if order > count {
		return 0, fmt.Errorf("該活動目前總共設置%d筆的活動介紹，如要更新簽約牆設置，簽約牆排序欄位請設置%d以下(包含)的數值", count, count)
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
		"create_time":              create,
		"contract_order":           order,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsOrderExist 檢查是否已經存在該簽約牆排序
func (m ContractModel) IsOrderExist(order int, activityid, id string) bool {
	if id == "" {
		check, _ := m.Table(m.TableName).Where("contract_order", "=", order).
			Where("activity_id", "=", activityid).First()
		return check != nil
	}
	check, _ := m.Table(m.TableName).
		Where("contract_order", "=", order).
		Where("activity_id", "=", activityid).
		Where("id", "!=", id).
		First()
	return check != nil
}
