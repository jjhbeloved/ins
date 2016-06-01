package wls12c
import (
	"io/ioutil"
	"fmt"
	"asiainfo.com/ins/logs"
	"os"
	"os/exec"
	"encoding/json"
	"path/filepath"
	"asiainfo.com/ins/utils"
	"time"
)

type Wls12c struct {
	WLS_HOME	string	`json:"wlsHome"`
	AdminAddr	string	`json:"adminAddr"`
	AdminPort	string	`json:"adminPort"`
	SrvName		string 	`json:"srvName"`
	ListenAddr	string	`json:"listenAddr"`
	ListenPort	string	`json:"listenPort"`
	UserName	string	`json:"userName"`
	PassWord	string	`json:"passWord"`
	Margs		string	`json:"margs"`
	DomainPath	string	`json:"domainPath"`
	Jars		[]string`json:"jars"`
	Option		string	`json:"option"`
	ConsolePath	string	`json:"consolePath"`
}

func (w *Wls12c) Json(bs []byte) error {
	return json.Unmarshal(bs, &w)
}

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Method */

func (w *Wls12c) Add() error {
	wlst := w.WLS_HOME + "/wlserver/common/bin/wlst.sh"
	conf := w.touchConf()
	//defer logs.Print(os.Remove(conf))

	// 根据模板生成domain
	cmd := exec.Command(wlst, conf)
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

	// 创建控制 srv 脚本
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
def editMode():
  edit()
  startEdit()
def editActivate():
  save()
  activate()
connect("%s", "%s", 't3://%s:%s')
editMode()
cd('/')
cmo.createServer('%s')
cd('/Servers/%s')
cmo.setListenAddress('%s')
cmo.setListenPort(int(%s))
cmo.setListenPortEnabled(true)
cmo.setExternalDNSName('%s')
cd('/Servers/%s/SSL/%s')
cmo.setEnabled(false)
cd('/Servers/%s/Log/%s')
cmo.setBufferSizeKB(32)
cmo.setNumberOfFilesLimited(true)
cmo.setFileCount(101)
cmo.setDomainLogBroadcastSeverity('Info')
cmo.setMemoryBufferSeverity('Info')
cmo.setLogFileSeverity('Info')
cmo.setLoggerSeverity('Info')
cmo.setStdoutSeverity('Info')
cmo.setMemoryBufferSeverity('Info')
cmo.setFileMinSize(10240)
cd('/Servers/%s/WebServer/%s/WebServerLog/%s')
cmo.setBufferSizeKB(32)
cmo.setNumberOfFilesLimited(true)
cmo.setFileCount(101)
cmo.setFileMinSize(10240)
cd('/Servers/%s/DataSource/%s/DataSourceLogFile/%s')
cmo.setBufferSizeKB(32)
cmo.setNumberOfFilesLimited(true)
cmo.setFileCount(101)
cmo.setFileMinSize(10240)
editActivate()
disconnect()
exit()
`
/**
 * touch conf local file
 */
func (w *Wls12c) touchConf() string {
	os.MkdirAll(utils.TMPD, 0750)
	file, err := ioutil.TempFile(utils.TMPD, "wls12c_server_")//在DIR目录下创建tmp为文件名前缀的文件，获得file文件指针，DIR必须存在，否则创建不成功
	defer file.Close()
	if err != nil {
		logs.Print(fmt.Errorf("create %s error.", file.Name()))
		panic(err)
		return ""
	}

	file.WriteString(
		fmt.Sprintf(
			template,
			w.UserName, w.PassWord,
			w.AdminAddr, w.AdminPort,
			w.SrvName, w.SrvName,
			w.ListenAddr, w.ListenPort, w.ListenAddr,
			w.SrvName, w.SrvName,
			w.SrvName, w.SrvName,
			w.SrvName, w.SrvName, w.SrvName,
			w.SrvName, w.SrvName, w.SrvName,
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
	adminPath := filepath.Join(w.DomainPath, "servers", w.SrvName)
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
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s starting..."
export _SERVER=%s
export CLASSPATH="${CLASSPATH}:%s"
export USER_MEM_ARGS="%s -XX:ErrorFile=${_SERVER}/logs/dump_error.log -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=${_SERVER}/logs/ -Dcatalina.base=%s"
nohup %s/bin/startManagedWebLogic.sh %s t3://%s:%s 1>>/dev/null 2>>%s/logs/failedServer.log &
echo "%s started, pls wating 30 sec..."
`
const templateStopConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s stoping..."
%s/bin/stopManagedWebLogic.sh %s t3://%s:%s
echo "%s stoped, pls wating 30 sec..."
`
const templateRestartConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s restarting..."
%s
%s
echo "%s restarted, pls wating 30 sec..."
`
/**
 * touch console file
 */
func (w *Wls12c) touchConsoleScript() {
	err := utils.MkdirConsolesPath(w.ConsolePath)
	if err != nil {
		logs.Print(err)
	}
	srvpath := filepath.Join(w.DomainPath, "servers", w.SrvName)
	start := filepath.Join(w.ConsolePath, "start", "start_" + w.SrvName + ".sh")
	stop := filepath.Join(w.ConsolePath, "stop", "stop_" + w.SrvName + ".sh")
	restart := filepath.Join(w.ConsolePath, "restart", "restart_" + w.SrvName + ".sh")
	now := time.Now().String()
	var jars string
	for _, jar := range w.Jars {
		jars += ":" + jar
	}
	// start
	logs.Print(ioutil.WriteFile(
		start,
		[]byte(fmt.Sprintf(
			templateStartConsole,
			now,
			w.SrvName, srvpath, jars,
			w.Margs, srvpath, w.DomainPath, w.SrvName,
			w.AdminAddr, w.AdminPort, srvpath, w.SrvName)),
		0750,
	))


	// stop
	logs.Print(ioutil.WriteFile(
		stop,
		[]byte(fmt.Sprintf(templateStopConsole,
			now,
			w.SrvName, w.DomainPath, w.SrvName, w.AdminAddr, w.AdminPort, w.SrvName)),
		0750,
	))

	// restart
	logs.Print(ioutil.WriteFile(
		restart,
		[]byte(fmt.Sprintf(templateRestartConsole,
			now,
			w.SrvName, stop, start, w.SrvName)),
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
	start := filepath.Join(w.ConsolePath, "start", "start_" + w.SrvName + ".sh")
	cmd := exec.Command(start)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	logs.Print(cmd.Run())
}
/* End Run */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

