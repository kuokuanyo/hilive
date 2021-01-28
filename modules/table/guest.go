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
	"strconv"
)

// GetGuestPanel 取得活動嘉賓頁面、表單資訊
func (s *SystemTable) GetGuestPanel(ctx *context.Context) (guestTable Table) {
	// 建立BaseTable
	guestTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := guestTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("嘉賓照片", "guest_picture", db.Varchar)
	info.AddField("嘉賓名稱", "guest_name", db.Varchar)
	info.AddField("嘉賓簡介", "guest_introduce", db.Varchar)
	info.AddField("嘉賓詳情", "guest_detail", db.Varchar)
	info.AddField("嘉賓排序", "guest_order", db.Int)

	info.SetTable("activity_guest").SetTitle("活動嘉賓").SetDescription("活動嘉賓管理").
		// 刪除函式
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_guest").WhereIn("id", ids).Delete()
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
	formList := guestTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
	SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("嘉賓照片", "guest_picture", db.Varchar, form.Text)
	formList.AddField("嘉賓名稱", "guest_name", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("嘉賓簡介", "guest_introduce", db.Varchar, form.Text)
	formList.AddField("嘉賓詳情", "guest_detail", db.Varchar, form.Text)
	formList.AddField("嘉賓排序", "guest_order", db.Int, form.Text).SetFieldMust().
		SetFieldHelpMsg(template.HTML("請輸入數字設置活動嘉賓的排序")).FieldNotAllowAdd()

	formList.SetTable("activity_guest").SetTitle("活動嘉賓").SetDescription("活動嘉賓管理")

	// 設置嘉賓新增函式
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "guest_name") {
			return errors.New("活動ID、嘉賓名稱、排序等欄位都不能為空")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 新增嘉賓資料
			_, err := models.DefaultGuestModel().SetTx(tx).SetConn(s.conn).AddGuest(
				values.Get("activity_id"), values.Get("guest_picture"), values.Get("guest_name"),
				values.Get("guest_introduce"), values.Get("guest_detail"))
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 設置嘉賓更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "guest_name", "guest_order") {
			return errors.New("活動ID、嘉賓名稱、排序等欄位都不能為空")
		}

		order, _ := strconv.Atoi(values.Get("guest_order"))

		guestModel := models.GetGuestModelAndID("activity_guest", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 更新用戶資料
			_, err := guestModel.SetTx(tx).UpdateGuest(values.Get("activity_id"), values.Get("guest_picture"), values.Get("guest_name"),
				values.Get("guest_introduce"), values.Get("guest_detail"), order)
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
