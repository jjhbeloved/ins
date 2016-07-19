# ins
Install with Redis, Memcached, Tomcat, Weblogic, JDK, Zookeeper, ActiveMQ ...(e.g.: Hdfs, Ceph or more).
With Xiaoxiao.

# web
Add new future about web ui and restful api for ins.

# Build
> - **install**
go build -O ${target}/install ${src}/asiainfo.com/ins/install/install.go
> - **domain**
go build -O ${target}/domain ${src}/asiainfo.com/ins/domain/domain.go
> - **server**
go build -O ${target}/server ${src}/asiainfo.com/ins/server/server.go
> - **deploy**
go build -O ${target}/deploy ${src}/asiainfo.com/ins/deploy/deploy.go
> - **restful**
go build -O ${target}/restful ${src}/asiainfo.com/ins/restful/main.go
>>> - 需要将 ${src}/asiainfo.com/ins/restful/conf 目录下载到和二进制 restful 文件同级目录[修改配置 runmode ="prod"]
>>> - sh restful 后可以使用 restful api 接口
>>> - 默认端口 8088

# ALL 3RD
3RD basic pkg install, Include:
> - **JDK**
> - **Memcached**
> - **ActiveMQ**
> - **Zookeeper**
> - **Redis**
> - **Storm**
> - **MySQL**
> - **Tomcat**
> - **Weblogic**

# JDK
- JDK insall support **root user** and **normal user**.
- root user, will change ***/etc/bashrc***, 保证每个普通用户登录都会加载到这个环境变量.
- normal user, will change ***${HOME}/.bashrc***, 这个修改只针对一个普通用户的环境变量.

---------------------------------------
> - **root user**

>> 1. user login
>> 2. download dependency ***jdk_1.7.0_17.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 3. create ***${app}/conf/install/jdk.json*** (e.g: required)
>> 4. ***${app}/bin/install***, wait..., JDK installed.
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
>> 4. ***${app}/bin/install***, wait..., JDK installed.
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
- Memcached 开源社区本身并 ***no support HA replication Mode***, the ha mode 由开源社区提供, 最新的patch维护在1.4.13版本. 其余版本的patch均为无效.
- New Memcached version up to 1.4.25, 此版本没有HA patch, 需要由应用本身做多点写保证HA. 

---------------------------------------
> - **Alone Mode**

>> 1. download dependency ***libevent-2.0.22-stable.tar.gz, Memcached-1.4.25.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/memcached.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait..., Memcached Alone installed.
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
>> 3. ***${app}/bin/install***, wait..., Memcached Alone installed.
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

# ActiveMQ
- ActiveMQ的安装分为单机模式和桥接多播发现模式.
- 单机模式下的ActiveMQ只提供了单服务, 无法保证负载和HA能力. ***基于测试环境推荐使用***.
- 网桥多播发现模式下的ActiveMQ在多个MQ服务之间通过多播方式实现自发现加入集群, 集群内部通过桥接实现内部自通信, 保证一个MQ闲置时, 可以帮助繁忙的MQ一起消费队列, 实现负载均衡

---------------------------------------
> - **Alone Mode**

>> 1. download dependency ***apache-activemq-5.13.3-bin.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/amq.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait..., ActiveMQ Alone installed.
>> 4. start/stop/restart script in ***consolePath***
```
{
  "amq_pkg": "/veris/odc/install/amq_src/apache-activemq-5.13.3-bin.tar.gz",
  "amq_home": "/veris/odc/install/amq",
  "heap_size": 2048,
  "mq_name": "q1",
  "listenAddr": "0.0.0.0:61616",
  "consolePath": "/veris/odc/install/console",
  "adminName": "admin",
  "adminPWD": "admin",
  "userName": "user",
  "userPWD": "user",
  "webPort": "8161"
}
```
>> - ***amq_pkg*** is amq app pkg absolute path
>> - ***amq_home*** is amq be installed home path
>> - ***heap_size*** is amq is app install heap size
>> - ***mq_name*** is app group only name
>> - ***listenAddr*** is amq listen address[include ip:port]
>> - ***adminName*** is amq login web ui admin username
>> - ***adminPWD*** is amq is log we ui admin password
>> - ***userName*** is amq login web ui normal username
>> - ***userPWD*** is amq login web ui normal password
>> - ***consolePath*** is be installed amq, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后memcached占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

> - **Network Multicast Mode**

