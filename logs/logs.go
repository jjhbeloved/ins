package logs

import (
	"asiainfo.com/ins/cli"
	"asiainfo.com/ins/utils"
	"os"
	"time"
)

func PrintInfoLog(fname, content string) error {
	fname = fname + "/ins_" + time.Now().Format("2016-01-20") + ".log"
	bs := []byte("[INFO] " + time.Now().Format("2016-01-20_11:11:11") + content)
	return print(fname, bs)
}

func PrintErrorLog(fname, content string) error {
	fname = fname + "/ins_" + time.Now().Format(utils.DATE_DIR) + "_error.log"
	bs := []byte("[ERROR] " + time.Now().Format(utils.DATE_FILE) + " " + content + "\n")
	return print(fname, bs)
}

func print(fname string, bs []byte) error {
	return utils.WriteFileA(fname, bs, os.ModePerm)
}

func Print(err error) {
	if err != nil {
		PrintErrorLog(cli.LOGS_PATH, err.Error())
	}
}
