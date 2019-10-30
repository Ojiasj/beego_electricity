package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"shanghaiyiqi/controllers"
)

func init() {
	beego.InsertFilter("/Article/*", beego.BeforeExec, Filfter)
	beego.Router("/", &controllers.MainController{}, "get:ShowGet;post:Post")

	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandlePost")

	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	//文章列表页访问
	beego.Router("/Article/showArticleList", &controllers.ArticleController{}, "get:ShowArticleList")
	//添加文章
	beego.Router("/Article/addArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")
	//显示文章详情
	beego.Router("/Article/showArticleDetail", &controllers.ArticleController{}, "get:ShowArticleDetail")
	//编辑文章
	beego.Router("/Article/updateArticle", &controllers.ArticleController{}, "get:ShowUpdateArticle;post:HandleUpdateArticle")
	//删除文章
	beego.Router("/Article/deleteArticle", &controllers.ArticleController{}, "get:DeleteArticle")
	//添加分类
	beego.Router("/Article/AddArticleType", &controllers.ArticleController{}, "get:ShowAddType;post:HandleAddType")
	//退出登录
	beego.Router("/Article/logout", &controllers.UserController{}, "get:Logout")
	//删除类型
	beego.Router("/Article/deleteType", &controllers.ArticleController{}, "get:DeleteType")

}

var Filfter = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}

}