>> 1. download dependency ***apache-activemq-5.13.3-bin.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/amq.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait..., ActiveMQ Network Multicast Mode installed.
>> 4. start/stop/restart script in ***consolePath***
```
{
  "amq_pkg": "/veris/odc/install/amq_src/apache-activemq-5.13.3-bin.tar.gz",
  "amq_home": "/veris/odc/install/amq",
  "heap_size": 2048,
  "mq_name": "q1",
  "multicastAddr": "224.1.2.3:6255",
  "groupName": "humbird",
  "listenAddr": "0.0.0.0:61616",
  "consolePath": "/veris/odc/install/console",
  "adminName": "admin",
  "adminPWD": "admin",
  "userName": "user",
  "userPWD": "user",
  "webPort": "8161"
}
```
>> - ***amq_pkg*** is amq app pkg absolute path
>> - ***amq_home*** is amq be installed home path
>> - ***heap_size*** is amq is app install heap size
>> - ***mq_name*** is app group only name
>> - ***multicastAddr*** is amq multicast address[include ip:port]
>> - ***groupName*** is amq multicast group name
>> - ***listenAddr*** is amq listen address[include ip:port]
>> - ***adminName*** is amq login web ui admin username
>> - ***adminPWD*** is amq is log we ui admin password
>> - ***userName*** is amq login web ui normal username
>> - ***userPWD*** is amq login web ui normal password
>> - ***consolePath*** is be installed amq, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后memcached占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @activemq.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/install/activemq.json
```


# Zookeeper
- Zookeeper安装分为 Alone Mode 和 Cluster Mode

---------------------------------------
> - **Alone Mode**

>> 1. download dependency ***zookeeper-3.4.8.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/zk.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait..., Zookeeper Alone Mode installed.
>> 4. start/stop/restart script in ***consolePath***
```
{
  "zk_name" : "z1",
  "zk_pkg" : "/veris/odc/install/zookeeper/zookeeper-3.4.8.tar.gz",
  "zk_home" : "/veris/odc/install/zk",
  "dataDir" : "/veris/odc/install/zk/db",
  "clientPort" : "2181",
  "maxClientCnxns" : "60",
  "snapRetainCount" : "3",
  "purgeInterval" : "1",
  "consolePath" : "/veris/odc/install/console",
  "jvm" : "-Xmx2048M -Xms2048M -XX:MaxPermSize=128m -XX:PermSize=32m"
}
```
>> - ***zk_pkg*** is zookeeper app pkg absolute path
>> - ***zk_home*** is zookeeper be installed home path
>> - ***zk_name*** is zookeeper alias name
>> - ***dataDir*** is zookeeper data directory
>> - ***clientPort*** is zookeeper listen port
>> - ***maxClientCnxns*** is zookeeper allow connection clients
>> - ***snapRetainCount*** is zookeeper data snap counts
>> - ***purgeInterval*** is zookeeper purge snap interval
>> - ***jvm*** is zookeeper start jvm
>> - ***consolePath*** is be installed amq, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后memcached占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

> - **Cluster Mode**

>> 1. download dependency ***zookeeper-3.4.8.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/zk.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait..., Zookeeper Cluster Mode installed.
>> 4. start/stop/restart script in ***consolePath***
```
{
  "zk_name" : "z1",
  "zk_pkg" : "/veris/odc/install/zookeeper/zookeeper-3.4.8.tar.gz",
  "zk_home" : "/veris/odc/install/zk",
  "id" : 1,
  "dataDir" : "/veris/odc/install/zk/db",
  "clientPort" : "2181",
  "maxClientCnxns" : "60",
  "clusters" : [
    {
      "name" : "server.1",
      "address" : "10.1.245.185:2888:3888"
    }
  ],
  "snapRetainCount" : "3",
  "purgeInterval" : "1",
  "consolePath" : "/veris/odc/install/console",
  "jvm" : "-Xmx2048M -Xms2048M -XX:MaxPermSize=128m -XX:PermSize=32m"
}
```
>> - ***zk_pkg*** is zookeeper app pkg absolute path
>> - ***zk_home*** is zookeeper be installed home path
>> - ***zk_name*** is zookeeper alias name
>> - ***id*** is zookeeper cluster 单前节点的编号[需要与cluseters里面的address匹配]
>> - ***dataDir*** is zookeeper data directory
>> - ***clientPort*** is zookeeper listen port
>> - ***maxClientCnxns*** is zookeeper allow connection clients
>> - ***clusters*** is zookeeper cluster 节点详细
>>> - ***name*** is zookeeper cluster 某个节点的别名, 必须是 {server}.{编号} 格式
>>> - ***address*** is zookeeper cluster 某个节点的地址
>> - ***snapRetainCount*** is zookeeper data snap counts
>> - ***purgeInterval*** is zookeeper purge snap interval
>> - ***jvm*** is zookeeper start jvm
>> - ***consolePath*** is be installed amq, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后memcached占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @zk.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/install/zk.json
```


# Redis
- Redis作为缓存容器, 有多种模式, 安装步骤分为install和domain两个阶段

---------------------------------------
> - **Install**

