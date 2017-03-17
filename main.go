package main

import (
	"github.com/astaxie/beego"
	_ "github.com/partikle/partikle/routers"
)

func main() {
	beego.SetStaticPath("/swagger", "swagger/")
	beego.Run()
}
