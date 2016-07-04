package controllers

import (
	"github.com/astaxie/beego"
	"asiainfo.com/ins/cli"
	"asiainfo.com/ins/restful/models"
	"fmt"
	"os"
	"asiainfo.com/ins/utils"
	"asiainfo.com/ins/domain/redis"
	"asiainfo.com/ins/domain/tomcat"
	"asiainfo.com/ins/domain/wls"
	"asiainfo.com/ins/domain/mysql"
)

type DomainController struct {
	beego.Controller
}

// @Title PostDomain
// @Description create basic domain
// @Success 200 {string} success
// @Failure 400 Invalid input
// @router /:type [post]
func (this *DomainController) Domain() {
	os.Mkdir(utils.TMPD, 0777)
	os.Mkdir(cli.LOGS_PATH, 0750)
	os.Mkdir(cli.PKG_PATH, 0750)

	typ := this.GetString(":type")
	body := this.Ctx.Input.RequestBody
	var dom Domain
	var option string
	var rsp map[string]*models.RSP
	this.Ctx.Output.Status = 400

	switch typ {
	case cli.WLS12CPREFIX:
		var wls12 wls.Wls12c
		e := wls12.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		dom = &wls12
		option = wls12.Option
	case cli.TOMCATPREFIX:
		var tomcat tomcat.Tomcat
		e := tomcat.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		dom = &tomcat
		option = tomcat.Option
	case cli.REDISPREFIX:
		var redis redis.Redis
		e := redis.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		dom = &redis
		option = redis.Option
	case cli.MYSQLPREFIX:
		var mysql mysql.Mysql
		e := mysql.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		dom = &mysql
		option = mysql.Option
	default:
		dom = nil
	}
	if dom != nil {
		var e error
		if option == "REMOVE" {
			e = dom.Remove()
		} else {
			e = dom.Add()
		}
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		rsp = models.HBD_200(fmt.Sprintf("Domain create %s Success", typ))
		this.Ctx.Output.Status = 200

	} else {
		rsp = models.HBD_xxx(500, "input request type not exist.", fmt.Errorf("%s is not found on the server.", typ))
		this.Ctx.Output.Status = 500
	}
	over:
	this.Data["json"] = rsp
	this.ServeJSON()
}

type Domain interface {
	Add() error
	Remove() error
}
