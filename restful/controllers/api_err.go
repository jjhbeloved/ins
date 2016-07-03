package controllers

import (
	"github.com/astaxie/beego"
)

type Error4xxController struct {
	beego.Controller
}

func (this *Error4xxController) Error404() {
	this.Data["json"] = map[string]interface{}{"error":404}
	this.ServeJSON()
}

