package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:DomainController"] = append(beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:DomainController"],
		beego.ControllerComments{
			"Domain",
			`/:type`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:InstallController"] = append(beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:InstallController"],
		beego.ControllerComments{
			"Install",
			`/:type`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"] = append(beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"],
		beego.ControllerComments{
			"List",
			`/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"] = append(beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"],
		beego.ControllerComments{
			"PostList",
			`/list`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"] = append(beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"],
		beego.ControllerComments{
			"Info",
			`/info/:id`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"] = append(beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"],
		beego.ControllerComments{
			"ToEdit",
			`/edit/:id`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"] = append(beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"],
		beego.ControllerComments{
			"Edit",
			`/edit/:id`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"] = append(beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewAppController"],
		beego.ControllerComments{
			"Active",
			`/active`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewDomainController"] = append(beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewDomainController"],
		beego.ControllerComments{
			"List",
			`/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewDomainController"] = append(beego.GlobalControllerRouter["asiainfo.com/ins/restful/controllers:ViewDomainController"],
		beego.ControllerComments{
			"PostList",
			`/list`,
			[]string{"post"},
			nil})

}