>> 1. download dependency ***redis-3.2.0.tar.gz,ruby-2.3.1.tar.gz,redis-3.3.0.gem*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/redis.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait..., Redis installed.
```
{
  "redis_pkg": "/veris/odc/install/redis_src/redis-3.2.0.tar.gz",
  "redis_home": "/veris/odc/install/redis",
  "ruby_pkg": "/veris/odc/install/redis_src/ruby-2.3.1.tar.gz",
  "ruby_home": "/veris/odc/install/redis/ruby",
  "gems_pkg": "/veris/odc/install/redis_src/redis-3.3.0.gem",
  "isCluster": true
}
```
>> - ***redis_pkg*** is redis app pkg absolute path
>> - ***redis_home*** is redis be installed home path
>> - ***ruby_pkg*** is ruby app pkg absolute path
>> - ***ruby_home*** is ruby installed home path
>> - ***gems_pkg*** is gems app pkg absolute path
>> - ***isCluster*** is redis是否被安装成集群模式, 决定是否会使用 ruby/gems

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @redis.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/install/redis.json
```

---------------------------------------
> - **Domain Singleton**

>> 1. create ***${app}/conf/domain/redis.json*** (e.g: required)
>> 2. ***${app}/bin/domain***, wait..., Redis Singleton created.
>> 3. start/stop/restart script in ***consolePath***
```
{
  "redis_name": "redis1",
  "redis_base": "/veris/odc/install/redis",
  "redis_home": "/veris/odc/install/redis",
  "daemonize": true,
  "ips": [
    "10.1.245.185"
  ],
  "port": 17000,
  "dbDir": "/veris/odc/install/redis/17000/db",
  "maxClientCon": 4096,
  "maxMemory": "419430400",
  "type": "singleton",
  "isCluster": false,
  "consolePath": "/veris/odc/install/console",
  "option": "ADD"
}
```
>> - ***redis_name*** is redis alias
>> - ***redis_base*** is redis base app installed absolute path
>> - ***redis_home*** is redis 域配置文件/pid 所在路径
>> - ***daemonize*** is redis 启动是否使用后台进程启动
>> - ***ips*** is redis 启动监听的网卡地址列表, 是一个字符串数组
>> - ***port*** is redis 启动的时候监听的端口
>> - ***dbDir*** is redis 启动后的持久化数据 和 增量数据存储的地址
>> - ***maxClientCon*** is redis 允许最多客户端连接个数
>> - ***maxMemory*** is redis 启动时占用的最大内存, 但是是 byte
>> - ***type*** is redis 启动的模式, 分为ha/cluster/singleton, 这里请选择singleton
>> - ***isCluster*** is redis 是否启动cluster模式的状态标识, ha/singleton模式下请选择false
>> - ***option*** choose ADD install option, 暂时没有一键卸载功能
>> - ***consolePath*** is be installed redis, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后memcached占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

> - **Domain Sentinel HA**

>> 1. create ***${app}/conf/domain/redis.json*** (e.g: required)
>> 2. ***${app}/bin/domain***, wait..., Redis Sentinel HA created.
>> 3. start/stop/restart script in ***consolePath***
```
{
  "redis_name": "redis1",
  "redis_base": "/veris/odc/install/redis",
  "redis_home": "/veris/odc/install/redis",
  "daemonize": true,
  "ips": [
    "10.1.245.185"
  ],
  "port": 17000,
  "dbDir": "/veris/odc/install/redis/17000/db",
  "maxClientCon": 4096,
  "maxMemory": "419430400",
  "isCluster": false,
  "type": "ha",
  "sentinel": {
    "isMaster" : true,
    "ip": "10.1.245.185",
    "port": 17001,
    "Dir": "/veris/odc/install/redis/17000/db",
    "masterName": "amd1",
    "masterIP": "10.1.245.185",
    "masterPort": 17000,
    "parallelSyncs": 1,
    "downTime": "10000",
    "failoverTime": "180000",
    "quorum": 1
  },
  "consolePath": "/veris/odc/install/console",
  "option": "ADD"
}
```
>> - ***redis_name*** is redis alias
>> - ***redis_base*** is redis base app installed absolute path
>> - ***redis_home*** is redis 域配置文件/pid 所在路径
>> - ***daemonize*** is redis 启动是否使用后台进程启动
>> - ***ips*** is redis 启动监听的网卡地址列表, 是一个字符串数组
>> - ***port*** is redis 启动的时候监听的端口
>> - ***dbDir*** is redis 启动后的持久化数据 和 增量数据存储的地址
>> - ***maxClientCon*** is redis 允许最多客户端连接个数
>> - ***maxMemory*** is redis 启动时占用的最大内存, 但是是 byte
>> - ***type*** is redis 启动的模式, 分为ha/cluster/singleton, 这里请选择ha
>> - ***isCluster*** is redis 是否启动cluster模式的状态标识, ha/singleton模式下请选择false
>> - ***sentinel*** is redis sentinel 进程的配置, sentinel实际上是redis 管理HA的一个进程, 需要额外进程和配置
>>> - ***isMaster*** is redis sentinel 当前redis进程是否作为master
>>> - ***ip*** is redis sentinel 进程监听的网卡地址
>>> - ***port*** is redis sentinel 进程监听的端口
>>> - ***Dir*** is redis sentinel 数据存储的目录
>>> - ***masterName*** is redis sentinel master节点的alias 别名, 各master/slave需要一致
>>> - ***masterIP*** is redis sentinel master节点所在的网口IP
>>> - ***masterPort*** is redis sentinel master节点所在的IP端口
>>> - ***parallelSyncs*** is redis sentinel ha模式下执行故障转移时, 最多可以从多少slave服务向新的master节点同步数据, 值越大, 同步越快, 但是在同步的slave节点暂时不可用
>>> - ***donwTime*** is redis sentinel ha模式下指定认为服务器已经断线所需要的毫秒数
>>> - ***failoverTime*** is redis sentinel ha模式下failover确认的时间
>>> - ***quorum*** is redis sentinel ha模式下服务集群做投票选举最小赞成票数
>> - ***option*** choose ADD install option, 暂时没有一键卸载功能
>> - ***consolePath*** is be installed redis, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后memcached占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

> - **Domain Cluster**

>> 1. create ***${app}/conf/domain/redis.json*** (e.g: required)
>> 2. ***${app}/bin/domain***, wait..., Redis Cluster created.
>> 3. start/stop/restart script in ***consolePath***
```
{
  "redis_name": "redis1",
  "redis_base": "/veris/odc/install/redis",
  "redis_home": "/veris/odc/install/redis",
  "daemonize": true,
  "ips": [
    "10.1.245.185"
  ],
  "port": 17000,
  "dbDir": "/veris/odc/install/redis/17000/db",
  "maxClientCon": 4096,
  "maxMemory": "419430400",
  "type": "cluster",
  "isCluster": true,
  "consolePath": "/veris/odc/install/console",
  "option": "ADD"
}
```
>> - ***redis_name*** is redis alias
>> - ***redis_base*** is redis base app installed absolute path
>> - ***redis_home*** is redis 域配置文件/pid 所在路径
>> - ***daemonize*** is redis 启动是否使用后台进程启动
>> - ***ips*** is redis 启动监听的网卡地址列表, 是一个字符串数组
>> - ***port*** is redis 启动的时候监听的端口
>> - ***dbDir*** is redis 启动后的持久化数据 和 增量数据存储的地址
>> - ***maxClientCon*** is redis 允许最多客户端连接个数
>> - ***maxMemory*** is redis 启动时占用的最大内存, 但是是 byte
>> - ***type*** is redis 启动的模式, 分为ha/cluster/singleton, 这里请选择cluster
>> - ***isCluster*** is redis 是否启动cluster模式的状态标识, ha/singleton模式下请选择true
>> - ***option*** choose ADD install option, 暂时没有一键卸载功能
>> - ***consolePath*** is be installed redis, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后redis占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @redis.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/domain/redis
```


