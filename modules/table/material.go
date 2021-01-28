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

// GetMaterialPanel 取得活動資料頁面、表單資訊
func (s *SystemTable) GetMaterialPanel(ctx *context.Context) (materialTable Table) {
	// 建立BaseTable
	materialTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := materialTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("資料名稱", "data_name", db.Varchar)
	info.AddField("資料說明", "data_introduce", db.Varchar)
	info.AddField("資料連結", "data_link", db.Varchar)
	info.AddField("資料排序", "data_order", db.Int)

	info.SetTable("activity_material").SetTitle("活動資料").SetDescription("活動資料管理").
		// 刪除函式
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_material").WhereIn("id", ids).Delete()
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
	formList := materialTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).SetFieldHelpMsg(template.HTML("活動辨別ID")).
		SetFieldMust().FieldNotAllowEdit()
	formList.AddField("資料名稱", "data_name", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("資料說明", "data_introduce", db.Varchar, form.Text)
	formList.AddField("資料連結", "data_link", db.Varchar, form.Text)
	formList.AddField("資料排序", "data_order", db.Int, form.Text).
		SetFieldHelpMsg(template.HTML("請輸入數字設置活動資料的排序")).SetFieldMust().FieldNotAllowAdd()

	formList.SetTable("activity_material").SetTitle("活動資料").SetDescription("活動資料管理")

	// 設置資料新增函式
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "data_name") {
			return errors.New("活動ID、資料名稱、排序等欄位都不能為空")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultMaterialModel().SetTx(tx).SetConn(s.conn).AddMaterial(
				values.Get("activity_id"), values.Get("data_name"), values.Get("data_introduce"),
				values.Get("data_link"))
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 設置資料更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "data_name", "data_order") {
			return errors.New("活動ID、資料名稱、排序等欄位都不能為空")
		}

		order, _ := strconv.Atoi(values.Get("data_order"))
		materialModel := models.GetMaterialModelAndID("activity_material", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 更新用戶資料
			_, err := materialModel.SetTx(tx).UpdateMaterial(values.Get("activity_id"),
				values.Get("data_name"), values.Get("data_introduce"),
				values.Get("data_link"), order)
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
