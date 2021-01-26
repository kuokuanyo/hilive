package table

import (
	"database/sql"
	"errors"
	"fmt"
	"hilive/context"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/template/form"
	"hilive/template/types"
	"html/template"
	"strconv"
	"strings"
)

// GetMessagePanel 取得訊息牆的頁面、表單資訊
func (s *SystemTable) GetMessagePanel(ctx *context.Context) (messageTable Table) {
	// 建立BaseTable
	messageTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊
	info := messageTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("允許圖片上訊息牆", "picture_message", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "允許"
			}
			return "禁止"
		})
	info.AddField("允許圖片自動上訊息牆", "picture_auto", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "允許"
			}
			return "禁止"
		})
	info.AddField("畫面自動刷新秒數", "refresh_second", db.Int)
	info.AddField("防刷屏", "prevent_status_update", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("跑馬燈訊息", "message", db.Varchar).
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			pathArr := strings.Split(model.Value, "\n")
			res := ""
			for i := 0; i < len(pathArr); i++ {
				if i == len(pathArr)-1 {
					res += string(template.HTML(fmt.Sprintf(`<span class="label label-success" style="background-color: ;">%s</span>`, pathArr[i])))
				} else {
					res += string(template.HTML(fmt.Sprintf(`<span class="label label-success" style="background-color: ;">%s</span>`, pathArr[i]) + "<br><br>"))
				}
			}
			return res
		})

	info.SetTable("activity_set_message").SetTitle("訊息牆").SetDescription("訊息牆管理").
		// 刪除函式
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_message").WhereIn("id", ids).Delete()
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
	formList := messageTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
	SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("允許圖片上訊息牆", "picture_message", db.Int, form.Radio).
	SetFieldOptions(types.FieldOptions{
		{Text: "允許", Value: "1"},
		{Text: "禁止", Value: "0"},
	}).SetFieldMust().
	SetDisplayFunc(func(value types.FieldModel) interface{} {
		var stauts []string
		if value.ID == "" {
			return []string{value.Value}
		}

		statusModel, _ := s.table("activity_set_message").Select("picture_message").FindByID(value.ID)
		stauts = append(stauts, strconv.FormatInt(statusModel["picture_message"].(int64), 10))
		return stauts
	}).SetFieldDefault("1")
	formList.AddField("允許圖片自動上訊息牆", "picture_auto", db.Int, form.Radio).
	SetFieldOptions(types.FieldOptions{
		{Text: "允許", Value: "1"},
		{Text: "禁止", Value: "0"},
	}).SetFieldMust().
	SetDisplayFunc(func(value types.FieldModel) interface{} {
		var stauts []string
		if value.ID == "" {
			return []string{value.Value}
		}

		statusModel, _ := s.table("activity_set_message").Select("picture_auto").FindByID(value.ID)
		stauts = append(stauts, strconv.FormatInt(statusModel["picture_auto"].(int64), 10))
		return stauts
	}).SetFieldDefault("1")
	formList.AddField("畫面自動刷新秒數", "refresh_second", db.Int, form.Number).SetFieldMust().SetFieldDefault("5")
	formList.AddField("防刷屏", "prevent_status_update", db.Int, form.Radio).
	SetFieldOptions(types.FieldOptions{
		{Text: "開啟", Value: "1"},
		{Text: "關閉", Value: "0"},
	}).SetFieldMust().
	SetDisplayFunc(func(value types.FieldModel) interface{} {
		var stauts []string
		if value.ID == "" {
			return []string{value.Value}
		}

		statusModel, _ := s.table("activity_set_message").Select("prevent_status_update").FindByID(value.ID)
		stauts = append(stauts, strconv.FormatInt(statusModel["prevent_status_update"].(int64), 10))
		return stauts
	}).SetFieldDefault("1")
	formList.AddField("跑馬燈訊息", "message", db.Varchar, form.TextArea).
		SetFieldHelpMsg(template.HTML("請一行設置一個跑馬燈訊息，若要輸入新的跑馬燈請換行"))

	formList.SetTable("activity_set_message").SetTitle("訊息牆").SetDescription("訊息牆管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "picture_message", "picture_auto", "refresh_second", "prevent_status_update") {
			return errors.New("活動ID、訊息牆資訊、刷新秒數等欄位都不能為空")
		}

		if models.DefaultMessageModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), "") {
			return errors.New("該活動已設置過訊息牆的基礎設定")
		}

		second, _ := strconv.Atoi(values.Get("refresh_second"))
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 新增活動資料
			_, err := models.DefaultMessageModel().SetTx(tx).SetConn(s.conn).AddMessage(
				values.Get("activity_id"), values.Get("picture_message"),
				values.Get("picture_auto"), values.Get("prevent_status_update"),
				values.Get("message"), second)
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
		if values.IsEmpty("activity_id", "picture_message", "picture_auto", "refresh_second", "prevent_status_update") {
			return errors.New("活動ID、訊息牆資訊、刷新秒數等欄位都不能為空")
		}

		second, _ := strconv.Atoi(values.Get("refresh_second"))
		messageModel := models.GetMessageModelAndID("activity_set_message", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := messageModel.SetTx(tx).UpdateActivityMessage(values.Get("activity_id"),
				values.Get("picture_message"), values.Get("picture_auto"),
				values.Get("prevent_status_update"), values.Get("message"), second)
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
