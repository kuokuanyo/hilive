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
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust()
	formList.AddField("介紹標題", "introduce_title", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("介紹內容", "introduce_content", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("介紹排序", "introduce_order", db.Int, form.Text).
		SetFieldHelpMsg(template.HTML("請輸入數字設置活動介紹的排序")).SetFieldMust()

	formList.SetTable("activity_introduce").SetTitle("活動介紹").SetDescription("活動介紹管理")
	// 設置活動介紹新增函式
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "introduce_title", "introduce_content", "introduce_order") {
			return errors.New("活動ID、介紹標題、內容、排序等欄位都不能為空")
		}

		order, _ := strconv.Atoi(values.Get("introduce_order"))
		if models.DefaultIntroduceModel().SetConn(s.conn).IsOrderExist(order, values.Get("activity_id"), "") {
			return errors.New("活動已在該排序中建立活動介紹，請設置其他排序")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 新增活動資料
			_, err := models.DefaultIntroduceModel().SetTx(tx).SetConn(s.conn).AddActivityIntroduce(
				values.Get("activity_id"), values.Get("introduce_title"), values.Get("introduct_content"), order)
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
		if models.DefaultIntroduceModel().SetConn(s.conn).IsOrderExist(order, values.Get("activity_id"), values.Get("id")) {
			return errors.New("活動已在該排序中建立活動介紹，請設置其他排序")
		}

		introduceModel := models.GetIntroduceModelAndID("activity_introduce", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 更新用戶資料
			_, err := introduceModel.SetTx(tx).UpdateActivityIntroduce(values.Get("activity_id"),
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
