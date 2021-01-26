package models

import (
	dbsql "database/sql"
	"errors"
	"fmt"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// GuestModel 活動嘉賓資料表欄位
type GuestModel struct {
	Base `json:"-"`

	ID             int64  `json:"id"`
	ActivityID     string `json:"activity_id"`
	GuestPicture   string `json:"guest_picture"`
	GuestName      string `json:"guest_name"`
	GuestIntroduce string `json:"guest_introduce"`
	GuestDetail    string `json:"guest_detail"`
	GuestOrder     int    `json:"guest_order"`
}

// DefaultGuestModel 預設GuestModel
func DefaultGuestModel() GuestModel {
	return GuestModel{Base: Base{TableName: "activity_guest"}}
}

// GetGuestModelAndID 設置GuestModel與ID
func GetGuestModelAndID(tablename, id string) GuestModel {
	idInt, _ := strconv.Atoi(id)
	return GuestModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (g GuestModel) SetConn(conn db.Connection) GuestModel {
	g.Conn = conn
	return g
}

// SetTx 設置Tx
func (g GuestModel) SetTx(tx *dbsql.Tx) GuestModel {
	g.Base.Tx = tx
	return g
}

// AddActivityGuest 增加活動嘉賓資料
func (g GuestModel) AddActivityGuest(activityid, picture, name, introduce, detail string, order int) (GuestModel, error) {
	// 檢查是否有該活動
	_, err := g.SetTx(g.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return g, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	// 判斷order設置
	res, _ := g.SetTx(g.Base.Tx).Table("activity_guest").Where("activity_id", "=", activityid).All()
	count := len(res)
	if order > count+1 {
		return g, fmt.Errorf("該活動目前總共設置%d筆的活動介紹，如要新增活動介紹，活動排序欄位請設置%d以下(包含)的數值", count, count+1)
	}

	id, err := g.SetTx(g.Base.Tx).Table(g.TableName).Insert(sql.Value{
		"activity_id":     activityid,
		"guest_picture":   picture,
		"guest_name":      name,
		"guest_introduce": introduce,
		"guest_detail":    detail,
		"guest_order":     order,
	})

	g.ID = id
	g.ActivityID = activityid
	g.GuestPicture = picture
	g.GuestName = name
	g.GuestIntroduce = introduce
	g.GuestDetail = detail
	g.GuestOrder = order
	return g, err
}

// UpdateActivityGuest 更新活動嘉賓資料
func (g GuestModel) UpdateActivityGuest(activityid, picture, name, introduce, detail string, order int) (int64, error) {
	_, err := g.SetTx(g.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	// 判斷order設置
	model, _ := g.SetTx(g.Base.Tx).Table("activity_guest").Where("id", "=", g.ID).First()

	res, _ := g.SetTx(g.Base.Tx).Table("activity_guest").Where("activity_id", "=", activityid).All()
	count := len(res)
	if fmt.Sprintf("%v", model["activity_id"]) == activityid {
		if order > count {
			return 0, fmt.Errorf("該活動目前總共設置%d筆的活動嘉賓，如要更新活動嘉賓，嘉賓排序欄位請設置%d以下(包含)的數值", count, count)
		}
	} else {
		if order > count+1 {
			return 0, fmt.Errorf("該活動目前總共設置%d筆的活動嘉賓，如要更新活動嘉賓，嘉賓排序欄位請設置%d以下(包含)的數值", count, count+1)
		}
	}

	fieldValues := sql.Value{
		"activity_id":     activityid,
		"guest_picture":   picture,
		"guest_name":      name,
		"guest_introduce": introduce,
		"guest_detail":    detail,
		"guest_order":     order,
	}

	return g.SetTx(g.Tx).Table(g.Base.TableName).
		Where("id", "=", g.ID).Update(fieldValues)
}

// IsOrderExist 檢查是否已經存在該嘉賓排序
func (g GuestModel) IsOrderExist(order int, activityid, id string) bool {
	if id == "" {
		check, _ := g.Table(g.TableName).Where("guest_order", "=", order).
			Where("activity_id", "=", activityid).First()
		return check != nil
	}
	check, _ := g.Table(g.TableName).
		Where("guest_order", "=", order).
		Where("activity_id", "=", activityid).
		Where("id", "!=", id).
		First()
	return check != nil
}
