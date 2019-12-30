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
	beego.Router("/goodsDetail", &controllers.GoodsController{}, "get:ShowGoodsDetail")
	// 获取商品列表页面
	beego.Router("/goodsList", &controllers.GoodsController{}, "get:ShowList")
	//商品搜索
	beego.Router("/goodsSearch", &controllers.GoodsController{}, "post:HandleSearch")
	// 购物车
	beego.Router("/user/addCart", &controllers.CartController{}, "post:HandleAddCart")
	// 购物车页面
	beego.Router("/user/cart", &controllers.CartController{}, "get:ShowCart")
	// 商品购物车
	beego.Router("/user/UpdateCart", &controllers.CartController{}, "post:HandleUpdateCart")
	// 购物车删除
	beego.Router("/user/deleteCart", &controllers.CartController{}, "post:DeleteCart")
	// 订单页面
	beego.Router("/user/showOrder", &controllers.OrderController{}, "post:ShowOrder")
	// 提交订单
	beego.Router("/user/addOrder", &controllers.OrderController{}, "post:AddOrder")
	//处理支付
	beego.Router("/user/pay", &controllers.OrderController{}, "get:HandlePay")
	//支付成功
	beego.Router("/user/payok", &controllers.OrderController{}, "get:PayOk")
}

var filterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
