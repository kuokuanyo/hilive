package models

import (
	dbsql "database/sql"
	"errors"
	"fmt"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// IntroduceModel 活動介紹資料表欄位
type IntroduceModel struct {
	Base `json:"-"`

	ID               int64  `json:"id"`
	ActivityID       string `json:"activity_id"`
	IntroduceTitle   string `json:"introduce_title"`
	IntroduceContent string `json:"introduce_content"`
	IntroduceOrder   int    `json:"introduce_order"`
}

// DefaultIntroduceModel 預設IntroduceModel
func DefaultIntroduceModel() IntroduceModel {
	return IntroduceModel{Base: Base{TableName: "activity_introduce"}}
}

// GetIntroduceModelAndID 設置ActivityModel與ID
func GetIntroduceModelAndID(tablename, id string) IntroduceModel {
	idInt, _ := strconv.Atoi(id)
	return IntroduceModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (i IntroduceModel) SetConn(conn db.Connection) IntroduceModel {
	i.Conn = conn
	return i
}

// SetTx 設置Tx
func (i IntroduceModel) SetTx(tx *dbsql.Tx) IntroduceModel {
	i.Base.Tx = tx
	return i
}

// AddActivityIntroduce 增加活動介紹資料
func (i IntroduceModel) AddActivityIntroduce(activityid, title, content string, order int) (IntroduceModel, error) {
	// 檢查是否有該活動
	_, err := i.SetTx(i.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return i, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	// 判斷order設置
	res, _ := i.SetTx(i.Base.Tx).Table("activity_introduce").Where("activity_id", "=", activityid).All()
	count := len(res)
	if order > count+1 {
		return i, fmt.Errorf("該活動目前總共設置%d筆的活動介紹，如要新增活動介紹，活動排序欄位請設置%d以下(包含)的數值", count, count+1)
	}

	id, err := i.SetTx(i.Base.Tx).Table(i.TableName).Insert(sql.Value{
		"activity_id":       activityid,
		"introduce_title":   title,
		"introduce_content": content,
		"introduce_order":   order,
	})

	i.ID = id
	i.ActivityID = activityid
	i.IntroduceTitle = title
	i.IntroduceContent = content
	return i, err
}

// UpdateActivityIntroduce 更新活動介紹資料
func (i IntroduceModel) UpdateActivityIntroduce(activityid, title, content string, order int) (int64, error) {
	_, err := i.SetTx(i.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	// 判斷order設置
	model, _ := i.SetTx(i.Base.Tx).Table("activity_introduce").Where("id", "=", i.ID).First()

	res, _ := i.SetTx(i.Base.Tx).Table("activity_introduce").Where("activity_id", "=", activityid).All()
	count := len(res)
	if fmt.Sprintf("%v", model["activity_id"]) == activityid {
		if order > count {
			return 0, fmt.Errorf("該活動目前總共設置%d筆的活動介紹，如要更新活動介紹，介紹排序欄位請設置%d以下(包含)的數值", count, count)
		}
	} else {
		if order > count+1 {
			return 0, fmt.Errorf("該活動目前總共設置%d筆的活動介紹，如要更新活動介紹，介紹排序欄位請設置%d以下(包含)的數值", count, count+1)
		}
	}

	fieldValues := sql.Value{
		"activity_id":       activityid,
		"introduce_title":   title,
		"introduce_content": content,
		"introduce_order":   order,
	}

	return i.SetTx(i.Tx).Table(i.Base.TableName).
		Where("id", "=", i.ID).Update(fieldValues)
}

// IsOrderExist 檢查是否已經存在該介紹排序
func (i IntroduceModel) IsOrderExist(order int, activityid, id string) bool {
	if id == "" {
		check, _ := i.Table(i.TableName).Where("introduce_order", "=", order).
			Where("activity_id", "=", activityid).First()
		return check != nil
	}
	check, _ := i.Table(i.TableName).
		Where("introduce_order", "=", order).
		Where("activity_id", "=", activityid).
		Where("id", "!=", id).
		First()
	return check != nil
}
