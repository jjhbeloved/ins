package main
import (
	"asiainfo.com/ins/cli"
	"fmt"
	"asiainfo.com/ins/utils"
	"io/ioutil"
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/install/wls"
	"path/filepath"
	"asiainfo.com/ins/install/tomcat"
)

func init()  {
	//cli.InitInstall()
}

/**
* go build -o install /veris/odc/app/go/3rd/src/asiainfo.com/ins/install/install.go
*/
func main()  {
	chs2()
}

/**
 * 遍历安装配置文档, 遍历安装
 */
func chs1()  {
	for key, val := range cli.Install{
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
func chs2()  {
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