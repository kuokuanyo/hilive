package models

import (
	dbsql "database/sql"
	"errors"
	"hilive/modules/db"
	"hilive/modules/db/sql"
	"strconv"
)

// GameLotteryModel 幸運轉盤的資料表欄位
type GameLotteryModel struct {
	Base `json:"-"`

	ID              int64  `json:"id"`
	ActivityID      string `json:"activity_id"`
	GameID          string `json:"game_id"`
	Title           string `json:"title"`
	Rule            string `json:"rule"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	GameStatus      string `json:"game_status"`
	ScreenOpen      string `json:"screen_open"`
	MaxWinTimes     int    `json:"max_win_times"`
	MaxPeople       int    `json:"max_people"`
	RaffleFrequency string `json:"raffle_frequency"`
	RaffleTimes     int    `json:"raffle_times"`
}

// DefaultGameLotteryModel 預設GameLotteryModel
func DefaultGameLotteryModel() GameLotteryModel {
	return GameLotteryModel{Base: Base{TableName: "activity_set_gamelottery"}}
}

// GetGameLotteryModelAndID 設置GameLotteryModel與ID
func GetGameLotteryModelAndID(tablename, id string) GameLotteryModel {
	idInt, _ := strconv.Atoi(id)
	return GameLotteryModel{Base: Base{TableName: tablename}, ID: int64(idInt)}
}

// SetConn 設定connection
func (m GameLotteryModel) SetConn(conn db.Connection) GameLotteryModel {
	m.Conn = conn
	return m
}

// SetTx 設置Tx
func (m GameLotteryModel) SetTx(tx *dbsql.Tx) GameLotteryModel {
	m.Base.Tx = tx
	return m
}

// AddGameLottery 增加幸運轉盤遊戲
func (m GameLotteryModel) AddGameLottery(activityid, gameid, title, rule, start, end, status,
	open, freq string, winTimes, people, raffle int) (GameLotteryModel, error) {
	// 檢查是否有該活動
	_, err := m.SetTx(m.Base.Tx).Table("activity").Select("id").Where("activity_id", "=", activityid).First()
	if err != nil {
		return m, errors.New("查詢不到此活動ID，請輸入正確活動ID")
	}

	id, err := m.SetTx(m.Base.Tx).Table(m.TableName).Insert(sql.Value{
		"activity_id":      activityid,
		"game_id":          gameid,
		"title":            title,
		"rule":             rule,
		"start_time":       start,
		"end_time":         end,
		"game_status":      status,
		"screen_open":      open,
		"max_win_times":    winTimes,
		"max_people":       people,
		"raffle_frequency": freq,
		"raffle_times":     raffle,
	})

	m.ID = id
	m.ActivityID = activityid
	m.GameID = gameid
	m.Title = title
	m.Rule = rule
	m.StartTime = start
	m.EndTime = end
	m.GameStatus = status
	m.ScreenOpen = open
	m.MaxWinTimes = winTimes
	m.MaxPeople = people
	m.RaffleFrequency = freq
	m.RaffleTimes = raffle
	return m, err
}

// UpdateGameLottery 更新幸運轉盤資料
func (m GameLotteryModel) UpdateGameLottery(activityid, gameid, title, rule, start, end, status,
	open, freq string, winTimes, people, raffle int) (int64, error) {
	model, err := m.SetTx(m.Base.Tx).Table(m.Base.TableName).Where("id", "=", m.ID).First()
	if err != nil {
		return 0, errors.New("查詢不到此活動")
	}
	if model["activity_id"] != activityid {
		return 0, errors.New("資料中的活動ID不符合，無法更新資料")
	}
	if model["game_id"] != gameid {
		return 0, errors.New("資料中的遊戲ID不符合，無法更新資料")
	}

	fieldValues := sql.Value{
		"title":            title,
		"rule":             rule,
		"start_time":       start,
		"end_time":         end,
		"game_status":      status,
		"screen_open":      open,
		"max_win_times":    winTimes,
		"max_people":       people,
		"raffle_frequency": freq,
		"raffle_times":     raffle,
	}

	return m.SetTx(m.Tx).Table(m.Base.TableName).
		Where("id", "=", m.ID).Update(fieldValues)
}