# Storm
- Storm安装分为supervisor cluster, nimbus HA模式

---------------------------------------
> - **Install Cluster And Nimbus HA**
>> 1. download dependency ***apache-storm-1.0.1.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/storm.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait..., Storm installed.
>> 4. ***${consolePath}***目录下会有2种脚本, nimbus和supervisor.[启动顺序是 nimbus->supervisor]
```
{
  "storm_name" : "storm1",
  "storm_pkg": "/veris/odc/install/storm_src/apache-storm-1.0.1.tar.gz",
  "storm_home": "/veris/odc/install/storm",
  "zks": [
    "10.1.234.149"
  ],
  "zkPort": 9201,
  "zkDir": "/teststorm",
  "mode": "distributed",
  "nimbusHA": [
    "10.1.245.185",
    "10.1.245.186"
  ],
  "stormDataDir": "/veris/odc/install/storm/data",
  "stormLogDir": "/veris/odc/install/storm/logs",
  "slotsPorts": [
    16701,
    16702
  ],
  "workerHeap": 512,
  "workerJVM": "-Xmx512M -Xms512M -XX:MaxPermSize=128m -XX:PermSize=32m",
  "uiHost": "10.1.245.185",
  "uiPort": 16710,
  "uiJVM": "-Xmx256M -Xms256M -XX:MaxPermSize=128m -XX:PermSize=32m",
  "topoMSGTime": 60,
  "rmLibs": [
    "log4j-over-slf4j-1.6.6.jar"
  ],
  "stormYaml": "",
  "option": "ADD",
  "consolePath": "/veris/odc/install/console"
}
```
>> - ***storm_name*** is storm alias
>> - ***storm_pkg*** is storm app pkg absolute path
>> - ***storm_home*** is storm be installed home path
>> - ***zks*** is storm 依赖 zk 的 IP 数组[storm这里很烂, 只能指定ip, 端口不能有差异]
>> - ***zkPort*** is storm 依赖 zk 的 Port
>> - ***zkDir*** is storm 依赖 zk, 将 metadata 存储在 zk 的指定目录上
>> - ***mode*** is storm 启动的模型[distributed][这里只允许 distributed]
>> - ***nimbusHA*** is storm 集群的master地址, 是一个数组[集群都有多个master地址]
>> - ***stormDataDir*** is storm 数据存储的目录
>> - ***stormLogDir*** is storm 日志存储的目录
>> - ***slotsPorts*** is storm worker进程启动的个数[每个worker对应一个端口, 要启动多少个worker就规划多少个端口]
>> - ***workerHeap*** is storm worker启动分配的堆大小
>> - ***workerJVM*** is storm worker启动分配的jvm参数
>> - ***uiHost*** is storm UI 监听的地址
>> - ***uiPort*** is storm UI 监听的端口
>> - ***uiJVM*** is storm UI 启动分配的jvm参数
>> - ***topoMSGTime*** is storm 停止的时候等待topo停止的时间
>> - ***rmLibs*** is storm 安装的时候需要删除默认的lib包, 主要是解决日志包有冲突的问题
>> - ***stormYaml*** is storm 使用自定义的配置文件, 而不用内置的参数, 如果这个有值, 之前填写有关配置项的参数都是无效的
>> - ***option*** choose ADD install option, 暂时没有一键卸载功能
>> - ***consolePath*** is be installed storm, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后storm占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @storm
.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/install/storm.json
```


# MySQL
- MySQL作为关系型数据库, 安装步骤分为install和domain两个阶段
- Install只会安装MySQL基础的组件
- Domain基于MySQL基础组件构建独立的实例, 默认生成的是调优后InnoDB类型

---------------------------------------
> - **Install**

>> 1. download dependency ***mysql-5.6.31.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/mysql.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait..., MySQL installed.
```
{
  "name": "mysql-5.6.31",
  "mysql_pkg": "/veris/odc/install/src/mysql-5.6.31.tar.gz",
  "mysql_home": "/veris/odc/mysql"
}
```
>> - ***name*** is mysql alias
>> - ***mysql_pkg*** is mysql app pkg absolute path
>> - ***mysql_home*** is mysql be installed home path

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @tomcat.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/install/mysql.json
```

