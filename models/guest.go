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

// AddGuest 增加活動嘉賓資料
func (g GuestModel) AddGuest(activityid, picture, name, introduce, detail string) (GuestModel, error) {
	// 檢查是否有該活動
	_, err := g.SetTx(g.Base.Tx).Table("activity").Where("activity_id", "=", activityid).First()
	if err != nil {
		return g, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	// 判斷order設置
	res, _ := g.SetTx(g.Base.Tx).Table("activity_guest").Where("activity_id", "=", activityid).All()
	count := len(res)

	id, err := g.SetTx(g.Base.Tx).Table(g.TableName).Insert(sql.Value{
		"activity_id":     activityid,
		"guest_picture":   picture,
		"guest_name":      name,
		"guest_introduce": introduce,
		"guest_detail":    detail,
		"guest_order":     count + 1,
	})

	g.ID = id
	g.ActivityID = activityid
	g.GuestPicture = picture
	g.GuestName = name
	g.GuestIntroduce = introduce
	g.GuestDetail = detail
	g.GuestOrder = count + 1
	return g, err
}

// UpdateGuest 更新活動嘉賓資料
func (g GuestModel) UpdateGuest(activityid, picture, name, introduce, detail string, order int) (int64, error) {
	// 檢查是否有該活動
	_, err := g.SetTx(g.Base.Tx).Table("activity").Where("activity_id", "=", activityid).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	model, err := g.SetTx(g.Base.Tx).Table(g.Base.TableName).Where("id", "=", g.ID).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動")
	}

	res, _ := g.SetTx(g.Base.Tx).Table(g.Base.TableName).Where("activity_id", "=", activityid).All()
	count := len(res)
	if order > count {
		return 0, fmt.Errorf("該活動目前總共設置%d筆的活動嘉賓，如要更新活動嘉賓，嘉賓排序欄位請設置%d以下(包含)的數值", count, count)
	}

	// 還沒更新前的order
	originalOrder := model["guest_order"]
	if fmt.Sprintf("%v", originalOrder) != strconv.Itoa(order) {
		_, err = g.SetTx(g.Tx).Table(g.Base.TableName).
			Where("activity_id", "=", activityid).Where("guest_order", "=", order).Update(sql.Value{
			"guest_order": originalOrder,
		})
		if err != nil {
			if err.Error() != "沒有影響任何資料" {
				return 0, errors.New("更新活動嘉賓order欄位發生錯誤")
			}
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
