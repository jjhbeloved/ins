package amq

import (
	"encoding/json"
	"os"
	"path/filepath"
	"os/exec"
	"time"
	"asiainfo.com/ins/utils"
	"io/ioutil"
	"asiainfo.com/ins/logs"
	"fmt"
	"strconv"
)

type AMQ struct {
	AMQ_PKG       string       `json:"amq_pkg"`
	AMQ_HOME      string       `json:"amq_home"`
	HEAP_SIZE     int          `json:"heap_size"`
	MQ_Name       string       `json:"mq_name"`
	MulticastAddr string       `json:"multicastAddr"`
	GroupName     string       `json:"groupName"`
	ListenAddr    string       `json:"listenAddr"`
	ConsolePath   string       `json:"consolePath"`
	AdminName     string       `json:"adminName"`
	AdminPWD      string       `json:"adminPWD"`
	UserName      string       `json:"userName"`
	UserPWD       string       `json:"userPWD"`
	WebPort       string       `json:"webPort"`
}

func (w *AMQ) Json(bs []byte) error {
	return json.Unmarshal(bs, &w)
}

const installAMQ = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
mkdir -p %s
tar xzf %s -C %s --strip-components=1
cd %s
`

func (amq *AMQ) Install() error {
	_, err := os.Stat(amq.AMQ_PKG)
	if err != nil {
		return err
	}
	now := time.Now().String()
	amqSh := filepath.Join(utils.TMPD, "amqInstall.sh")
	defer os.Remove(amqSh)
	logs.Print(ioutil.WriteFile(
		amqSh,
		[]byte(fmt.Sprintf(installAMQ,
			now,
			amq.AMQ_HOME,
			amq.AMQ_PKG, amq.AMQ_HOME,
			amq.AMQ_HOME,
		)),
		0750,
	))
	// 根据模板生成domain
	cmd := exec.Command(amqSh)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	amq.touchConf()
	amq.touchConsoleScript()

	return nil

}


/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Configuration File */

const wrapperConf = `
#wrapper.debug=TRUE
set.default.ACTIVEMQ_HOME=../..
set.default.ACTIVEMQ_BASE=../..
set.default.ACTIVEMQ_CONF=%ACTIVEMQ_BASE%/conf
set.default.ACTIVEMQ_DATA=%ACTIVEMQ_BASE%/data
wrapper.working.dir=.

# Java Application
wrapper.java.command=java

wrapper.java.mainclass=org.tanukisoftware.wrapper.WrapperSimpleApp

wrapper.java.classpath.1=%ACTIVEMQ_HOME%/bin/wrapper.jar
wrapper.java.classpath.2=%ACTIVEMQ_HOME%/bin/activemq.jar

wrapper.java.library.path.1=%ACTIVEMQ_HOME%/bin/linux-x86-64/

wrapper.java.additional.1=-Dactivemq.home=%ACTIVEMQ_HOME%
wrapper.java.additional.2=-Dactivemq.base=%ACTIVEMQ_BASE%
wrapper.java.additional.3=-Djavax.net.ssl.keyStorePassword=password
wrapper.java.additional.4=-Djavax.net.ssl.trustStorePassword=password
wrapper.java.additional.5=-Djavax.net.ssl.keyStore=%ACTIVEMQ_CONF%/broker.ks
wrapper.java.additional.6=-Djavax.net.ssl.trustStore=%ACTIVEMQ_CONF%/broker.ts
wrapper.java.additional.7=-Dcom.sun.management.jmxremote
wrapper.java.additional.8=-Dorg.apache.activemq.UseDedicatedTaskRunner=true
wrapper.java.additional.9=-Djava.util.logging.config.file=logging.properties
wrapper.java.additional.10=-Dactivemq.conf=%ACTIVEMQ_CONF%
wrapper.java.additional.11=-Dactivemq.data=%ACTIVEMQ_DATA%
wrapper.java.additional.12=-Djava.security.auth.login.config=%ACTIVEMQ_CONF%/login.config

# Uncomment to enable jmx
#wrapper.java.additional.n=-Dcom.sun.management.jmxremote.port=1616
#wrapper.java.additional.n=-Dcom.sun.management.jmxremote.authenticate=false
#wrapper.java.additional.n=-Dcom.sun.management.jmxremote.ssl=false

# Uncomment to enable YourKit profiling
#wrapper.java.additional.n=-Xrunyjpagent

