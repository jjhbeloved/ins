package main

import (
	"asiainfo.com/ins/cli"
	"asiainfo.com/ins/deploy/wls12c"
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/utils"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

/**
* go build -o deploy /veris/odc/app/go/3rd/src/asiainfo.com/ins/deploy/deploy.go
 */
func main() {
	chs2()
}

func chs1() {
	//ins_path := cli.CONF_HOME + "/server"
	//fs, err := utils.GetAllFiles(ins_path)
	//if err != nil {
	//	fmt.Println(err)
	//	panic(err)
	//}

	//for _, f := range fs {
	//	key := f.Name()
	//	fn := ins_path + "/" + key
	//	var srv Server
	//	var option string
	//}
}

func chs2() {
	ins_path := filepath.Join(cli.CONF_HOME, "deploy")
	fs, err := utils.GetAllFiles(ins_path)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for _, f := range fs {
		key := f.Name()
		fn := filepath.Join(ins_path, key)
		var dpy Deploy
		var option string
		var wls12 wls12c.Wls12c
		bs, _ := ioutil.ReadFile(fn)
		e := wls12.Json(bs)
		if e != nil {
			logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
			continue
		}
		dpy = &wls12
		option = wls12.Option
		if dpy != nil {
			var e error
			if option == "REMOVE" {
				e = dpy.Remove()
			} else {
				e = dpy.Add()
			}
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
		}
	}
}

type Deploy interface {
	Add() error
	Remove() error
}
