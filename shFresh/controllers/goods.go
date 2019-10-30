// hello,jiajia
package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"shFresh/models"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"math"
)

type GoodsController struct {
	beego.Controller
}

// 用户名获取
func GetUser(this *beego.Controller) string {
	userName := this.GetSession("userName")
	if userName == nil {
		//this.Redirect("/login", 302)
		this.Data["userName"] = ""
		return ""
	} else {
		this.Data["userName"] = userName
	}
	return userName.(string)
}

// 展示首页
func (this *GoodsController) ShowIndex() {
	_ = GetUser(&this.Controller)
	// 获取类型数据
	o := orm.NewOrm()
	var goodsTypes [] models.GoodsType
	o.QueryTable("GoodsType").All(&goodsTypes)
	this.Data["goodsTypes"] = goodsTypes
	// 获取轮播数据
	var IndexGoodsBanner []models.IndexGoodsBanner
	o.QueryTable("IndexGoodsBanner").OrderBy("Index").All(&IndexGoodsBanner)
	this.Data["IndexGoodsBanner"] = IndexGoodsBanner
	//获取促销商品数据
	var promotionGoods []models.IndexPromotionBanner
	o.QueryTable("IndexPromotionBanner").OrderBy("Index").All(&promotionGoods)
	this.Data["promotionGoods"] = promotionGoods

	// 首页展示商品数据
	goods := make([]map[string]interface{}, len(goodsTypes))
	// 向切片interface插入类型数据
	for index, value := range goodsTypes {
		// 商品获取对应的类型的首页展示商品
		temp := make(map[string]interface{})
		//goods[index] = value
		temp["type"] = value
		goods[index] = temp
	}

	//商品数据

	for _, value := range goods {
		var textGoods []models.IndexTypeGoodsBanner
		// 获取文字商品数据
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").OrderBy("Index").Filter("GoodsType", value["type"]).Filter("DisplayType", 0).All(&textGoods)
		var imgGoods []models.IndexTypeGoodsBanner
		// 获取图片商品数据
		o.QueryTable("IndexTypeGoodsBanner").RelatedSel("GoodsType", "GoodsSKU").OrderBy("Index").Filter("GoodsType", value["type"]).Filter("DisplayType", 1).All(&imgGoods)
		value["textGoods"] = textGoods
		value["imgGoods"] = imgGoods
	}
	this.Data["goods"] = goods

	// 购物车
	cartCount := GetCartCount(&this.Controller)
	this.Data["cartCount"] = cartCount


	this.TplName = "index.html"
}

func ShowLayout(this *beego.Controller) {
	//查询类型
	o := orm.NewOrm()
	var types []models.GoodsType
	o.QueryTable("GoodsType").All(&types)
	this.Data["types"] = types
	//获取用户信息
	GetUser(this)
	//制定Layout
	this.Layout = "goodsLayout.html"
}

// 商品详情
func (this *GoodsController) ShowGoodsDetail() {
	// 获取数据
	id, err := this.GetInt("id")
	// 校验数据
	if err != nil {
		beego.Info("浏览器请求错误", err)
		this.Redirect("/", 302)
		return
	}
	// 处理数据
	o := orm.NewOrm()
	var goodsSku models.GoodsSKU
	goodsSku.Id = id
	//o.Read(&goodsSku)
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType", "Goods").Filter("Id", id).One(&goodsSku)

	// 获取同类型时间靠前的商品数据
	var goodsNew []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType", goodsSku.GoodsType).OrderBy("Time").Limit(2, 0).All(&goodsNew)
	this.Data["goodsNew"] = goodsNew
	// 返回数据
	this.Data["goodsSku"] = goodsSku

	// 添加历史访问记录
	//1.判断用户是否登陆
	userName := this.GetSession("userName")
	//2.添加历史记录
	if userName != nil {
		// 查询用户信息
		o := orm.NewOrm()
		var user models.User
		user.Name = userName.(string)
		o.Read(&user, "Name")
		//添加记录，用redis储存
		conn, err := redis.Dial("tcp", "192.168.88.130:6379")
		defer conn.Close()
		if err != nil {
			beego.Info("redis连接错误", err)
		}
		// 把以前相同的商品删除
		conn.Do("lrem", "history_"+strconv.Itoa(user.Id), 0, id)
		// 添加新的
		conn.Do("lpush", "history_"+strconv.Itoa(user.Id), id)
	}

	ShowLayout(&this.Controller)
	cartCount := GetCartCount(&this.Controller)
	this.Data["cartCount"] = cartCount
	this.TplName = "detail.html"
}

