package memcached

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

const (
	TAR_GZ string = "tar.gz"
	ZIP    string = "zip"
)

type Memcached struct {
	Memcached_PKG     string `json:"memcached_pkg"`
	Libevent_PKG      string `json:"libevent_pkg"`
	Repcached_PATCH   string `json:"repcached_patch"`
	Memcached_HOME    string `json:"memcached_home"`
	Libevent_HOME     string `json:"libevent_home"`
	Port              string `json:"port"`
	Memory            string `json:"memory"`
	Connections       string `json:"connections"`
	RepecachedPort    string `json:"repecachedPort"`
	RepecachedAddress string `json:"repecachedAddress"`
	Option            string `json:"option"`
	ConsolePath       string `json:"consolePath"`
}

func (w *Memcached) Json(bs []byte) error {
	return json.Unmarshal(bs, &w)
}

func (mem *Memcached) Install() error {
	_, err := os.Stat(mem.Libevent_PKG)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	now := time.Now().String()
	memcached := filepath.Join(utils.TMPD, "memcacheCreate.sh")
	defer os.Remove(memcached)
	if len(mem.Repcached_PATCH) > 0 {
		logs.Print(ioutil.WriteFile(
			memcached,
			[]byte(fmt.Sprintf(installMemcachedWithRepcached,
				now,
				mem.Libevent_HOME,
				mem.Libevent_PKG, mem.Libevent_HOME,
				mem.Libevent_HOME,
				mem.Libevent_HOME,
				mem.Libevent_HOME,
				mem.Memcached_HOME,
				mem.Memcached_PKG, mem.Memcached_HOME,
				mem.Memcached_HOME,
				mem.Repcached_PATCH,
				mem.Memcached_HOME, mem.Libevent_HOME,
				mem.Memcached_HOME,
			)),
			0750,
		))
	} else {
		logs.Print(ioutil.WriteFile(
			memcached,
			[]byte(fmt.Sprintf(installMemcached,
				now,
				mem.Libevent_HOME,
				mem.Libevent_PKG, mem.Libevent_HOME,
				mem.Libevent_HOME,
				mem.Libevent_HOME,
				mem.Libevent_HOME,
				mem.Memcached_HOME,
				mem.Memcached_PKG, mem.Memcached_HOME,
				mem.Memcached_HOME,
				mem.Memcached_HOME, mem.Libevent_HOME,
				mem.Memcached_HOME,
			)),
			0750,
		))
	}
	defer os.Remove(memcached)

	// 根据模板生成domain
	cmd = exec.Command("sh", memcached)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	mem.touchConsoleScript()
	return err
}

const installMemcached = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
sh configure -prefix=%s
make -j4 && make install
rm -rf %s-SRC

mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
sh configure -prefix=%s -with-libevent=%s
make -j4 && make install
rm -rf %s-SRC
`

const installMemcachedWithRepcached = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
sh configure -prefix=%s
make -j4 && make install
rm -rf %s-SRC

mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
patch -p1 -i %s
sh configure -prefix=%s -with-libevent=%s --enable-replication
make -j4 && make install
rm -rf %s-SRC
`

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Console File */
const templateStartConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "starting memcached...."
%s/bin/memcached -d -m %d -c %d -p %d -U 0 -P %s/memcached.pid
echo "memcache started."
`
const templateStartConsoleWithRepcached = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "starting memcached...."
%s/bin/memcached -d -m %d -c %d -p %d -U 0 -P %s/memcached.pid -L -x %s -X %s
echo "memcache started."
`
const templateStopConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "memcached stopping..."
killall memcached
echo "memcached stopped, pls wating 30 sec..."
`
const templateRestartConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "memcached restarting..."
%s
%s
echo "memcached restarted, pls wating 30 sec..."
`

/**
 * touch console file
 */
func (mem *Memcached) touchConsoleScript() {
	err := utils.MkdirConsolesPath(mem.ConsolePath)
	if err != nil {
		logs.Print(err)
	}
	start := filepath.Join(mem.ConsolePath, "start", "start_memcached.sh")
	stop := filepath.Join(mem.ConsolePath, "stop", "stop_memcached.sh")
	restart := filepath.Join(mem.ConsolePath, "restart", "restart_memcached.sh")
	now := time.Now().String()
	vm := 512
	if len(mem.Memory) > 0 {
		vm, _ = strconv.Atoi(mem.Memory)
	}
	port := 11211
	if len(mem.Port) > 0 {
		port, _ = strconv.Atoi(mem.Port)
	}
	conns := 2048
	if len(mem.Port) > 0 {
		conns, _ = strconv.Atoi(mem.Connections)
	}
	if len(mem.RepecachedAddress) > 0 && len(mem.RepecachedPort) > 0 {
		// start
		logs.Print(ioutil.WriteFile(
			start,
			[]byte(fmt.Sprintf(
				templateStartConsoleWithRepcached,
				now,
				mem.Memcached_HOME, vm, conns, port, mem.Memcached_HOME, mem.RepecachedAddress, mem.RepecachedPort,
			)),
			0750,
		))
	} else {
		// start
		logs.Print(ioutil.WriteFile(
			start,
			[]byte(fmt.Sprintf(
				templateStartConsole,
				now,
				mem.Memcached_HOME, vm, conns, port, mem.Memcached_HOME,
			)),
			0750,
		))
	}

	// stop
	logs.Print(ioutil.WriteFile(
		stop,
		[]byte(fmt.Sprintf(templateStopConsole,
			now,
		)),
		0750,
	))

	// restart
	logs.Print(ioutil.WriteFile(
		restart,
		[]byte(fmt.Sprintf(templateRestartConsole,
			now,
			stop, start,
		)),
		0750,
	))
}

/* End Console File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
