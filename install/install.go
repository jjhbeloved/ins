package main

import (
	"asiainfo.com/ins/cli"
	"asiainfo.com/ins/install/amq"
	"asiainfo.com/ins/install/jdk"
	"asiainfo.com/ins/install/memcached"
	"asiainfo.com/ins/install/redis"
	"asiainfo.com/ins/install/storm"
	"asiainfo.com/ins/install/tomcat"
	"asiainfo.com/ins/install/wls"
	"asiainfo.com/ins/install/zookeeper"
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/utils"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"asiainfo.com/ins/install/mysql"
)

func init() {
	//cli.InitInstall()
}

/**
* go build -o install /veris/odc/app/go/3rd/src/asiainfo.com/ins/install/install.go
* go build -o install /install_apps/server/go/3rd/src/asiainfo.com/ins/install/install.go
 */
func main() {
	chs2()
}

/**
 * 遍历安装配置文档, 遍历安装
 */
func chs1() {
	for key, val := range cli.Install {
		switch key {
		case cli.WLS12C:
			wls12 := &wls.Wls12c{}
			err := wls12.Builder(val)
			print(err)
			print(wls12.Install())
		default:
			continue
		}
	}
}

/**
 * 遍历配置目录, 遍历安装符合条件的文件
 */
func chs2() {
	os.Mkdir(utils.TMPD, 0777)
	ins_path := filepath.Join(cli.CONF_HOME, "install")
	fs, err := utils.GetAllFiles(ins_path)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for _, f := range fs {
		key := f.Name()
		fn := filepath.Join(ins_path, key)
		var ins Installer
		switch key {
		case cli.WLS12CCONF:
			var wls12 wls.Wls12c
			bs, _ := ioutil.ReadFile(fn)
			e := wls12.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			ins = &wls12
		case cli.TOMCATCONF:
			var tomcat tomcat.Tomcat
			bs, _ := ioutil.ReadFile(fn)
			e := tomcat.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			ins = &tomcat
		case cli.JDKCONF:
			var jdk jdk.Jdk
			bs, _ := ioutil.ReadFile(fn)
			e := jdk.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			ins = &jdk
		case cli.MEMCACHEDCONF:
			var mem memcached.Memcached
			bs, _ := ioutil.ReadFile(fn)
			e := mem.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			ins = &mem
		case cli.ACTIVEMQCONF:
			var amq amq.AMQ
			bs, _ := ioutil.ReadFile(fn)
			e := amq.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			ins = &amq
		case cli.ZKCONF:
			var zk zookeeper.ZK
			bs, _ := ioutil.ReadFile(fn)
			e := zk.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			ins = &zk
		case cli.REDISCONF:
			var redis redis.Redis
			bs, _ := ioutil.ReadFile(fn)
			e := redis.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			ins = &redis
		case cli.STORMCONF:
			var storm storm.Storm
			bs, _ := ioutil.ReadFile(fn)
			e := storm.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			ins = &storm
		case cli.MYSQLCONF:
			var mysql mysql.Mysql
			bs, _ := ioutil.ReadFile(fn)
			e := mysql.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			ins = &mysql
		default:
			ins = nil
			continue
		}
		if ins != nil {
			e := ins.Install()
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
		}
	}
}

type Installer interface {
	Install() error
}
