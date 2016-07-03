package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

/**
domain 域
 */
type Domain struct {
	DomainCode string           `json:"domain_code" orm:"pk"`
	Name       string           `json:"name" orm:"unique"`
	Device     []*Device        `json:"device" orm:"reverse(many); column(domain_host)"` // 设置多对多的反向关系
}

/**
创建 domain 域 M2M
 */
//func (domain *Domain) CreateM2M(device *Device) error {
//	_, err := orm.NewOrm().QueryM2M(domain, "Device").Add(device)
//	if err != nil {
//		return err
//	}
//	return nil
//}

/**
查询 domain 域, 指定查询的 where col 字段
 */
func (domain *Domain) Query(col ...string) error {
	o := orm.NewOrm()
	return o.Read(domain, col...)
}

/**
创建 domain 域
 */
func (domain *Domain) Create() (int64, error) {
	o := orm.NewOrm()
	return o.Insert(domain)
}

/**
更新 domain 域, 可以指定更新的 col 字段
 */
func (domain *Domain) Update(col ...string) (int64, error) {
	o := orm.NewOrm()
	return o.Update(domain, col...)
}

/**
删除 domain 域
 */
func (domain *Domain) Delete() (int64, error) {
	o := orm.NewOrm()
	//return o.QueryTable(domain).Filter("domain_code", domain.DomainCode).Delete()
	return o.Delete(domain)
}

/**
TODO
on_delete(cascade) 有bug, 好像是 string 类型会被解析成 char 无法删除
 */
type Device struct {
	DomainHost string          `json:"domain_host" orm:"pk"`
	Host       string          `json:"host"`
	Domain     *Domain         `json:"domain" orm:"rel(fk);column(domain_code);on_delete(do_nothing)"` // 设置多对多的反向关系
	Ip         []*Ip           `json:"ip" orm:"reverse(many);column(device_ipv4)"`                     // 设置多对多的反向关系
	Memory     int             `json:"memory"`
	Cpu        int             `json:"cpu"`
	Storage    string          `json:"storage"`
	OsType     string          `json:"osType"`
	User       []*User         `json:"user" orm:"reverse(many);column(device_uid)"`
	Group      []*Group        `json:"group" orm:"reverse(many);column(device_gid)"`
}

/**
查询 device 设备, 指定查询的 where col 字段
 */
func (device *Device) Query(col ...string) error {
	o := orm.NewOrm()
	return o.Read(device, col...)
}

/**
创建 device 设备
 */
func (device *Device) Create() (int64, error) {
	o := orm.NewOrm()
	device.DomainHost = device.Domain.DomainCode + "_" + device.Host
	return o.Insert(device)
}

/**
更新 device 设备, 可以指定更新的 col 字段
 */
func (device *Device) Update(col ...string) (int64, error) {
	o := orm.NewOrm()
	return o.Update(device, col...)
}

/**
删除 device 设备
 */
func (device *Device) Delete() (int64, error) {
	o := orm.NewOrm()
	device.DomainHost = device.Domain.DomainCode + "_" + device.Host
	return o.Delete(device)
}

type Ip struct {
	DeviceIpv4 string        `json:"device_ipv4" orm:"pk"`
	Device     *Device       `json:"device" orm:"rel(fk);column(domain_host);on_delete(do_nothing)"`
	Ipv4       string        `json:"ipv4"`
	Ipv6       string        `json:"ipv6" orm:"null"`
}

/**
查询 ip IP, 指定查询的 where col 字段
 */
func (ip *Ip) Query(col ...string) error {
	o := orm.NewOrm()
	ip.DeviceIpv4 = ip.Device.DomainHost + "_" + ip.Ipv4
	return o.Read(ip, col...)
}

/**
创建 ip IP
 */
func (ip *Ip) Create() (int64, error) {
	o := orm.NewOrm()
	ip.DeviceIpv4 = ip.Device.DomainHost + "_" + ip.Ipv4
	return o.Insert(ip)
}

/**
更新 ip IP, 可以指定更新的 col 字段
 */
func (ip *Ip) Update(col ...string) (int64, error) {
	o := orm.NewOrm()
	ip.DeviceIpv4 = ip.Device.DomainHost + "_" + ip.Ipv4
	return o.Update(ip, col...)
}

/**
删除 ip IP
 */
func (ip *Ip) Delete() (int64, error) {
	o := orm.NewOrm()
	ip.DeviceIpv4 = ip.Device.DomainHost + "_" + ip.Ipv4
	return o.Delete(ip)
}

