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
	"time"
)

// GetPicturePanel 取得圖片牆頁面、表單資訊
func (s *SystemTable) GetPicturePanel(ctx *context.Context) (pictureTable Table) {
	// 建立BaseTable
	pictureTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := pictureTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("活動開始時間", "start_time", db.Datetime)
	info.AddField("活動結束時間", "end_time", db.Datetime)
	info.AddField("屏幕間隔秒數", "switch_second", db.Int)
	info.AddField("屏幕播放順序", "play_order", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "順序播放"
			}
			return "隨機播放"
		})
	info.AddField("圖片路徑", "picture_path", db.Varchar).
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
	info.AddField("屏幕背景", "background", db.Varchar)

	info.SetTable("activity_set_picture").SetTitle("圖片牆").SetDescription("圖片管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_picture").WhereIn("id", ids).Delete()
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
	formList := pictureTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
	SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("活動開始時間", "start_time", db.Datetime, form.Datetime).SetFieldMust()
	formList.AddField("活動結束時間", "end_time", db.Datetime, form.Datetime).SetFieldMust()
	formList.AddField("屏幕間隔秒數", "switch_second", db.Int, form.Number).SetFieldMust().SetFieldDefault("5")
	formList.AddField("屏幕播放順序", "play_order", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Text: "順序播放", Value: "1"},
			{Text: "隨機播放", Value: "0"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var status []string
			if value.ID == "" {
				return []string{value.Value}
			}

			model, _ := s.table("activity_set_picture").Select("play_order").FindByID(value.ID)
			status = append(status, strconv.FormatInt(model["play_order"].(int64), 10))
			return status
		}).SetFieldDefault("1")
	formList.AddField("圖片路徑", "picture_path", db.Varchar, form.TextArea).
		SetFieldHelpMsg(template.HTML("請一行設置一個圖片路徑，若要輸入新路徑請換行輸入"))
	formList.AddField("屏幕背景", "background", db.Varchar, form.Text)

	formList.SetTable("activity_set_picture").SetTitle("圖片牆").SetDescription("圖片管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "start_time", "end_time", "switch_second", "play_order") {
			return errors.New("活動ID、時間、秒數等欄位都不能為空")
		}

		if models.DefaultPictureModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), "") {
			return errors.New("該活動已設置過圖片牆的基礎設定")
		}

		// 時間判斷
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("start_time"), time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("end_time"), time.Local)
		boolTime := end.After(start) && start.Before(end)
		if boolTime == false {
			return errors.New("時間設置發生錯誤，請重新設置(結束時間在開始時間之後)")
		}

		second, _ := strconv.Atoi(values.Get("switch_second"))
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultPictureModel().SetTx(tx).SetConn(s.conn).AddPicture(
				values.Get("activity_id"), values.Get("start_time"),
				values.Get("end_time"), second, values.Get("play_order"),
				values.Get("picture_path"), values.Get("background"))
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
		if values.IsEmpty("activity_id", "start_time", "end_time", "switch_second", "play_order") {
			return errors.New("活動ID、時間、秒數等欄位都不能為空")
		}

		// 時間判斷
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("start_time"), time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("end_time"), time.Local)
		boolTime := end.After(start) && start.Before(end)
		if boolTime == false {
			return errors.New("時間設置發生錯誤，請重新設置(結束時間在開始時間之後)")
		}

		second, _ := strconv.Atoi(values.Get("switch_second"))
		model := models.GetPictureModelAndID("activity_set_picture", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).UpdatePicture(values.Get("activity_id"), values.Get("start_time"),
				values.Get("end_time"), second, values.Get("play_order"),
				values.Get("picture_path"), values.Get("background"))
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
