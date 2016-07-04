package jdk

import (
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"asiainfo.com/ins/cli"
)

const (
	TAR_GZ string = "tar.gz"
	ZIP    string = "zip"
)

type Jdk struct {
	Pkg      string `json:"pkg"`
	JDK_HOME string `json:"jdk_home"`
	IsROOT   bool   `json:"isRoot"`
}

func (w *Jdk) Json(bs []byte) error {
	return json.Unmarshal(bs, &w)
}

func (jdk *Jdk) Install() error {
	pkg, err := utils.DownloadToDir(jdk.Pkg, cli.PKG_PATH)
	if err != nil {
		return err
	}
	jdk.Pkg = pkg

	fileinfo, err := os.Stat(jdk.Pkg)
	if err != nil {
		return err
	}

	suffix := filepath.Base(fileinfo.Name())
	var cmd *exec.Cmd
	uncompress := filepath.Join(utils.TMPD, "uncompress.sh")
	defer os.Remove(uncompress)
	if strings.HasSuffix(suffix, TAR_GZ) {
		//name, err = utils.UnTarGz(tomcat.Pkg, dir, tomcat.IsRemove, false)
		// copy
		logs.Print(ioutil.WriteFile(
			uncompress,
			[]byte(fmt.Sprintf(tarJDKPkg,
				jdk.JDK_HOME,
				jdk.Pkg, jdk.JDK_HOME,
			)),
			0750,
		))
	} else if strings.HasSuffix(suffix, ZIP) {
		//name, err = utils.Unzip(tomcat.Pkg, dir, tomcat.IsRemove)
		// copy
		logs.Print(ioutil.WriteFile(
			uncompress,
			[]byte(fmt.Sprintf(unzipJDKPkg,
				jdk.JDK_HOME,
				jdk.Pkg, jdk.JDK_HOME,
			)),
			0750,
		))
	} else {
		return fmt.Errorf("%s is not exists.", jdk.Pkg)
	}
	if err != nil {
		return err
	}
	// 根据模板生成domain
	cmd = exec.Command("sh", uncompress)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	if jdk.IsROOT && os.Getuid() == 0 {
		bashrc := filepath.Join("/etc/bashrc")
		utils.WriteFileA(
			bashrc,
			[]byte(fmt.Sprintf(javaEnv,
				jdk.JDK_HOME,
			)),
			0644,
		)
	} else {
		bashrc := filepath.Join(os.Getenv("HOME"), ".bashrc")
		utils.WriteFileA(
			bashrc,
			[]byte(fmt.Sprintf(javaEnv,
				time.Now().String(),
				jdk.JDK_HOME,
			)),
			0644,
		)
	}
	return err
}

const tarJDKPkg = `#!/bin/bash
mkdir -p %s
tar xzf %s -C %s --strip-components=1
`

const unzipJDKPkg = `#!/bin/bash
mkdir -p %s
unzip %s -d %s
`
const javaEnv = `
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
JAVA_HOME=%s
PATH=${JAVA_HOME}/bin:${PATH}
export PATH JAVA_HOME
`
