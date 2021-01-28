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
)

// GetTopicPanel 取得主題牆的頁面、表單資訊
func (s *SystemTable) GetTopicPanel(ctx *context.Context) (topicTable Table) {
	// 建立BaseTable
	topicTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊
	info := topicTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("背景圖片", "background", db.Varchar)

	info.SetTable("activity_set_topic").SetTitle("主題牆").SetDescription("主題牆管理").
		// 刪除函式
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_topic").WhereIn("id", ids).Delete()
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
	formList := topicTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("背景圖片", "background", db.Varchar, form.Text)

	formList.SetTable("activity_set_topic").SetTitle("主題牆").SetDescription("主題牆管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("活動ID不能為空")
		}

		if models.DefaultTopicModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), "") {
			return errors.New("該活動已設置過主題牆的基礎設定")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultTopicModel().SetTx(tx).SetConn(s.conn).AddTopic(
				values.Get("activity_id"), values.Get("background"))
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 主題牆基礎設置更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id") {
			return errors.New("活動ID不能為空")
		}

		if models.DefaultTopicModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), values.Get("id")) {
			return errors.New("該活動已設置過主題牆的基礎設定")
		}

		topicModel := models.GetTopicModelAndID("activity_set_topic", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := topicModel.SetTx(tx).UpdateTopic(values.Get("activity_id"),
				values.Get("background"))
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
