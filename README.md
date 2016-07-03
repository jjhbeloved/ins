# ins
Install with Redis, Memcached, Tomcat, Weblogic, JDK, Zookeeper, ActiveMQ ...(e.g.: Hdfs, Ceph or more).
With Xiaoxiao.

# web
Add new future about web ui and restful api for ins.

# ALL 3RD
3RD basic pkg install, Include:
> - **JDK**
> - **Memcached**
> - **ActiveMQ**
> - **Zookeeper**
> - **Redis**
> - **Tomcat**
> - **Weblogic**

# JDK
- JDK insall support **root user** and **normal user**.
- root user，will change ***/etc/bashrc***, 保证每个普通用户登录都会加载到这个环境变量.
- normal user，will change ***${HOME}/.bashrc***, 这个修改只针对一个普通用户的环境变量.

---------------------------------------
> - **root user**

>> 1. user login
>> 2. download dependency ***jdk_1.7.0_17.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 3. create ***${app}/conf/install/jdk.json*** (e.g: required)
>> 4. ***${app}/bin/install***$，wait...，JDK installed.
>> 5. jdk 安装作为基础环境变量生效, 具体设置了哪些请参考 ***/etc/bashrc*** 文件
```
{
  "pkg" : "/veris/odc/install/src/jdk_1.7.0_17.tar.gz",
  "jdk_home" : "/usr/bin/jdk7",
  "isRoot" : true
}
```
>> - ***pkg*** is JDK app pkg absolute path
>> - ***jdk_home*** is ${JAVA_HOME} path
>> - ***isRoot*** choose root=true/normal user=false mode

> - **normal user**

>> 1. user login
>> 2. download dependency ***jdk_1.7.0_17.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 3. create ***${app}/conf/install/jdk.json*** (e.g: required)
>> 4. ***${app}/bin/install***$，wait...，JDK installed.
>> 5. jdk 安装作为基础环境变量生效, 具体设置了哪些请参考 ***${HOME}/.bashrc*** 文件
```
{
  "pkg" : "/veris/odc/install/src/jdk_1.7.0_17.tar.gz",
  "jdk_home" : "/veris/odc/app/jdk7",
  "isRoot" : false
}
```
>> - ***pkg*** is JDK app pkg absolute path
>> - ***jdk_home*** is ${JAVA_HOME} path
>> - ***isRoot*** choose root=true/normal user=false mode

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @jdk.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/install/jdk.json
```

# Memcached
- Memcached install have ***Alone Mode*** and ***HA*** Mode.
- Memcached 开源社区本身并 ***no support HA replication Mode***, the ha mode 由开源社区提供，最新的patch维护在1.4.13版本. 其余版本的patch均为无效.
- New Memcached version up to 1.4.25，此版本没有HA patch，需要由应用本身做多点写保证HA. 

---------------------------------------
> - **Alone Mode**

>> 1. download dependency ***libevent-2.0.22-stable.tar.gz, Memcached-1.4.25.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/memcached.json*** (e.g: required)
>> 3. ***${app}/bin/install***$，wait...，Memcached Alone installed.
>> 4. start/stop/restart script in ***consolePath***
```
{
  "memcached_pkg": "/veris/odc/install/src/memcached-1.4.25.tar.gz",
  "libevent_pkg": "/veris/odc/install/src/libevent-2.0.22-stable.tar.gz",
  "memcached_home": "/veris/odc/memcached",
  "libevent_home": "/veris/odc/memcached/libevent",
  "port": "9999",
  "memory" : "512",
  "connections" : "2048",
  "option" : "ADD",
  "consolePath" : "/veris/odc/console"
}
```
>> - ***memcached_pkg*** is memcached app pkg absolute path
>> - ***libevent_pkg*** is libevent app pkg absolute path
>> - ***memcached_home*** is memcached be installed home path
>> - ***libevent_home*** is libevent be installed home path
>> - ***port*** is memcached LISTEN port
>> - ***memory*** is memcached 启动后占用的内存块大小, 单位是M
>> - ***conections*** is memcached 启动后允许的最大连接数
>> - ***option*** choose ADD install option, 暂时没有一键卸载功能
>> - ***consolePath*** is be installed memcached, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后memcached占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

> - **HA Mode**

>> 1. download dependency ***libevent-2.0.22-stable.tar.gz, Memcached-1.4.13.tar.gz, repcached-2.3.1-1.4.13.patch*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/memcached.json*** (e.g: required)
>> 3. ***${app}/bin/install***$，wait...，Memcached Alone installed.
>> 4. start/stop/restart script in ***consolePath***
```
{
  "memcached_pkg": "/veris/odc/install/src/memcached-1.4.13.tar.gz",
  "libevent_pkg": "/veris/odc/install/src/libevent-2.0.22-stable.tar.gz",
  "repcached_patch": "/veris/odc/install/src/repcached-2.3.1-1.4.13.patch",
  "memcached_home": "/veris/odc/memcached",
  "libevent_home": "/veris/odc/memcached/libevent",
  "port": "9999",
  "memory" : "512",
  "connections" : "2048",
  "repecachedPort" : "9998",
  "repecachedAddress" : "10.1.245.185",
  "option" : "ADD",
  "consolePath" : "/veris/odc/console"
}
```
>> - ***memcached_pkg*** is memcached app pkg absolute path
>> - ***libevent_pkg*** is libevent app pkg absolute path
>> - ***repcached_patch*** is memcached ha repached patch pkg absolute path
>> - ***memcached_home*** is memcached be installed home path
>> - ***libevent_home*** is libevent be installed home path
>> - ***port*** is memcached LISTEN port
>> - ***memory*** is memcached 启动后占用的内存块大小, 单位是M
>> - ***conections*** is memcached 启动后允许的最大连接数
>> - ***repecachedPort*** is be HA两个Memcached直接进行replication的通信端口
>> - ***repecachedAddress*** is other memcached server address
>> - ***option*** choose ADD install option, 暂时没有一键卸载功能
>> - ***consolePath*** is be installed memcached, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后memcached占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @memcached.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/install/memcached.json
```


