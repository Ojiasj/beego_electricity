package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"shFreshManagement/controllers"
)

func init() {
	beego.InsertFilter("/Commodity/*", beego.BeforeExec, Filfter)
	// 注册
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandlePost")
	// 登录
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	// 退出登录
	beego.Router("/Commodity/logout", &controllers.UserController{}, "get:Logout")
	// 商品列表页访问
	beego.Router("/Commodity/showCommodityList", &controllers.CommodityController{}, "get:ShowCommodityList")
	//查看商品详细
	beego.Router("/Commodity/showCommodityDetail", &controllers.CommodityController{}, "get:ShowCommodityDetail")
	// 添加商品
	beego.Router("/Commodity/AddCommodity", &controllers.CommodityController{}, "get:ShowAddCommodity;post:HandleAddCommodity")
	// 编辑商品
	beego.Router("/Commodity/EditCommodity", &controllers.CommodityController{}, "get:ShowEditCommodity;post:HandleEditCommodity")
	// 删除商品
	beego.Router("/Commodity/DelCommodity", &controllers.CommodityController{}, "get:DelCommodity")
	// 添加分类
	beego.Router("/Commodity/AddCommodityType", &controllers.CommodityController{}, "get:ShowAddType;post:HandleAddType")
	// 分类详细
	beego.Router("/Commodity/TypeDeail", &controllers.CommodityController{}, "get:ShowTypeDetail")
	// 删除分类
	beego.Router("/Commodity/DelType", &controllers.CommodityController{}, "get:DelType")
	// 添加商品SPU
	beego.Router("/Commodity/AddCommoditySPU", &controllers.CommodityController{}, "get:ShowAddCommoditySPU;post:HandleAddCommoditySPU")
	// 商品SPU详细
	beego.Router("/Commodity/CommoditySPUDetail",&controllers.CommodityController{},"get:ShowCommoditySPUDetail")
	// 删除商品SPU
	beego.Router("/Commodity/DelCommoditySPU",&controllers.CommodityController{},"get:DelCommoditySPU")

	// 测试
	beego.Router("/test", &controllers.CommodityController{}, "get:TestFunc")

}

var Filfter = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}

}