# Uncomment to enable remote debugging
#wrapper.java.additional.n=-Xdebug -Xnoagent -Djava.compiler=NONE
#wrapper.java.additional.n=-Xrunjdwp:transport=dt_socket,server=y,suspend=n,address=5005

# Initial Java Heap Size (in MB)
#wrapper.java.initmemory=3

# Application parameters.  Add parameters as needed starting from 1
wrapper.app.parameter.1=org.apache.activemq.console.Main
wrapper.app.parameter.2=start

# Format of output for the console.  (See docs for formats)
wrapper.console.format=PM

# Log Level for console output.  (See docs for log levels)
wrapper.console.loglevel=INFO

# Log file to use for wrapper output logging.
wrapper.logfile=%ACTIVEMQ_DATA%/wrapper.log

# Format of output for the log file.  (See docs for formats)
wrapper.logfile.format=LPTM

# Log Level for log file output.  (See docs for log levels)
wrapper.logfile.loglevel=INFO

# Maximum size that the log file will be allowed to grow to before
#  the log is rolled. Size is specified in bytes.  The default value
#  of 0, disables log rolling.  May abbreviate with the 'k' (kb) or
#  'm' (mb) suffix.  For example: 10m = 10 megabytes.
wrapper.logfile.maxsize=50

# Maximum number of rolled log files which will be allowed before old
#  files are deleted.  The default value of 0 implies no limit.
wrapper.logfile.maxfiles=3

# Log Level for sys/event log output.  (See docs for log levels)
wrapper.syslog.loglevel=NONE

# Title to use when running as a console
wrapper.console.title=ActiveMQ

# Name of the service
wrapper.ntservice.name=ActiveMQ

# Display name of the service
wrapper.ntservice.displayname=ActiveMQ

# Description of the service
wrapper.ntservice.description=ActiveMQ Broker

# Service dependencies.  Add dependencies as needed starting from 1
wrapper.ntservice.dependency.1=

# Mode in which the service is installed.  AUTO_START or DEMAND_START
wrapper.ntservice.starttype=AUTO_START

# Allow the service to interact with the desktop.
wrapper.ntservice.interactive=false

`
const wrapperConfA = `
# Maximum Java Heap Size (in MB)
wrapper.java.maxmemory=%d
`

const amqMXml = `
<beans
  xmlns="http://www.springframework.org/schema/beans"
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://www.springframework.org/schema/beans http://www.springframework.org/schema/beans/spring-beans.xsd
  http://activemq.apache.org/schema/core http://activemq.apache.org/schema/core/activemq-core.xsd">

    <bean class="org.springframework.beans.factory.config.PropertyPlaceholderConfigurer">
        <property name="locations">
            <value>file:${activemq.conf}/credentials.properties</value>
        </property>
    </bean>

    <bean id="logQuery" class="io.fabric8.insight.log.log4j.Log4jLogQuery"
          lazy-init="false" scope="singleton"
          init-method="start" destroy-method="stop">
    </bean>
    <broker xmlns="http://activemq.apache.org/schema/core" brokerName="%s" dataDirectory="${activemq.data}">

        <destinationPolicy>
            <policyMap>
              <policyEntries>
                <policyEntry topic=">" >
                  <pendingMessageLimitStrategy>
                    <constantPendingMessageLimitStrategy limit="1000"/>
                  </pendingMessageLimitStrategy>
                </policyEntry>
              </policyEntries>
            </policyMap>
        </destinationPolicy>
        <managementContext>
            <managementContext createConnector="false"/>
        </managementContext>
        <networkConnectors>
           <networkConnector uri="multicast://%s?group=%s"
            dynamicOnly="true"
            networkTTL="3"
            prefetchSize="1"
            decreaseNetworkConsumerPriority="true" />
        </networkConnectors>
        <persistenceAdapter>
            <kahaDB directory="${activemq.data}/kahadb" indexCacheSize="100000" indexWriteBatchSize="1000" enableJournalDiskSyncs="false"  journalMaxFileLength="128mb" concurrentStoreAndDispatchQueues="true" concurrentStoreAndDispatchTopics="true"/>
        </persistenceAdapter>
          <systemUsage>
            <systemUsage>
                <memoryUsage>
                    <memoryUsage percentOfJvmHeap="70" />
                </memoryUsage>
                <storeUsage>
                    <storeUsage limit="100 gb"/>
                </storeUsage>
                <tempUsage>
                    <tempUsage limit="50 gb"/>
                </tempUsage>
            </systemUsage>
        </systemUsage>

        <transportConnectors>
            <transportConnector name="openwire" uri="tcp://%s?maximumConnections=1000&amp;wireFormat.maxFrameSize=104857600" discoveryUri="multicast://%s?group=%s" enableStatusMonitor="true" updateClusterClients="true" rebalanceClusterClients="true" updateClusterClientsOnRemove="true"/>
        </transportConnectors>
        <shutdownHooks>
            <bean xmlns="http://www.springframework.org/schema/beans" class="org.apache.activemq.hooks.SpringContextHook" />
        </shutdownHooks>

    </broker>
    <import resource="jetty.xml"/>
