package models

import (
	"asiainfo.com/ins/utils"
)

type AppView struct {
	DomainCode string        `json:"domain_code"`
	HostName   string        `json:"hostname"`
	IP         []string      `json:"ip"`
	User       string        `json:"user"`
	AppId      int           `json:"app_id"`
	AppName    string        `json:"app_name"`
	AppVersion string        `json:"app_version"`
	AppPorts   string        `json:"app_ports"`
	IsActivity bool          `json:"isActivity"`
}

func QueryAppAndIP(id int) (*App, []*Ip) {
	app := &App{Id: id}
	err := QueryApp(app)
	utils.PanicIf(err)
	ips, err := QueryIPs(app.User.Device.DomainHost)
	//utils.PanicIf(err)
	return app, ips
}

func findDomain(device *Device, domains []*Domain) {
	//for _, domain := range domains {
	//	if device.Domain.DomainCode ==  domain.DomainCode {
	//		device.Ip = append(device.Ip, ip)
	//	}
	//}
}

func QueryAllAppView() ([]*AppView, string, string, string) {
	ips, err := QueryAllIP(-1, 0)
	if err != nil {
		return nil, "", "", ""
	}
	apps, err := QueryAllApp(-1, 0)
	if err != nil {
		return nil, "", "", ""
	}
	for _, app := range apps {
		findIP(app.User.Device, ips)
	}
	return appCastToView(apps)
}

func QueryAppView(domainCode, deviceName, userName string) ([]*AppView, string, string, string) {
	apps, _ := QueryAppByUser(domainCode + "_" + deviceName + "_" + userName)
	if len(apps) > 0 {
		ips, _ := QueryIPs(apps[0].User.Device.DomainHost)
		if len(ips) > 0 {
			for _, app := range apps {
				findIP(app.User.Device, ips)
			}
		}
	}
	return appCastToView(apps)
}

func appCastToView(apps []*App) ([]*AppView, string, string, string) {
	var views []*AppView
	var code, hostname, username string
	for n, app := range apps {
		if n == 0 {
			code = app.User.Device.Domain.DomainCode
			hostname = app.User.Device.Host
			username = app.User.Name
		}
		view := &AppView{
			DomainCode: app.User.Device.Domain.DomainCode,
			HostName: app.User.Device.Host,
			User: app.User.Name,
			AppId: app.Id,
			AppName: app.Name,
			AppVersion: app.Version,
			AppPorts: app.Ports,
			IsActivity: app.IsActivity,
		}
		for _, ip := range app.User.Device.Ip {
			view.IP = append(view.IP, ip.Ipv4)
		}
		views = append(views, view)
	}
	return views, code, hostname, username
}


func findIP(device *Device, ips []*Ip) {
	for _, ip := range ips {
		if ip.Device.DomainHost == device.DomainHost {
			device.Ip = append(device.Ip, ip)
		}
	}
}