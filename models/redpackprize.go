package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// RedpackPrizeModel 搖紅包獎品的資料表欄位
type RedpackPrizeModel struct {
	Base `json:"-"`

	ID             int64  `json:"id"`
	GameID         string `json:"game_id"`
	PrizeName      string `json:"prize_name"`
	PrizeType      string `json:"prize_type"`
	Picture        string `json:"picture"`
	Amount         int    `json:"amount"`
	Remain         int    `json:"remain"`
	Price          int    `json:"price"`
	RedeemMethod   string `json:"redeem_method"`
	RedeemPassword string `json:"redeem_password"`
}

// DefaultRedpackPrizeModel 預設RedpackPrizeModel
func DefaultRedpackPrizeModel() RedpackPrizeModel {
	return RedpackPrizeModel{Base: Base{TableName: "activity_redpack_prize"}}
}

// GetRedpackPrizeModelAndID 設置RedpackPrizeModel與ID
func GetRedpackPrizeModelAndID(tablename, id string) RedpackPrizeModel {
	idInt, _ := strconv.Atoi(id)
	return RedpackPrizeModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m RedpackPrizeModel) SetConn(conn db.Connection) RedpackPrizeModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m RedpackPrizeModel) SetTx(tx *dbsql.Tx) RedpackPrizeModel {
	m.Base.Tx = tx
	return m
}

// AddRedpackPrize 增加搖紅包獎品
func (m RedpackPrizeModel) AddRedpackPrize(gameid, name, prizeType, method,
	password, picture string, amount, price int) (RedpackPrizeModel, error) {
	// 檢查是否有該遊戲
	_, err := m.SetTx(m.Base.Tx).Table("activity_set_redpack").Where("game_id", "=", gameid).First()
	if err != nil {
		return m, errors.New("查詢不到搖紅包遊戲ID，請輸入正確遊戲ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"game_id":         gameid,
		"prize_name":      name,
		"prize_type":      prizeType,
		"amount":          amount,
		"remain":          amount,
		"redeem_method":   method,
		"redeem_password": password,
		"price":           price,
		"picture":         picture,
	})

	m.ID = id
	m.PrizeName = name
	m.PrizeType = prizeType
	m.Amount = amount
	m.RedeemMethod = method
	m.RedeemPassword = password
	m.Price = price
	m.Remain = amount
	m.Picture = picture
	return m, err
}

// UpdateRedpackPrize 更新搖紅包獎品資料
func (m RedpackPrizeModel) UpdateRedpackPrize(gameid, name, prizeType, method,
	password, picture string, amount, price, remain int) (int64, error) {
	// 檢查是否有該遊戲
	_, err := m.SetTx(m.Base.Tx).Table("activity_set_redpack").Where("game_id", "=", gameid).First()
	if err != nil {
		return 0, errors.New("查詢不到搖紅包遊戲ID，請輸入正確遊戲ID")
	}
	if amount < remain {
		return 0, errors.New("獎品剩餘數量不可大於獎品總數量，請重新設置")
	}

	fieldValues := sql.Value{
		"game_id":         gameid,
		"prize_name":      name,
		"prize_type":      prizeType,
		"amount":          amount,
		"redeem_method":   method,
		"redeem_password": password,
		"price":           price,
		"remain":          remain,
		"picture":         picture,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}
