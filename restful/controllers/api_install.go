package controllers

import (
	"github.com/astaxie/beego"
	"asiainfo.com/ins/install/jdk"
	"asiainfo.com/ins/install/memcached"
	"asiainfo.com/ins/install/amq"
	"asiainfo.com/ins/install/zookeeper"
	"asiainfo.com/ins/install/wls"
	"asiainfo.com/ins/install/tomcat"
	"asiainfo.com/ins/install/mysql"
	"asiainfo.com/ins/install/storm"
	"asiainfo.com/ins/install/redis"
	"asiainfo.com/ins/cli"
	"asiainfo.com/ins/restful/models"
	"fmt"
	"os"
	"asiainfo.com/ins/utils"
)

type InstallController struct {
	beego.Controller
}

// @Title PostInstall
// @Description install basic model
// @Success 200 {string} success
// @Failure 400 Invalid input
// @router /:type [post]
func (this *InstallController) Install() {
	os.Mkdir(utils.TMPD, 0777)
	os.Mkdir(cli.LOGS_PATH, 0750)
	os.Mkdir(cli.PKG_PATH, 0750)

	typ := this.GetString(":type")
	body := this.Ctx.Input.RequestBody
	var ins Installer
	var rsp map[string]*models.RSP
	this.Ctx.Output.Status = 400

	switch typ {
	case cli.WLS12CCONF:
		var wls12 wls.Wls12c
		e := wls12.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		ins = &wls12
	case cli.TOMCATCONF:
		var tomcat tomcat.Tomcat
		e := tomcat.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		ins = &tomcat
	case cli.JDKCONF:
		var jdk jdk.Jdk
		e := jdk.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		ins = &jdk
	case cli.MEMCACHEDCONF:
		var mem memcached.Memcached
		e := mem.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		ins = &mem
	case cli.ACTIVEMQCONF:
		var amq amq.AMQ
		e := amq.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		ins = &amq
	case cli.ZKCONF:
		var zk zookeeper.ZK
		e := zk.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		ins = &zk
	case cli.REDISCONF:
		var redis redis.Redis
		e := redis.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		ins = &redis
	case cli.STORMCONF:
		var storm storm.Storm
		e := storm.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		ins = &storm
	case cli.MYSQLCONF:
		var mysql mysql.Mysql
		e := mysql.Json(body)
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		ins = &mysql
	default:
		ins = nil
	}
	if ins != nil {
		e := ins.Install()
		if e != nil {
			rsp = models.HBD_xxx(400, "input request body have error.", e)
			goto over
		}
		rsp = models.HBD_200(fmt.Sprintf("Install %s Success", typ))
		this.Ctx.Output.Status = 200

	} else {
		rsp = models.HBD_xxx(500, "input request type not exist.", fmt.Errorf("%s is not found on the server.", typ))
		this.Ctx.Output.Status = 500
	}
	over:
	this.Data["json"] = rsp
	this.ServeJSON()
}

type Installer interface {
	Install() error
}