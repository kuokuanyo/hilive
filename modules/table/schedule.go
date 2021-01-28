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
	"html/template"
	"time"
)

// GetSchedulePanel 取得活動行程的頁面、表單資訊
func (s *SystemTable) GetSchedulePanel(ctx *context.Context) (scheduleTable Table) {
	// 建立BaseTable
	scheduleTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊
	info := scheduleTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("行程名稱", "schedule_name", db.Varchar)
	info.AddField("行程內容", "schedule_content", db.Varchar)
	info.AddField("行程開始時間", "start_time", db.Datetime)
	info.AddField("行程結束時間", "end_time", db.Datetime)

	info.SetTable("activity_schedule").SetTitle("活動行程").SetDescription("活動行程管理").
		// 刪除函式
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_schedule").WhereIn("id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除活動資料發生錯誤"), nil
					}
				}
				return nil, nil
			})
			return txErr
		})

	// 增加表單資訊
	formList := scheduleTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).SetFieldMust().
	SetFieldHelpMsg(template.HTML("活動辨別ID")).FieldNotAllowEdit()
	formList.AddField("行程名稱", "schedule_name", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("行程內容", "schedule_content", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("行程開始時間", "start_time", db.Datetime, form.Datetime).SetFieldMust()
	formList.AddField("行程結束時間", "end_time", db.Datetime, form.Datetime).SetFieldMust()

	formList.SetTable("activity_schedule").SetTitle("活動行程").SetDescription("活動行程管理")

	// 設置活動新增函式
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "schedule_name", "schedule_content", "start_time", "end_time") {
			return errors.New("活動ID、行程名稱內容、行程時間等欄位都不能為空")
		}

		// 時間判斷
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("start_time"), time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("end_time"), time.Local)
		boolTime := end.After(start) && start.Before(end)
		if boolTime == false {
			return errors.New("時間設置發生錯誤，請重新設置(結束時間在開始時間之後)")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 新增活動資料
			_, err := models.DefaultScheduleModel().SetTx(tx).SetConn(s.conn).AddSchedule(
				values.Get("activity_id"), values.Get("schedule_name"), values.Get("schedule_content"),
				values.Get("start_time"), values.Get("end_time"))
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 設置活動行程更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "schedule_name", "schedule_content", "start_time", "end_time") {
			return errors.New("活動ID、行程名稱內容、行程時間等欄位都不能為空")
		}

		// 時間判斷
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("start_time"), time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("end_time"), time.Local)
		boolTime := end.After(start) && start.Before(end)
		if boolTime == false {
			return errors.New("時間設置發生錯誤，請重新設置(結束時間在開始時間之後)")
		}

		scheduleModel := models.GetScheduleModelAndID("activity_schedule", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 更新活動行程資料
			_, err := scheduleModel.SetTx(tx).UpdateSchedule(values.Get("activity_id"), values.Get("schedule_name"),
				values.Get("schedule_content"), values.Get("start_time"), values.Get("end_time"))
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
