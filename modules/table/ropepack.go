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
)

// GetRopepackPanel 取得套紅包頁面、表單資訊
func (s *SystemTable) GetRopepackPanel(ctx *context.Context) (gameTable Table) {
	// 建立BaseTable
	gameTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := gameTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("遊戲專屬ID", "game_id", db.Varchar).FieldFilterable()
	info.AddField("遊戲標題", "title", db.Varchar)
	info.AddField("中獎機率(%)", "percent", db.Int)
	info.AddField("限時時間(秒)", "second", db.Int)
	info.AddField("允許重複搖中", "allow_repeat", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "允許"
			}
			return "禁止"
		})
	info.AddField("遊戲狀態", "game_status", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})

	info.SetTable("activity_set_ropepack").SetTitle("套紅包").SetDescription("遊戲管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_ropepack").WhereIn("id", ids).Delete()
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
	formList := gameTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("遊戲專屬ID", "game_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("遊戲辨別ID")).SetFieldMust().
		FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("遊戲標題", "title", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("中獎機率(%)", "percent", db.Int, form.Number).
		SetFieldMust().SetFieldDefault("20").
		SetFieldHelpMsg(template.HTML("數值請設置0~100區間內"))
	formList.AddField("限時時間(秒)", "second", db.Int, form.Number).
		SetFieldMust().SetFieldDefault("30")
	formList.AddField("允許重複搖中", "allow_repeat", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "允許"},
			{Value: "0", Text: "禁止"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_ropepack").Select("allow_repeat").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["allow_repeat"].(int64), 10))
			return stauts
		}).SetFieldDefault("0")
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

			statusModel, _ := s.table("activity_set_ropepack").Select("game_status").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["game_status"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")

	formList.SetTable("activity_set_ropepack").SetTitle("套紅包").SetDescription("遊戲管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "title", "second", "allow_repeat",
			"game_status", "percent") {
			return errors.New("活動ID、標題、遊戲設置等欄位都不能為空")
		}

		second, _ := strconv.Atoi(values.Get("second"))
		percent, _ := strconv.Atoi(values.Get("percent"))
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultRopepackModel().SetTx(tx).SetConn(s.conn).AddRopepack(
				values.Get("activity_id"), utils.UUID(8), values.Get("title"),
				values.Get("allow_repeat"), values.Get("game_status"), second, percent)
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
		if values.IsEmpty("activity_id", "game_id", "title", "second", "allow_repeat",
			"game_status", "percent") {
			return errors.New("活動ID、標題、遊戲設置等欄位都不能為空")
		}

		second, _ := strconv.Atoi(values.Get("second"))
		percent, _ := strconv.Atoi(values.Get("percent"))
		model := models.GetRopepackModelAndID("activity_set_ropepack", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).UpdateRopepack(values.Get("activity_id"), values.Get("game_id"),
				values.Get("title"), values.Get("allow_repeat"),
				values.Get("game_status"), second, percent)
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
