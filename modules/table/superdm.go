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

// GetSuperdmPanel 取得超級彈幕頁面、表單資訊
func (s *SystemTable) GetSuperdmPanel(ctx *context.Context) (superdmTable Table) {
	// 建立BaseTable
	superdmTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := superdmTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("超級彈幕訊息審核", "message_check", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("醒目彈幕價格", "eye_catching_price", db.Int)
	info.AddField("大號彈幕價格", "large_danmu_price", db.Int)
	info.AddField("刷屏彈幕價格", "status_update_price", db.Int)
	info.AddField("圖片彈幕價格", "picture_danmu_price", db.Int)

	info.SetTable("activity_set_superdm").SetTitle("超級彈幕").SetDescription("超級彈幕管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_superdm").WhereIn("id", ids).Delete()
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
	formList := superdmTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
	SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("超級彈幕訊息審核", "message_check", db.Int, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "開啟", Value: "1"},
			{Text: "關閉", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var open []string
			if value.ID == "" {
				return []string{value.Value}
			}

			openModel, _ := s.table("activity_set_superdm").Select("message_check").FindByID(value.ID)
			open = append(open, strconv.FormatInt(openModel["message_check"].(int64), 10))
			return open
		}).SetFieldDefault("1")
	formList.AddField("醒目彈幕價格", "eye_catching_price", db.Int, form.Number).SetFieldMust().SetFieldDefault("5")
	formList.AddField("大號彈幕價格", "large_danmu_price", db.Int, form.Number).SetFieldMust().SetFieldDefault("8")
	formList.AddField("刷屏彈幕價格", "status_update_price", db.Int, form.Number).SetFieldMust().SetFieldDefault("8")
	formList.AddField("圖片彈幕價格", "picture_danmu_price", db.Int, form.Number).SetFieldMust().SetFieldDefault("10")

	formList.SetTable("activity_set_superdm").SetTitle("超級彈幕").SetDescription("超級彈幕管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "message_check", "eye_catching_price",
			"large_danmu_price", "status_update_price", "picture_danmu_price") {
			return errors.New("活動ID、審核、價格等欄位都不能為空")
		}

		if models.DefaultSuperdmModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), "") {
			return errors.New("此活動已建立超級彈幕基礎設置")
		}

		eye, _ := strconv.Atoi(values.Get("eye_catching_price"))
		large, _ := strconv.Atoi(values.Get("large_danmu_price"))
		statusUpdate, _ := strconv.Atoi(values.Get("status_update_price"))
		picture, _ := strconv.Atoi(values.Get("picture_danmu_price"))

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultSuperdmModel().SetTx(tx).SetConn(s.conn).AddSuperdm(
				values.Get("activity_id"), values.Get("message_check"), eye, large,
				statusUpdate, picture)
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
		if values.IsEmpty("activity_id", "message_check", "eye_catching_price",
			"large_danmu_price", "status_update_price", "picture_danmu_price") {
			return errors.New("活動ID、審核、價格等欄位都不能為空")
		}

		eye, _ := strconv.Atoi(values.Get("eye_catching_price"))
		large, _ := strconv.Atoi(values.Get("large_danmu_price"))
		statusUpdate, _ := strconv.Atoi(values.Get("status_update_price"))
		picture, _ := strconv.Atoi(values.Get("picture_danmu_price"))

		model := models.GetSuperdmModelAndID("activity_set_superdm", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).UpdateSuperdm(values.Get("activity_id"), values.Get("message_check"), eye, large,
				statusUpdate, picture)
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
