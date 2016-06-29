package wls

import (
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

/*
WLS12C = jdkHome:/veris/odc/app/jdk, invLoc:/veris/odc/app/oraInventory, rspFile:/veris/odc/console/ins/conf/weblogic.rsp, jarLoc:/veris/odc/app/fmw_12.1.3.0.0_wls.jar
*/

const (
	jarLoc  string = "jarLoc"
	invLoc  string = "invLoc"
	rspFile string = "rspFile"
	jdkHome string = "jdkHome"
)

type Wls12c struct {
	JARLoc   string `json:"jarLoc"`
	InvLoc   string `json:"invLoc"`
	RspFile  string `json:"rspFile"`
	JDK_HOME string `json:"jdk_home"`
}

func (w *Wls12c) Builder(str string) error {
	strs := strings.Split(str, ",")
	for _, val := range strs {
		vals := strings.Split(val, ":")
		a := utils.TrimLeftRightSpace(vals[0])
		b := utils.TrimLeftRightSpace(vals[1])
		if len(b) <= 0 {
			return fmt.Errorf("%s can not be empyt.", a)
		}
		switch a {
		case jarLoc:
			w.JARLoc = b
		case invLoc:
			w.InvLoc = b
		case rspFile:
			w.RspFile = b
		case jdkHome:
			w.JDK_HOME = b
		default:
			continue
		}
	}
	return nil
}

func (w *Wls12c) Json(bs []byte) error {
	return json.Unmarshal(bs, &w)
}

func (w *Wls12c) Install() error {
	java := filepath.Join(w.JDK_HOME, "bin", "java")
	jvms := "-Xms512m"
	jvmx := "-Xmx512m"
	user, _ := user.Current()

	inv := touchInv(w.InvLoc, user.Gid)
	defer logs.Print(os.Remove(inv))

	cmd := exec.Command(java, jvmx, jvms, "-jar", w.JARLoc, "-silent", "-responseFile", w.RspFile, "-invPtrLoc", inv)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

/**
 * touch inv local file
 */
func touchInv(inv_loc, gname string) string {
	os.MkdirAll(utils.TMPD, 0750)
	file, err := ioutil.TempFile(utils.TMPD, "install_wls12c_") //在DIR目录下创建tmp为文件名前缀的文件，获得file文件指针，DIR必须存在，否则创建不成功
	defer file.Close()
	if err != nil {
		logs.Print(fmt.Errorf("create %s error.", file.Name()))
		panic(err)
		return ""
	}
	file.WriteString("inventory_loc=" + inv_loc)
	file.WriteString("\ninst_group=" + gname)
	return file.Name()
}
