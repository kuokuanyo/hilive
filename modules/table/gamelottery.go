package table

import (
	"database/sql"
	"errors"
	"hilive/context"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/modules/utils"
	"hilive/template/form"
	"hilive/template/types"
	"html/template"
	"strconv"
	"time"
)

// GetGameLottery 取得轉盤遊戲頁面、表單資訊
func (s *SystemTable) GetGameLottery(ctx *context.Context) (gametable Table) {
	// 建立BaseTable
	gametable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := gametable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("小遊戲專屬ID", "game_id", db.Varchar).FieldFilterable()
	info.AddField("遊戲標題", "title", db.Varchar)
	info.AddField("遊戲規則", "rule", db.Varchar)
	info.AddField("遊戲開始時間", "start_time", db.Datetime)
	info.AddField("遊戲結束時間", "end_time", db.Datetime)
	info.AddField("遊戲狀態", "game_status", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("大屏幕狀態", "screen_open", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("每人最多中獎次數", "max_win_times", db.Int)
	info.AddField("最多參與抽獎人數", "max_people", db.Int)
	info.AddField("抽獎頻率", "raffle_frequency", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "每天"
			}
			return "最多"
		})
	info.AddField("抽獎次數", "raffle_times", db.Int)

	info.SetTable("activity_set_gamelottery").SetTitle("幸運轉盤").SetDescription("遊戲管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_gamelottery").WhereIn("id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除活動資料發生錯誤"), nil
					}
				}
				return nil, nil
			})
			return txErr
		})

	// 取得FormPanel
	formList := gametable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("小遊戲專屬ID", "game_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("小遊戲辨別ID")).SetFieldMust().FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("遊戲標題", "title", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("遊戲規則", "rule", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("遊戲開始時間", "start_time", db.Datetime, form.Datetime).SetFieldMust()
	formList.AddField("遊戲結束時間", "end_time", db.Datetime, form.Datetime).SetFieldMust()
	formList.AddField("遊戲狀態", "game_status", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "開啟"},
			{Value: "0", Text: "關閉"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_gamelottery").Select("game_status").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["game_status"].(int64), 10))
			return stauts
		}).SetFieldDefault("0")
	formList.AddField("大屏幕狀態", "screen_open", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "開啟"},
			{Value: "0", Text: "關閉"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_gamelottery").Select("screen_open").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["screen_open"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("每人最多中獎次數", "max_win_times", db.Int, form.Number).SetFieldMust().SetFieldDefault("2")
	formList.AddField("最多參與抽獎人數", "max_people", db.Int, form.Number).SetFieldMust().SetFieldDefault("30")
	formList.AddField("抽獎頻率", "raffle_frequency", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "每天"},
			{Value: "0", Text: "最多"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_gamelottery").Select("raffle_frequency").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["raffle_frequency"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("抽獎次數", "raffle_times", db.Int, form.Number).SetFieldMust().SetFieldDefault("1")

	formList.SetTable("activity_set_gamelottery").SetTitle("幸運轉盤").SetDescription("遊戲管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "title", "rule", "start_time",
			"end_time", "game_status", "screen_open", "max_win_times", "max_people",
			"raffle_frequency", "raffle_times") {
			return errors.New("活動ID、標題、規則、時間等欄位都不能為空")
		}

		// 時間判斷
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("start_time"), time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("end_time"), time.Local)
		boolTime := end.After(start) && start.Before(end)
		if boolTime == false {
			return errors.New("時間設置發生錯誤，請重新設置(結束時間在開始時間之後)")
		}

		win, _ := strconv.Atoi(values.Get("max_win_times"))
		people, _ := strconv.Atoi(values.Get("max_people"))
		raffle, _ := strconv.Atoi(values.Get("raffle_times"))
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultGameLotteryModel().SetTx(tx).SetConn(s.conn).AddGameLottery(
				values.Get("activity_id"), utils.UUID(8), values.Get("title"), values.Get("rule"),
				values.Get("start_time"), values.Get("end_time"), values.Get("game_status"),
				values.Get("screen_open"), values.Get("raffle_frequency"), win, people, raffle)
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id", "title", "rule", "start_time",
			"end_time", "game_status", "screen_open", "max_win_times", "max_people",
			"raffle_frequency", "raffle_times") {
			return errors.New("活動ID、遊戲ID、標題、規則、時間等欄位都不能為空")
		}

		// 時間判斷
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("start_time"), time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("end_time"), time.Local)
		boolTime := end.After(start) && start.Before(end)
		if boolTime == false {
			return errors.New("時間設置發生錯誤，請重新設置(結束時間在開始時間之後)")
		}

		win, _ := strconv.Atoi(values.Get("max_win_times"))
		people, _ := strconv.Atoi(values.Get("max_people"))
		raffle, _ := strconv.Atoi(values.Get("raffle_times"))
		model := models.GetGameLotteryModelAndID("activity_set_gamelottery", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).UpdateGameLottery(values.Get("activity_id"), values.Get("game_id"),
				values.Get("title"), values.Get("rule"),
				values.Get("start_time"), values.Get("end_time"), values.Get("game_status"),
				values.Get("screen_open"), values.Get("raffle_frequency"), win, people, raffle)
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	return
}
