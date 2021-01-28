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

// GetPrizeStaffPanel 取得參加中獎人員頁面、表單資訊
func (s *SystemTable) GetPrizeStaffPanel(ctx *context.Context) (staffTable Table) {
	// 建立BaseTable
	staffTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := staffTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("用戶ID", "user_id", db.Varchar).FieldFilterable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("遊戲專屬ID", "game_id", db.Varchar).FieldFilterable()
	info.AddField("中獎時間", "win_time", db.Datetime)
	info.AddField("獎品名稱", "prize_name", db.Varchar)
	info.AddField("兌獎方式", "redeem_method", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "現場兌獎"
			}
			return "郵寄兌獎"
		})
	info.AddField("兌獎狀態", "redeem_status", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "已領獎"
			}
			return "未領獎"
		})
	info.AddField("兌獎密碼", "redeem_password", db.Varchar)

	info.SetTable("activity_prize_staff").SetTitle("中獎人員").SetDescription("人員管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_prize_staff").WhereIn("id", ids).Delete()
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
	formList := staffTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("用戶ID", "user_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("用戶辨別ID")).SetFieldMust().
		FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().
		FieldNotAllowEdit()
	formList.AddField("遊戲專屬ID", "game_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("遊戲辨別ID")).SetFieldMust().
		FieldNotAllowEdit()
	formList.AddField("中獎時間", "win_time", db.Datetime, form.Datetime).SetFieldMust()
	formList.AddField("獎品名稱", "prize_name", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("兌獎方式", "redeem_method", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "現場兌獎"},
			{Value: "0", Text: "郵寄兌獎"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_prize_staff").
				Select("redeem_method").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["redeem_method"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("兌獎狀態", "redeem_status", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "已領獎"},
			{Value: "0", Text: "未領獎"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_prize_staff").
				Select("redeem_status").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["redeem_status"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("兌獎密碼", "redeem_password", db.Varchar, form.Text)

	formList.SetTable("activity_prize_staff").SetTitle("中獎人員").SetDescription("人員管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("user_id", "activity_id", "game_id", "win_time",
			"prize_name", "redeem_method", "redeem_status") {
			return errors.New("人員、活動、遊戲ID、時間、兌獎狀態等欄位都不能為空")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultPrizeStaffModel().SetTx(tx).SetConn(s.conn).AddStaff(
				values.Get("user_id"), values.Get("activity_id"), values.Get("game_id"),
				values.Get("win_time"), values.Get("prize_name"), values.Get("redeem_method"),
				values.Get("redeem_status"), values.Get("redeem_password"))
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
		if values.IsEmpty("user_id", "activity_id", "game_id", "win_time",
			"prize_name", "redeem_method", "redeem_status") {
			return errors.New("人員、活動、遊戲ID、時間、兌獎狀態等欄位都不能為空")
		}

		model := models.GetPrizeStaffModelAndID("activity_prize_staff", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).UpdateStaff(values.Get("user_id"), values.Get("activity_id"), values.Get("game_id"),
				values.Get("win_time"), values.Get("prize_name"), values.Get("redeem_method"),
				values.Get("redeem_status"), values.Get("redeem_password"))
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