</beans>
`

const amqSXml = `
<beans
  xmlns="http://www.springframework.org/schema/beans"
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://www.springframework.org/schema/beans http://www.springframework.org/schema/beans/spring-beans.xsd
  http://activemq.apache.org/schema/core http://activemq.apache.org/schema/core/activemq-core.xsd">

    <bean class="org.springframework.beans.factory.config.PropertyPlaceholderConfigurer">
        <property name="locations">
            <value>file:${activemq.conf}/credentials.properties</value>
        </property>
    </bean>

    <bean id="logQuery" class="io.fabric8.insight.log.log4j.Log4jLogQuery"
          lazy-init="false" scope="singleton"
          init-method="start" destroy-method="stop">
    </bean>
    <broker xmlns="http://activemq.apache.org/schema/core" brokerName="%s" dataDirectory="${activemq.data}">

        <destinationPolicy>
            <policyMap>
              <policyEntries>
                <policyEntry topic=">" >
                  <pendingMessageLimitStrategy>
                    <constantPendingMessageLimitStrategy limit="1000"/>
                  </pendingMessageLimitStrategy>
                </policyEntry>
              </policyEntries>
            </policyMap>
        </destinationPolicy>
        <managementContext>
            <managementContext createConnector="false"/>
        </managementContext>
        <persistenceAdapter>
            <kahaDB directory="${activemq.data}/kahadb" indexCacheSize="100000" indexWriteBatchSize="1000" enableJournalDiskSyncs="false"  journalMaxFileLength="128mb" concurrentStoreAndDispatchQueues="true" concurrentStoreAndDispatchTopics="true"/>
        </persistenceAdapter>
          <systemUsage>
            <systemUsage>
                <memoryUsage>
                    <memoryUsage percentOfJvmHeap="70" />
                </memoryUsage>
                <storeUsage>
                    <storeUsage limit="100 gb"/>
                </storeUsage>
                <tempUsage>
                    <tempUsage limit="50 gb"/>
                </tempUsage>
            </systemUsage>
        </systemUsage>

        <transportConnectors>
            <transportConnector name="openwire" uri="tcp://%s?maximumConnections=1000&amp;wireFormat.maxFrameSize=104857600"/>
        </transportConnectors>
        <shutdownHooks>
            <bean xmlns="http://www.springframework.org/schema/beans" class="org.apache.activemq.hooks.SpringContextHook" />
        </shutdownHooks>

    </broker>
    <import resource="jetty.xml"/>
