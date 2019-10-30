package main

import (
	_ "shFresh/routers"
	"github.com/astaxie/beego"
	_ "shFresh/models"
)

func main() {
	//beego.SetStaticPath()["/static"] = "static"
	beego.Run()
}

