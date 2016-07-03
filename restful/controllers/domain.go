package controllers

import (
	"asiainfo.com/ins/restful/models"
	"fmt"
)

// domain
type DomainController struct {
	BaseController
}

// @Title List
// @Description List All Domain
// @Success 200 {string} success
// @Failure 400 Invalid input
// @router /list [get]
func (this *DomainController) List() {
	domains, _ := models.QueryAllDomain(0, -1)
	this.Data["domains"] = domains
	if len(domains) > 0 {
		this.Data["domain_code"] = domains[0].DomainCode
		this.Data["domain_name"] = domains[0].Name
	}
	this.TplName = "domain/list.html"
	this.Render()
}

// @Title List
// @Description List Some Domain
// @Success 200 {string} success
// @Failure 400 Invalid input
// @router /list [post]
func (this *DomainController) PostList() {
	domainOption := this.GetString("DomainOption")
	searchValue := this.GetString("SearchValue")
	fmt.Println(domainOption)
	fmt.Println(searchValue)
}