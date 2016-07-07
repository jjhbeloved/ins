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
> - **Storm**
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
>> 4. ***${app}/bin/install***, wait...，JDK installed.
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
>> 4. ***${app}/bin/install***, wait...，JDK installed.
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
>> 3. ***${app}/bin/install***, wait...，Memcached Alone installed.
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
>> 3. ***${app}/bin/install***, wait...，Memcached Alone installed.
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
- 网桥多播发现模式下的ActiveMQ在多个MQ服务之间通过多播方式实现自发现加入集群, 集群内部通过桥接实现内部自通信，保证一个MQ闲置时，可以帮助繁忙的MQ一起消费队列，实现负载均衡

---------------------------------------
> - **Alone Mode**

>> 1. download dependency ***apache-activemq-5.13.3-bin.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/amq.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait...，ActiveMQ Alone installed.
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
>> 3. ***${app}/bin/install***, wait...，ActiveMQ Network Multicast Mode installed.
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
>> 3. ***${app}/bin/install***, wait...，Zookeeper Alone Mode installed.
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
>> 3. ***${app}/bin/install***, wait...，Zookeeper Cluster Mode installed.
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
>> 3. ***${app}/bin/install***, wait...，Redis installed.
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
>> 2. ***${app}/bin/domain***, wait...，Redis Singleton created.
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
>> 2. ***${app}/bin/domain***, wait...，Redis Sentinel HA created.
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
>> 2. ***${app}/bin/domain***, wait...，Redis Cluster created.
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
>> - ***consolePath*** is be installed redis, 自动启停脚本生产的路径, 此路径下会有start/restart/stop三个目录, 启动后memcached占用的资源是通过启停脚本自动生成的, 如果不通过启停脚本启动, 资源需要自己重新分配

---------------------------------------
**Restful API:**
```
curl -X POST --data-binary @redis.json -H "Content-Type:application/json;charset=UTF-8" http://localhost:8080/v1/domain/redis
```


# Tomcat
- Tomcat作为web容器，我们只依赖于它启动我们的应用，安装步骤分为install和domain两个阶段
- Install只会安装tomcat基础的组件，涉及tomcat, apr, jdk, openssl.(e.g.: apr的安装保证作为server, 可以对外提供高效的apr高性能服务, 默认tomcat为bio模式, 可以直接选择nio模式, apr为第三种可选的高效模式, 基于apache, 但是需要手动编译和安装, 取决于操作系统编译环境)
- Domain会构建server域，一个域下可以安装多个server.并且分为BIO, NIO,APR三种模式


---------------------------------------
> - **Install**

>> 1. download dependency ***apache-tomcat-7.0.69.tar.gz, apr-1.5.2.tar.gz, apr-util-1.5.4.tar.gz, openssl-1.0.1t.tar.gz, tomcat-native.tar.gz*** (e.g.: 这个可以根据开源社区下载最新版)
>> 2. create ***${app}/conf/install/tomcat.json*** (e.g: required)
>> 3. ***${app}/bin/install***, wait...，Redis installed.
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


