package controllers

import (
	"asiainfo.com/ins/restful/models"
	"asiainfo.com/ins/utils"
)

// app
type ViewAppController struct {
	BaseController
}

// @Title List
// @Description List All App
// @Success 200 {string} success
// @Failure 400 Invalid input
// @router /list [get]
func (this *ViewAppController) List() {
	apps, code, device, user := models.QueryAllAppView()
	this.Data["apps"] = apps
	this.Data["domain_code"] = code
	this.Data["device_host"] = device
	this.Data["user_name"] = user
	this.TplName = "app/list.html"
	this.Render()
}

// @Title List
// @Description List Some App
// @Success 200 {string} success
// @Failure 400 Invalid input
// @router /list [post]
func (this *ViewAppController) PostList() {
	domainCode := this.GetString("DomainCode")
	deviceName := this.GetString("DeviceHost")
	userName := this.GetString("UserName")
	apps, code, device, user := models.QueryAppView(domainCode, deviceName, userName)
	this.Data["apps"] = apps
	this.Data["domain_code"] = code
	this.Data["device_host"] = device
	this.Data["user_name"] = user
	this.TplName = "app/list.html"
	this.Render()
}

// @Title Info
// @Description App Detail Info
// @Success 200 {string} success
// @Failure 400 Invalid input
// @router /info/:id [get]
func (this *ViewAppController)Info() {
	id, err := this.GetInt(":id")
	utils.PanicIf(err)
	app, ips := models.QueryAppAndIP(id)
	this.Data["app"] = app
	this.Data["ips"] = ips
	this.TplName = "app/info.html"
	this.Render()
}

// @Title Info
// @Description App Detail Info
// @Success 200 {string} success
// @Failure 400 Invalid input
// @router /edit/:id [get]
func (this *ViewAppController)ToEdit() {
	id, err := this.GetInt(":id")
	utils.PanicIf(err)
	app, ips := models.QueryAppAndIP(id)
	this.Data["app"] = app
	this.Data["ips"] = ips
	this.TplName = "app/edit.html"
	this.Render()
}

// @Title Info
// @Description App Detail Info
// @Success 200 {string} success
// @Failure 400 Invalid input
// @router /edit/:id [put]
func (this *ViewAppController)Edit() {
	id, err := this.GetInt(":id")
	utils.PanicIf(err)
	app, ips := models.QueryAppAndIP(id)
	this.Data["app"] = app
	this.Data["ips"] = ips
	this.TplName = "app/edit.html"
	this.Render()
}

// TODO
// @Title Active
// @Description Active App
// @Success 200 {string} success
// @Failure 400 Invalid input
// @router /active [post]
func (this *ViewAppController) Active() {
	jsonRet := &models.RSP{}
	_, err := this.GetInt("Id")
	active, err := this.GetBool("active")
	if err != nil {
		jsonRet.Message = err.Error()
	} else {
		jsonRet.Success = true
		if active == true {
			jsonRet.Message = "start success."
		} else {
			jsonRet.Message = "stop success."
		}
		jsonRet.Data = active
	}
	this.Data["json"] = jsonRet
	this.ServeJSON()
}