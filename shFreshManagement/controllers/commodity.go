package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"encoding/gob"
	"bytes"
	"shFreshManagement/models"
	"github.com/gomodule/redigo/redis"
	"path"
	"github.com/weilaihui/fdfs_client"
)

type CommodityController struct {
	beego.Controller
}

// 封装上传文件函数
func UploadFile(this *beego.Controller, filePath string) string {

	//通过GetFile获取文件信息
	f, h, err := this.GetFile(filePath)
	if err != nil {
		beego.Info("获取文件信息失败:", err)
	}
	defer f.Close()
	//然后对上传的文件进行格式和大小判断
	//1.判断文件格式
	ext := path.Ext(h.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		beego.Info("上传文件格式不正确")
		return "error"
	}

	//2.文件大小
	if h.Size > 5000000 {
		beego.Info("文件太大，不允许上传")
		return "error"
	}
	//3.上传文件
	//先获取一个[]byte
	fileBuffer := make([]byte, h.Size)
	//把文件数据读入到fileBuffer中
	f.Read(fileBuffer)
	beego.Info("获取文件成功")
	//获取client对象
	client, err := fdfs_client.NewFdfsClient("/etc/fdfs/client.conf")
	if err != nil {
		beego.Info("获取client对象失败:", err)
	}
	//上传
	fdfsresponse, _ := client.UploadByBuffer(fileBuffer, ext)
	//返回文件ID
	return fdfsresponse.RemoteFileId
}

// 展示商品列表页
func (this *CommodityController) ShowCommodityList() {
	//session判断

	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}

	o := orm.NewOrm()
	qs := o.QueryTable("GoodsSKU") //queryseter

	//查询总记录数
	typeName := this.GetString("select")
	var count int64

	//获取总页数
	pageSize := 6

	//获取页码
	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}

	//获取数据
	//作用就是获取数据库部分数据,第一个参数，获取几条,第二个参数，从那条数据开始获取,返回值还是querySeter
	//起始位置计算
	start := (pageIndex - 1) * pageSize

	//qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)

	if typeName == "" {
		count, _ = qs.Limit(pageSize, start).RelatedSel("GoodsType").Filter("GoodsType__Name", "新鲜水果").Count()
	} else {
		count, _ = qs.Limit(pageSize, start).RelatedSel("GoodsType").Filter("GoodsType__Name", typeName).Count()
	}
	pageCount := math.Ceil(float64(count) / float64(pageSize))

	//根据选中的类型查询相应类型商品
	var goods []models.GoodsSKU

	if typeName == "" {
		qs.Limit(pageSize, start).All(&goods)
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name", "新鲜水果").All(&goods)
	} else {
		o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("GoodsType__Name", typeName).All(&goods)
		qs.Limit(pageSize, start).RelatedSel("GoodsType").Filter("GoodsType__Name", typeName).All(&goods)
	}
	this.Data["goods"] = goods

	//获取商品类型
	var types []models.GoodsType
	//获取数据
	conn, err := redis.Dial("tcp", "192.168.88.130:6379")
	if err != nil {
		beego.Info("redis连接失败")
		return
	}
	//尝试从redis中获取数据
	//解码
	rep, err := conn.Do("get", "types")
	data, _ := redis.Bytes(rep, err)
	//解码
	dec := gob.NewDecoder(bytes.NewReader(data)) //获取解码器
	dec.Decode(&types)
	if len(types) == 0 {
		//从redis中获取数据不成功，从mysql获取数据
		o.QueryTable("GoodsType").All(&types)
		//把获取到的数据存储到redis中
		//编码操作
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer) //获取编码器
		enc.Encode(&types)             //编码
		//数据存储
		conn.Do("set", "types", buffer.Bytes())
		beego.Info("从mysql中获取数据")
	}
	//解码操作
	o.QueryTable("GoodsType").All(&types)
	this.Data["types"] = types

	//传递数据
	this.Data["userName"] = userName.(string)
	this.Data["typeName"] = typeName
	this.Data["pageIndex"] = pageIndex
	this.Data["pageCount"] = int(pageCount)
	this.Data["count"] = count

	//指定试图布局
	this.Layout = "layout.html"
	this.TplName = "index.html"
}

// 查看商品详细
func (this *CommodityController) ShowCommodityDetail() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	goodsId, err := this.GetInt("id")
	if err != nil {
		beego.Info("获取商品id失败", err)
		return
	}
	this.Data["userName"] = userName

	var goods models.GoodsSKU
	o := orm.NewOrm()
	o.QueryTable("GoodsSKU").RelatedSel("GoodsType").Filter("Id", goodsId).One(&goods)
	this.Data["goods"] = goods
	this.Layout = "layout.html"
	this.TplName = "goodsContent.html"
}

