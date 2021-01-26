package table

import (
	"database/sql"
	"errors"
	"hilive/context"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/template/form"
	"hilive/template/types"
	"html/template"
	"strconv"
)

// GetActivityOverviewPanel 取得活動總覽頁面、表單資訊
func (s *SystemTable) GetActivityOverviewPanel(ctx *context.Context) (overviewTable Table) {
	// 建立BaseTable
	overviewTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := overviewTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("遊戲類型", "type_name", db.Varchar).FieldFilterable().
		FieldJoin(types.Join{
			JoinTable: "activity_game",
			JoinField: "id",
			BaseTable: "activity_game_open",
			Field:     "game_id",
		}).
		FieldJoin(types.Join{
			JoinTable: "activity_game_type",
			JoinField: "type_id",
			BaseTable: "activity_game",
			Field:     "game_type",
		})
	info.AddField("遊戲名稱", "game_name", db.Varchar).FieldFilterable().
		FieldJoin(types.Join{
			JoinTable: "activity_game",
			JoinField: "id",
			BaseTable: "activity_game_open",
			Field:     "game_id",
		})
	info.AddField("是否開啟", "open", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})

	info.SetTable("activity_game_open").SetTitle("活動總覽").SetDescription("總覽管理").
		// 刪除函式
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_game_open").WhereIn("id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除活動資料發生錯誤"), nil
					}
				}
				return nil, nil
			})
			return txErr
		})

	// 增加表單資訊欄位
	formList := overviewTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
	SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("遊戲名稱", "game_id", db.Varchar, form.SelectSingle).SetFieldMust().
		SetFieldOptionFromTable("activity_game", "game_name", "id").
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			var activityGame []string
			if model.ID == "" {
				return activityGame
			}

			gameModel, _ := s.table("activity_game_open").Select("game_id").FindByID(model.ID)
			activityGame = append(activityGame, strconv.FormatInt(gameModel["game_id"].(int64), 10))
			return activityGame
		})
	formList.AddField("是否開啟遊戲", "open", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "開啟", Value: "1"},
			{Text: "關閉", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var open []string
			if value.ID == "" {
				return []string{value.Value}
			}

			openModel, _ := s.table("activity_game_open").Select("open").FindByID(value.ID)
			open = append(open, strconv.FormatInt(openModel["open"].(int64), 10))
			return open
		}).SetFieldDefault("1")
		
	formList.SetTable("activity_game_open").SetTitle("活動總覽").SetDescription("總覽管理")
	// 設置活動新增函式
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id", "open") {
			return errors.New("活動ID、遊戲名稱、是否開啟等欄位都不能為空")
		}

		if models.DefaultOverviewModel().SetConn(s.conn).IsGameExist(values.Get("game_id"), values.Get("activity_id"), "") {
			return errors.New("此活動已建立該遊戲")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 新增活動資料
			_, err := models.DefaultOverviewModel().SetTx(tx).SetConn(s.conn).AddActivityOverview(
				values.Get("activity_id"), values.Get("game_id"), values.Get("open"))
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 設置活動更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "game_id", "open") {
			return errors.New("活動ID、遊戲名稱、是否開啟等欄位都不能為空")
		}
		if models.DefaultOverviewModel().SetConn(s.conn).IsGameExist(values.Get("game_id"), values.Get("activity_id"), values.Get("id")) {
			return errors.New("此活動已建立該遊戲")
		}

		gameModel := models.GetOverviewModelAndID("activity_game_open", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 更新活動總覽資料
			_, err := gameModel.SetTx(tx).UpdateActivityOverview(values.Get("activity_id"), values.Get("game_id"), values.Get("open"))
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