---------------------------------------

> - **Domain**

>> 1. create ***${app}/conf/domain/mysql_xxx.json*** (e.g: required, 只要是mysql开头的文件即可)
>> 2. ***${app}/bin/domain***, wait..., MySQL domain created.
>> 3. start/stop/restart script in ***consolePath***
```
{
  "mysql_name" : "o2p",
  "mysql_home" : "/veris/odc/mysql",
  "mysql_base" : "/veris/odc/mysql",
  "port" : 3401,
  "consolePath" : "/veris/odc/mysql/3401/db/scripts",
  "option" : "ADD"
}
```
>> - ***mysql_name*** is mysql domain alias
>> - ***mysql_home*** is mysql be installed home path
>> - ***mysql_base*** is mysql domain 安装的路径
>> - ***port*** is mysql domain 启动后监听的端口
>> - ***option*** choose ADD install option, 暂时没有一键卸载功能
>> - ***consolePath*** is be installed mysql domain, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后 mysql domain占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @mysql_3401.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/domain/mysql
```


# Tomcat
- Tomcat作为web容器, 我们只依赖于它启动我们的应用, 安装步骤分为install和domain两个阶段
- Install只会安装tomcat基础的组件, 涉及tomcat, apr, jdk, openssl.(e.g.: apr的安装保证作为server, 可以对外提供高效的apr高性能服务, 默认tomcat为bio模式, 可以直接选择nio模式, apr为第三种可选的高效模式, 基于apache, 但是需要手动编译和安装, 取决于操作系统编译环境)
- Domain会构建server域, 一个域下可以安装多个server.并且分为BIO, NIO,APR三种模式

---------------------------------------
> - **Install**

>> 1. download dependency ***apache-tomcat-7.0.69.tar.gz, apr-1.5.2.tar.gz, apr-util-1.5.4.tar.gz, openssl-1.0.1t.tar.gz, tomcat-native.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/tomcat.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait..., Redis installed.
```
{
    "tomcat_home":"/veris/odc/app/tomcat7",
    "pkg":"/veris/odc/install/tomcat7/apache-tomcat-7.0.69.tar.gz",
    "isRemove":false,
    "tomcat_native_pkg":"/veris/odc/install/tomcat7/tomcat-native.tar.gz",
    "apr_pkg":"/veris/odc/install/tomcat7/apr-1.5.2.tar.gz",
    "apr_util_pkg":"/veris/odc/install/tomcat7/apr-util-1.5.4.tar.gz",
    "openssl_pkg":"/veris/odc/install/tomcat7/openssl-1.0.1t.tar.gz",
    "apr_home":"/veris/odc/install/tomcat-native/apr",
    "apr_util_home":"/veris/odc/install/tomcat-native/apr-util",
    "openssl_home":"/veris/odc/install/tomcat-native/openssl",
    "tomcat_native_home":"/veris/odc/install/tomcat-native",
    "java_home":"/veris/odc/app/jdk"
}
```
>> - ***tomcat_home*** is tomcat be installed home path
>> - ***pkg*** is tomcat app pkg absolute path
>> - ***isRemove*** is tomcat removed after installed
>> - ***tomcat_native_pkg*** is tomcat_native app pkg absolute path
>> - ***apr_pkg*** is apr app pkg absolute path
>> - ***apr_util_pkg*** is apr-util app pkg absolute path
>> - ***openssl_pkg*** is openssl app pkg absolute path
>> - ***apr_home*** is apr be installed home path
>> - ***apr_util_home*** is apr-util be installed home path
>> - ***openssl_home*** is openssl be installed home path
>> - ***tomcat_native_home*** is tomcat_native be installed home path
>> - ***java_home*** is ${JAVA_HOME} path

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @tomcat.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/install/tomcat.json
```

