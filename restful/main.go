package main

import (
	oos "asiainfo.com/ins/restful/models"
	_ "asiainfo.com/ins/restful/docs"
	_ "asiainfo.com/ins/restful/routers"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"

	_ "github.com/mattn/go-sqlite3"
	"runtime"
	"asiainfo.com/ins/restful/controllers"
)

const (
	DRIVER_NAME = "sqlite3"
	DATA_SOURCE = "file:/tmp/c.c"
	MAX_IDLE_CONN = 5
	MAX_OPEN_CONN = 30
)

func registerDB() {
	orm.Debug = true
	orm.RegisterDataBase("default", DRIVER_NAME, DATA_SOURCE, MAX_IDLE_CONN, MAX_OPEN_CONN)
	orm.RegisterModel(
		new(oos.User), new(oos.App),
		new(oos.Device), new(oos.Domain),
		new(oos.Ip), new(oos.Group),
		new(oos.GroupName), new(oos.AppType),
	)
	orm.RunCommand()
	orm.RunSyncdb("default", false, true)
}

func init() {
	// set CPU
	runtime.GOMAXPROCS(runtime.NumCPU())
	registerDB()
}

func main() {

	//http.Handle("/", http.FileServer(http.Dir("/tmp")))
	//http.ListenAndServe(":8123", nil)

	//te2st()
	beeGo()
}

func beeGo() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.SetStaticPath("static", "static")
		beego.SetStaticPath("public", "public")
	}
	beego.ErrorController(&controllers.Error4xxController{})
	beego.Run()
}

func te2st() {
	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库
	domain := new(oos.Domain)
	domain.DomainCode = "22"
	domain.Name = "HU"

	gname := new(oos.GroupName)
	gname.Name = "root"
	group := new(oos.Group)
	group.Name = "root"
	group.Gid = 0
	group.SubGroup = append(group.SubGroup, gname)


	app := new(oos.App)
	app.Name = "agrldcz"
	app.Version = "v1"
	app.Home = "/tmp/webapps/serviceAgent"
	app.Ips = "127.0.0.1"
	app.Ports = "8280;8281"
	app.CheckUrl = "TCP"
	app.CheckRet = "200"

	appType := new(oos.AppType)
	appType.Name = "basic"
	appType.App = append(appType.App, app)

	user := new(oos.User)
	user.Name = "root"
	user.Password = "123"
	user.Uid = 0

	ip := new(oos.Ip)
	ip.Ipv4 = "127.0.0.1"

	device := new(oos.Device)
	device.Cpu = 4
	device.Host = "intapp1"
	device.Memory = 16
	device.Storage = "128G"
	device.OsType = "redhat 7.1"
	device.DomainHost = domain.DomainCode + "_" + device.Host
	device.Domain = domain
	device.Ip = append(device.Ip, ip)
	device.Group = append(device.Group, group)
	device.Ip = append(device.Ip, ip)
	device.Ip = append(device.Ip, ip)

	domain.Device = append(domain.Device, device)
	user.Device = device
	group.Device = device
	ip.Device = device

	gname.Group = group
	user.Group = group

	group.DeviceGid = group.Device.DomainHost + "_" + group.Name
	user.DeviceUid = user.Device.DomainHost + "_" + user.Name

	app.User = user
	user.Apps = append(user.Apps, app)

	appType.Create()
	app.Type = appType
	domain.Create()
	device.Create()
	gname.Create()
	group.Create()
	ip.Create()
	_, err := app.Create()
	if err != nil {
		panic(err)
	}
	user.Create()
	appType.Name = "server"
	appType.Create()
	domain.Query("domain_code")
}