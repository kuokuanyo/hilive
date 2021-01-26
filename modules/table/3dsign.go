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

// Get3DSignPanel 取得3D簽到牆頁面、表單資訊
func (s *SystemTable) Get3DSignPanel(ctx *context.Context) (signTable Table) {
	// 建立BaseTable
	signTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊欄位
	info := signTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("頭像形狀", "avatar_shape", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "圓形"
			}
			return "方形"
		})
	info.AddField("是否顯示簽到人數", "display", db.Tinyint).
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			if value.Value == "1" {
				return "開啟"
			}
			return "關閉"
		})
	info.AddField("簽到牆背景", "background", db.Varchar)

	info.SetTable("activity_set_3d").SetTitle("3D簽到牆").SetDescription("3D簽到牆管理").
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity_set_3d").WhereIn("id", ids).Delete()
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
	formList := signTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Text).
		SetFieldHelpMsg(template.HTML("活動辨別ID")).SetFieldMust().FieldNotAllowEdit()
	formList.AddField("頭像形狀", "avatar_shape", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "圓形"},
			{Value: "0", Text: "方形"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_3d").Select("avatar_shape").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["avatar_shape"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("是否顯示簽到人數", "display", db.Tinyint, form.Radio).
		SetFieldOptions(types.FieldOptions{
			{Value: "1", Text: "開啟"},
			{Value: "0", Text: "關閉"},
		}).SetFieldMust().
		SetDisplayFunc(func(value types.FieldModel) interface{} {
			var stauts []string
			if value.ID == "" {
				return []string{value.Value}
			}

			statusModel, _ := s.table("activity_set_3d").Select("display").FindByID(value.ID)
			stauts = append(stauts, strconv.FormatInt(statusModel["display"].(int64), 10))
			return stauts
		}).SetFieldDefault("1")
	formList.AddField("簽到牆背景", "background", db.Varchar, form.Text)

	formList.SetTable("activity_set_3d").SetTitle("3D簽到牆").SetDescription("3D簽到牆管理")

	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_id", "display", "avatar_shape") {
			return errors.New("活動ID、頭像形狀、顯示人數等欄位都不能為空")
		}

		if models.DefaultSign3DModel().SetConn(s.conn).IsActivityExist(values.Get("activity_id"), "") {
			return errors.New("該活動已設置過3D簽到牆的基礎設定")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := models.DefaultSign3DModel().SetTx(tx).SetConn(s.conn).Add3DSign(
				values.Get("activity_id"), values.Get("avatar_shape"),
				values.Get("display"), values.Get("background"))
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
		if values.IsEmpty("activity_id", "display", "avatar_shape") {
			return errors.New("活動ID、頭像形狀、顯示人數等欄位都不能為空")
		}

		model := models.GetSign3DModelAndID("activity_set_3d", values.Get("id")).SetConn(s.conn)
		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			_, err := model.SetTx(tx).Update3DSign(values.Get("activity_id"), values.Get("avatar_shape"),
				values.Get("display"), values.Get("background"))
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