---------------------------------------

> - **Domain**

>> 1. create ***${app}/conf/domain/tomcat_xxx.json*** (e.g: required, 只要是tomcat开头的文件即可)
>> 2. ***${app}/bin/domain***, wait..., Tomcat domain created.
>> 3. start/stop/restart script in ***consolePath***
```
{
  "jdkHome": "/veris/odc/app/jdk",
  "tomcatHome": "/veris/odc/app/tomcat7",
  "nativeHome": "/veris/odc/install/tomcat-native",
  "protocol": "org.apache.coyote.http11.Http11AprProtocol",
  "sharedLoader": "/veris/odc/properties",
  "domainPath": "/veris/odc/install",
  "aliasName": "cache",
  "apps" : [
    {
      "appName" : "cache",
      "app_home" : "/veris/odc/install/webapps/cache"
    }, {
      "appName" : "task",
      "app_home" : "/veris/odc/install/webapps/task"
    },
  ],
  "servers": [
    {
      "version": "1",
      "listenPort": "8315",
      "shutdownPort": "8316"
    },
    {
      "version": "2",
      "listenPort": "8317",
      "shutdownPort": "8318"
    }
  ],
  "envs" : [
    {
        "name": "PROPERTY_NAME_SPACE",
        "value": "/cfg/DMZ/T/1/agfpdmz"
    },
    {
        "name": "CONNECT_STRING",
        "value": "10.1.236.154:9201,10.1.236.155:9201,10.1.236.155:9201"
    }
  ],
  "timeout": "10",
  "jvm": "-Xmx512m -Xms512m -XX:MaxPermSize=128m -XX:PermSize=128m -Xmn128M",
  "option": "ADD",
  "consolePath": "/veris/odc/install/console"
}
```
>> - ***jdkHome*** is ${JAVA_HOME} path
>> - ***tomcatHome*** is tomcat base app installed absolute path
>> - ***nativeHome*** is tomcat-native base app installed absolute path
>> - ***protocol*** is tomcat 启动的模式, 可选3种模式:
>>> - ***BIO模式*** HTTP/1.1
>>> - ***NIO模式*** org.apache.coyote.http11.Http11NioProtocol
>>> - ***APR模式*** org.apache.coyote.http11.Http11AprProtocol
>> - ***sharedLoader*** is tomcat domain 应用配置文件所在目录
>> - ***domainPath*** is tomcat domain 域所在的根目录[我们约定一个域下面有多个server, 每个服务同事提供多种app]
>> - ***aliasName*** is tomcat domain 下面要安装的server别名, 一份配置文件只能安装一个server
>> - ***apps*** is tomcat domain 下面安装的server 启动后能提供的 app 应用[一个/context就是一个应用], 旗下是一个数组
>>> - ***appName*** 对外提供 app 应用的 context 上下文
>>> - ***app_home*** 程序应用[这里一般约定是war包解压后的目录]
>> - ***servers*** is 需要基于以上配置, 安装***多少套相同配置***的服务, 本身是数组[会在${domainPath}/${aliasName}/${version}构建对应配置]
>>> - ***version*** is 在同一个 ${domainPath}/${aliasName}/${version} 来区别多个同类服务的目录, 适用于scale out
>>> - ***listenPort*** is tomcat domain server app listen port
>>> - ***shutdownPort*** is tomcat domain server app listen shutdown port[建议是listenPort + 1]
>> - ***envs*** is 启动的 tomcat domain 要临时生效的用户环境变量
>>> - ***name*** is 相当于 export JAVA_HOME=/hihi 的 ***JAVA_HOME***
>>> - ***value*** is 相当于 export JAVA_HOME=/hihi 的 ***/hihi***
>> - ***timeout*** is tomcat domain server停止之前等待的时长
>> - ***jvm*** is tomcat domain server启动的时候的jvm参数
>> - ***option*** choose ADD install option, 暂时没有一键卸载功能
>> - ***consolePath*** is be installed tomcat domain, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后tomcat domain占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @tomcat_cache.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/domain/tomcat
```


# Weblogic
- weblogic作为web容器，我们依赖于它启动我们的应用，安装步骤分为install, domain, server, deploy四个阶段度.
- 和tomcat不同, tomcat需要人为的规划域, weblogic设计上就是通过domain划分域, 域下面可以启动多个server. 而每个server需要启动哪些app根据deploy来取决, 但是每个server下面的app共享server的jvm环境变量.
- 创建顺序是install->domain->server->deploy.

---------------------------------------
> - **Install**

>> 1. download dependency ***fmw_12.1.3.0.0_wls.jar*** (e.g.: 这个可以根据开源社区下载最新版, 暂时***只支持12c版本***), 并且手动生成wls12c.rsp文件, 文件内容如下
```
[ENGINE]
#DO NOT CHANGE THIS.
Response File Version=1.0.0.0.0
[GENERIC]
#The oracle home location. This can be an existing Oracle Home or a new Oracle Home
ORACLE_HOME=/veris/odc/app/weblogic
#Set this variable value to the Installation Type selected. e.g. WebLogic Server, Coherence, Complete with Examples.
INSTALL_TYPE=WebLogic Server
#Provide the My Oracle Support Username. If you wish to ignore Oracle Configuration Manager configuration provide empty string for user name.
MYORACLESUPPORT_USERNAME=
#Provide the My Oracle Support Password
MYORACLESUPPORT_PASSWORD=<SECURE VALUE>
#Set this to true if you wish to decline the security updates. Setting this to true and providing empty string for My Oracle Support username will ignore the Oracle Configuration Manager configuration
DECLINE_SECURITY_UPDATES=true
#Set this to true if My Oracle Support Password is specified
SECURITY_UPDATES_VIA_MYORACLESUPPORT=false
#Provide the Proxy Host
PROXY_HOST=
#Provide the Proxy Port
PROXY_PORT=
#Provide the Proxy Username
PROXY_USER=
#Provide the Proxy Password
PROXY_PWD=<SECURE VALUE>
#Type String (URL format) Indicates the OCM Repeater URL which should be of the format [scheme[Http/Https]]://[repeater host]:[repeater port]
COLLECTOR_SUPPORTHUB_URL=
```
>>> - ***ORACLE_HOME*** 需要修改, 表示weblogic被安装的目录地址
>>> - 由于版本早起没有考虑到这块在应用内部生成, 需要使用者手工生成
>> 2. create ***${app}/conf/install/wls12c.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait..., MySQL installed.
```
{
    "jarLoc":"/veris/odc/install/wls12c/fmw_12.1.3.0.0_wls.jar",
    "invLoc":"/veris/odc/app/oraInventory",
    "rspFile":"/veris/odc/install/wls12c/wls12c.rsp",
    "jdk_home":"/veris/odc/app/jdk"
}
```
>> - ***jarLoc*** is wls12c app pkg absolute path
>> - ***invLoc*** is wls12c安装后的一些patch依赖目录
>> - ***jarLoc*** is wls12c静默安装需要的配置文件, 参考步骤1生成的文件目录
>> - ***jarLoc*** is ${JAVA_HOME} path

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @wls12c.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/install/wls12c.json
```