// 编辑商品
func (this *CommodityController) ShowEditCommodity() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	//获取商品id
	goodsId, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取商品id失败", err)
	}
	var goodsSku models.GoodsSKU
	goodsSku.Id = goodsId
	o := orm.NewOrm()
	o.Read(&goodsSku)
	this.Data["goods"] = goodsSku
	this.Layout = "layout.html"
	this.TplName = "update.html"
}

// 处理编辑商品
func (this *CommodityController) HandleEditCommodity() {
	goodsName := this.GetString("goodsName")
	desc := this.GetString("desc")
	goodsStock, err := this.GetInt("goodsStock")
	goodsPrice, err := this.GetInt("goodsPrice")
	goodsId, err := this.GetInt("id")
	if err != nil {
		beego.Info("获取商品数据失败:", err)
	}
	uploadname := UploadFile(&this.Controller, "uploadname")
	var goodsSku models.GoodsSKU
	goodsSku.Id = goodsId
	o := orm.NewOrm()
	o.Read(&goodsSku)
	if uploadname == "error" {

		beego.Error(goodsSku)
		uploadname = goodsSku.Image
		beego.Info(uploadname)
	}
	selectTypeId := goodsSku.GoodsType.Id
	selectGoodsSPUId := goodsSku.Goods.Id

	goodsSku.Name = goodsName
	goodsSku.Desc = desc
	goodsSku.Stock = goodsStock
	goodsSku.Price = goodsPrice
	goodsSku.Image = uploadname
	goodsSku.Unite = "1"

	var goodsType models.GoodsType
	o.QueryTable("GoodsType").Filter("Id", selectTypeId).One(&goodsType)
	var goodsSpu models.Goods
	o.QueryTable("Goods").Filter("Id", selectGoodsSPUId).One(&goodsSpu)
	goodsSku.GoodsType = &goodsType
	goodsSku.Goods = &goodsSpu
	beego.Info(goodsSku)
	o.Update(&goodsSku)

	this.ShowCommodityList()

}

// 删除商品
func (this *CommodityController) DelCommodity() {
	goodsId, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取id失败", err)
	}
	var goodsSku models.GoodsSKU
	goodsSku.Id = goodsId
	o := orm.NewOrm()
	o.Delete(&goodsSku, "Id")
	this.ShowCommodityList()
}

// 添加商品
func (this *CommodityController) ShowAddCommodity() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	this.Data["userName"] = userName

	// 商品类型查询
	var types []models.GoodsType
	o := orm.NewOrm()
	o.QueryTable("GoodsType").All(&types)
	this.Data["types"] = types
	var goodsSPU []models.Goods
	o.QueryTable("Goods").All(&goodsSPU)
	this.Data["goodsSPU"] = goodsSPU

	this.Layout = "layout.html"
	this.TplName = "add.html"
}

// 处理添加商品
func (this *CommodityController) HandleAddCommodity() {
	goodsName := this.GetString("goodsName")
	selectType := this.GetString("selectType")
	selectGoodsSPU := this.GetString("selectGoodsSPU")
	desc := this.GetString("desc")
	goodsStock, err := this.GetInt("goodsStock")
	goodsPrice, err := this.GetInt("goodsPrice")
	displayType := this.GetString("displayType")
	unite := this.GetString("unite")
	if err != nil {
		beego.Info("获取商品数据失败:", err)
	}

	uploadname := UploadFile(&this.Controller, "uploadname")
	o := orm.NewOrm()
	var goodsSku models.GoodsSKU
	goodsSku.Unite = unite
	goodsSku.Name = goodsName
	goodsSku.Desc = desc
	goodsSku.Stock = goodsStock
	goodsSku.Price = goodsPrice
	goodsSku.Image = uploadname
	var goodsType models.GoodsType
	o.QueryTable("GoodsType").Filter("Name", selectType).One(&goodsType)
	var goodsSpu models.Goods
	o.QueryTable("Goods").Filter("Name", selectGoodsSPU).One(&goodsSpu)
	beego.Info("----", goodsType)
	beego.Info("------", goodsSpu)
	goodsSku.GoodsType = &goodsType
	goodsSku.Goods = &goodsSpu
	o.Insert(&goodsSku)
	var indexTypeGoodsBanner models.IndexTypeGoodsBanner
	indexTypeGoodsBanner.GoodsType = &goodsType
	indexTypeGoodsBanner.GoodsSKU = &goodsSku
	// 判断是什么类型商品
	if displayType == "文字商品" {
		indexTypeGoodsBanner.DisplayType = 0
	} else {
		indexTypeGoodsBanner.DisplayType = 1
	}
	// 类型商品的index
	count, err := o.QueryTable("GoodsType").Filter("Name", selectType).Count()
	if err != nil {
		beego.Error(err)
	}
	if count == 0 {
		indexTypeGoodsBanner.Index = int(count)
	} else {
		indexTypeGoodsBanner.Index = int(count) + 1
	}
	o.Insert(&indexTypeGoodsBanner)

	this.Redirect("/Commodity/AddCommodity", 302)
}

