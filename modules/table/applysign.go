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

// GetApplysignPanel 取得報名簽到的頁面、表單資訊
func (s *SystemTable) GetApplysignPanel(ctx *context.Context) (applysignTable Table) {
	// 建立BaseTable
	applysignTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊
	info := applysignTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("使用者ID", "user_id", db.Varchar).FieldFilterable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("使用者名稱", "user_name", db.Varchar).FieldFilterable()
	info.AddField("使用者頭像", "user_avater", db.Varchar)
	info.AddField("簽到狀態", "status", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "簽到完成"
			}
			return "未簽到"
		})
	info.AddField("簽到時間", "sign_time", db.Datetime)

	info.SetTable("activity_applysign").SetTitle("報名簽到").SetDescription("簽到管理").
		// 刪除函式
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_applysign").WhereIn("id", ids).Delete()
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
	formList := applysignTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("使用者ID", "user_id", db.Varchar, form.Text).SetFieldHelpMsg(template.HTML("使用者辨別ID")).SetFieldMust()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust()
	formList.AddField("使用者名稱", "user_name", db.Varchar, form.Text).SetFieldMust()
	formList.AddField("使用者頭像", "user_avater", db.Varchar, form.Text)
	formList.AddField("簽到狀態", "status", db.Tinyint, form.SelectSingle).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "簽到完成"},
			{Value: "0", Text: "未簽到"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var status []string
			if value.ID == "" {
				return status
			}

			statusModel, _ := s.table("activity_applysign").Select("status").FindByID(value.ID)
			status = append(status, strconv.FormatInt(statusModel["status"].(int64), 10))
			return status
		})

	formList.SetTable("activity_applysign").SetTitle("報名簽到").SetDescription("簽到管理")
	// 設置報名簽到新增函式
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("user_id", "activity_id", "user_name", "status") {
			return errors.New("活動ID、使用者資訊、簽到狀態等欄位都不能為空")
		}

		if models.DefaultApplysignModel().SetConn(s.conn).IsSignExist(values.Get("activity_id"), "", values.Get("id")) {
			return errors.New("該用戶已報名簽到此活動")
		}

		status, _ := strconv.Atoi(values.Get("status"))
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 新增活動資料
			_, err := models.DefaultApplysignModel().SetTx(tx).SetConn(s.conn).AddApplysign(
				values.Get("user_id"), values.Get("activity_id"), values.Get("user_name"),
				values.Get("user_avater"), status)
			if err != nil {
				if err.Error() != "沒有影響任何資料" {
					return err, nil
				}
			}
			return nil, nil
		})
		return txErr
	})

	// 報名簽到更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("user_id", "activity_id", "user_name", "status") {
			return errors.New("活動ID、使用者資訊、簽到狀態等欄位都不能為空")
		}

		if models.DefaultApplysignModel().SetConn(s.conn).IsSignExist(values.Get("activity_id"), values.Get("user_id"), values.Get("id")) {
			return errors.New("該用戶已報名簽到此活動")
		}

		status, _ := strconv.Atoi(values.Get("status"))
		applyModel := models.GetApplysignModelAndID("activity_applysign", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 更新報名簽到資料
			_, err := applyModel.SetTx(tx).UpdateActivityApplysign(values.Get("user_id"),
				values.Get("activity_id"), values.Get("user_name"),
				values.Get("user_avater"), status)
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