---------------------------------------

> - **Domain**

>> 1. create ***${app}/conf/domain/wls12c_xxx.json*** (e.g: required, 只要是 wls12c 开头的文件即可)
>> 2. ***${app}/bin/domain***, wait..., Wls12c domain created.
>> 3. start/stop/restart script in ***consolePath***(e.g: start_admin.sh/restart_admin.sh/stop_admin.sh)
```
{
    "wlsHome":      "/veris/odc/app/weblogic",
    "listenAddr":   "",
    "listenPort":   "7001",
    "userName":     "dev_odc",
    "passWord":     "69@Ailk!@#",
    "mode":         "prod",
    "jdkHome":      "/veris/odc/app/jdk",
    "domainPath":   "/veris/odc/dev_domain",
    "option":       "ADD",
    "consolePath":  "/veris/odc/console"
}
```
>> - ***jdkHome*** is ${JAVA_HOME} path
>> - ***wlsHome*** is wls12c base app installed absolute path
>> - ***domainPath*** is wls12c domain installed absolute path
>> - ***listenAddr*** is wls12c domain 被安装后, admin端点监听的 addr, weblogic有admin/servers的概念. admin是主控. 默认为空, 监听是0.0.0.0
>> - ***listenPort*** is wls12c domain 被安装后, admin端点监听的 port
>> - ***wlsHome*** is wls12c base app installed absolute path
>> - ***userName*** is wls12c domain 被安装后, admin端点web界面登录的管理员用户名
>> - ***passWord*** is wls12c domain 被安装后, admin端点web界面登录的管理员密码
>> - ***mode*** is wls12c domain安装时候选择的模型, 有生产模式(prod)和开发模式(这个模式我也忘了是啥变量了)
>> - ***option*** choose ADD install option, 暂时没有一键卸载功能
>> - ***consolePath*** is be installed tomcat domain, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后tomcat domain占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @wls12c.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/donmain/wls12c
```

---------------------------------------

> - **Server**

>> 1. create ***${app}/conf/server/xxx.json*** (e.g: required, 根据规划创建服务配置文件)
>> 2. ***${app}/bin/server***, wait..., Wls12c server created.
>> 3. start/stop/restart script in ***consolePath***(e.g: start_xxx.sh/restart_xxx.sh/stop_xxx.sh)
```
{
  "wlsHome": "/data/odc/app/weblogic",
  "adminAddr": "10.1.234.149",
  "adminPort": "7001",
  "userName": "dev_odc",
  "passWord": "69@Ailk!@#",
  "listenAddr": "10.1.234.149",
  "listenPort": "8135",
  "srvName": "cache",
  "margs": "-Xmx512m -Xms512m -XX:PermSize=256m -XX:MaxPermSize=256m",
  "jars": [
    "$HOME/properties/jar/resoure.jar"
  ],
  "domainPath": "/data/odc/dev_domain",
  "option": "ADD",
  "consolePath": "/data/odc/console"
}
```
>> - ***wlsHome*** is wls12c base app installed absolute path
>> - ***domainPath*** is wls12c domain installed absolute path
>> - ***adminAddr*** is wls12c domain 被安装后, admin端点监听的 addr, weblogic有admin/servers的概念. admin是主控. 默认为空, 监听是0.0.0.0
>> - ***adminPort*** is wls12c domain 被安装后, admin端点监听的 port
>> - ***userName*** is wls12c domain 被安装后, admin端点web界面登录的管理员用户名
>> - ***passWord*** is wls12c domain 被安装后, admin端点web界面登录的管理员密码
>> - ***listenAddr*** is wls12c server 被安装后, server 监听的地址
>> - ***listenPort*** is wls12c server 被安装后, server 监听的端口
>> - ***srvName*** is wls12c server 被安装后, server 的 alias 别名
>> - ***margs*** is wls12c server 被安装后, server 启动时 分配的 jvm 参数
>> - ***jars*** is wls12c server 被安装后, 启动时要导入的 配置第三方jar包
>> - ***envs*** 本版本暂时不支持
>> - ***option*** choose ADD install option, 暂时没有一键卸载功能
>> - ***consolePath*** is be installed tomcat domain, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后tomcat domain占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

---------------------------------------
**Restful API:** 暂时不支持
```
```

---------------------------------------

> - **Deploy**

>> 1. create ***${app}/conf/deploy/xxx.json*** (e.g: required, 根据规划创建服务配置文件)
>> 2. ***${app}/bin/deploy***, wait..., Wls12c app created.
>> 3. start/stop/restart script in ***consolePath***(e.g: 共deploy/undeploy/redeploy/application_start/application_stop五个目录. 目录下脚本分别以pkgName找到即可)
```
{
    "wlsHome":      "/veris/odc/app/weblogic",
    "adminAddr":    "10.1.245.185",
    "adminPort":    "7001",
    "userName":     "dev_odc",
    "passWord":     "69@Ailk!@#",
    "srvName":      "cache",
    "pkgName":      "cache",
    "pkgPath":      "/veris/odc/webapps/wls/cache",
    "option":       "ADD",
    "consolePath":  "/veris/odc/console"
}
```
>> - ***wlsHome*** is wls12c base app installed absolute path
>> - ***adminAddr*** is wls12c domain 被安装后, admin端点监听的 addr, weblogic有admin/servers的概念. admin是主控. 默认为空, 监听是0.0.0.0
>> - ***adminPort*** is wls12c domain 被安装后, admin端点监听的 port
>> - ***userName*** is wls12c domain 被安装后, admin端点web界面登录的管理员用户名
>> - ***passWord*** is wls12c domain 被安装后, admin端点web界面登录的管理员密码
>> - ***srvName*** is wls12c server 被安装后, server 的 alias 别名
>> - ***pkgName*** is wls12c server 被安装后, app 的 alias 别名, 一个 server 下可以有多个 app
>> - ***pkgPath*** is 基于 wls12c server 使用的 war 包的绝对路径, 需要自己 解压缩包.
>> - ***option*** choose ADD install option, 暂时没有一键卸载功能
>> - ***consolePath*** is be installed tomcat domain, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后tomcat domain占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

---------------------------------------
**Restful API:** 暂时不支持
```
```