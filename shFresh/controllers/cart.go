package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"github.com/astaxie/beego/orm"
	"shFresh/models"
	"strconv"
)

type CartController struct {
	beego.Controller
}

func (this *CartController) HandleAddCart() {
	// 获取数据
	skuid, err1 := this.GetInt("skuid")
	count, err2 := this.GetInt("count")
	resp := make(map[string]interface{})
	defer this.ServeJSON()

	// 校验数据
	if err1 != nil || err2 != nil {
		resp["code"] = 1 //请求数据不正确
		resp["msg"] = "请求数据不正确"
		this.Data["json"] = resp
		return
	}
	userName := this.GetSession("userName")
	if userName == nil {
		resp["code"] = 1 //没有登陆
		resp["msg"] = "没有登陆"
		this.Data["json"] = resp
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")
	// 处理数据
	// 购物车数据，存在redis中
	conn, err := redis.Dial("tcp", "192.168.88.130:6379")
	if err != nil {
		beego.Info("redis数据库连接错误")
		return
	}
	//先获取原来的数量，然后给数量加起来
	preCount, err := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), skuid))
	beego.Info("preCount", preCount)
	finaCount := count + preCount
	beego.Info("finaCount", finaCount)

	conn.Do("hset", "cart_"+strconv.Itoa(user.Id), skuid, finaCount)

	rep, err := conn.Do("hlen", "cart_"+strconv.Itoa(user.Id))

	// 回复助手函数
	cartCount, _ := redis.Int(rep, err)
	beego.Info("cartCount", cartCount)

	resp["code"] = 5 // 成功
	resp["msg"] = "ok"
	resp["cartCount"] = cartCount

	this.Data["json"] = resp
	// 返回json数据
}

// 获取购物车数量
func GetCartCount(this *beego.Controller) int {
	// 从redis中获取数据
	userName := this.GetSession("userName")
	if userName == nil {
		return 0
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")
	conn, err := redis.Dial("tcp", "192.168.88.130:6379")
	if err != nil {
		return 0
	}
	defer conn.Close()

	rep, err := conn.Do("hlen", "cart_"+strconv.Itoa(user.Id))
	// 回复助手函数
	cartCount, _ := redis.Int(rep, err)
	return cartCount
}

func (this *CartController) ShowCart() {
	userName := GetUser(&this.Controller)
	cartCount := GetCartCount(&this.Controller)
	conn, err := redis.Dial("tcp", "192.168.88.130:6379")
	if err != nil {
		beego.Info("redis连接错误", err)
		return
	}
	defer conn.Close()
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	o.Read(&user, "Name")

	goodsMap, _ := redis.IntMap(conn.Do("hgetall", "cart_"+strconv.Itoa(user.Id))) // 返回切片
	goods := make([]map[string]interface{}, len(goodsMap))
	i := 0
	totalPrice := 0
	totalCount := 0
	for index, value := range goodsMap {
		skuid, _ := strconv.Atoi(index)
		var goodsSku models.GoodsSKU
		goodsSku.Id = skuid
		o.Read(&goodsSku)

		temp := make(map[string]interface{})
		temp["goods"] = goodsSku
		temp["count"] = value
		temp["addPrice"] = goodsSku.Price * value

		totalPrice += goodsSku.Price * value
		totalCount += value

		goods[i] = temp
		i += 1
	}

	this.Data["userName"] = userName
	this.Data["cartCount"] = cartCount
	this.Data["goodsSku"] = goods
	this.Data["totalCount"] = totalCount
	this.Data["totalPrice"] = totalPrice

	this.TplName = "cart.html"

}

func (this *CartController) HandleUpdateCart() {
	// 获取数据
	skuid, err1 := this.GetInt("skuid")
	count, err2 := this.GetInt("count")
	beego.Info(skuid,count)
	resp := make(map[string]interface{})
	defer this.ServeJSON()
	// 校验数据
	if err1 != nil || err2 != nil {
		resp["code"] = 1
		resp["errmsg"] = "请求失败"
		this.Data["json"] = resp
		return
	}
	//处理数据
	conn, err := redis.Dial("tcp", "192.168.88.130:6379")
	if err != nil {
		resp["code"] = 2
		resp["errmsg"] = "redis连接失败"
		this.Data["json"] = resp
		return
	}
	userName := this.GetSession("userName")
	if userName == nil {
		resp["code"] = 3
		resp["errmsg"] = "当前用户未登陆"
		this.Data["json"] = resp
		return
	}
	o := orm.NewOrm()
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")
	defer conn.Close()

	conn.Do("hset", "cart_"+strconv.Itoa(user.Id), skuid, count)
	resp["code"] = 5
	resp["errmsg"] = "OK"
	this.Data["json"] = resp
}