</beans>
`

const jettyXml = `
<beans xmlns="http://www.springframework.org/schema/beans" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
    xsi:schemaLocation="http://www.springframework.org/schema/beans http://www.springframework.org/schema/beans/spring-beans.xsd">

    <bean id="securityLoginService" class="org.eclipse.jetty.security.HashLoginService">
        <property name="name" value="ActiveMQRealm" />
        <property name="config" value="${activemq.conf}/jetty-realm.properties" />
    </bean>

    <bean id="securityConstraint" class="org.eclipse.jetty.util.security.Constraint">
        <property name="name" value="BASIC" />
        <property name="roles" value="user,admin" />
        <!-- set authenticate=false to disable login -->
        <property name="authenticate" value="true" />
    </bean>
    <bean id="adminSecurityConstraint" class="org.eclipse.jetty.util.security.Constraint">
        <property name="name" value="BASIC" />
        <property name="roles" value="admin" />
         <!-- set authenticate=false to disable login -->
        <property name="authenticate" value="true" />
    </bean>
    <bean id="securityConstraintMapping" class="org.eclipse.jetty.security.ConstraintMapping">
        <property name="constraint" ref="securityConstraint" />
        <property name="pathSpec" value="/api/*,/admin/*,*.jsp" />
    </bean>
    <bean id="adminSecurityConstraintMapping" class="org.eclipse.jetty.security.ConstraintMapping">
        <property name="constraint" ref="adminSecurityConstraint" />
        <property name="pathSpec" value="*.action" />
    </bean>

    <bean id="rewriteHandler" class="org.eclipse.jetty.rewrite.handler.RewriteHandler">
        <property name="rules">
            <list>
                <bean id="header" class="org.eclipse.jetty.rewrite.handler.HeaderPatternRule">
                  <property name="pattern" value="*"/>
                  <property name="name" value="X-FRAME-OPTIONS"/>
                  <property name="value" value="SAMEORIGIN"/>
                </bean>
            </list>
        </property>
    </bean>

	<bean id="secHandlerCollection" class="org.eclipse.jetty.server.handler.HandlerCollection">
		<property name="handlers">
			<list>
   	            <ref bean="rewriteHandler"/>
				<bean class="org.eclipse.jetty.webapp.WebAppContext">
					<property name="contextPath" value="/admin" />
					<property name="resourceBase" value="${activemq.home}/webapps/admin" />
					<property name="logUrlOnStart" value="true" />
				</bean>
				<bean class="org.eclipse.jetty.webapp.WebAppContext">
					<property name="contextPath" value="/api" />
					<property name="resourceBase" value="${activemq.home}/webapps/api" />
					<property name="logUrlOnStart" value="true" />
				</bean>
				<bean class="org.eclipse.jetty.server.handler.ResourceHandler">
					<property name="directoriesListed" value="false" />
					<property name="welcomeFiles">
						<list>
							<value>index.html</value>
						</list>
					</property>
					<property name="resourceBase" value="${activemq.home}/webapps/" />
				</bean>
				<bean id="defaultHandler" class="org.eclipse.jetty.server.handler.DefaultHandler">
					<property name="serveIcon" value="false" />
				</bean>
			</list>
		</property>
	</bean>
    <bean id="securityHandler" class="org.eclipse.jetty.security.ConstraintSecurityHandler">
        <property name="loginService" ref="securityLoginService" />
        <property name="authenticator">
            <bean class="org.eclipse.jetty.security.authentication.BasicAuthenticator" />
        </property>
        <property name="constraintMappings">
            <list>
                <ref bean="adminSecurityConstraintMapping" />
                <ref bean="securityConstraintMapping" />
            </list>
        </property>
        <property name="handler" ref="secHandlerCollection" />
    </bean>

    <bean id="contexts" class="org.eclipse.jetty.server.handler.ContextHandlerCollection">
    </bean>

    <bean id="jettyPort" class="org.apache.activemq.web.WebConsolePort" init-method="start">
             <!-- the default port number for the web console -->
        <property name="host" value="0.0.0.0"/>
        <property name="port" value="%d"/>
    </bean>

    <bean id="Server" depends-on="jettyPort" class="org.eclipse.jetty.server.Server"
        destroy-method="stop">

        <property name="handler">
            <bean id="handlers" class="org.eclipse.jetty.server.handler.HandlerCollection">
                <property name="handlers">
                    <list>
                        <ref bean="contexts" />
                        <ref bean="securityHandler" />
                    </list>
                </property>
            </bean>
        </property>

    </bean>

    <bean id="invokeConnectors" class="org.springframework.beans.factory.config.MethodInvokingFactoryBean">
    	<property name="targetObject" ref="Server" />
    	<property name="targetMethod" value="setConnectors" />
    	<property name="arguments">
    	<list>
           	<bean id="Connector" class="org.eclipse.jetty.server.ServerConnector">
           		<constructor-arg ref="Server" />
                   <property name="host" value="#{systemProperties['jetty.host']}" />
                   <property name="port" value="#{systemProperties['jetty.port']}" />
               </bean>
            </list>
    	</property>
    </bean>
	<bean id="configureJetty" class="org.springframework.beans.factory.config.MethodInvokingFactoryBean">
		<property name="staticMethod" value="org.apache.activemq.web.config.JspConfigurer.configureJetty" />
		<property name="arguments">
			<list>
				<ref bean="Server" />
				<ref bean="secHandlerCollection" />
			</list>
		</property>
	</bean>
    <bean id="invokeStart" class="org.springframework.beans.factory.config.MethodInvokingFactoryBean"
    	depends-on="configureJetty, invokeConnectors">
    	<property name="targetObject" ref="Server" />
    	<property name="targetMethod" value="start" />
    </bean>
