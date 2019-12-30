package controllers

import (
	"github.com/astaxie/beego"
	"strconv"
	"github.com/astaxie/beego/orm"
	"shFresh/models"
	"github.com/gomodule/redigo/redis"
	"time"
	"strings"
	"github.com/smartwalle/alipay"
	"fmt"
)

type OrderController struct {
	beego.Controller
}

func (this *OrderController) ShowOrder() {
	userName := GetUser(&this.Controller)
	// 获取数据
	skuids := this.GetStrings("skuid")
	//beego.Info(skuids)
	if len(skuids) == 0 {
		beego.Info("购物车请求数据错误")
		this.Redirect("/user/cart", 302)
	}
	// 处理数据
	o := orm.NewOrm()
	conn, _ := redis.Dial("tcp", "192.168.88.130:6379")
	// 获取用户数据
	var user models.User
	user.Name = userName
	o.Read(&user, "Name")
	defer conn.Close()
	goodsBuffer := make([]map[string]interface{}, len(skuids))
	totalPrice := 0
	totalCount := 0
	for index, skuid := range skuids {
		temp := make(map[string]interface{})
		id, _ := strconv.Atoi(skuid)
		// 查询商品数据
		var goodsSku models.GoodsSKU
		goodsSku.Id = id
		o.Read(&goodsSku)
		temp["goods"] = goodsSku
		// 获取商品数量
		count, _ := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), id))
		temp["count"] = count

		// 计算小计
		amount := goodsSku.Price * count
		temp["amount"] = amount
		goodsBuffer[index] = temp

		// 计算总金额和总件数
		totalPrice += amount
		totalCount += count

	}
	this.Data["userName"] = userName
	this.Data["goodsBuffer"] = goodsBuffer
	// 获取地址数据
	var addrs []models.Address
	o.QueryTable("Address").RelatedSel("User").Filter("User__id", user.Id).All(&addrs)
	this.Data["addrs"] = addrs
	this.Data["totalCount"] = totalCount
	this.Data["totalPrice"] = totalPrice
	transferPrice := 10
	this.Data["transferPrice"] = transferPrice
	this.Data["realyPrice"] = totalPrice + transferPrice
	// 传递所有商品id
	this.Data["skuids"] = skuids

	this.TplName = "place_order.html"

}

