package main

import (
	"asiainfo.com/ins/cli"
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/server/wls12c"
	"asiainfo.com/ins/utils"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

/**
* go build -o server /veris/odc/app/go/3rd/src/asiainfo.com/ins/server/server.go
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
	ins_path := filepath.Join(cli.CONF_HOME, "server")
	fs, err := utils.GetAllFiles(ins_path)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for _, f := range fs {
		key := f.Name()
		fn := filepath.Join(ins_path, key)
		var srv Server
		var option string
		var wls12 wls12c.Wls12c
		bs, _ := ioutil.ReadFile(fn)
		e := wls12.Json(bs)
		if e != nil {
			logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
			continue
		}
		srv = &wls12
		option = wls12.Option
		if srv != nil {
			var e error
			if option == "REMOVE" {
				e = srv.Remove()
			} else {
				e = srv.Add()
			}
			if e != nil {
				logs.PrintErrorLog(cli.LOGS_PATH, e.Error())
				continue
			}
		}
	}
}

type Server interface {
	Add() error
	Remove() error
}
