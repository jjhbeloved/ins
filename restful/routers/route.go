// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"asiainfo.com/ins/restful/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSGet("/", func(ctx *context.Context) {
			ctx.Output.Body([]byte("notAllowed"))
		}),
		beego.NSNamespace("/err",
			beego.NSInclude(
				&controllers.Error4xxController{},
			),
		),
		beego.NSNamespace("/install",
			beego.NSInclude(
				&controllers.InstallController{},
			),
		),
		beego.NSNamespace("/domain",
			beego.NSInclude(
				&controllers.DomainController{},
			),
		),
	)
	app := beego.NewNamespace("/app",
		beego.NSInclude(
			&controllers.ViewAppController{},
		),
	)
	domain := beego.NewNamespace("/domain",
		beego.NSInclude(
			&controllers.ViewDomainController{},
		),
	)
	beego.Router("/", &controllers.IndexController{}, "*:Index")
	beego.AddNamespace(ns, app, domain)
}