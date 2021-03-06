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

// GetIntroducePanel 取得活動介紹的頁面、表單資訊
func (s *SystemTable) GetIntroducePanel(ctx *context.Context) (introduceTable Table) {
	// 建立BaseTable
	introduceTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊
	info := introduceTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("介紹標題", "introduce_title", db.Varchar)
	info.AddField("介紹內容", "introduce_content", db.Varchar)
	info.AddField("介紹排序", "introduce_order", db.Int)

	info.SetTable("activity_introduce").SetTitle("活動介紹").SetDescription("活動介紹管理").
		// 刪除函式
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_introduce").WhereIn("id", ids).Delete()
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
	formList := introduceTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("介紹標題", "introduce_title", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("介紹內容", "introduce_content", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("介紹排序", "introduce_order", db.Int, form.Text).
		SetFieldHelpMsg(template.HTML("請輸入數字設置活動介紹的排序")).FieldNotAllowAdd()

	formList.SetTable("activity_introduce").SetTitle("活動介紹").SetDescription("活動介紹管理")
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "introduce_title", "introduce_content") {
			return errors.New("活動ID、介紹標題、內容等欄位都不能為空")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultIntroduceModel().SetTx(tx).SetConn(s.conn).AddIntroduce(
				values.Get("activity_id"), values.Get("introduce_title"), values.Get("introduce_content"))
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 設置活動介紹更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "introduce_title", "introduce_content", "introduce_order") {
			return errors.New("活動ID、介紹標題、內容、排序等欄位都不能為空")
		}

		order, _ := strconv.Atoi(values.Get("introduce_order"))
		introduceModel := models.GetIntroduceModelAndID("activity_introduce", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := introduceModel.SetTx(tx).UpdateIntroduce(values.Get("activity_id"),
				values.Get("introduce_title"), values.Get("introduce_content"), order)
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
