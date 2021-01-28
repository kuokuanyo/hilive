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

// GetCountdownPanel 取得倒數計時頁面、表單資訊
func (s *SystemTable) GetCountdownPanel(ctx *context.Context) (countdowntable Table) {
	// 建立BaseTable
	countdowntable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := countdowntable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("倒數計時秒數", "second", db.Int)
	info.AddField("倒數計時後進入頁面", "index_url", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "當前頁面"
			}
			return "3D簽到牆"
		})
	info.AddField("頭像形狀", "avatar_shape", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "圓形"
			}
			return "方形"
		})

	info.SetTable("activity_set_countdown").SetTitle("倒數計時").SetDescription("倒數計時管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_countdown").WhereIn("id", ids).Delete()
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
	formList := countdowntable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
	SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("倒數計時秒數", "second", db.Int, form.Number).SetFieldMust().SetFieldDefault("5")
	formList.AddField("倒數計時後進入頁面", "index_url", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "當前頁面"},
			{Value: "0", Text: "3D簽到牆"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_countdown").Select("index_url").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["index_url"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("頭像形狀", "avatar_shape", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "圓形"},
			{Value: "0", Text: "方形"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_countdown").Select("avatar_shape").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["avatar_shape"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")

	formList.SetTable("activity_set_countdown").SetTitle("倒數計時").SetDescription("倒數計時管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "second", "index_url", "avatar_shape") {
			return errors.New("活動ID、秒數、頭像形狀、導向頁面等欄位都不能為空")
		}

		if models.DefaultCountdownModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), "") {
			return errors.New("該活動已設置過倒數計時的基礎設定")
		}

		second, _ := strconv.Atoi(values.Get("second"))
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultCountdownModel().SetTx(tx).SetConn(s.conn).AddCountdown(
				values.Get("activity_id"), values.Get("index_url"),
				values.Get("avatar_shape"), second)
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
		if values.IsEmpty("activity_id", "second", "index_url", "avatar_shape") {
			return errors.New("活動ID、秒數、頭像形狀、導向頁面等欄位都不能為空")
		}

		if models.DefaultCountdownModel().SetConn(s.conn).
		IsActivityExist(values.Get("activity_id"), values.Get("id")) {
			return errors.New("該活動已設置過倒數計時的基礎設定")
		}

		second, _ := strconv.Atoi(values.Get("second"))
		model := models.GetCountdownModelAndID("activity_set_countdown", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).UpdateCountdown(values.Get("activity_id"), values.Get("index_url"),
				values.Get("avatar_shape"), second)
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