// 展示添加分类页面
func (this *CommodityController) ShowAddType() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	this.Data["userName"] = userName

	o := orm.NewOrm()
	var types []models.GoodsType
	o.QueryTable("GoodsType").All(&types)

	//传递数据
	this.Data["types"] = types
	this.Layout = "layout.html"
	this.TplName = "addType.html"

}

// 处理添加分类数据
func (this *CommodityController) HandleAddType() {

	//获取数据
	typeName := this.GetString("typeName")
	logoPath := UploadFile(&this.Controller, "uploadlogo")
	typeImage := UploadFile(&this.Controller, "uploadTypeImage")
	//校验数据
	beego.Info(typeName, logoPath, typeImage)
	if typeName == "" || logoPath == "" || typeImage == "" {
		beego.Info("信息不完整，请重新输入")
		return
	}
	//处理数据
	//插入操作
	o := orm.NewOrm()
	var goodsType models.GoodsType
	goodsType.Name = typeName
	goodsType.Logo = logoPath
	goodsType.Image = typeImage
	o.Insert(&goodsType)

	//返回视图
	this.Redirect("/Commodity/AddCommodityType", 302)

}

// 分类详情页面
func (this *CommodityController) ShowTypeDetail() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}

	typeId, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取id失败", err)
	}
	o := orm.NewOrm()
	var goodsType models.GoodsType
	o.QueryTable("GoodsType").Filter("Id", typeId).One(&goodsType)

	this.Data["goodsType"] = goodsType
	this.Layout = "layout.html"
	this.TplName = "typeContent.html"
}

// 删除分类
func (this *CommodityController) DelType() {
	typeId, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取分类id失败", err)
		return
	}
	var goodsType models.GoodsType
	o := orm.NewOrm()
	o.QueryTable("GoodsType").Filter("Id", typeId).One(&goodsType)
	o.Delete(&goodsType)

	var types []models.GoodsType
	o.QueryTable("GoodsType").All(&types)

	//传递数据
	this.Data["types"] = types

	this.Layout = "layout.html"
	this.TplName = "addType.html"
}

// 添加商品SPU
func (this *CommodityController) ShowAddCommoditySPU() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	this.Data["userName"] = userName

	var goodsSpus []models.Goods
	o := orm.NewOrm()
	o.QueryTable("Goods").All(&goodsSpus)

	this.Data["goodsSpus"] = goodsSpus
	this.Layout = "layout.html"
	this.TplName = "addGoodsSPU.html"
}

// 处理添加商品SPU
func (this *CommodityController) HandleAddCommoditySPU() {
	spuName := this.GetString("spuName")
	spuDetail := this.GetString("spuDetail")
	o := orm.NewOrm()
	var goodsSpu models.Goods
	goodsSpu.Name = spuName
	goodsSpu.Detail = spuDetail
	o.Insert(&goodsSpu)

	this.Redirect("/Commodity/AddCommoditySPU", 302)

}

// 处理展示商品SPU
func (this *CommodityController) ShowCommoditySPUDetail() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	this.Data["userName"] = userName

	spuId, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取商品SPU的id失败", err)
		return
	}
	o := orm.NewOrm()
	var goods models.Goods
	o.QueryTable("Goods").Filter("Id", spuId).One(&goods)
	this.Data["goods"] = goods

	this.Layout = "layout.html"
	this.TplName = "goodsSPUContent.html"
}

// 删除商品SPU
func (this *CommodityController) DelCommoditySPU() {
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	this.Data["userName"] = userName

	// 获取id
	spuId, err := this.GetInt("id")
	if err != nil {
		beego.Error("获取SPU的id失败", err)
		return
	}
	var goods models.Goods
	goods.Id = spuId
	o := orm.NewOrm()
	o.Delete(&goods)

	var goodsSpus []models.Goods
	o.QueryTable("Goods").All(&goodsSpus)

	this.Data["goodsSpus"] = goodsSpus
	this.Layout = "layout.html"
	this.TplName = "addGoodsSPU.html"
}

func (this *CommodityController) TestFunc() {
	o := orm.NewOrm()
	//var types models.GoodsType
	count, err := o.QueryTable("GoodsType").Filter("Name", "AJ").Count()
	if err != nil {
		beego.Error(err)
	}
	beego.Error(count)
	this.TplName = "test.html"
}