// 添加订单
func (this *OrderController) AddOrder() {
	//获取数据
	addrid, _ := this.GetInt("addrid")
	payId, _ := this.GetInt("payId")
	skuid := this.GetString("skuids")
	ids := skuid[1:len(skuid)-1]

	skuids := strings.Split(ids, " ")

	beego.Error(skuids)
	//totalPrice,_ := this.GetInt("totalPrice")
	totalCount, _ := this.GetInt("totalCount")
	transferPrice, _ := this.GetInt("transferPrice")
	realyPrice, _ := this.GetInt("realyPrice")

	resp := make(map[string]interface{})
	defer this.ServeJSON()
	//校验数据
	if len(skuids) == 0 {
		resp["code"] = 1
		resp["errmsg"] = "数据库链接错误"
		this.Data["json"] = resp
		return
	}
	//处理数据
	//向订单表中插入数据
	o := orm.NewOrm()

	o.Begin() //标识事务的开始

	userName := this.GetSession("userName")
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")

	var order models.OrderInfo
	order.OrderId = time.Now().Format("2006010215030405") + strconv.Itoa(user.Id)
	order.User = &user
	order.Orderstatus = 1
	order.PayMethod = payId
	order.TotalCount = totalCount
	order.TotalPrice = realyPrice
	order.TransitPrice = transferPrice
	//查询地址
	var addr models.Address
	addr.Id = addrid
	o.Read(&addr)

	order.Address = &addr

	//执行插入操作
	o.Insert(&order)

	//想订单商品表中插入数据
	conn, _ := redis.Dial("tcp", "192.168.88.130:6379")

	for _, skuid := range skuids {
		id, _ := strconv.Atoi(skuid)

		var goods models.GoodsSKU
		goods.Id = id
		i := 3

		for i > 0 {
			o.Read(&goods)

			var orderGoods models.OrderGoods

			orderGoods.GoodsSKU = &goods
			orderGoods.OrderInfo = &order

			count, _ := redis.Int(conn.Do("hget", "cart_"+strconv.Itoa(user.Id), id))

			if count > goods.Stock {
				resp["code"] = 2
				resp["errmsg"] = "商品库存不足"
				this.Data["json"] = resp
				o.Rollback() //标识事务的回滚
				return
			}

			preCount := goods.Stock

			time.Sleep(time.Second * 5)
			beego.Info(preCount, user.Id)

			orderGoods.Count = count

			orderGoods.Price = count * goods.Price

			o.Insert(&orderGoods)

			goods.Stock -= count
			goods.Sales += count

			updateCount, _ := o.QueryTable("GoodsSKU").Filter("Id", goods.Id).Filter("Stock", preCount).Update(orm.Params{"Stock": goods.Stock, "Sales": goods.Sales})
			if updateCount == 0 {
				if i > 0 {
					i -= 1
					continue
				}
				resp["code"] = 3
				resp["errmsg"] = "商品库存改变,订单提交失败"
				this.Data["json"] = resp
				o.Rollback() //标识事务的回滚
				return
			} else {
				conn.Do("hdel", "cart_"+strconv.Itoa(user.Id), goods.Id)
				break
			}
		}

	}

	//返回数据
	o.Commit() //提交事务
	resp["code"] = 5
	resp["errmsg"] = "ok"
	this.Data["json"] = resp

}

