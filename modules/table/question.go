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

// GetQuestionPanel 取得提問牆的頁面、表單資訊
func (s *SystemTable) GetQuestionPanel(ctx *context.Context) (questionTable Table) {
	// 建立BaseTable
	questionTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊
	info := questionTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("提問訊息審核", "message_check", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("匿名提問", "anonymous", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("隱藏已解答問題", "hide_answered", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("二維碼", "qrcode", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})

	info.SetTable("activity_set_question").SetTitle("提問牆").SetDescription("提問牆管理").
		// 刪除函式
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_question").WhereIn("id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除活動資料發生錯誤"), nil
					}
				}
				return nil, nil
			})
			return txErr
		})

	// 增加表單欄位資訊
	formList := questionTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust()
	formList.AddField("提問訊息審核", "message_check", db.Tinyint, form.SelectSingle).SetFieldMust().
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "開啟"},
			{Value: "0", Text: "關閉"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var allow []string
			if value.ID == "" {
				return allow
			}

			allowModel, _ := s.table("activity_set_question").Select("message_check").FindByID(value.ID)
			allow = append(allow, strconv.FormatInt(allowModel["message_check"].(int64), 10))
			return allow
		})
	formList.AddField("匿名提問", "anonymous", db.Tinyint, form.SelectSingle).SetFieldMust().
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "開啟"},
			{Value: "0", Text: "關閉"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var allow []string
			if value.ID == "" {
				return allow
			}

			allowModel, _ := s.table("activity_set_question").Select("anonymous").FindByID(value.ID)
			allow = append(allow, strconv.FormatInt(allowModel["anonymous"].(int64), 10))
			return allow
		})
	formList.AddField("隱藏已解答問題", "hide_answered", db.Tinyint, form.SelectSingle).SetFieldMust().
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "開啟"},
			{Value: "0", Text: "關閉"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var allow []string
			if value.ID == "" {
				return allow
			}

			allowModel, _ := s.table("activity_set_question").Select("hide_answered").FindByID(value.ID)
			allow = append(allow, strconv.FormatInt(allowModel["hide_answered"].(int64), 10))
			return allow
		})
	formList.AddField("二維碼", "qrcode", db.Tinyint, form.SelectSingle).SetFieldMust().
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "開啟"},
			{Value: "0", Text: "關閉"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var allow []string
			if value.ID == "" {
				return allow
			}

			allowModel, _ := s.table("activity_set_question").Select("qrcode").FindByID(value.ID)
			allow = append(allow, strconv.FormatInt(allowModel["qrcode"].(int64), 10))
			return allow
		})

	formList.SetTable("activity_set_question").SetTitle("提問牆").SetDescription("提問牆管理")
	// 提問牆基礎設置新增函式
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "message_check", "anonymous", "hide_answered", "qrcode") {
			return errors.New("活動ID、提問牆資訊、qrcode等欄位都不能為空")
		}

		if models.DefaultQuestionModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), "") {
			return errors.New("該活動已設置過提問牆的基礎設定")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultQuestionModel().SetTx(tx).SetConn(s.conn).AddQuestion(
				values.Get("activity_id"), values.Get("message_check"),
				values.Get("anonymous"), values.Get("hide_answered"),
				values.Get("qrcode"))
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 提問牆基礎設置更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "message_check", "anonymous", "hide_answered", "qrcode") {
			return errors.New("活動ID、提問牆資訊、qrcode等欄位都不能為空")
		}

		if models.DefaultQuestionModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), values.Get("id")) {
			return errors.New("該活動已設置過提問牆的基礎設定")
		}

		questionModel := models.GetQuestionModelAndID("activity_set_question", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := questionModel.SetTx(tx).UpdateQuestion(values.Get("activity_id"), values.Get("message_check"),
				values.Get("anonymous"), values.Get("hide_answered"),
				values.Get("qrcode"))
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
