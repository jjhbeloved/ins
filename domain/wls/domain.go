package wls12c

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os/exec"
	"os"
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/utils"
	"path/filepath"
)

type Wls12c struct {
	WLS_HOME    string        `json:"wlsHome"`
	ListenAddr  string        `json:"listenAddr"`
	ListenPort  string        `json:"listenPort"`
	UserName    string        `json:"userName"`
	PassWord    string        `json:"passWord"`
	Mode        string        `json:"mode"`
	DomainPath  string        `json:"domainPath"`
	JDK_HOME    string        `json:"jdkHome"`
	Option      string        `json:"option"`
	ConsolePath string        `json:"consolePath"`
}

func (w *Wls12c) Builder(str string) error {
	return nil
}

func (w *Wls12c) Json(bs []byte) error {
	return json.Unmarshal(bs, &w)
}

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Method */

func (w *Wls12c) Add() error {
	wlst := filepath.Join(w.WLS_HOME, "wlserver", "common", "bin", "wlst.sh")
	conf := w.touchConf()

	// 根据模板生成domain
	cmd := exec.Command(wlst, conf)
	//defer logs.Print(os.Remove(conf))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// 创建无密码
	err := w.touchNonPasswdRun()
	if err != nil {
		return err
	}

	// 创建控制 admin 脚本
	w.touchConsoleScript()

	// 直接启动 admin
	w.run()
	return nil
}

func (w *Wls12c) Remove() error {
	return nil
}
/* End Method */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Configuration File */

const template = `
readTemplate('%s/wlserver/common/templates/wls/wls.jar')
cd('/Servers/AdminServer')
set('ListenAddress', '%s')
set('ListenPort', %s)
cd('/Security/base_domain/User/weblogic')
cmo.setName('%s')
cmo.setUserPassword('%s')
setOption('ServerStartMode', '%s')
setOption('JavaHome', '%s')
setOption('OverwriteDomain', 'true')
writeDomain('%s')
closeTemplate()
`
/**
 * touch conf local file
 */
func (w *Wls12c) touchConf() string {
	os.MkdirAll(utils.TMPD, 0750)
	file, err := ioutil.TempFile(utils.TMPD, "wls12c_domain_")//在DIR目录下创建tmp为文件名前缀的文件，获得file文件指针，DIR必须存在，否则创建不成功
	defer file.Close()
	if err != nil {
		logs.Print(fmt.Errorf("create %s error.", file.Name()))
		panic(err)
		return ""
	}

	file.WriteString(
		fmt.Sprintf(
			template,
			w.WLS_HOME, w.ListenAddr, w.ListenPort,
			w.UserName, w.PassWord,
			w.Mode, w.JDK_HOME, w.DomainPath,
		))

	return file.Name()
}
/* End Configuration File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Security File */

const templateSec = `
username=%s
password=%s
`
/**
 * touch non password run file
 */
func (w *Wls12c) touchNonPasswdRun() error {
	adminPath := filepath.Join(w.DomainPath, "servers", "AdminServer")
	adminSec := filepath.Join(adminPath ,"security")
	adminLog := filepath.Join(adminPath ,"logs")
	adminBoot := filepath.Join(adminPath ,"security", "boot.properties")

	os.MkdirAll(adminSec, 0750)
	os.MkdirAll(adminLog, 0755)

	return ioutil.WriteFile(
		adminBoot,
		[]byte(fmt.Sprintf(templateSec, w.UserName, w.PassWord)),
		0600,
	)
}
/* End Security File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Console File */

const templateStartConsole = `#!/bin/bash

echo "%s admin starting..."
nohup %s/bin/startWebLogic.sh 1>>%s/servers/AdminServer/logs/admin_auto.log 2>>%s/servers/AdminServer/logs/failedServer.log &
echo "%s admin started, pls wating 30 sec..."
`
const templateStopConsole = `#!/bin/bash

echo "%s admin stoping..."
%s/bin/stopWebLogic.sh
echo "%s admin stoped, pls wating 30 sec..."
`
const templateRestartConsole = `#!/bin/bash

echo "%s admin restarting..."
%s
%s
echo "%s admin restarted, pls wating 30 sec..."
`
/**
 * touch console file
 */
func (w *Wls12c) touchConsoleScript() {
	err := utils.MkdirConsolesPath(w.ConsolePath)
	if err != nil {
		logs.Print(err)
	}
	domainName := utils.GetFileName(w.DomainPath)

	start := filepath.Join(w.ConsolePath, "start", "start_" + domainName + ".sh")
	stop := filepath.Join(w.ConsolePath, "stop", "stop_" + domainName + ".sh")
	restart := filepath.Join(w.ConsolePath, "restart", "restart_" + domainName + ".sh")
	// start
	logs.Print(ioutil.WriteFile(
		start,
		[]byte(fmt.Sprintf(templateStartConsole, domainName, w.DomainPath, w.DomainPath, w.DomainPath, domainName)),
		0750,
	))


	// stop
	logs.Print(ioutil.WriteFile(
		stop,
		[]byte(fmt.Sprintf(templateStopConsole, domainName, domainName)),
		0750,
	))

	// restart
	logs.Print(ioutil.WriteFile(
		restart,
		[]byte(fmt.Sprintf(templateRestartConsole, domainName, stop, start, domainName)),
		0750,
	))
}
/* End Console File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Run */
/**
 * run wls12c
 */
func (w *Wls12c) run() {
	domainName := utils.GetFileName(w.DomainPath)
	start := filepath.Join(w.ConsolePath, "start", "start_" + domainName + ".sh")
	cmd := exec.Command(start)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logs.Print(cmd.Run())
}
/* End Run */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
