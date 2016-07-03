package mysql

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

type Mysql struct {
	Name            string    `json:"name"`
	MYSQL_PKG       string    `json:"mysql_pkg"`
	MYSQL_HOME      string    `json:"mysql_home"`
}

func (mysql *Mysql) Json(bs []byte) error {
	return json.Unmarshal(bs, &mysql)
}

const installMysql = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
mkdir -p %s-src
tar xzf %s -C %s-src --strip-components=1
cd %s-src
cmake  . -DCMAKE_INSTALL_PREFIX=%s
make -j4
make install
rm -rf %s-src
mkdir -p %s/sbin
cp %s/bin/mysqld %s/sbin
`

func (mysql *Mysql) Install() error {
	_, err := os.Stat(mysql.MYSQL_PKG)
	if err != nil {
		return err
	}
	now := time.Now().String()
	mySh := filepath.Join(utils.TMPD, "mysqlInstall.sh")
	defer os.Remove(mySh)
	logs.Print(ioutil.WriteFile(
		mySh,
		[]byte(fmt.Sprintf(installMysql,
			now,
			mysql.MYSQL_HOME,
			mysql.MYSQL_PKG, mysql.MYSQL_HOME,
			mysql.MYSQL_HOME,
			mysql.MYSQL_HOME,
			mysql.MYSQL_HOME,
			mysql.MYSQL_HOME,
			mysql.MYSQL_HOME, mysql.MYSQL_HOME,
		)),
		0750,
	))
	// 根据模板生成domain
	cmd := exec.Command(mySh)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil

}
