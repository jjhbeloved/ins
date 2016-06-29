package redis

import (
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Redis struct {
	Redis_Name     string        `json:"redis_name"`
	Redis_BASE     string        `json:"redis_base"`
	Redis_HOME     string        `json:"redis_home"`
	Daemonize      bool          `json:"daemonize"`
	IPs            []string      `json:"ips"`
	Port           int           `json:"port"`
	DbDir          string        `json:"dbDir"`
	SentinelServer Sentinel      `json:"sentinel"`
	IsCluster      bool          `json:"isCluster"`
	Type           string        `json:"type"`
	Masters        []RedisServer `json:"masters"`
	Slaves         []RedisServer `json:"slaves"`
	MaxClientCon   int           `json:"maxClientCon"`
	MaxMemory      string        `json:"maxMemory"`
	ConsolePath    string        `json:"consolePath"`
	Option         string        `json:"option"`
}

type RedisServer struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Flag bool   `json:"flag"`
}

type Sentinel struct {
	IsMaster      bool   `json:"isMaster"`
	IP            string `json:"ip"`
	Port          int    `json:"port"`
	Dir           string `json:"dir"`
	MasterName    string `json:"masterName"`
	MasterIP      string `json:"masterIP"`
	MasterPort    int    `json:"masterPort"`
	ParallelSyncs int    `json:"parallelSyncs"`
	DownTime      string `json:"downTime"`
	FailoverTime  string `json:"failoverTime"`
	Quorum        int    `json:"quorum"`
}

func (w *Redis) Json(bs []byte) error {
	return json.Unmarshal(bs, &w)
}

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Method */

func (redis *Redis) Add() error {
	redis.touchConf()
	redis.touchConsoleScript()
	return nil
}

func (redis *Redis) Remove() error {
	return nil
}

/* End Method */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Configuration File */
const redisConf = `########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################

################################## INCLUDES ###################################

# include /path/to/local.conf
# include /path/to/other.conf

################################## NETWORK #####################################

bind %s
protected-mode yes
port %d
tcp-backlog 511
# Close the connection after a client is idle for N seconds (0 to disable)
timeout 0
tcp-keepalive 60

################################ GENERAL  ################################
daemonize %s
supervised no
pidfile %s
loglevel notice
logfile "%s"
# syslog-enabled no
# syslog-ident redis
# Specify the syslog facility. Must be USER or between LOCAL0-LOCAL7.
# syslog-facility local0
databases 16

################################ SNAPSHOTTING  ################################

save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename %s
dir %s

################################# REPLICATION #################################

# slaveof <masterip> <masterport>
%s
# masterauth <master-password>
slave-serve-stale-data yes
slave-read-only yes

# 1) Disk-backed: The Redis master creates a new process that writes the RDB
#                 file on disk. Later the file is transferred by the parent
#                 process to the slaves incrementally.
# 2) Diskless: The Redis master creates a new process that directly writes the
#              RDB file to slave sockets, without touching the disk at all.
# With slow disks and fast (large bandwidth) networks, diskless replication
# works better.
repl-diskless-sync no
repl-diskless-sync-delay 5
repl-ping-slave-period 10
repl-timeout 60
# If you select "no" the delay for data to appear on the slave side will
# be reduced but more bandwidth will be used for replication.
# By default we optimize for low latency, but in very high traffic conditions
# or when the master and slaves are many hops away, turning this to "yes" may
# be a good idea.
repl-disable-tcp-nodelay yes
# 设置repl-backlog-size 64mb. 默认值是1M, 当写入量很大时, backlog溢出会导致增量复制不成功
repl-backlog-size 64mb
repl-backlog-ttl 3600
# 适用Sentinel, 当master失效后, Sentinel将会从slave列表中找到权重值最低(>0)的slave, 并提升为master
slave-priority 100

# By default min-slaves-to-write is set to 0 (feature disabled) and
# min-slaves-max-lag is set to 10.
min-slaves-to-write 0
min-slaves-max-lag 10

################################## SECURITY ###################################

# requirepass foobared
# rename-command CONFIG ""

################################### LIMITS ####################################

maxclients %d
maxmemory %s
# volatile-lru -> remove the key with an expire set using an LRU algorithm
# allkeys-lru -> remove any key according to the LRU algorithm
# volatile-random -> remove a random key with an expire set
# allkeys-random -> remove a random key, any key
# volatile-ttl -> remove the key with the nearest expire time (minor TTL)
# noeviction -> don't expire at all, just return an error on write operations
maxmemory-policy allkeys-lru
# The default of 5 produces good enough results. 10 Approximates very closely
# true LRU but costs a bit more CPU. 3 is very fast but not very accurate.
maxmemory-samples 5

############################## APPEND ONLY MODE ###############################
#  但是这样会造成 appendonly.aof 文件过大，所以 redis 还支持了 BGREWRITEAOF 指令，对 appendonly.aof 进行重新整理。
#  你可以同时开启 asynchronous dumps 和  AOF
appendonly yes
appendfilename appendonly.aof
# no: don't fsync, just let the OS flush the data when it wants. Faster.
# always: fsync after every write to the append only log. Slow, Safest.
# everysec: fsync only one time every second. Compromise.
appendfsync everysec
# If you have latency problems turn this to "yes". Otherwise leave it as
# "no" that is the safest pick from the point of view of durability.
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
aof-load-truncated yes

################################ LUA SCRIPTING  ###############################

lua-time-limit 5000

################################ REDIS CLUSTER  ###############################
#
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
# WARNING EXPERIMENTAL: Redis Cluster is considered to be stable code, however
# in order to mark it as "mature" we need to wait for a non trivial percentage
# of users to deploy it in production.
# ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
cluster-enabled %s
cluster-config-file nodes.conf
cluster-node-timeout 15000
# 设为0, 从节点会一直尝试启动FailOver.
# 设为正数,失联大于一定时间(factor*节点TimeOut), 不再进行FailOver.
cluster-slave-validity-factor 0
# Default is 1 (slaves migrate only if their masters remain with at least
# one slave). To disable migration just set it to a very large value.
# A value of 0 can be set but is useful only for debugging and dangerous
# in production.
cluster-migration-barrier 1
# 设置为no, 集群丢失Key的情况下仍提供查询服务
cluster-require-full-coverage no

################################## SLOW LOG ###################################
# The following time is expressed in microseconds, so 1000000 is equivalent
# to one second. Note that a negative number disables the slow log, while
# a value of zero forces the logging of every command.
slowlog-log-slower-than 10000
slowlog-max-len 128

################################ LATENCY MONITOR ##############################

# By default latency monitoring is disabled since it is mostly not needed
# if you don't have latency issues, and collecting data has a performance
# impact, that while very small, can be measured under big load. Latency
# monitoring can easily be enabled at runtime using the command
# "CONFIG SET latency-monitor-threshold <milliseconds>" if needed.
latency-monitor-threshold 0
notify-keyspace-events ""

############################### ADVANCED CONFIG ###############################

# Hashes are encoded using a memory efficient data structure when they have a
# small number of entries, and the biggest entry does not exceed a given
# threshold. These thresholds can be configured using the following directives.
hash-max-ziplist-entries 512
hash-max-ziplist-value 64

# The highest performing option is usually -2 (8 Kb size) or -1 (4 Kb size),
# but if your use case is unique, adjust the settings as necessary.
list-max-ziplist-size -2
list-compress-depth 0

set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64

hll-sparse-max-bytes 3000
# 影响实时性
activerehashing yes
# 解决 output-buffer RAM 超过 client request 导致 OOM
# client-output-buffer-limit normal 0 0 0
client-output-buffer-limit normal 256mb 128mb 60
# client-output-buffer-limit slave 256mb 64mb 60
client-output-buffer-limit slave 512mb 256mb 180
client-output-buffer-limit pubsub 32mb 8mb 60
hz 10
aof-rewrite-incremental-fsync yes
`

