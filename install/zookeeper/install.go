package zookeeper

import (
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type ZK struct {
	ZK_Name         string    `json:"zk_name"`
	ZK_PKG          string    `json:"zk_pkg"`
	ZK_HOME         string    `json:"zk_home"`
	ID              int       `json:"id"`
	DataDir         string    `json:"dataDir"`
	ClientPort      string    `json:"clientPort"`
	MaxClientCnxns  string    `json:"maxClientCnxns"`
	Clusters        []Cluster `json:"clusters"`
	SnapRetainCount string    `json:"snapRetainCount"`
	PurgeInterval   string    `json:"purgeInterval"`
	ConsolePath     string    `json:"consolePath"`
	JVM             string    `json:"jvm"`
}

type Cluster struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (zk *ZK) Json(bs []byte) error {
	return json.Unmarshal(bs, &zk)
}

const installZK = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
mkdir -p %s
tar xzf %s -C %s --strip-components=1
cd %s
mkdir -p dump logs %s
echo %d > %s
`

func (zk *ZK) Install() error {
	_, err := os.Stat(zk.ZK_PKG)
	if err != nil {
		return err
	}
	now := time.Now().String()
	zkSh := filepath.Join(utils.TMPD, "zkInstall.sh")
	defer os.Remove(zkSh)
	logs.Print(ioutil.WriteFile(
		zkSh,
		[]byte(fmt.Sprintf(installZK,
			now,
			zk.ZK_HOME,
			zk.ZK_PKG, zk.ZK_HOME,
			zk.ZK_HOME,
			zk.DataDir,
			zk.ID, filepath.Join(zk.DataDir, "myid"),
		)),
		0750,
	))
	// 根据模板生成domain
	cmd := exec.Command(zkSh)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	zk.touchConf()
	zk.touchConsoleScript()

	return nil

}

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Configuration File */

const zkConf = `########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
# The number of milliseconds of each tick
tickTime=2000
# The number of ticks that the initial
# synchronization phase can take
initLimit=10
# The number of ticks that can pass between
# sending a request and getting an acknowledgement
syncLimit=5
# the directory where the snapshot is stored.
# do not use /tmp for storage, /tmp here is just
# example sakes.
dataDir=%s
# the port at which the clients will connect
clientPort=%d

%s

# the maximum number of client connections.
# increase this if you need to handle more clients
maxClientCnxns=%d
# The number of snapshots to retain in dataDir
autopurge.snapRetainCount=%d
# Purge task interval in hours
# Set to "0" to disable auto purge feature
autopurge.purgeInterval=%d
`

/**
 * touch conf local file
 */
func (zk *ZK) touchConf() {
	zk_conf := filepath.Join(zk.ZK_HOME, "conf", "zoo.cfg")
	now := time.Now().String()

	clientPort := 2181
	maxClientCnxns := 60
	snapRetainCount := 3
	purgeInterval := 1
	if len(zk.ClientPort) > 0 {
		clientPort, _ = strconv.Atoi(zk.ClientPort)
	}
	if len(zk.MaxClientCnxns) > 0 {
		maxClientCnxns, _ = strconv.Atoi(zk.MaxClientCnxns)
	}
	if len(zk.SnapRetainCount) > 0 {
		snapRetainCount, _ = strconv.Atoi(zk.SnapRetainCount)
	}
	if len(zk.PurgeInterval) > 0 {
		purgeInterval, _ = strconv.Atoi(zk.PurgeInterval)
	}
	var clusters string
	for _, server := range zk.Clusters {
		clusters += server.Name + "=" + server.Address + "\n"
	}

	// zk_conf
	logs.Print(ioutil.WriteFile(
		zk_conf,
		[]byte(fmt.Sprintf(zkConf,
			now,
			zk.DataDir, clientPort,
			clusters,
			maxClientCnxns, snapRetainCount, purgeInterval,
		)),
		0750,
	))

}

/* End Configuration File */
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
export ZK_HOME=%s
if [ ! -f %s ]; then
  echo %d > %s
fi
export JVM="%s"
export ZOO_LOG4J_PROP=INFO,CONSOLE
export ZOO_LOG_DIR=%s
export JVMFLAGS="${JVM} -XX:ErrorFile=${ZK_HOME}/dump/err_zk.log -Xloggc:${ZK_HOME}/dump/gc_zk.log -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=${ZK_HOME}/dump/ -Dfile.encoding=UTF-8 -Dencoding=UTF-8"
cd ${ZOO_LOG_DIR}
%s start
sleep 5s
echo "%s started, pls wating 30 sec..."
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

const templateStopConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s stopping..."
export ZOO_LOG_DIR=%s
cd ${ZOO_LOG_DIR}
%s stop
sleep 5s
echo "%s stopped, pls wating 30 sec..."
`

/**
 * touch console file
 */
func (zk *ZK) touchConsoleScript() {
	err := utils.MkdirConsolesPath(zk.ConsolePath)
	if err != nil {
		logs.Print(err)
	}
	start := filepath.Join(zk.ConsolePath, "start", "start_"+zk.ZK_Name+".sh")
	stop := filepath.Join(zk.ConsolePath, "stop", "stop_"+zk.ZK_Name+".sh")
	restart := filepath.Join(zk.ConsolePath, "restart", "restart_"+zk.ZK_Name+".sh")
	myid := filepath.Join(zk.DataDir, "myid")
	zkLogs := filepath.Join(zk.ZK_HOME, "logs")
	zkConsole := filepath.Join(zk.ZK_HOME, "bin", "zkServer.sh")
	jvm := "-Xmx2048M -Xms2048M -XX:MaxPermSize=128m -XX:PermSize=32m"
	if len(zk.JVM) > 0 {
		jvm = zk.JVM
	}
	now := time.Now().String()
	// start
	logs.Print(ioutil.WriteFile(
		start,
		[]byte(fmt.Sprintf(templateStartConsole,
			now,
			zk.ZK_Name,
			zk.ZK_HOME, myid, zk.ID, myid,
			jvm, zkLogs, zkConsole,
			zk.ZK_Name,
		)),
		0750,
	))

	// stop
	logs.Print(ioutil.WriteFile(
		stop,
		[]byte(fmt.Sprintf(templateStopConsole,
			now,
			zk.ZK_Name,
			zkLogs,
			zkConsole,
			zk.ZK_Name,
		)),
		0750,
	))

	// restart
	logs.Print(ioutil.WriteFile(
		restart,
		[]byte(fmt.Sprintf(templateRestartConsole,
			now,
			zk.ZK_Name, stop, start, zk.ZK_Name)),
		0750,
	))
}

/* End Console File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
