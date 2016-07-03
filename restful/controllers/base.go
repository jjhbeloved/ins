package controllers

import (
	"github.com/astaxie/beego"
	"log"
)

type BaseController struct {

	beego.Controller

}

func (this *BaseController)Prepare()  {
	log.Println("Prepare.....")

	var m map[string]string = make(map[string]string)
	m["header"] = "layout/header.html"
	m["footer"] = "layout/footer.html"
	m["nav"] = "layout/nav.html"
	this.LayoutSections = m
	this.Layout = "layout/layout.html"
}