const redisSentinelConf = `########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
port %d
sentinel announce-ip %s
dir %s
sentinel monitor %s %s %d %d
sentinel failover-timeout %s %s
sentinel parallel-syncs %s %d
sentinel down-after-milliseconds %s %s
`

/**
 * touch conf local file
 */
func (redis *Redis) touchConf() {
	os.MkdirAll(utils.TMPD, 0750)
	_redisHome := filepath.Join(redis.Redis_HOME, strconv.Itoa(redis.Port))
	_redisLogDir := filepath.Join(_redisHome, "logs")
	os.MkdirAll(_redisHome, 0755)
	os.MkdirAll(_redisLogDir, 0755)
	os.MkdirAll(redis.DbDir, 0755)
	_redisConf := filepath.Join(_redisHome, "redis.conf")
	_sentinelConf := filepath.Join(_redisHome, "sentinel.conf")
	_type := redis.Type
	now := time.Now().String()
	var ips string
	for _, ip := range redis.IPs {
		ips += ip + " "
	}
	var _daemonize, _clusterEnabled string
	slaveof := ""
	if redis.Daemonize {
		_daemonize = "yes"
	} else {
		_daemonize = "no"
	}
	if redis.IsCluster {
		_clusterEnabled = "yes"
	} else {
		_clusterEnabled = "no"
	}
	switch _type {
	case "ha":
		_sentinelServer := redis.SentinelServer
		if !_sentinelServer.IsMaster {
			slaveof = fmt.Sprintf("slaveof %s %d", _sentinelServer.MasterIP, _sentinelServer.MasterPort)
		}
		logs.Print(ioutil.WriteFile(
			_sentinelConf,
			[]byte(fmt.Sprintf(redisSentinelConf,
				now,
				_sentinelServer.Port, _sentinelServer.IP,
				_sentinelServer.Dir,
				_sentinelServer.MasterName, _sentinelServer.MasterIP, _sentinelServer.MasterPort, _sentinelServer.Quorum,
				_sentinelServer.MasterName, _sentinelServer.FailoverTime,
				_sentinelServer.MasterName, _sentinelServer.ParallelSyncs,
				_sentinelServer.MasterName, _sentinelServer.DownTime,
			)),
			0750,
		))
	case "cluster":
		break
	case "singleton":
		break
	default:

	}
	// _redisConf
	logs.Print(ioutil.WriteFile(
		_redisConf,
		[]byte(fmt.Sprintf(redisConf,
			now,
			ips, redis.Port, _daemonize, filepath.Join(_redisHome, redis.Redis_Name+".pid"), filepath.Join(_redisLogDir, redis.Redis_Name+".log"),
			redis.Redis_Name+".rdb", redis.DbDir,
			slaveof,
			redis.MaxClientCon, redis.MaxMemory,
			_clusterEnabled,
		)),
		0750,
	))
}

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Console File */

const templateStartConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s starting..."
%s %s
echo "%s started, pls wating 30 sec..."
`
const templateStopConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s stoping..."
%s -h %s -p %d shutdown
echo "%s stoped, pls wating 30 sec..."
`
const templateRestartConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s admin restarting..."
%s
%s
echo "%s admin restarted, pls wating 30 sec..."
`

const templateStartSenConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s starting..."
%s %s >> %s &
echo "%s started, pls wating 30 sec..."
`
const templateStopSenConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s stoping..."
ps -ef | grep -v grep | grep sentinel | grep %d |awk '{print $2}' | xargs kill -15
echo "%s stoped, pls wating 30 sec..."
`
const templateRestartSenConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
echo "%s admin restarting..."
%s
%s
echo "%s admin restarted, pls wating 30 sec..."
`

/**
 * touch console file
 */
func (redis *Redis) touchConsoleScript() {
	err := utils.MkdirConsolesPath(redis.ConsolePath)
	if err != nil {
		logs.Print(err)
	}

	_type := redis.Type
	_sentinelServer := redis.SentinelServer
	_redisHome := filepath.Join(redis.Redis_HOME, strconv.Itoa(redis.Port))
	_redisConf := filepath.Join(_redisHome, "redis.conf")
	start := filepath.Join(redis.ConsolePath, "start", "start_"+redis.Redis_Name+"_"+strconv.Itoa(redis.Port)+".sh")
	stop := filepath.Join(redis.ConsolePath, "stop", "stop_"+redis.Redis_Name+"_"+strconv.Itoa(redis.Port)+".sh")
	restart := filepath.Join(redis.ConsolePath, "restart", "restart_"+redis.Redis_Name+"_"+strconv.Itoa(redis.Port)+".sh")
	now := time.Now().String()
	// start
	logs.Print(ioutil.WriteFile(
		start,
		[]byte(fmt.Sprintf(templateStartConsole,
			now,
			redis.Redis_Name,
			filepath.Join(redis.Redis_BASE, "bin", "redis-server"), _redisConf,
			redis.Redis_Name,
		)),
		0750,
	))
	// stop
	logs.Print(ioutil.WriteFile(
		stop,
		[]byte(fmt.Sprintf(templateStopConsole,
			now,
			redis.Redis_Name,
			filepath.Join(redis.Redis_BASE, "bin", "redis-cli"), redis.IPs[0], redis.Port,
			redis.Redis_Name,
		)),
		0750,
	))
	// restart
	logs.Print(ioutil.WriteFile(
		restart,
		[]byte(fmt.Sprintf(templateRestartConsole,
			now,
			redis.Redis_Name, stop, start, redis.Redis_Name)),
		0750,
	))

	switch _type {
	case "ha":
		os.MkdirAll(_sentinelServer.Dir, 0755)
		_sentinelConf := filepath.Join(_redisHome, "sentinel.conf")
		_sentinelLog := filepath.Join(_redisHome, "logs", "sentinel_"+redis.Redis_Name+".log")
		startSen := filepath.Join(redis.ConsolePath, "start", "start_"+redis.Redis_Name+"_sentinel_"+strconv.Itoa(redis.Port)+".sh")
		stopSen := filepath.Join(redis.ConsolePath, "stop", "stop_"+redis.Redis_Name+"_sentinel_"+strconv.Itoa(redis.Port)+".sh")
		restartSen := filepath.Join(redis.ConsolePath, "restart", "restart_"+redis.Redis_Name+"_sentinel_"+strconv.Itoa(redis.Port)+".sh")
		// start sen
		logs.Print(ioutil.WriteFile(
			startSen,
			[]byte(fmt.Sprintf(templateStartSenConsole,
				now,
				redis.Redis_Name,
				filepath.Join(redis.Redis_BASE, "bin", "redis-sentinel"), _sentinelConf, _sentinelLog,
				redis.Redis_Name,
			)),
			0750,
		))
		// stop sen
		logs.Print(ioutil.WriteFile(
			stopSen,
			[]byte(fmt.Sprintf(templateStopSenConsole,
				now,
				redis.Redis_Name,
				_sentinelServer.Port,
				redis.Redis_Name,
			)),
			0750,
		))
		// restart sen
		logs.Print(ioutil.WriteFile(
			restartSen,
			[]byte(fmt.Sprintf(templateRestartSenConsole,
				now,
				redis.Redis_Name, stopSen, startSen, redis.Redis_Name)),
			0750,
		))
	case "cluster":
		break
	case "singleton":
		break
	default:

	}
}

/* End Console File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
