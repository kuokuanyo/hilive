package models

import (
	dbsql "database/sql"
	"errors"
	"fmt"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// MaterialModel 活動資料的資料表欄位
type MaterialModel struct {
	Base `json:"-"`

	ID            int64  `json:"id"`
	ActivityID    string `json:"activity_id"`
	DataName      string `json:"data_name"`
	DataIntroduce string `json:"data_introduce"`
	DataLink      string `json:"data_link"`
	DataOrder     int    `json:"data_order"`
}

// DefaultMaterialModel 預設MaterialModel
func DefaultMaterialModel() MaterialModel {
	return MaterialModel{Base: Base{TableName: "activity_material"}}
}

// GetMaterialModelAndID 設置MaterialModel與ID
func GetMaterialModelAndID(tablename, id string) MaterialModel {
	idInt, _ := strconv.Atoi(id)
	return MaterialModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m MaterialModel) SetConn(conn db.Connection) MaterialModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m MaterialModel) SetTx(tx *dbsql.Tx) MaterialModel {
	m.Base.Tx = tx
	return m
}

// AddActivityMaterial 增加活動資料
func (m MaterialModel) AddActivityMaterial(activityid, name, introduce, link string, order int) (MaterialModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	// 判斷order設置
	res, err := m.SetTx(m.Base.Tx).Table("activity_material").Select("id").Where("activity_id", "=", activityid).All()
	count := len(res)
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}
	if order > count+1 {
		return m, fmt.Errorf("該活動目前總共設置%d筆的活動資料，如要新增活動資料，活動排序欄位請設置%d以下(包含)的數值", count, count+1)
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":    activityid,
		"data_name":      name,
		"data_introduce": introduce,
		"data_link":      link,
		"data_order":     order,
	})

	m.ID = id
	m.ActivityID = activityid
	m.DataName = name
	m.DataIntroduce = introduce
	m.DataLink = link
	m.DataOrder = order
	return m, err
}

// UpdateActivityMaterial 更新活動資料
func (m MaterialModel) UpdateActivityMaterial(activityid, name, introduce, link string, order int) (int64, error) {
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	// 判斷order設置
	res, err := m.SetTx(m.Base.Tx).Table("activity_material").Select("id").Where("activity_id", "=", activityid).All()
	count := len(res)
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}
	if order > count {
		return 0, fmt.Errorf("該活動目前總共設置%d筆的活動資料，如要更新活動資料，活動排序欄位請設置%d以下(包含)的數值", count, count)
	}

	fieldValues := sql.Value{
		"activity_id":    activityid,
		"data_name":      name,
		"data_introduce": introduce,
		"data_link":      link,
		"data_order":     order,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsOrderExist 檢查是否已經存在該資料排序
func (m MaterialModel) IsOrderExist(order int, activityid, id string) bool {
	check, _ := m.Table(m.TableName).
		Where("data_order", "=", order).
		Where("activity_id", "=", activityid).
		First()
	if check != nil {
		model, _ := m.Table(m.TableName).
			Where("id", "=", id).First()
		if fmt.Sprintf("%v", model["data_order"]) == fmt.Sprintf("%v", order) {
			return false
		}
	}
	return check != nil
}
