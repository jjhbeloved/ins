package main
import (
	"asiainfo.com/ins/cli"
	"asiainfo.com/ins/utils"
	"fmt"
	"asiainfo.com/ins/domain/wls"
	"io/ioutil"
	"asiainfo.com/ins/logs"
	"path/filepath"
	"asiainfo.com/ins/domain/tomcat"
	"asiainfo.com/ins/domain/redis"
	"strings"
)

/**
* go build -o domain /veris/odc/app/go/3rd/src/asiainfo.com/ins/domain/domain.go
*/
func main()  {
	chs1()
}

func chs1() {
	ins_path := filepath.Join(cli.CONF_HOME, "domain")
	fs, err := utils.GetAllFiles(ins_path)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for _, f := range fs {
		key := f.Name()
		fn := filepath.Join(ins_path, key)
		var dom Domain
		var option string
		if strings.HasPrefix(key, cli.WLS12CPREFIX) {
			var wls12 wls12c.Wls12c
			bs, _ := ioutil.ReadFile(fn)
			e := wls12.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			dom = &wls12
			option = wls12.Option
		} else if strings.HasPrefix(key, cli.TOMCATPREFIX) {
			var tomcat tomcat.Tomcat
			bs, _ := ioutil.ReadFile(fn)
			e := tomcat.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			dom = &tomcat
			option = tomcat.Option
		} else if strings.HasPrefix(key, cli.REDISPREFIX) {
			var redis redis.Redis
			bs, _ := ioutil.ReadFile(fn)
			e := redis.Json(bs)
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
			dom = &redis
			option = redis.Option
		} else {
			dom = nil
			continue
		}
		if dom != nil {
			var e error
			if option == "REMOVE" {
				e = dom.Remove()
			} else {
				e = dom.Add()
			}
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
		}
	}
}

type Domain interface {
	Add() error
	Remove() error
}
