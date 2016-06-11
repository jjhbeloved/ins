package cli
import (
	"asiainfo.com/ins/utils"
)

const (
	WLS12C string = "WLS12C"
	WLS11G string = "WLS11G"
	TOMCAT string = "TOMCAT"
	JDK string = "JDK"
	MEMCACHED string = "MEMCACHED"
	REDIS string = "REDIS"
	ACTIVEMQ string = "ACTIVEMQ"
	ZK string = "ZK"

	WLS12CCONF string = "wls12c.json"
	WLS11GCONF string = "wls11g.json"
	TOMCATCONF string = "tomcat.json"
	JDKCONF string = "jdk.json"
	MEMCACHEDCONF string = "memcached.json"
	REDISCONF string = "redis.json"
	ACTIVEMQCONF string = "amq.json"
	ZKCONF string = "zk.json"

	WLS12CPREFIX string = "wls12c"
	WLS11GPREFIX string = "wls11g"
	TOMCATPREFIX string = "tomcat"
	JDKPREFIX string = "jdk"
	MEMCACHEDPREFIX string = "memcached"
	REDISPREFIX string = "redis"
	ACTIVEMQPREFIX string = "amq"
	ZKPREFIX string = "zk"
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
