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

// GetHoldScreenPanel 取得霸屏頁面、表單資訊
func (s *SystemTable) GetHoldScreenPanel(ctx *context.Context) (holdscreenTable Table) {
	// 建立BaseTable
	holdscreenTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := holdscreenTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("霸屏每秒價格", "holdscreen_price", db.Int)
	info.AddField("霸屏消息審核", "message_check", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("只允許以霸屏形式發送圖片", "only_picture", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("霸屏最低秒數", "minimum_second", db.Int)
	info.AddField("生日主題", "birthday_topic", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("表白主題", "confess_topic", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("求婚主題", "propose_topic", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("祝福主題", "bless_topic", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("女神主題", "goddess_topic", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})

	info.SetTable("activity_set_holdscreen").SetTitle("霸屏").SetDescription("霸屏管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_holdscreen").WhereIn("id", ids).Delete()
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
	formList := holdscreenTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust()
	formList.AddField("霸屏每秒價格", "holdscreen_price", db.Int, form.Number).SetFieldMust().SetFieldDefault("2")
	formList.AddField("霸屏消息審核", "message_check", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "開啟", Value: "1"},
			{Text: "關閉", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var open []string
			if value.ID == "" {
				return []string{value.Value}
			}

			openModel, _ := s.table("activity_set_holdscreen").Select("message_check").FindByID(value.ID)
			open = append(open, strconv.FormatInt(openModel["message_check"].(int64), 10))
			return open
		}).SetFieldDefault("0")
	formList.AddField("只允許以霸屏形式發送圖片", "only_picture", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "開啟", Value: "1"},
			{Text: "關閉", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var status []string
			if value.ID == "" {
				return []string{value.Value}
			}

			model, _ := s.table("activity_set_holdscreen").Select("only_picture").FindByID(value.ID)
			status = append(status, strconv.FormatInt(model["only_picture"].(int64), 10))
			return status
		}).SetFieldDefault("0")
	formList.AddField("霸屏最低秒數", "minimum_second", db.Int, form.Number).SetFieldMust().SetFieldDefault("30")
	formList.AddField("生日主題", "birthday_topic", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "開啟", Value: "1"},
			{Text: "關閉", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var status []string
			if value.ID == "" {
				return []string{value.Value}
			}

			model, _ := s.table("activity_set_holdscreen").Select("birthday_topic").FindByID(value.ID)
			status = append(status, strconv.FormatInt(model["birthday_topic"].(int64), 10))
			return status
		}).SetFieldDefault("1")
	formList.AddField("表白主題", "confess_topic", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "開啟", Value: "1"},
			{Text: "關閉", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var status []string
			if value.ID == "" {
				return []string{value.Value}
			}

			model, _ := s.table("activity_set_holdscreen").Select("confess_topic").FindByID(value.ID)
			status = append(status, strconv.FormatInt(model["confess_topic"].(int64), 10))
			return status
		}).SetFieldDefault("1")
	formList.AddField("求婚主題", "propose_topic", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "開啟", Value: "1"},
			{Text: "關閉", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var status []string
			if value.ID == "" {
				return []string{value.Value}
			}

			model, _ := s.table("activity_set_holdscreen").Select("propose_topic").FindByID(value.ID)
			status = append(status, strconv.FormatInt(model["propose_topic"].(int64), 10))
			return status
		}).SetFieldDefault("1")
	formList.AddField("祝福主題", "bless_topic", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "開啟", Value: "1"},
			{Text: "關閉", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var status []string
			if value.ID == "" {
				return []string{value.Value}
			}

			model, _ := s.table("activity_set_holdscreen").Select("bless_topic").FindByID(value.ID)
			status = append(status, strconv.FormatInt(model["bless_topic"].(int64), 10))
			return status
		}).SetFieldDefault("1")
	formList.AddField("女神主題", "goddess_topic", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "開啟", Value: "1"},
			{Text: "關閉", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var status []string
			if value.ID == "" {
				return []string{value.Value}
			}

			model, _ := s.table("activity_set_holdscreen").Select("goddess_topic").FindByID(value.ID)
			status = append(status, strconv.FormatInt(model["goddess_topic"].(int64), 10))
			return status
		}).SetFieldDefault("1")

	formList.SetTable("activity_set_holdscreen").SetTitle("霸屏").SetDescription("霸屏管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "holdscreen_price", "message_check", "only_picture", "minimum_second",
			"birthday_topic", "confess_topic", "propose_topic", "bless_topic", "goddess_topic") {
			return errors.New("活動ID、秒數、主題等欄位都不能為空")
		}

		if models.DefaultHoldScreenModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), "") {
			return errors.New("該活動已設置過霸屏的基礎設定")
		}

		price, _ := strconv.Atoi(values.Get("holdscreen_price"))
		second, _ := strconv.Atoi(values.Get("minimum_second"))
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultHoldScreenModel().SetTx(tx).SetConn(s.conn).AddHoldScreen(
				values.Get("activity_id"), price, values.Get("message_check"),
				values.Get("only_picture"), second, values.Get("birthday_topic"),
				values.Get("confess_topic"), values.Get("propose_topic"), values.Get("bless_topic"),
				values.Get("goddess_topic"))
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
		if values.IsEmpty("activity_id", "holdscreen_price", "message_check", "only_picture", "minimum_second",
			"birthday_topic", "confess_topic", "propose_topic", "bless_topic", "goddess_topic") {
			return errors.New("活動ID、秒數、主題等欄位都不能為空")
		}

		if models.DefaultHoldScreenModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), values.Get("id")) {
			return errors.New("該活動已設置過霸屏的基礎設定")
		}

		price, _ := strconv.Atoi(values.Get("holdscreen_price"))
		second, _ := strconv.Atoi(values.Get("minimum_second"))
		model := models.GetHoldScreenModelAndID("activity_set_holdscreen", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).UpdateHoldScreen(values.Get("activity_id"), price, values.Get("message_check"),
				values.Get("only_picture"), second, values.Get("birthday_topic"),
				values.Get("confess_topic"), values.Get("propose_topic"), values.Get("bless_topic"),
				values.Get("goddess_topic"))
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