type User struct {
	DeviceUid string        `json:"device_uid" orm:"pk"`
	Device    *Device       `json:"device" orm:"rel(fk);column(domain_host);on_delete(do_nothing)"`
	Name      string        `json:"name"`
	Uid       int           `json:"uid"`
	Group     *Group        `json:"group" orm:"rel(one)"`
	Password  string        `json:"password" orm:"null"`
	CertPath  string        `json:"cert_path" orm:"null"`
	Apps      []*App        `json:"apps" orm:"reverse(many);column(id)"`
	Updated   time.Time     `json:"last_date" orm:"auto_now;type(datetime)"`
}

/**
查询 user 用户, 指定查询的 where col 字段
 */
func (user *User) Query(col ...string) error {
	o := orm.NewOrm()
	user.DeviceUid = user.Device.DomainHost + "_" + user.Name
	return o.Read(user, col...)
}

/**
创建 user 用户
 */
func (user *User) Create() (int64, error) {
	o := orm.NewOrm()
	user.DeviceUid = user.Device.DomainHost + "_" + user.Name
	return o.Insert(user)
}

/**
更新 user 用户, 可以指定更新的 col 字段
 */
func (user *User) Update(col ...string) (int64, error) {
	o := orm.NewOrm()
	user.DeviceUid = user.Device.DomainHost + "_" + user.Name
	return o.Update(user, col...)
}

/**
删除 user 用户
 */
func (user *User) Delete() (int64, error) {
	o := orm.NewOrm()
	user.DeviceUid = user.Device.DomainHost + "_" + user.Name
	return o.Delete(user)
}

type Group struct {
	DeviceGid string        `json:"device_gid" orm:"pk"`
	Device    *Device       `json:"device" orm:"rel(fk);column(domain_host);on_delete(do_nothing)"`
	Name      string        `json:"name"`
	Gid       int           `json:"gid"`
	SubGroup  []*GroupName  `json:"subGroup" orm:"null;reverse(many);column(device_group)"`
}

/**
查询 group 属组, 指定查询的 where col 字段
 */
func (group *Group) Query(col ...string) error {
	o := orm.NewOrm()
	group.DeviceGid = group.Device.DomainHost + "_" + group.Name
	return o.Read(group, col...)
}

/**
创建 group 属组
 */
func (group *Group) Create() (int64, error) {
	o := orm.NewOrm()
	group.DeviceGid = group.Device.DomainHost + "_" + group.Name
	return o.Insert(group)
}

/**
更新 group 属组, 可以指定更新的 col 字段
 */
func (group *Group) Update(col ...string) (int64, error) {
	o := orm.NewOrm()
	group.DeviceGid = group.Device.DomainHost + "_" + group.Name
	return o.Update(group, col...)
}

/**
删除 group 属组
 */
func (group *Group) Delete() (int64, error) {
	o := orm.NewOrm()
	group.DeviceGid = group.Device.DomainHost + "_" + group.Name
	return o.Delete(group)
}

type GroupName struct {
	DeviceGroup string        `json:"device_group" orm:"pk"`
	Name        string        `json:"name"`
	Group       *Group        `json:"group" orm:"rel(fk);column(device_gid);on_delete(do_nothing)"`
}

/**
查询 group_name 属组名, 指定查询的 where col 字段
 */
func (gname *GroupName) Query(col ...string) error {
	o := orm.NewOrm()
	gname.DeviceGroup = gname.Group.Device.DomainHost + "_" + gname.Name
	return o.Read(gname, col...)
}

/**
创建 group_name 属组名
 */
func (gname *GroupName) Create() (int64, error) {
	o := orm.NewOrm()
	gname.DeviceGroup = gname.Group.Device.DomainHost + "_" + gname.Name
	return o.Insert(gname)
}

/**
更新 group_name 属组名, 可以指定更新的 col 字段
 */
func (gname *GroupName) Update(col ...string) (int64, error) {
	o := orm.NewOrm()
	gname.DeviceGroup = gname.Group.Device.DomainHost + "_" + gname.Name
	return o.Update(gname, col...)
}

/**
删除 group_name 属组名
 */
func (gname *GroupName) Delete() (int64, error) {
	o := orm.NewOrm()
	gname.DeviceGroup = gname.Group.Device.DomainHost + "_" + gname.Name
	return o.Delete(gname)
}

type AppType struct {
	Id   int        `json:"app_id" orm:"pk:auto"`
	Name string     `json:"name" orm:"unique"`
	App  []*App     `json:"app" orm:"null;reverse(many)"`
}

