package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// LotteryOtherModel 幸運轉盤謝謝參與獎的資料表欄位
type LotteryOtherModel struct {
	Base `json:"-"`

	ID     int64  `json:"id"`
	GameID string `json:"game_id"`
	Amount int    `json:"amount"`
}

// DefaultLotteryOtherModel 預設LotteryOtherModel
func DefaultLotteryOtherModel() LotteryOtherModel {
	return LotteryOtherModel{Base: Base{TableName: "activity_gamelottery_other_prize"}}
}

// GetLotteryOtherModelAndID 設置LotteryOtherModel與ID
func GetLotteryOtherModelAndID(tablename, id string) LotteryOtherModel {
	idInt, _ := strconv.Atoi(id)
	return LotteryOtherModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m LotteryOtherModel) SetConn(conn db.Connection) LotteryOtherModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m LotteryOtherModel) SetTx(tx *dbsql.Tx) LotteryOtherModel {
	m.Base.Tx = tx
	return m
}

// AddLotteryOther 增加幸運轉盤謝謝參與獎資料
func (m LotteryOtherModel) AddLotteryOther(gameid string, amount int) (LotteryOtherModel, error) {
	// 檢查是否有該遊戲
	_, err := m.SetTx(m.Base.Tx).Table("activity_set_gamelottery").
		Where("game_id", "=", gameid).First()
	if err != nil {
		return m, errors.New("查詢不到幸運轉盤遊戲ID，請輸入正確遊戲ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"game_id": gameid,
		"amount":  amount,
	})

	m.ID = id
	m.GameID = gameid
	m.Amount = amount
	return m, err
}

// UpdateLotteryOther 更新幸運轉盤謝謝參與獎資料
func (m LotteryOtherModel) UpdateLotteryOther(gameid string, amount int) (int64, error) {
	// 檢查是否有該遊戲
	_, err := m.SetTx(m.Base.Tx).Table("activity_set_gamelottery").
		Where("game_id", "=", gameid).First()
	if err != nil {
		return 0, errors.New("查詢不到幸運轉盤遊戲ID，請輸入正確遊戲ID")
	}

	fieldValues := sql.Value{
		"game_id": gameid,
		"amount":  amount,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}

// IsGameExist 檢查遊戲是否已經設置謝謝參與獎
func (m LotteryOtherModel) IsGameExist(gameid, id string) bool {
	if id == "" {
		check, _ := m.Table(m.TableName).Where("game_id", "=", gameid).First()
		return check != nil
	}
	check, _ := m.Table(m.TableName).
		Where("game_id", "=", gameid).
		Where("id", "!=", id).
		First()
	return check != nil
}
