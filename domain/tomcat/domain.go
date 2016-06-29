package tomcat

import (
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type Tomcat struct {
	TOMCAT_HOME  string   `json:"tomcatHome"`
	NATIVE_HOME  string   `json:"nativeHome"`
	JDK_HOME     string   `json:"jdkHome"`
	ServerLoader string   `json:"serverLoader"`
	SharedLoader string   `json:"sharedLoader"`
	DomainPath   string   `json:"domainPath"`
	AliasName    string   `json:"aliasName"`
	Apps         []App    `json:"apps"`
	Servers      []Server `json:"servers"`
	Protocol     string   `json:"protocol"`
	JVM          string   `json:"jvm"`
	Envs         []Env    `json:"envs"`
	Timeout      string   `json:"timeout"`
	Option       string   `json:"option"`
	ConsolePath  string   `json:"consolePath"`
}

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Server struct {
	Version      string `json:"version"`
	ListenAddr   string `json:"listenAddr"`
	ListenPort   string `json:"listenPort"`
	ShutdownPort string `json:"shutdownPort"`
}

type App struct {
	AppName  string `json:"appName"`
	APP_HOME string `json:"app_home"`
}

func (w *Tomcat) Json(bs []byte) error {
	return json.Unmarshal(bs, &w)
}

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Method */

func (tomcat *Tomcat) Add() error {
	err := tomcat.shell()
	tomcat.touchConsoleScript()
	return err
}

func (tomcat *Tomcat) Remove() error {
	return nil
}

/* End Method */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Configuration File */

const templateCatalinaProperties = `package.access=sun.,org.apache.catalina.,org.apache.coyote.,org.apache.jasper.,org.apache.tomcat.
package.definition=sun.,java.,org.apache.catalina.,org.apache.coyote.,\
org.apache.jasper.,org.apache.naming.,org.apache.tomcat.
common.loader=${catalina.base}/lib,${catalina.base}/lib/*.jar,${catalina.home}/lib,${catalina.home}/lib/*.jar
server.loader=%s
shared.loader=%s
tomcat.util.scan.StandardJarScanFilter.jarsToSkip=\
bootstrap.jar,commons-daemon.jar,tomcat-juli.jar,\
annotations-api.jar,el-api.jar,jsp-api.jar,servlet-api.jar,websocket-api.jar,\
catalina.jar,catalina-ant.jar,catalina-ha.jar,catalina-storeconfig.jar,\
catalina-tribes.jar,\
jasper.jar,jasper-el.jar,ecj-*.jar,\
tomcat-api.jar,tomcat-util.jar,tomcat-util-scan.jar,tomcat-coyote.jar,\
tomcat-dbcp.jar,tomcat-jni.jar,tomcat-websocket.jar,\
tomcat-i18n-en.jar,tomcat-i18n-es.jar,tomcat-i18n-fr.jar,tomcat-i18n-ja.jar,\
tomcat-juli-adapters.jar,catalina-jmx-remote.jar,catalina-ws.jar,\
tomcat-jdbc.jar,\
tools.jar,\
commons-beanutils*.jar,commons-codec*.jar,commons-collections*.jar,\
commons-dbcp*.jar,commons-digester*.jar,commons-fileupload*.jar,\
commons-httpclient*.jar,commons-io*.jar,commons-lang*.jar,commons-logging*.jar,\
commons-math*.jar,commons-pool*.jar,\
jstl.jar,taglibs-standard-spec-*.jar,\
geronimo-spec-jaxrpc*.jar,wsdl4j*.jar,\
ant.jar,ant-junit*.jar,aspectj*.jar,jmx.jar,h2*.jar,hibernate*.jar,httpclient*.jar,\
jmx-tools.jar,jta*.jar,log4j*.jar,mail*.jar,slf4j*.jar,\
xercesImpl.jar,xmlParserAPIs.jar,xml-apis.jar,\
junit.jar,junit-*.jar,ant-launcher.jar,\
cobertura-*.jar,asm-*.jar,dom4j-*.jar,icu4j-*.jar,jaxen-*.jar,jdom-*.jar,\
jetty-*.jar,oro-*.jar,servlet-api-*.jar,tagsoup-*.jar,xmlParserAPIs-*.jar,\
xom-*.jar
tomcat.util.scan.StandardJarScanFilter.jarsToScan=\
log4j-core*.jar,log4j-taglib*.jar,log4javascript*.jar,slf4j-taglib*.jar
# String cache configuration.
tomcat.util.buf.StringCache.byte.enabled=true
`

const templateServerXml = `<?xml version='1.0' encoding='utf-8'?>
<Server port="%s" shutdown="SHUTDOWN">
  <Listener className="org.apache.catalina.startup.VersionLoggerListener" />
  <Listener className="org.apache.catalina.core.AprLifecycleListener" SSLEngine="on" />
  <Listener className="org.apache.catalina.core.JreMemoryLeakPreventionListener" />
  <Listener className="org.apache.catalina.mbeans.GlobalResourcesLifecycleListener" />
  <Listener className="org.apache.catalina.core.ThreadLocalLeakPreventionListener" />
  <GlobalNamingResources>
    <Resource name="UserDatabase" auth="Container"
              type="org.apache.catalina.UserDatabase"
              description="User database that can be updated and saved"
              factory="org.apache.catalina.users.MemoryUserDatabaseFactory"
              pathname="conf/tomcat-users.xml" />
  </GlobalNamingResources>
  <Service name="Catalina">
    <Connector port="%s" protocol="%s" connectionTimeout="20000" maxConnections="2000" maxHttpHeaderSize="8192" maxThreads="300" minSpareThreads="100" maxSpareThreads="500"  maxProcessors="500" minProcessors="50" acceptorThreadCount="10" enableLookups="false" acceptCount="200"  disableUploadTimeout="true" URIEncoding="UTF-8" compression="on" compressionMinSize="10240" noCompressionUserAgents="gozilla, traviata" redirectPort="8443" />
    <Engine name="Catalina" defaultHost="localhost">

      <Realm className="org.apache.catalina.realm.LockOutRealm">
        <Realm className="org.apache.catalina.realm.UserDatabaseRealm"
               resourceName="UserDatabase"/>
      </Realm>
      <Host name="localhost"  appBase="webapps"
            unpackWARs="true" autoDeploy="true">
%s
      </Host>
    </Engine>
  </Service>
</Server>
`

const copyDomain = `#!/bin/bash

mkdir -p %s %s %s %s %s %s
cp -r %s %s
`

/**
 * touch conf local file
 */
func (tomcat *Tomcat) touchConf(server Server) {

	fullName := filepath.Join(tomcat.DomainPath, tomcat.AliasName, server.Version)
	catalinaProperties := filepath.Join(fullName, "conf", "catalina.properties")
	serverXml := filepath.Join(fullName, "conf", "server.xml")
	protocol := "org.apache.coyote.http11.Http11NioProtocol"
	if len(tomcat.Protocol) > 0 {
		protocol = tomcat.Protocol
	}

	// catalinaProperties
	logs.Print(ioutil.WriteFile(
		catalinaProperties,
		[]byte(fmt.Sprintf(templateCatalinaProperties,
			tomcat.ServerLoader,
			tomcat.SharedLoader,
		)),
		0750,
	))
	var servers string
	for _, app := range tomcat.Apps {
		servers += fmt.Sprintf("<Context path=\"%s\" docBase=\"%s\" reloadable=\"false\" crossContext=\"true\" allowLinking=\"true\"/>\n", "/"+app.AppName, app.APP_HOME)
	}
	// serverXml
	logs.Print(ioutil.WriteFile(
		serverXml,
		[]byte(fmt.Sprintf(templateServerXml,
			server.ShutdownPort,
			server.ListenPort, protocol,
			servers,
		)),
		0750,
	))
}

/**
 * touch conf local file
 */
func (tomcat *Tomcat) shell() error {
	for _, server := range tomcat.Servers {
		os.MkdirAll(utils.TMPD, 0750)
		tmp := filepath.Join(utils.TMPD, "exec.sh")
		sourceConf := filepath.Join(tomcat.TOMCAT_HOME, "conf")
		fullName := filepath.Join(tomcat.DomainPath, tomcat.AliasName, server.Version)
		targetConf := filepath.Join(fullName, "conf")
		targetTemp := filepath.Join(fullName, "temp")
		targetLogs := filepath.Join(fullName, "logs")
		targetWork := filepath.Join(fullName, "work")
		targetDump := filepath.Join(fullName, "dump")
		targetWebapps := filepath.Join(fullName, "webapps")
		// copy
		logs.Print(ioutil.WriteFile(
			tmp,
			[]byte(fmt.Sprintf(copyDomain,
				targetConf, targetTemp, targetLogs, targetWork, targetDump, targetWebapps,
				sourceConf, fullName,
			)),
			0750,
		))
		defer os.Remove(tmp)
		// 根据模板生成domain
		cmd := exec.Command(tmp)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
		tomcat.touchConf(server)
	}
	return nil
}

/* End Configuration File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Console File */
const templateStartConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s starting..."
%s
CATALINA_HOME="%s"
CATALINA_BASE="%s"
JAVA_HOME="%s"
export CATALINA_HOME CATALINA_BASE JAVA_HOME
export JVM="%s"
export JAVA_OPTS="${JVM} -Djava.library.path=%s -Dbalance_job_group_first_wait_time=20 -XX:ErrorFile=${CATALINA_BASE}/dump/err.log -XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=${CATALINA_BASE}/dump/  -Dfile.encoding=UTF-8 -Dencoding=UTF-8"
export CATALINA_OPTS="-DTOMCAT_NAME=%s -XX:+UseConcMarkSweepGC -XX:+CMSConcurrentMTEnabled -XX:CMSInitiatingOccupancyFraction=59 -XX:SurvivorRatio=4 -XX:ParallelGCThreads=8"
CATALINA_TMPDIR=${CATALINA_BASE}/temp
export CATALINA_TMPDIR
cd ${CATALINA_TMPDIR}
export CATALINA_PID=%s.pid
$CATALINA_HOME/bin/startup.sh
echo "%s started, pls wating 30 sec..."
`
const templateStopConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s stopping..."
CATALINA_HOME=%s
CATALINA_BASE=%s
JAVA_HOME=%s
export CATALINA_BASE CATALINA_HOME JAVA_HOME
CATALINA_TMPDIR=${CATALINA_BASE}/temp
export CATALINA_TMPDIR
cd ${CATALINA_TMPDIR}
export CATALINA_PID=%s.pid
$CATALINA_HOME/bin/catalina.sh stop %d -force
echo "%s stopped, pls wating 30 sec..."
`
const templateRestartConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s restarting..."
%s
%s
echo "%s restarted, pls wating 30 sec..."
`

/**
 * touch console file
 */
func (tomcat *Tomcat) touchConsoleScript() {
	err := utils.MkdirConsolesPath(tomcat.ConsolePath)
	if err != nil {
		logs.Print(err)
	}
	for _, server := range tomcat.Servers {
		//srvpath := filepath.Join(tomcat.DomainPath, tomcat.ServerName, tomcat.Version)
		simpleName := tomcat.AliasName + server.Version
		start := filepath.Join(tomcat.ConsolePath, "start", "start_"+simpleName+".sh")
		stop := filepath.Join(tomcat.ConsolePath, "stop", "stop_"+simpleName+".sh")
		restart := filepath.Join(tomcat.ConsolePath, "restart", "restart_"+simpleName+".sh")
		fullName := filepath.Join(tomcat.DomainPath, tomcat.AliasName, server.Version)
		nativeLib := filepath.Join(tomcat.NATIVE_HOME, "lib")
		var envs string
		now := time.Now().String()

		timeout := 10
		if len(tomcat.Timeout) > 0 {
			timeout, _ = strconv.Atoi(tomcat.Timeout)
		}
		jvm := "-Xmx512m -Xms512m -XX:MaxPermSize=128m -XX:PermSize=128m"
		if len(tomcat.JVM) > 0 {
			jvm = tomcat.JVM
		}
		if len(tomcat.Envs) > 0 {
			for _, env := range tomcat.Envs {
				envs += fmt.Sprintf("export %s=%s\n", env.Name, env.Value)
			}
		}
		// start
		logs.Print(ioutil.WriteFile(
			start,
			[]byte(fmt.Sprintf(
				templateStartConsole,
				now,
				simpleName,
				envs,
				tomcat.TOMCAT_HOME, fullName, tomcat.JDK_HOME,
				jvm, nativeLib,
				simpleName,
				simpleName,
				simpleName,
			)),
			0750,
		))

		// stop
		logs.Print(ioutil.WriteFile(
			stop,
			[]byte(fmt.Sprintf(templateStopConsole,
				now,
				simpleName,
				tomcat.TOMCAT_HOME, fullName, tomcat.JDK_HOME,
				simpleName, timeout,
				simpleName,
			)),
			0750,
		))

		// restart
		logs.Print(ioutil.WriteFile(
			restart,
			[]byte(fmt.Sprintf(templateRestartConsole,
				now,
				simpleName, stop, start, simpleName)),
			0750,
		))
	}
}

/* End Console File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
