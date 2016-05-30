package cli
import (
	"asiainfo.com/ins/utils"
)

const (
	WLS12C string = "WLS12C"
	WLS11G string = "WLS11G"
	TOMCAT string = "TOMCAT"
	JDK7 string = "JDK7"
	REDIS3 string = "REDIS3"
	ACTIVEMQ512 string = "ACTIVEMQ512"

	WLS12CCONF string = "wls12c.conf"
	WLS11GCONF string = "wls11g.conf"
	TOMCATCONF string = "tomcat.conf"
	JDK7CONF string = "jdk7.conf"
	REDIS3CONF string = "redis3.conf"
	ACTIVEMQ512CONF string = "activemq512.conf"
)

var CONF_PATH = utils.GetParentDirectory(utils.GetCurrPath()) + "/conf/install.conf"
var CONF_HOME = utils.GetParentDirectory(utils.GetCurrPath()) + "/conf"
var LOGS_PATH = utils.GetParentDirectory(utils.GetCurrPath()) + "/logs"

var Install = make(map[string]string)

/**
 * 初始化安装配置文档
 */
//func InitInstall()  {
//
//	f, err := os.Open(CONF_PATH)
//	defer f.Close()
//	if err != nil {
//		e := logs.PrintErrorLog(LOGS_PATH, err.Error())
//		if e != nil {
//			fmt.Println(e)
//		}
//		panic(err)
//	}
//
//	rd := bufio.NewReader(f)
//	for {
//		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
//		if err != nil || io.EOF == err {
//			break
//		}
//		//相当于PHP的trim
//		line = strings.TrimSpace(line)
//		if len(line) == 0 {
//			continue
//		}
//		// 定义切割后的长度
//		lines := strings.SplitN(line, "=", 2)
//		key := strings.ToUpper(utils.TrimLeftRightSpace(lines[0]))
//		val := utils.TrimLeftRightSpace(lines[1])
//		Install[key] = val
//	}
//}