// 处理支付
func (this *OrderController) HandlePay() {
	var privateKey = "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCAcnzLAjnQiWxjiOckHLhWyvSOKWeojZEakOCIcKtXKI70PSZhqPc+H973+TFnf8OJ6faULXg8KeTMhSeDvGqcq25nXKVcUEqILga5bb7qYD64lYV0rAdiM8FN2ksUtJLqILHgwuv4Da10x8ykprpJiqyFvSHlz3GpcPrbwp1Ng2fXWCPIZyLxPYUujzhaH9KZUR77gErcSJBf9WsL2aFUwBKz95Dn0sX0i7C4KpDaGiw4zWnOegQEbhtuk0sS4k04O3oWRF9b3pTSET0VorVBqrDY1o6PQGh1SSjkU1nKV4kqeHGgkEwFji22V/EL0KuWbMXxMCb6Lp8KcTKMCqpBAgMBAAECggEAYDNe+7zDCEikgTe7xfQsq/R9jSu9kyPUFY2EXSvlZ/Xce1iBMouWAVVlbVuZgObT0KsGHpyffI/z6Kmhmqy3PHe4zHk68BTCfZPbPW3Qm0FSOHTj6yocrZQNpv1jVjKYBfpOvwO+L54u11P9FSQ6SXIvHEs25COmAT245Hax0acN4IO8h1UlkJl8EEZSxCSlH5i2IHcC8I5ZjO0m1UeZmDmY5hIwqj6Jj123V1HaxCkIvFKFrfjV0ekgF2wh+pcwqCgz+ZxdrN4jAx6Dq/bZ+PZXWiuBH/yUshMzmtmLbqR5HC39D2kuug0Al2YK6otacWY9hpCZOlI8vOzqaMoZAQKBgQDmYXbvee3BxjEK2/w+TWO+oK7CJf+4z3wEaHYnKTdxXsDBKrpwG7oOjRD26Ut5Z4euIpFIthpFomhP+HPcMcXkWFnWatExVbzasMd1D2I0ZKE+xTVRIaX6IcH6hjsFdAzwepF8eNo4dCcy5/8Y/LhIHnrHJ0EMW0qy9O3eTBhtKQKBgQCOuyeWgn7nQuDZ/Ma9qWRYi0PZD+5pvxn56hGSsKKRmM9kZbHBHvchGmBiM8voifBeI8c/so2plv9L0tqkpdsyyELLPGC4cCzePGauievFmvb1xuzuglTUxESUijZr870m7BLNsomWS2Pg2fQb9lyaT7LUhI5d+C7tPF7i5YbfWQKBgQCe3a0PnOwYiNw+2c5LBED5LoL0fRGn1uR1gbegb6q88hNH4XgpYOVfrWV6bwxNztfgfBPlqTXauRDnvLPgry4AtfBVjNlvBOmzgN46Wf5llNMgnwuSQ/rZzzed07yVmR5nIo564DfTYD27vAHMsFE/4kzWUrpnk/iiCYaSmbTqaQKBgBisT6ab/lX54KauJKjqnFcWE+906pDLITFrfggukpu6n7dKQRUSRkQprBmtvIUBO1T6uFnRgt2bJZy04Wju5tp7ddhuIoYflhIAvPtVCjXJmQFZluGQFBwHpZdL4SZ+JueQjZuTTmd1ttlKtAEVcGKYfmBwXa/u1CXcdsagSNVZAoGBAIOPO6GTZ/l7fIikJXmTD1fuX0+dnI3+g4IpnefbS7w96YqrleKz+HcwgulfcJG7hiyE8HW+0ae0Zh/+4u6MTWdsTffwonhQq7s8AlTDegKQrOMBhkyyQwOWqzhmmGT2iC83e4leE4wBvSrsy6Y/AlgTsQs1FGntHl0BZHzfs8ML" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	appId := "2016101600703864"
	var client, err = alipay.New(appId, privateKey, false)

	err3 := client.LoadAppPublicCertFromFile("../key/appCertPublicKey_2016101600703864.crt")     // 加载应用公钥证书
	err4 := client.LoadAliPayRootCertFromFile("../key/alipayRootCert.crt")           			  // 加载支付宝根证书
	err5 := client.LoadAliPayPublicCertFromFile("../key/alipayCertPublicKey_RSA2.crt")			 // 加载支付宝公钥证书
	if err3 != nil || err4 != nil || err5 != nil {
		beego.Info("err3", err3)
		beego.Info("err4", err4)
		beego.Info("err5", err5)
		beego.Info("++++++++++++++++++++++++++++++++")
	}
	orderId := this.GetString("orderId")
	totalPrice := this.GetString("totalPrice")

	beego.Error(orderId, totalPrice)

	// 将 key 的验证调整到初始化阶段
	if err != nil {
		fmt.Println(err)
		return
	}

	var p = alipay.TradePagePay{}
	//p.NotifyURL = "http://xxx"
	p.ReturnURL = "http://192.168.88.130:8080/user/payok"
	p.Subject = "我的购物平台"
	p.OutTradeNo = orderId
	p.TotalAmount = totalPrice
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"
  	var url, err1 = client.TradePagePay(p)
	if err1 != nil {
		fmt.Println(err)
		beego.Info("wo")
	}

	var payURL = url.String()
	//fmt.Println(payURL)
	this.Redirect(payURL, 302)
}

//支付成功
func (this *OrderController) PayOk() {
	//获取数据
	//out_trade_no=999998888777
	orderId := this.GetString("out_trade_no")

	//校验数据
	if orderId == "" {
		beego.Info("支付返回数据错误")
		this.Redirect("/user/userCenterOrder", 302)
		return
	}

	//操作数据

	o := orm.NewOrm()
	count, _ := o.QueryTable("OrderInfo").Filter("OrderId", orderId).Update(orm.Params{"Orderstatus": 2})
	if count == 0 {
		beego.Info("更新数据失败")
		this.Redirect("/user/userCenterOrder", 302)
		return
	}

	//返回视图
	this.Redirect("/user/userCenterOrder", 302)
}
