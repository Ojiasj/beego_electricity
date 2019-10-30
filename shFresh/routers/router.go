package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"shFresh/controllers"
)

func init() {
	beego.InsertFilter("/user/*", beego.BeforeExec, filterFunc)
	///////////////////////////////////////////////用户///////////////////////////////////////////////
	beego.Router("/register", &controllers.UserController{}, "get:ShowReg;post:HandleReg")
	//激活用户
	beego.Router("/active", &controllers.UserController{}, "get:ActiveUser")
	//用户登录
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	// 退出登录
	beego.Router("/user/logout", &controllers.UserController{}, "get:Logout")
	// 用户中心信息页
	beego.Router("/user/userCenterInfo", &controllers.UserController{}, "get:ShowUserCenterInfo")
	// 用户中心订单页
	beego.Router("/user/userCenterOrder", &controllers.UserController{}, "get:ShowUserCenterOrder")
	// 用户地址显示页
	beego.Router("/user/userCenterSite", &controllers.UserController{}, "get:ShowUserCenterSite;post:HandUserCenterSite")
	//////////////////////////////////////////////商品页///////////////////////////////////////////////
	// 跳转首页
	beego.Router("/", &controllers.GoodsController{}, "get:ShowIndex")
	// 商品详情显示
	beego.Router("/goodsDetail",&controllers.GoodsController{},"get:ShowGoodsDetail")
	// 获取商品列表页面
	beego.Router("/goodsList",&controllers.GoodsController{},"get:ShowList")
	//商品搜索
	beego.Router("/goodsSearch",&controllers.GoodsController{},"post:HandleSearch")

}

var filterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
