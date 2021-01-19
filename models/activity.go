package models

import (
	dbsql "database/sql"
	"errors"
	"fmt"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// ActivityModel 活動資料表欄位
type ActivityModel struct {
	Base `json:"-"`

	ID                   int64  `json:"id"`
	ActivityID           string `json:"activity_id"`
	ActivityName         string `json:"activity_name"`
	ActivityType         string `json:"activity_type"`
	ExpectedParticipants string `json:"expected_participants"`
	Participants         string `json:"participants"`
	City                 string `json:"city"`
	Town                 string `json:"town"`
	StartTime            string `json:"start_time"`
	EndTime              string `json:"end_time"`
}

// DefaultActivityModel 預設ActivityModel
func DefaultActivityModel() ActivityModel {
	return ActivityModel{Base: Base{TableName: "activity"}}
}

// GetActivityModelAndID 設置ActivityModel與ID
func GetActivityModelAndID(tablename, id string) ActivityModel {
	idInt, _ := strconv.Atoi(id)
	return ActivityModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (a ActivityModel) SetConn(conn db.Connection) ActivityModel {
	a.Conn = conn
	return a
}

// SetTx 設置Tx
func (a ActivityModel) SetTx(tx *dbsql.Tx) ActivityModel {
	a.Base.Tx = tx
	return a
}

// AddActivity 增加活動資料
func (a ActivityModel) AddActivity(activityid, activityName, activityType, expected, city, town, start, end string) (ActivityModel, error) {
	cityModel, err := a.SetTx(a.Base.Tx).Table("activity_town").Select("city_id", "town_id").Where("town", "=", town).First()
	if err != nil {
		return a, errors.New("輸入活動區域名稱發生錯誤，原因: 無此區域")
	}
	if fmt.Sprintf("%v", cityModel["city_id"]) != city {
		return a, errors.New("該縣市無此區域，請重新輸入活動區域")
	}

	id, err := a.SetTx(a.Base.Tx).Table(a.TableName).Insert(sql.Value{
		"activity_id":           activityid,
		"activity_name":         activityName,
		"activity_type":         activityType,
		"expected_participants": expected,
		"participants":          "0",
		"city":                  city,
		"town":                  cityModel["town_id"],
		"start_time":            start,
		"end_time":              end,
	})

	a.ID = id
	a.ActivityID = activityid
	a.ActivityName = activityName
	a.ActivityType = activityType
	a.ExpectedParticipants = expected
	a.Participants = "0"
	a.City = city
	a.Town = town
	a.StartTime = start
	a.EndTime = end
	return a, err
}

// Update 更新活動資料
func (a ActivityModel) Update(activityid, activityName, activityType, expected, city, town, start, end string) (int64, error) {
	cityModel, err := a.SetTx(a.Base.Tx).Table("activity_town").Select("city_id", "town_id").Where("town", "=", town).First()
	if err != nil {
		return 0, errors.New("輸入活動區域名稱發生錯誤，原因: 無此區域")
	}
	if fmt.Sprintf("%v", cityModel["city_id"]) != city {
		return 0, errors.New("該縣市無此區域，請重新輸入活動區域")
	}

	fieldValues := sql.Value{
		"activity_id":           activityid,
		"activity_name":         activityName,
		"activity_type":         activityType,
		"expected_participants": expected,
		"city":                  city,
		"town":                  cityModel["town_id"],
		"start_time":            start,
		"end_time":              end,
	}

	return a.SetTx(a.Tx).Table(a.Base.TableName).
		Where("id", "=", a.ID).Update(fieldValues)
}
