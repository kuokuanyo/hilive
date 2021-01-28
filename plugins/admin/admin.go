package admin

import (
	"hilive/controller"
	"hilive/guard"
	"hilive/modules/config"
	"hilive/modules/service"
	"hilive/modules/table"
	"hilive/plugins"
)

// Admin 也屬於Plugin(interface)所有方法，設置handler與guardian(api方法)
type Admin struct {
	*plugins.Base
	config  *config.Config
	list    table.List // 放置所有頁面及表單資訊
	guard   *guard.Guard
	handler *controller.Handler
}

// NewAdmin 設置一個Admin(struct)
func NewAdmin(tableCfg ...table.List) *Admin {
	return &Admin{
		list:    make(table.List).CombineAll(tableCfg),
		Base:    &plugins.Base{PlugName: "admin"},
		handler: controller.NewHandler(),
	}
}

// -----plugin的所有方法-----start

// InitPlugin 設置admin(struct)以及設置api路徑、功能
func (admin *Admin) InitPlugin(services service.List) {
	// InitBase 設置Base.Conn、Base.Services
	admin.Base.InitBase(services)

	// 從Services中取得config
	c := config.GetService(services.Get("config"))

	// 將參數設置至SystemTable(struct)
	st := table.NewSystemTable(admin.Conn, c)

	admin.list.Combine(table.List{
		"manager":                 st.GetManagerPanel,
		"roles":                   st.GetRolesPanel,
		"permission":              st.GetPermissionPanel,
		"activity":                st.GetActivityPanel,
		"activity_overview":       st.GetActivityOverviewPanel,
		"activity_introduce":      st.GetIntroducePanel,
		"activity_schedule":       st.GetSchedulePanel,
		"activity_guest":          st.GetGuestPanel,
		"activity_material":       st.GetMaterialPanel,
		"activity_applysign":      st.GetApplysignPanel,
		"message":                 st.GetMessagePanel,
		"topic":                   st.GetTopicPanel,
		"question":                st.GetQuestionPanel,
		"danmu":                   st.GetDanmuPanel,
		"super_danmu":             st.GetSuperdmPanel,
		"picture":                 st.GetPicturePanel,
		"holdscreen":              st.GetHoldScreenPanel,
		"contract":                st.GetContractPanel,
		"sign":                    st.GetSignPanel,
		"3Dsign":                  st.Get3DSignPanel,
		"countdown":               st.GetCountdownPanel,
		"gamelottery":             st.GetGameLottery,
		"gamelottery_prize":       st.GetLotteryPrizePanel,
		"gamelottery_other_prize": st.GetLotteryOtherPrize,
		"redpack":                 st.GetRedpackPanel,
		"redpack_prize":           st.GetRedpackPrizePanel,
		"ropepack":                st.GetRopepackPanel,
		"ropepack_prize":          st.GetRopepackPrizePanel,
		"game_staff":              st.GetGameStaffPanel,
		"prize_staff":             st.GetPrizeStaffPanel,
	})

	// 設置Config
	admin.config = c

	// 設置admin.guardian
	admin.guard = guard.NewGuard(admin.Services, admin.Conn, admin.list, c)

	// 設置admin.handler
	admin.handler.NewHandler(c, admin.Services, admin.Conn, admin.list)

	// ***************放置api的地方*****************
	admin.initRouter()

	// 設置services，services是service套件中的全域變數
	table.SetServices(services)
}

// -----plugin的所有方法-----end
