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
func (m MaterialModel) AddActivityMaterial(activityid, name, introduce, link string) (MaterialModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	// 判斷order設置
	res, _ := m.SetTx(m.Base.Tx).Table("activity_material").Where("activity_id", "=", activityid).All()
	count := len(res)

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":    activityid,
		"data_name":      name,
		"data_introduce": introduce,
		"data_link":      link,
		"data_order":     count + 1,
	})

	m.ID = id
	m.ActivityID = activityid
	m.DataName = name
	m.DataIntroduce = introduce
	m.DataLink = link
	m.DataOrder = count + 1
	return m, err
}

// UpdateActivityMaterial 更新活動資料
func (m MaterialModel) UpdateActivityMaterial(activityid, name, introduce, link string, order int) (int64, error) {
	model, err := m.SetTx(m.Base.Tx).Table(m.Base.TableName).Where("id", "=", m.ID).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動")
	}
	if model["activity_id"] != activityid {
		return 0, errors.New("資料中的活動ID不符合，無法更新資料")
	}

	res, _ := m.SetTx(m.Base.Tx).Table(m.Base.TableName).Where("activity_id", "=", activityid).All()
	count := len(res)
	if order > count {
		return 0, fmt.Errorf("該活動目前總共設置%d筆的活動資料，如要更新活動資料，資料排序欄位請設置%d以下(包含)的數值", count, count)
	}

	// 還沒更新前的order
	originalOrder := model["data_order"]
	if fmt.Sprintf("%v", originalOrder) != strconv.Itoa(order) {
		_, err = m.SetTx(m.Tx).Table(m.Base.TableName).
			Where("activity_id", "=", activityid).Where("data_order", "=", order).Update(sql.Value{
			"data_order": originalOrder,
		})
		if err != nil {
			if err.Error() != "沒有影響任何資料" {
				return 0, errors.New("更新活動資料order欄位發生錯誤")
			}
		}
	}

	fieldValues := sql.Value{
		"data_name":      name,
		"data_introduce": introduce,
		"data_link":      link,
		"data_order":     order,
	}
	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}