</beans>
`

const jettyReal = `
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
%s: %s, admin
%s: %s, user
`

/**
 * touch conf local file
 */
func (amq *AMQ) touchConf() {
	_wrapperConf := filepath.Join(amq.AMQ_HOME, "bin", "linux-x86-64", "wrapper.conf")
	_amqXml := filepath.Join(amq.AMQ_HOME, "conf", "activemq.xml")
	_jettyXml := filepath.Join(amq.AMQ_HOME, "conf", "jetty.xml")
	_jettyReal := filepath.Join(amq.AMQ_HOME, "conf", "jetty-realm.properties")
	now := time.Now().String()

	userName := "user"
	userPWD := "user"
	if len(amq.UserName) > 0 {
		userName = amq.UserName
	}
	if len(amq.UserPWD) > 0 {
		userPWD = amq.UserPWD
	}
	adminName := "admin"
	adminPWD := "admin"
	if len(amq.AdminName) > 0 {
		adminName = amq.AdminName
	}
	if len(amq.AdminPWD) > 0 {
		adminPWD = amq.AdminPWD
	}
	mqName := "humbird"
	group := "humbird"
	addr := "0.0.0.0:61616"
	webPort := 8161

	if len(amq.GroupName) > 0 {
		group = amq.GroupName
	}
	if len(amq.MQ_Name) > 0 {
		mqName = amq.MQ_Name
	}
	if len(amq.ListenAddr) > 0 {
		addr = amq.ListenAddr
	}
	if len(amq.WebPort) > 0 {
		webPort, _ = strconv.Atoi(amq.WebPort)
	}

	// wrapperConf
	logs.Print(ioutil.WriteFile(
		_wrapperConf,
		[]byte(wrapperConf + fmt.Sprintf(wrapperConfA,
			amq.HEAP_SIZE,
		)),
		0750,
	))

	multicast := amq.MulticastAddr
	if len(multicast) > 0 {
		// multicast := "224.1.2.3:6255"
		// amqXml
		logs.Print(ioutil.WriteFile(
			_amqXml,
			[]byte(fmt.Sprintf(amqMXml,
				mqName, multicast, group,
				addr, multicast, group,
			)),
			0750,
		))
	} else {
		// amqXml
		logs.Print(ioutil.WriteFile(
			_amqXml,
			[]byte(fmt.Sprintf(amqSXml,
				mqName,
				addr,
			)),
			0750,
		))
	}

	// jettyXml
	logs.Print(ioutil.WriteFile(
		_jettyXml,
		[]byte(fmt.Sprintf(jettyXml,
			webPort,
		)),
		0750,
	))

	// jettyReal
	logs.Print(ioutil.WriteFile(
		_jettyReal,
		[]byte(fmt.Sprintf(jettyReal,
			now,
			adminName, adminPWD,
			userName, userPWD,
		)),
		0750,
	))
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
echo "%s started, pls wating 30 sec..."
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

const templateStopConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s stopping..."
%s
echo "%s stopped, pls wating 30 sec..."
`
/**
 * touch console file
 */
func (amq *AMQ) touchConsoleScript() {
	err := utils.MkdirConsolesPath(amq.ConsolePath)
	if err != nil {
		logs.Print(err)
	}
	start := filepath.Join(amq.ConsolePath, "start", "start_" + amq.MQ_Name + ".sh")
	stop := filepath.Join(amq.ConsolePath, "stop", "stop_" + amq.MQ_Name + ".sh")
	restart := filepath.Join(amq.ConsolePath, "restart", "restart_" + amq.MQ_Name + ".sh")
	now := time.Now().String()
	// start
	logs.Print(ioutil.WriteFile(
		start,
		[]byte(fmt.Sprintf(templateStartConsole,
			now,
			amq.MQ_Name,
			filepath.Join(amq.AMQ_HOME, "bin", "linux-x86-64", "activemq") + " start",
			amq.MQ_Name,
		)),
		0750,
	))

	// stop
	logs.Print(ioutil.WriteFile(
		stop,
		[]byte(fmt.Sprintf(templateStopConsole,
			now,
			amq.MQ_Name,
			filepath.Join(amq.AMQ_HOME, "bin", "linux-x86-64", "activemq") + " stop",
			amq.MQ_Name,
		)),
		0750,
	))

	// restart
	logs.Print(ioutil.WriteFile(
		restart,
		[]byte(fmt.Sprintf(templateRestartConsole,
			now,
			amq.MQ_Name, stop, start, amq.MQ_Name)),
		0750,
	))
}

/* End Console File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
