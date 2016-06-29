package wls12c

import (
	"asiainfo.com/ins/logs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Wls12c struct {
	WLS_HOME    string `json:"wlsHome"`
	AdminAddr   string `json:"adminAddr"`
	AdminPort   string `json:"adminPort"`
	SrvName     string `json:"srvName"`
	UserName    string `json:"userName"`
	PassWord    string `json:"passWord"`
	PkgName     string `json:"pkgName"`
	PkgPath     string `json:"pkgPath"`
	Option      string `json:"option"`
	ConsolePath string `json:"consolePath"`
}

func (w *Wls12c) Json(bs []byte) error {
	return json.Unmarshal(bs, &w)
}

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Method */

func (w *Wls12c) Add() error {
	wlst := w.WLS_HOME + "/wlserver/common/bin/wlst.sh"
	dpy, err := w.touchDeploy(wlst)
	if err != nil {
		return err
	}
	//defer logs.Print(os.Remove(conf))

	// 根据模板生成domain
	cmd := exec.Command(dpy)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

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

const deployTemplate = `
def editMode():
  edit()
  startEdit()
def editActivate():
  save()
  activate()
connect("%s", "%s", 't3://%s:%s')
progress=deploy(appName='%s',path='%s',targets='%s', stageMode='nostage', timeout=120000)
progress.printStatus()
disconnect()
exit()
`

const undeployTemplate = `
def editMode():
  edit()
  startEdit()
def editActivate():
  save()
  activate()
connect("%s", "%s", 't3://%s:%s')
undeploy(appName='%s',path='%s',targets='%s', forceUndeployTimeout=120000)
disconnect()
exit()
`

const redeployTemplate = `
def editMode():
  edit()
  startEdit()
def editActivate():
  save()
  activate()
connect("%s", "%s", 't3://%s:%s')
progress=redeploy(appName='%s',appPath='%s', timeout=120000)
progress.getState()
disconnect()
exit()
`

const startApplicationsTemplate = `
def editMode():
  edit()
  startEdit()
def editActivate():
  save()
  activate()
connect("%s", "%s", 't3://%s:%s')
progress=startApplication(appName='%s',stageMode='NOSTAGE', adminMode='false', timeout=120000)
progress.getState()
disconnect()
exit()
`

const stopApplicationsTemplate = `
def editMode():
  edit()
  startEdit()
def editActivate():
  save()
  activate()
connect("%s", "%s", 't3://%s:%s')
progress=stopApplication(appName='%s', timeout=120000)
progress.getState()
disconnect()
exit()
`

const updateApplicationsTemplate = `
def editMode():
  edit()
  startEdit()
def editActivate():
  save()
  activate()
connect("%s", "%s", 't3://%s:%s')
progress=updateApplication(appName='%s',adminMode='false', timeout=120000)
progress.getState()
disconnect()
exit()
`

/**
 * touch conf and run file
 */
func (w *Wls12c) touchDeploy(wlst string) (string, error) {
	deployConfDir := filepath.Join(w.ConsolePath, "deploy", "conf")
	undeployConfDir := filepath.Join(w.ConsolePath, "undeploy", "conf")
	redeployConfDir := filepath.Join(w.ConsolePath, "redeploy", "conf")
	startApplicationsConfDir := filepath.Join(w.ConsolePath, "application_start", "conf")
	stopApplicationsConfDir := filepath.Join(w.ConsolePath, "application_stop", "conf")
	//updateApplicationsConfDir := filepath.Join(w.ConsolePath, "application_update", "conf")

	os.MkdirAll(deployConfDir, 0750)
	os.MkdirAll(undeployConfDir, 0750)
	os.MkdirAll(redeployConfDir, 0750)
	os.MkdirAll(startApplicationsConfDir, 0750)
	os.MkdirAll(stopApplicationsConfDir, 0750)
	//os.MkdirAll(updateApplicationsConfDir, 0750)

	dname := filepath.Join(deployConfDir, w.PkgName+".py")
	uname := filepath.Join(undeployConfDir, w.PkgName+".py")
	rname := filepath.Join(redeployConfDir, w.PkgName+".py")
	startAname := filepath.Join(startApplicationsConfDir, w.PkgName+".py")
	stopAname := filepath.Join(stopApplicationsConfDir, w.PkgName+".py")
	//updateAname := filepath.Join(updateApplicationsConfDir, w.PkgName + ".py")
	err := ioutil.WriteFile(
		dname,
		[]byte(fmt.Sprintf(
			deployTemplate,
			w.UserName, w.PassWord,
			w.AdminAddr, w.AdminPort,
			w.PkgName, w.PkgPath, w.SrvName,
		)),
		0750,
	)
	if err != nil {
		rm(dname, uname, rname, startAname, stopAname)
		return "", err
	}
	err = ioutil.WriteFile(
		uname,
		[]byte(fmt.Sprintf(
			undeployTemplate,
			w.UserName, w.PassWord,
			w.AdminAddr, w.AdminPort,
			w.PkgName, w.PkgPath, w.SrvName,
		)),
		0750,
	)
	if err != nil {
		rm(dname, uname, rname, startAname, stopAname)
		return "", err
	}
	err = ioutil.WriteFile(
		rname,
		[]byte(fmt.Sprintf(
			redeployTemplate,
			w.UserName, w.PassWord,
			w.AdminAddr, w.AdminPort,
			w.PkgName, w.PkgPath,
		)),
		0750,
	)
	if err != nil {
		rm(dname, uname, rname, startAname, stopAname)
		return "", err
	}
	err = ioutil.WriteFile(
		startAname,
		[]byte(fmt.Sprintf(
			startApplicationsTemplate,
			w.UserName, w.PassWord,
			w.AdminAddr, w.AdminPort,
			w.SrvName,
		)),
		0750,
	)
	if err != nil {
		rm(dname, uname, rname, startAname, stopAname)
		return "", err
	}
	err = ioutil.WriteFile(
		stopAname,
		[]byte(fmt.Sprintf(
			stopApplicationsTemplate,
			w.UserName, w.PassWord,
			w.AdminAddr, w.AdminPort,
			w.SrvName,
		)),
		0750,
	)
	if err != nil {
		rm(dname, uname, rname, startAname, stopAname)
		return "", err
	}
	//err = ioutil.WriteFile(
	//	updateAname,
	//	[]byte(fmt.Sprintf(
	//		updateApplicationsTemplate,
	//		w.UserName, w.PassWord,
	//		w.AdminAddr, w.AdminPort,
	//		w.SrvName,
	//	)),
	//	0750,
	//)
	//if err != nil {
	//	rm(dname, uname, rname, startAname, stopAname)
	//	return "", err
	//}
	_dname, _, _, _, _ := w.touchConsoleScript(wlst, dname, uname, rname, startAname, stopAname)
	return _dname, nil
}

func rm(names ...string) {
	for _, name := range names {
		os.Remove(name)
	}
}

/* End Configuration File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Console File */

const templateDeployConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s deploy..."
%s %s
echo "%s deployed, pls wating 30 sec..."
`
const templateUndeployConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s undeploy..."
%s %s
echo "%s undeployed, pls wating 30 sec..."
`
const templateRedeployConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s redeploy..."
%s %s
echo "%s redeployed, pls wating 30 sec..."
`

const templateStartApplicationsConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s start applications..."
%s %s
echo "%s started applications, pls wating 30 sec..."
`
const templateStopApplicationsConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s stop applications..."
%s %s
echo "%s stoped applications, pls wating 30 sec..."
`
const templateUpdateApplicationsConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s update applications..."
%s %s
echo "%s updated applications, pls wating 30 sec..."
`

/**
 * touch console file
 */
func (w *Wls12c) touchConsoleScript(wlst, dname, uname, rname, startAname, stopAname string) (string, string, string, string, string) {
	_dname := filepath.Join(filepath.Join(w.ConsolePath, "deploy"), "deploy_"+w.PkgName+".sh")
	_uname := filepath.Join(filepath.Join(w.ConsolePath, "undeploy"), "undeploy_"+w.PkgName+".sh")
	_rname := filepath.Join(filepath.Join(w.ConsolePath, "redeploy"), "redeploy_"+w.PkgName+".sh")
	_startAname := filepath.Join(filepath.Join(w.ConsolePath, "application_start"), w.PkgName+"_start.sh")
	_stopAname := filepath.Join(filepath.Join(w.ConsolePath, "application_stop"), w.PkgName+"_stop.sh")
	//_updateAname := filepath.Join(filepath.Join(w.ConsolePath, "application_update"), w.PkgName + "_update.sh")
	now := time.Now().String()
	logs.Print(ioutil.WriteFile(
		_dname,
		[]byte(fmt.Sprintf(templateDeployConsole,
			now,
			w.PkgName, wlst, dname, w.PkgName)),
		0750,
	))
	logs.Print(ioutil.WriteFile(
		_uname,
		[]byte(fmt.Sprintf(templateUndeployConsole,
			now,
			w.PkgName, wlst, uname, w.PkgName)),
		0750,
	))
	logs.Print(ioutil.WriteFile(
		_rname,
		[]byte(fmt.Sprintf(templateRedeployConsole,
			now,
			w.PkgName, wlst, rname, w.PkgName)),
		0750,
	))
	logs.Print(ioutil.WriteFile(
		_startAname,
		[]byte(fmt.Sprintf(templateStartApplicationsConsole,
			now,
			w.PkgName, wlst, startAname, w.PkgName)),
		0750,
	))
	logs.Print(ioutil.WriteFile(
		_stopAname,
		[]byte(fmt.Sprintf(templateStopApplicationsConsole,
			now,
			w.PkgName, wlst, stopAname, w.PkgName)),
		0750,
	))
	//logs.Print(ioutil.WriteFile(
	//	_updateAname,
	//	[]byte(fmt.Sprintf(templateUpdateApplicationsConsole, w.PkgName, wlst, updateAname, w.PkgName)),
	//	0750,
	//))
	return _dname, _uname, _rname, _startAname, _stopAname
}

/* End Console File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