/**
查询 app_type 应用类型, 指定查询的 where col 字段
 */
func (this *AppType) Query(col ...string) error {
	o := orm.NewOrm()
	return o.Read(this, col...)
}

/**
创建 app_type 应用类型
 */
func (this *AppType) Create() (int64, error) {
	o := orm.NewOrm()
	return o.Insert(this)
}

/**
更新 app_type 应用类型, 可以指定更新的 col 字段
 */
func (this *AppType) Update(col ...string) (int64, error) {
	o := orm.NewOrm()
	return o.Update(this, col...)
}

/**
删除 app_type 应用类型
 */
func (this *AppType) Delete() (int64, error) {
	o := orm.NewOrm()
	return o.Delete(this)
}

type App struct {
	Id         int           `json:"app_id" orm:"pk:auto"`
	Name       string        `json:"name"`
	Type       *AppType      `json:"type" orm:"rel(fk);on_delete(cascade)"`
	Version    string        `json:"version"`
	Home       string        `json:"home"`
	Ips        string        `json:"ips"`
	Ports      string        `json:"ports"`
	CheckUrl   string        `json:"checkUrl"`
	CheckRet   string        `json:"checkRet"`
	IsActivity bool          `json:"isActivity"`
	Updated    time.Time     `json:"last_date" orm:"auto_now;type(datetime)"`
	User       *User         `json:"user" orm:"rel(fk);column(device_user);on_delete(cascade)"`
}

// 多字段唯一键
func (app *App) TableUnique() [][]string {
	return [][]string{
		[]string{"Name", "Version", "User"},
	}
}

/**
查询 app 应用, 指定查询的 where col 字段
 */
func (app *App) Query(col ...string) error {
	o := orm.NewOrm()
	return o.Read(app, col...)
}

/**
创建 app 应用
 */
func (app *App) Create() (int64, error) {
	o := orm.NewOrm()
	return o.Insert(app)
}

/**
更新 app 应用, 可以指定更新的 col 字段
 */
func (app *App) Update(col ...string) (int64, error) {
	o := orm.NewOrm()
	return o.Update(app, col...)
}

/**
删除 app 应用
 */
func (app *App) Delete() (int64, error) {
	o := orm.NewOrm()
	return o.Delete(app)
}

// #######################################

/**
查询所有 domain 域, 指定查询的 limit, offset 字段
 */
func QueryAllDomain(limit, offset int) ([]*Domain, error) {
	var domains []*Domain
	o := orm.NewOrm()
	_, err := o.QueryTable("domain").Limit(limit, offset).All(&domains)
	return domains, err
}

/**
查询所有 device 设备, 指定查询的 limit, offset 字段
 */
func QueryAllDevice(limit, offset int) ([]*Device, error) {
	var devices []*Device
	o := orm.NewOrm()
	_, err := o.QueryTable("device").Limit(limit, offset).RelatedSel().All(&devices)
	return devices, err
}

/**
查询所有 ip IP, 指定查询的 limit, offset 字段
 */
func QueryAllIP(limit, offset int) ([]*Ip, error) {
	var ips []*Ip
	o := orm.NewOrm()
	_, err := o.QueryTable("ip").Limit(limit, offset).All(&ips)
	return ips, err
}

/**
查询 foreign key 相关 ip IP, 指定查询的 limit, offset 字段
 */
func QueryIPs(deviceId string) ([]*Ip, error) {
	var ips []*Ip
	o := orm.NewOrm()
	_, err := o.QueryTable("ip").Filter("domain_host", deviceId).All(&ips)
	return ips, err
}

/**
查询所有 app App, 指定查询的 limit, offset 字段
 */
func QueryAllApp(limit, offset int, tab ...interface{}) ([]*App, error) {
	var apps []*App
	o := orm.NewOrm()
	_, err := o.QueryTable("app").Limit(limit, offset).RelatedSel(tab...).All(&apps)
	return apps, err
}

/**
查询 primary key 相关 app App, 指定查询的 limit, offset 字段
 */
func QueryApp(app *App, tab ...interface{}) error {
	o := orm.NewOrm()
	err := o.QueryTable("app").Filter("id", app.Id).RelatedSel(tab...).One(app)
	return err
}

/**
查询 user 相关 app App, 指定查询的 limit, offset 字段
 */
func QueryAppByUser(device_user string, tab ...interface{}) ([]*App, error) {
	var apps []*App
	o := orm.NewOrm()
	_, err := o.QueryTable("app").Filter("device_user", device_user).RelatedSel(tab...).All(&apps)
	return apps, err
}
