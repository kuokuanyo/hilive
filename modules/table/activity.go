package table

import (
	"database/sql"
	"errors"
	"hilive/context"
	"hilive/models"
	"hilive/modules/config"
	"hilive/modules/db"
	form2 "hilive/modules/form"
	"hilive/modules/utils"
	"hilive/template/form"
	"hilive/template/types"
	"html/template"
	"strconv"
	"time"
)

// GetActivityPanel 取得建立活動的頁面、表單資訊
func (s *SystemTable) GetActivityPanel(ctx *context.Context) (activityTable Table) {
	// 建立BaseTable
	activityTable = DefaultBaseTable(DefaultConfigTableByDriver(config.GetDatabaseDriver()))

	// 增加頁面資訊
	info := activityTable.GetInfo()
	info.AddField("ID", "id", "INT").FieldSortable()
	info.AddField("活動專屬ID", "activity_id", db.Varchar).FieldFilterable()
	info.AddField("活動主題", "activity_name", db.Varchar).FieldFilterable()
	info.AddField("活動類型", "activity_type", db.Varchar).FieldFilterable().
		FieldJoin(types.Join{
			JoinTable: "activity_type",
			JoinField: "type_id",
			BaseTable: "activity",
			Field:     "activity_type",
		})
	info.AddField("預計參加人數", "participant_interval", db.Varchar).
		FieldJoin(types.Join{
			JoinTable: "activity_participant",
			JoinField: "participant_id",
			BaseTable: "activity",
			Field:     "expected_participants",
		})
	info.AddField("已參加人數", "participants", db.Int)
	info.AddField("活動縣市", "city", db.Varchar).FieldFilterable().
		FieldJoin(types.Join{
			JoinTable: "activity_city",
			JoinField: "city_id",
			BaseTable: "activity",
			Field:     "city",
		})
	info.AddField("活動地區", "town", db.Varchar).FieldFilterable().
		FieldJoin(types.Join{
			JoinTable: "activity_town",
			JoinField: "town_id",
			BaseTable: "activity",
			Field:     "town",
		})
	info.AddField("活動開始時間", "start_time", db.Datetime)
	info.AddField("活動結束時間", "end_time", db.Datetime)

	info.SetTable("activity").SetTitle("活動").SetDescription("活動管理").
		// 刪除函式
		SetDeleteFunc(func(idArr []string) error {
			var ids = interfaces(idArr)

			_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
				err := s.connection().SetTx(tx).
					Table("activity").WhereIn("id", ids).Delete()
				if err != nil {
					if err.Error() != "沒有影響任何資料" {
						return errors.New("刪除活動資料發生錯誤"), nil
					}
				}
				return nil, nil
			})
			return txErr
		})

	// 取得FormPanel
	formList := activityTable.GetFormPanel()
	formList.AddField("ID", "id", "INT", form.Default).FieldNotAllowAdd().FieldNotAllowEdit()
	formList.AddField("活動專屬ID", "activity_id", db.Varchar, form.Default).FieldNotAllowAdd().FieldNotAllowEdit().SetFieldHelpMsg(template.HTML("活動辨別ID"))
	formList.AddField("活動主題", "activity_name", db.Varchar, form.Text).SetFieldMust().SetFieldHelpMsg(template.HTML("活動主題名稱不能超過20個字"))
	formList.AddField("活動類型", "activity_type", db.Varchar, form.SelectSingle).SetFieldMust().
		SetFieldOptionFromTable("activity_type", "activity_type", "type_id").
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			var activityType []string
			if model.ID == "" {
				return activityType
			}

			typeModel, _ := s.table("activity").Select("activity_type").FindByID(model.ID)
			activityType = append(activityType, strconv.FormatInt(typeModel["activity_type"].(int64), 10))
			return activityType
		})
	formList.AddField("預計參加人數", "expected_participants", db.Varchar, form.SelectSingle).SetFieldMust().
		SetFieldOptionFromTable("activity_participant", "participant_interval", "participant_id").
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			var activityParticipate []string
			if model.ID == "" {
				return activityParticipate
			}

			partModel, _ := s.table("activity").Select("expected_participants").FindByID(model.ID)
			activityParticipate = append(activityParticipate, strconv.FormatInt(partModel["expected_participants"].(int64), 10))
			return activityParticipate
		})
	formList.AddField("活動縣市", "city", db.Varchar, form.SelectSingle).SetFieldMust().
		SetFieldOptionFromTable("activity_city", "city", "city_id").
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			var activityCity []string
			if model.ID == "" {
				return activityCity
			}

			cityModel, _ := s.table("activity").Select("city").FindByID(model.ID)
			activityCity = append(activityCity, strconv.FormatInt(cityModel["city"].(int64), 10))
			return activityCity
		})
	formList.AddField("活動地區", "town", db.Varchar, form.Text).SetFieldMust().
		SetDisplayFunc(func(model types.FieldModel) interface{} {
			var (
				town interface{}
			)
			if model.ID == "" {
				return town
			}

			// 取得activity的town(代號)
			townModel, _ := s.table("activity").Select("town").FindByID(model.ID)

			townModel2, _ := s.table("activity_town").Select("town").Where("town_id", "=", townModel["town"]).First()
			return townModel2["town"]
		}).SetFieldHelpMsg(template.HTML("ex:信義區"))
	formList.AddField("活動開始時間", "start_time", db.Datetime, form.Datetime).SetFieldMust()
	formList.AddField("活動結束時間", "end_time", db.Datetime, form.Datetime).SetFieldMust()

	formList.SetTable("activity").SetTitle("活動").SetDescription("活動管理")

	// 設置活動新增函式
	formList.SetInsertFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_name", "activity_type", "expected_participants",
			"city", "town", "start_time", "end_time") {
			return errors.New("活動主題、類型、地點、時間等欄位都不能為空")
		}

		// 時間判斷
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("start_time"), time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("end_time"), time.Local)
		if !end.After(start) && start.Before(end) {
			return errors.New("時間設置發生錯誤，請重新設置(結束時間在開始時間之後)")
		}

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 新增活動資料
			_, err := models.DefaultActivityModel().SetTx(tx).SetConn(s.conn).AddActivity(
				utils.UUID(8), values.Get("activity_name"), values.Get("activity_type"),
				values.Get("expected_participants"), values.Get("city"), values.Get("town"),
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

	// 設置活動更新函式
	formList.SetUpdateFunc(func(values form2.Values) error {
		if values.IsEmpty("activity_name", "activity_type", "expected_participants",
			"city", "town", "start_time", "end_time") {
			return errors.New("活動主題、類型、地點、時間等欄位都不能為空")
		}

		// 時間判斷
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("start_time"), time.Local)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", values.Get("end_time"), time.Local)
		boolTime := end.After(start) && start.Before(end)
		if boolTime == false {
			return errors.New("時間設置發生錯誤，請重新設置(結束時間在開始時間之後)")
		}

		activity := models.GetActivityModelAndID("activity", values.Get("id")).SetConn(s.conn)

		_, txErr := s.connection().WithTransaction(func(tx *sql.Tx) (e error, i map[string]interface{}) {
			// 更新用戶資料
			_, err := activity.SetTx(tx).Update(
				values.Get("activity_id"), values.Get("activity_name"), values.Get("activity_type"),
				values.Get("expected_participants"), values.Get("city"), values.Get("town"),
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
	return
}
