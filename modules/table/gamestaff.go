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

// GetGameStaffPanel 取得參加活動人員頁面、表單資訊
func (s *SystemTable) GetGameStaffPanel(ctx *context.Context) (staffTable Table) {
	// 建立BaseTable
	staffTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := staffTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("用戶ID", "user_id", db.Varchar).FieldFilterable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("遊戲專屬ID", "game_id", db.Varchar).FieldFilterable()
	info.AddField("報名狀態", "apply_status", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "報名成功"
			}
			return "報名失敗"
		})
	info.AddField("報名活動時間", "apply_time", db.Datetime)

	info.SetTable("activity_apply_game_staff").SetTitle("參加遊戲人員").SetDescription("人員管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_apply_game_staff").WhereIn("id", ids).Delete()
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
	formList.AddField("報名狀態", "apply_status", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "報名成功"},
			{Value: "0", Text: "報名失敗"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_apply_game_staff").
				Select("apply_status").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["apply_status"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")

	formList.SetTable("activity_apply_game_staff").SetTitle("參加遊戲人員").SetDescription("人員管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("user_id", "activity_id", "game_id", "apply_status") {
			return errors.New("人員、活動、遊戲ID、報名狀態等欄位都不能為空")
		}

		if models.DefaultGameStaffModel().SetConn(s.conn).IsStaffExist(
			values.Get("user_id"), values.Get("game_id"), "") {
			return errors.New("人員已報名該活動")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultGameStaffModel().SetTx(tx).SetConn(s.conn).AddStaff(
				values.Get("user_id"), values.Get("activity_id"), values.Get("game_id"),
				values.Get("apply_status"))
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
		if values.IsEmpty("user_id", "activity_id", "game_id", "apply_status") {
			return errors.New("人員、活動、遊戲ID、報名狀態等欄位都不能為空")
		}

		if models.DefaultGameStaffModel().SetConn(s.conn).IsStaffExist(
			values.Get("user_id"), values.Get("game_id"), values.Get("id")) {
			return errors.New("人員已報名該活動")
		}

		model := models.GetGameStaffModelAndID("activity_apply_game_staff", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).UpdateStaff(values.Get("user_id"), values.Get("activity_id"), values.Get("game_id"),
				values.Get("apply_status"))
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
