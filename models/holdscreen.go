package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// HoldScreenModel 霸屏的資料表欄位
type HoldScreenModel struct {
	Base `json:"-"`

	ID              int64  `json:"id"`
	ActivityID      string `json:"activity_id"`
	HoldscreenPrice int    `json:"holdscreen_price"`
	MessageCheck    string `json:"message_check"`
	OnlyPicture     string `json:"only_picture"`
	MinimumSecond   int    `json:"minimum_second"`
	BirthdayTopic   string `json:"birthday_topic"`
	ConfessTopic    string `json:"confess_topic"`
	ProposeTopic    string `json:"propose_topic"`
	BlessTopic      string `json:"bless_topic"`
	GoddessTopic    string `json:"goddess_topic"`
}

// DefaultHoldScreenModel 預設HoldScreenModel
func DefaultHoldScreenModel() HoldScreenModel {
	return HoldScreenModel{Base: Base{TableName: "activity_set_holdscreen"}}
}

// GetHoldScreenModelAndID 設置HoldScreenModel與ID
func GetHoldScreenModelAndID(tablename, id string) HoldScreenModel {
	idInt, _ := strconv.Atoi(id)
	return HoldScreenModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m HoldScreenModel) SetConn(conn db.Connection) HoldScreenModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m HoldScreenModel) SetTx(tx *dbsql.Tx) HoldScreenModel {
	m.Base.Tx = tx
	return m
}

// AddHoldScreen 增加霸屏資料
func (m HoldScreenModel) AddHoldScreen(activityid string, holdPrice int, check string, only string,
	second int, biryhday, confess, propose, bless, goddess string) (HoldScreenModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":      activityid,
		"holdscreen_price": holdPrice,
		"message_check":    check,
		"only_picture":     only,
		"minimum_second":   second,
		"birthday_topic":   biryhday,
		"confess_topic":    confess,
		"propose_topic":    propose,
		"bless_topic":      bless,
		"goddess_topic":    goddess,
	})

	m.ID = id
	m.ActivityID = activityid
	m.HoldscreenPrice = holdPrice
	m.MessageCheck = check
	m.OnlyPicture = only
	m.MinimumSecond = second
	m.BirthdayTopic = biryhday
	m.ConfessTopic = confess
	m.ProposeTopic = propose
	m.BlessTopic = bless
	m.GoddessTopic = goddess
	return m, err
}

// UpdateHoldScreen 更新霸屏資料
func (m HoldScreenModel) UpdateHoldScreen(activityid string, holdPrice int, check string, only string,
	second int, biryhday, confess, propose, bless, goddess string) (int64, error) {
	model, err := m.SetTx(m.Base.Tx).Table(m.Base.TableName).Where("id", "=", m.ID).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動")
	}
	if model["activity_id"] != activityid {
		return 0, errors.New("資料中的活動ID不符合，無法更新資料")
	}

	fieldValues := sql.Value{
		"holdscreen_price": holdPrice,
		"message_check":    check,
		"only_picture":     only,
		"minimum_second":   second,
		"birthday_topic":   biryhday,
		"confess_topic":    confess,
		"propose_topic":    propose,
		"bless_topic":      bless,
		"goddess_topic":    goddess,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsActivityExist 檢查活動是否已經設置霸屏
func (m HoldScreenModel) IsActivityExist(activityid, id string) bool {
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
