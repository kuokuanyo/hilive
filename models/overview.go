package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// OverviewModel 活動總覽資料表欄位
type OverviewModel struct {
	Base `json:"-"`

	ID         int64  `json:"id"`
	ActivityID string `json:"activity_id"`
	GameID     string `json:"game_id"`
	Open       string `json:"open"`
}

// DefaultOverviewModel 預設ActivityModel
func DefaultOverviewModel() OverviewModel {
	return OverviewModel{Base: Base{TableName: "activity_game_open"}}
}

// GetOverviewModelAndID 設置ActivityModel與ID
func GetOverviewModelAndID(tablename, id string) OverviewModel {
	idInt, _ := strconv.Atoi(id)
	return OverviewModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (o OverviewModel) SetConn(conn db.Connection) OverviewModel {
	o.Conn = conn
	return o
}

// SetTx 設置Tx
func (o OverviewModel) SetTx(tx *dbsql.Tx) OverviewModel {
	o.Base.Tx = tx
	return o
}

// AddActivityOverview 增加活動總覽資料
func (o OverviewModel) AddActivityOverview(activityid, game, open string) (OverviewModel, error) {
	_, err := o.SetTx(o.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return o, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := o.SetTx(o.Base.Tx).Table(o.TableName).Insert(sql.Value{
		"activity_id": activityid,
		"game_id":     game,
		"open":        open,
	})

	o.ID = id
	o.ActivityID = activityid
	o.GameID = game
	o.Open = open
	return o, err
}

// UpdateActivityOverview 更新活動總覽資料
func (o OverviewModel) UpdateActivityOverview(activityid, game, open string) (int64, error) {
	model, err := o.SetTx(o.Base.Tx).Table(o.Base.TableName).Select("activity_id").Where("id", "=", o.ID).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動")
	}
	if model["activity_id"] != activityid {
		return 0, errors.New("資料中的活動ID不符合，無法更新資料")
	}

	fieldValues := sql.Value{
		"game_id": game,
		"open":    open,
	}

	return o.SetTx(o.Tx).Table(o.Base.TableName).
		Where("id", "=", o.ID).Update(fieldValues)
}

// IsGameExist 檢查該活動是否已經創建過相同遊戲
func (o OverviewModel) IsGameExist(game, activityid, id string) bool {
	if id == "" {
		check, _ := o.Table(o.TableName).Where("game_id", "=", game).
			Where("activity_id", "=", activityid).First()
		return check != nil
	}
	check, _ := o.Table(o.TableName).
		Where("game_id", "=", game).
		Where("activity_id", "=", activityid).
		Where("id", "!=", id).
		First()
	return check != nil
}

// Delete 刪除活動總覽資料
func (o OverviewModel) Delete(id string) error {
	return o.SetTx(o.Tx).Table("activity_game_open").
		Where("id", "=", id).Delete()
}
