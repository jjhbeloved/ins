package hdfs

import (
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Hdfs struct {
	Hdfs_Name   string   `json:"hdfs_name"`
	Hdfs_PKG    string   `json:"hdfs_pkg"`
	Hdfs_HOME   string   `json:"hdfs_home"`
	JAVA_HOME   string   `json:"java_home"`
	ZKs         []string `json:"zks"`
	ZKPort      int      `json:"zkPort"`
	ZKDir       string   `json:"zkDir"`
	Option      string   `json:"option"`
	ConsolePath string   `json:"consolePath"`
}

func (hdfs *Hdfs) Json(bs []byte) error {
	return json.Unmarshal(bs, &hdfs)
}

const installStorm = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
mkdir -p %s
tar xzf %s -C %s --strip-components=1
cd %s
%s
rm -rf ../examples
mkdir -p %s %s
`

func (storm *Storm) Install() error {
	_, err := os.Stat(storm.Storm_PKG)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	now := time.Now().String()
	stormSh := filepath.Join(utils.TMPD, "stormCreate.sh")
	defer os.Remove(stormSh)
	var rms string
	for _, name := range storm.RMlibs {
		rms += fmt.Sprintf("mv %s _%s\n", name, name)
	}
	logs.Print(ioutil.WriteFile(
		stormSh,
		[]byte(fmt.Sprintf(installStorm,
			now,
			storm.Storm_HOME,
			storm.Storm_PKG, storm.Storm_HOME,
			filepath.Join(storm.Storm_HOME, "lib"),
			rms,
			storm.StormDataDir, storm.StormLogDir,
		)),
		0750,
	))

	// 根据模板生成domain
	cmd = exec.Command("sh", stormSh)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	storm.touchConf()
	storm.touchConsoleScript()

	return err
}

const stormYaml = `
storm.zookeeper.servers:
%s
storm.zookeeper.port: %d
storm.zookeeper.root: "%s"

storm.local.dir: "%s"
storm.log.dir:  "%s"

nimbus.seeds: [%s]

supervisor.slots.ports:
%s

storm.cluster.mode: "%s"

worker.heap.memory.mb: %d
worker.childopts: "%s"

ui.host:  "%s"
ui.port:  %d
ui.childopts: "%s"

topology.message.timeout.secs: %d

storm.health.check.dir: "healthchecks"
storm.health.check.timeout.ms: 5000
`

/**
 * touch conf local file
 */
func (storm *Storm) touchConf() {
	_stormYaml := filepath.Join(storm.Storm_HOME, "conf", "storm.yaml")
	if len(storm.StormYaml) > 0 {
		utils.CopyFile(_stormYaml, storm.StormYaml)
	} else {
		var zks, seeds, slots string
		for _, zk := range storm.ZKs {
			zks += fmt.Sprintf(" - \"%s\"\n", zk)
		}
		for i, seed := range storm.NimbusHA {
			if i == 0 {
				seeds += fmt.Sprintf(`"%s"`, seed)
			} else {
				seeds += fmt.Sprintf(`, "%s"`, seed)
			}
		}
		for _, slot := range storm.SlotsPorts {
			slots += fmt.Sprintf(" - %d\n", slot)
		}
		logs.Print(ioutil.WriteFile(
			_stormYaml,
			[]byte(fmt.Sprintf(stormYaml,
				zks, storm.ZKPort, storm.ZKDir,
				storm.StormDataDir, storm.StormLogDir,
				seeds, slots,
				storm.Mode,
				storm.WorkerHeap, storm.WorkerJVM,
				storm.UIHost, storm.UIPort, storm.UIJVM,
				storm.TopoMSGTime,
			)),
			0750,
		))
	}
}

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Console File */
const templateStartSupervisorConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "starting storm supervisor..."
%s supervisor >/dev/null 2>&1 &
echo " started storm supervisor..."
`

const templateRestartSupervisorConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s restarting..."
%s
%s
echo "%s restarted, pls wating 30 sec..."
`

const templateStopSupervisorConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "stopping storm supervisor"
ps -ef | grep storm| grep -v grep | grep supervisor | grep %s |awk '{print $2}' | xargs kill -15
echo "stopped storm supervisor"
`

const templateStartNimbusConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "starting storm nimbus..."
%s nimbus >/dev/null 2>&1 &
%s ui >/dev/null 2>&1 &
echo " started storm nimbus..."
`

const templateRestartNimbusConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "restarting..."
%s
%s
echo "restarted, pls wating 30 sec..."
`

const templateStopNimbusConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "stopping storm nimbus"
ps -ef | grep storm| grep -v grep | grep ui | awk '{print $2}' | xargs kill -15
ps -ef | grep storm| grep -v grep | grep nimbus | awk '{print $2}' | xargs kill -15
echo "stopped storm nimbus"
`

/**
 * touch console file
 */
func (storm *Storm) touchConsoleScript() {
	err := utils.MkdirConsolesPath(storm.ConsolePath)
	if err != nil {
		logs.Print(err)
	}
	startSupervisor := filepath.Join(storm.ConsolePath, "start", "start_supervisor_"+storm.Storm_Name+".sh")
	stopSupervisor := filepath.Join(storm.ConsolePath, "stop", "stop_supervisor_"+storm.Storm_Name+".sh")
	restartSupervisor := filepath.Join(storm.ConsolePath, "restart", "restart_supervisor_"+storm.Storm_Name+".sh")
	startNimbus := filepath.Join(storm.ConsolePath, "start", "start_nimbus_"+storm.Storm_Name+".sh")
	stopNimbus := filepath.Join(storm.ConsolePath, "stop", "stop_nimbus_"+storm.Storm_Name+".sh")
	restartNimbus := filepath.Join(storm.ConsolePath, "restart", "restart_nimbus_"+storm.Storm_Name+".sh")
	now := time.Now().String()
	// start
	logs.Print(ioutil.WriteFile(
		startSupervisor,
		[]byte(fmt.Sprintf(templateStartSupervisorConsole,
			now,
			filepath.Join(storm.Storm_HOME, "bin", "storm"),
		)),
		0750,
	))
	// stop
	logs.Print(ioutil.WriteFile(
		stopSupervisor,
		[]byte(fmt.Sprintf(templateStopSupervisorConsole,
			now,
		)),
		0750,
	))
	// restart
	logs.Print(ioutil.WriteFile(
		restartSupervisor,
		[]byte(fmt.Sprintf(templateRestartSupervisorConsole,
			now,
			stopSupervisor, startSupervisor,
		)),
		0750,
	))
	// start
	logs.Print(ioutil.WriteFile(
		startNimbus,
		[]byte(fmt.Sprintf(templateStartNimbusConsole,
			now,
			filepath.Join(storm.Storm_HOME, "bin", "storm"),
			filepath.Join(storm.Storm_HOME, "bin", "storm"),
		)),
		0750,
	))
	// stop
	logs.Print(ioutil.WriteFile(
		stopNimbus,
		[]byte(fmt.Sprintf(templateStopNimbusConsole,
			now,
		)),
		0750,
	))
	// restart
	logs.Print(ioutil.WriteFile(
		restartNimbus,
		[]byte(fmt.Sprintf(templateRestartNimbusConsole,
			now,
			stopNimbus, startNimbus,
		)),
		0750,
	))
}

/* End Console File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
