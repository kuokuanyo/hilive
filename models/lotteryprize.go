package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// LotteryPrizeModel 幸運轉盤獎品的資料表欄位
type LotteryPrizeModel struct {
	Base `json:"-"`

	ID             int64  `json:"id"`
	GameID         string `json:"game_id"`
	PrizeName      string `json:"prize_name"`
	PrizeType      string `json:"prize_type"`
	Picture        string `json:"picture"`
	Percent        int    `json:"percent"`
	Amount         int    `json:"amount"`
	Remain         int    `json:"remain"`
	Price          int    `json:"price"`
	AllowWinning   string `json:"allow_winning"`
	RedeemMethod   string `json:"redeem_method"`
	RedeemPassword string `json:"redeem_password"`
}

// DefaultLotteryPrizeModel 預設LotteryPrizeModel
func DefaultLotteryPrizeModel() LotteryPrizeModel {
	return LotteryPrizeModel{Base: Base{TableName: "activity_gamelottery_prize"}}
}

// GetLotteryPrizeModelAndID 設置LotteryPrizeModel與ID
func GetLotteryPrizeModelAndID(tablename, id string) LotteryPrizeModel {
	idInt, _ := strconv.Atoi(id)
	return LotteryPrizeModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m LotteryPrizeModel) SetConn(conn db.Connection) LotteryPrizeModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m LotteryPrizeModel) SetTx(tx *dbsql.Tx) LotteryPrizeModel {
	m.Base.Tx = tx
	return m
}

// AddLotteryPrize 增加幸運轉盤獎品
func (m LotteryPrizeModel) AddLotteryPrize(gameid, name, prizeType, allow, method,
	password, picture string, percent, amount, price int) (LotteryPrizeModel, error) {
	// 檢查是否有該遊戲
	_, err := m.SetTx(m.Base.Tx).Table("activity_set_gamelottery").Where("game_id", "=", gameid).First()
	if err != nil {
		return m, errors.New("查詢不到幸運轉盤遊戲ID，請輸入正確遊戲ID")
	}
	if percent > 100 || percent < 0 {
		return m, errors.New("中獎機率設置區間為1~100，請重新設置")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"game_id":         gameid,
		"prize_name":      name,
		"prize_type":      prizeType,
		"percent":         percent,
		"amount":          amount,
		"remain":          amount,
		"allow_winning":   allow,
		"redeem_method":   method,
		"redeem_password": password,
		"price":           price,
		"picture":         picture,
	})

	m.ID = id
	m.PrizeName = name
	m.PrizeType = prizeType
	m.Percent = percent
	m.Amount = amount
	m.AllowWinning = allow
	m.RedeemMethod = method
	m.RedeemPassword = password
	m.Price = price
	m.Remain = amount
	m.Picture = picture
	return m, err
}

// UpdateLotteryPrize 更新幸運轉盤獎品資料
func (m LotteryPrizeModel) UpdateLotteryPrize(gameid, name, prizeType, allow, method,
	password, picture string, percent, amount, price, left int) (int64, error) {
	// 檢查是否有該遊戲
	_, err := m.SetTx(m.Base.Tx).Table("activity_set_gamelottery").Where("game_id", "=", gameid).First()
	if err != nil {
		return 0, errors.New("查詢不到幸運轉盤遊戲ID，請輸入正確遊戲ID")
	}

	if amount < left {
		return 0, errors.New("獎品剩餘數量不可大於獎品總數量，請重新設置")
	}
	if percent > 100 || percent < 0 {
		return 0, errors.New("中獎機率設置區間為1~100，請重新設置")
	}

	fieldValues := sql.Value{
		"game_id":         gameid,
		"prize_name":      name,
		"prize_type":      prizeType,
		"percent":         percent,
		"amount":          amount,
		"allow_winning":   allow,
		"redeem_method":   method,
		"redeem_password": password,
		"price":           price,
		"remain":          left,
		"picture":         picture,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}