//分页
func PageTool(pageCount, PageIndex int) []int {
	var pages []int
	if pageCount <= 5 {
		pages = make([]int, pageCount)
		for i, _ := range pages {
			pages[i] = i + 1
			beego.Info(pages)
		}
	} else if PageIndex <= 3 {
		pages = []int{1, 2, 3, 4, 5}

	} else if PageIndex > pageCount-3 {
		pages = []int{pageCount - 4, pageCount - 3, pageCount - 2, pageCount - 1, pageCount}

	} else {
		pages = []int{PageIndex - 2, PageIndex - 1, PageIndex, PageIndex + 1, PageIndex + 2}

	}
	return pages
}

// 展示商品列表页
func (this *GoodsController) ShowList() {
	//获取数据
	id, err := this.GetInt("typeId")
	//校验数据
	if err != nil {
		beego.Info("请求路径错误")
		this.Redirect("/", 302)
		return
	}
	//处理数据
	ShowLayout(&this.Controller)
	//获取新品
	o := orm.NewOrm()
	var goodsNew []models.GoodsSKU
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__id", id).OrderBy("Time").Limit(2, 0).All(&goodsNew)

	//获取商品
	var goods []models.GoodsSKU
	//o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__id", id).All(&goods)
	//返回视图
	this.Data["goodsNew"] = goodsNew

	//分页实现
	//获取总页码pageCount
	count, err := o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__id", id).Count()
	pageSize := 3
	pageCount := math.Ceil(float64(count) / float64(pageSize))
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}

	pages := PageTool(int(pageCount), int(pageIndex))
	this.Data["pages"] = pages
	this.Data["typeId"] = id
	this.Data["pageIndex"] = pageIndex

	start := (pageIndex - 1) * pageSize
	// 获取上一页代码
	PrePage := pageIndex - 1
	if PrePage <= 1 {
		PrePage = 1
	}
	NextPage := pageIndex + 1
	if NextPage >= int(pageCount) {
		NextPage = int(pageCount)
	}
	this.Data["PrePage"] = PrePage
	this.Data["NextPage"] = NextPage

	//获取排序
	sort := this.GetString("sort")
	if sort == "" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__id", id).Limit(pageSize, start).All(&goods)
		this.Data["sort"] = sort
		this.Data["goods"] = goods
	} else if sort == "price" {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__id", id).OrderBy("Price").Limit(pageSize, start).All(&goods)
		this.Data["sort"] = sort
		this.Data["goods"] = goods
	} else {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__id", id).OrderBy("Sales").Limit(pageSize, start).All(&goods)
		this.Data["sort"] = sort
		this.Data["goods"] = goods
	}

	this.TplName = "list.html"
}

//商品搜索
func (this *GoodsController) HandleSearch() {
	//获取数据
	goodsName := this.GetString("goodsName")
	o := orm.NewOrm()
	var goods []models.GoodsSKU
	//校验数据
	if goodsName == "" {
		o.QueryTable("GoodsSKU").All(&goods)
		this.Data["goods"] = goods
		ShowLayout(&this.Controller)
		this.TplName = "search.html"
		return
	}
	//处理数据
	o.QueryTable("GoodsSKU").Filter("Name__icontains", goodsName).All(&goods)
	//返回数据
	this.Data["goods"] = goods
	ShowLayout(&this.Controller)
	this.TplName = "search.html"
}

