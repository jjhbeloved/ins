package mysql

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

type Mysql struct {
	Mysql_Name  string        `json:"mysql_name"`
	Mysql_HOME  string        `json:"mysql_home"`
	Mysql_Base  string        `json:"mysql_base"`
	Port        int           `json:"port"`
	Root_PWD    string        `json:"rootPwd"`
	ConsolePath string        `json:"consolePath"`
	Option      string        `json:"option"`
}

func (w *Mysql) Json(bs []byte) error {
	return json.Unmarshal(bs, &w)
}

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Method */

func (mysql *Mysql) Add() error {
	mysql.touchConf()
	mysql.touchConsoleScript()
	return nil
}

func (mysq *Mysql) Remove() error {
	return nil
}

/* End Method */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

/* --------------------------------------------------------- */
/* --------------------------------------------------------- */
/* Begin Configuration File */
const mysqlConf = `########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
[mysqld]
lower_case_table_names=1
server-id  = %d
user       = mysql
port       = %d
basedir    = %s
datadir    = %s/data
tmpdir     = %s/tmp
socket     = %s/proc/my_%d.sock
pid_file   = %s/proc/my_%d.pid

#repl
binlog-format=ROW
binlog-rows-query-log-events=0
max_binlog_cache_size = 256M
max_binlog_size = 512M
sync_binlog=1
log-slave-updates=true

gtid-mode=on
enforce-gtid-consistency=true
master-info-repository=TABLE
relay-log-info-repository=TABLE
sync-master-info=1
slave-parallel-workers=8
binlog-checksum=CRC32
master-verify-checksum=1
slave-sql-verify-checksum=1

# Logging
log_bin             =  %s/logs/binlog/master_%d_bin
log-error           =  %s/logs/error/error_%d.log
slow_query_log_file =  %s/logs/slow/slow_%d.log
general_log_file    =  %s/logs/general/general_%d.log
relay_log           =  %s/logs/relay/relay_%d.log
relay_log_info_file =  %s/logs/relay/relay_log_%d.info
expire_logs_days=3

slow_query_log=ON
long_query_time=1
event_scheduler=1

#audit
#plugin-dir=/veris/usr/lib64/mysql/plugin
#plugin-load=audit_log.so
#audit_log_file=%s/logs/audit/audit_%d.log
#audit_log_rotate_on_size=268435456
#audit_log_flush=ON
#audit_log_policy=ALL

character-set-server = utf8
init_connect = 'SET NAMES utf8'

#common
transaction_isolation = read-committed
lower_case_table_names = 1
symbolic-links=0
open_files_limit        = 30000
max_connections         = 2000
max_user_connections = 2040
back_log                        = 200
skip-name-resolve       = ON
max_allowed_packet = 128M
explicit_defaults_for_timestamp=true
sql_mode = NO_ENGINE_SUBSTITUTION,STRICT_TRANS_TABLES
log_bin_trust_function_creators=1
# innodb_max_dirty_pages_pct=20

# InnoDB
innodb_data_home_dir    = %s/innodb/data
innodb_data_file_path   = ibdata1:1024M;ibdata2:1024M:autoextend
innodb_log_group_home_dir = %s/innodb/logs

innodb_buffer_pool_size = 10240M
innodb_flush_method = O_DIRECT
innodb_thread_concurrency = 64
innodb_io_capacity = 8000
innodb_read_io_threads =8
innodb_write_io_threads = 64
innodb_change_buffering = inserts
innodb_log_buffer_size = 32M
innodb_log_files_in_group = 3
innodb_log_file_size = 1G
# innodb_flush_log_at_trx_commit = 2
# sync_binlog=0

#  init; Antelope  Compact;
innodb_file_format=barracuda
innodb_file_format_max=barracuda
innodb_strict_mode=1
innodb_print_all_deadlocks = 1

# MyISAM
key_buffer_size = 128M

# Other
query_cache_size = 0
# query_cache_type = 0
tmp_table_size = 32M
max_heap_table_size = 32M
thread_cache_size = 128
bulk_insert_buffer_size = 64M

# table_open_cache  = 10262
# table_definition_cache = 10000
table_definition_cache=200
table_open_cache=128
net_write_timeout = 300
net_read_timeout  = 300

# Remove leading # to set options mainly useful for reporting servers.
# The server defaults are faster for transactions and fast SELECTs.
# Adjust sizes as needed, experiment to find the optimal values.
join_buffer_size = 2M
sort_buffer_size = 2M
read_buffer_size = 2M
# read_rnd_buffer_size = 4M
`

/**
 * touch conf local file
 */
func (mysql *Mysql) touchConf() {
	os.MkdirAll(utils.TMPD, 0750)
	_mysqlBase := filepath.Join(mysql.Mysql_Base, strconv.Itoa(mysql.Port), "db")
	_mysqlProc := filepath.Join(_mysqlBase, "proc")
	_mysqlTmp := filepath.Join(_mysqlBase, "tmp")
	_mysqlData := filepath.Join(_mysqlBase, "data")
	_mysqlScripts := filepath.Join(_mysqlBase, "scripts")
	_mysqlEtc := filepath.Join(_mysqlBase, "etc")
	_mysqlBak := filepath.Join(_mysqlBase, "bak")
	_mysqlLogsAudit := filepath.Join(_mysqlBase, "logs", "audit")
	_mysqlLogsBinlog := filepath.Join(_mysqlBase, "logs", "binlog")
	_mysqlLogsGeneral := filepath.Join(_mysqlBase, "logs", "general")
	_mysqlLogsError := filepath.Join(_mysqlBase, "logs", "error")
	_mysqlLogsSlow := filepath.Join(_mysqlBase, "logs", "slow")
	_mysqlLogsRelay := filepath.Join(_mysqlBase, "logs", "relay")
	_mysqlInnodbLogs := filepath.Join(_mysqlBase, "innodb", "logs")
	_mysqlInnodbData := filepath.Join(_mysqlBase, "innodb", "data")
	os.MkdirAll(_mysqlProc, 0755)
	os.MkdirAll(_mysqlTmp, 0755)
	os.MkdirAll(_mysqlData, 0755)
	os.MkdirAll(_mysqlLogsAudit, 0755)
	os.MkdirAll(_mysqlLogsBinlog, 0755)
	os.MkdirAll(_mysqlLogsGeneral, 0755)
	os.MkdirAll(_mysqlLogsError, 0755)
	os.MkdirAll(_mysqlLogsSlow, 0755)
	os.MkdirAll(_mysqlLogsRelay, 0755)
	os.MkdirAll(_mysqlInnodbLogs, 0755)
	os.MkdirAll(_mysqlInnodbData, 0755)
	os.MkdirAll(_mysqlScripts, 0755)
	os.MkdirAll(_mysqlBak, 0755)
	os.MkdirAll(_mysqlEtc, 0755)
	_mysqlCnf := filepath.Join(_mysqlEtc, "my_" + strconv.Itoa(mysql.Port) + ".cnf")
	now := time.Now().String()
	// _redisConf
	logs.Print(ioutil.WriteFile(
		_mysqlCnf,
		[]byte(fmt.Sprintf(mysqlConf,
			now,
			mysql.Port, mysql.Port,
			mysql.Mysql_HOME,
			_mysqlBase, _mysqlBase, _mysqlBase, mysql.Port, _mysqlBase, mysql.Port,
			_mysqlBase, mysql.Port,
			_mysqlBase, mysql.Port,
			_mysqlBase, mysql.Port,
			_mysqlBase, mysql.Port,
			_mysqlBase, mysql.Port,
			_mysqlBase, mysql.Port,
			_mysqlBase, mysql.Port,
			_mysqlBase,
			_mysqlBase,
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
#!/bin/sh
%s/bin/mysqld_safe --defaults-file=%s/etc/my_%d.cnf  --ledir=%s/sbin  &
`
const templateStopConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
v_port=%d;
v_pid=%s

if (kill -0 $v_pid 2>/dev/null)
then
        echo -e "Shutting down MySQL"
        kill $v_pid

        while true
        do
                v_count=0;
                v_count=%s
                if [ "x$v_count" == "x1" ]
                then
                        echo -e "Db port=[$v_port]. Waiting for shutdown..."
                else
                        break;
                fi
                sleep 1;
        done
        echo -e "SUCCESS!"
else
        echo  "MySQL server process [PORT=$v_port] is not running!"
fi
`
const templateRestartConsole = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################

`

/**
 * touch console file
 */
func (mysql *Mysql) touchConsoleScript() {
	//err := utils.MkdirConsolesPath(mysql.ConsolePath)
	err := os.MkdirAll(mysql.ConsolePath, 0750)
	if err != nil {
		logs.Print(err)
	}
	_mysqlBase := filepath.Join(mysql.Mysql_Base, strconv.Itoa(mysql.Port), "db")
	_vpid := fmt.Sprintf("`ps -ef|grep mysqld|grep '%s/sbin/mysqld'|grep '\\--datadir = '|grep '\\--socket ='|grep $v_port|awk '{print $2}'`", mysql.Mysql_HOME)
	_vcount := fmt.Sprintf("`ps -ef|grep mysqld|grep $v_pid|wc -l`")
	start := filepath.Join(mysql.ConsolePath, "start" + ".sh")
	stop := filepath.Join(mysql.ConsolePath, "stop" + ".sh")
	restart := filepath.Join(mysql.ConsolePath, "restart" + ".sh")
	now := time.Now().String()
	// start
	logs.Print(ioutil.WriteFile(
		start,
		[]byte(fmt.Sprintf(templateStartConsole,
			now,
			mysql.Mysql_HOME, _mysqlBase, mysql.Port, mysql.Mysql_HOME,
		)),
		0750,
	))
	// stop
	logs.Print(ioutil.WriteFile(
		stop,
		[]byte(fmt.Sprintf(templateStopConsole,
			now,
			mysql.Port,
			_vpid, _vcount,
		)),
		0750,
	))
	// restart
	logs.Print(ioutil.WriteFile(
		restart,
		[]byte(fmt.Sprintf(templateRestartConsole,
			now,
			stop, start,
		)),
		0750,
	))

}

/* End Console File */
/* --------------------------------------------------------- */
/* --------------------------------------------------------- */

/**
* 1.	初始化数据块
* $HOME/mysql/scripts/mysql_install_db --defaults-file=/veris/odc/mysql/3401/db/etc/mysql_3401.cnf  --basedir=$HOME/mysql
* 2.	设置 ~/.bashrc
* alias ls='ls --color --show-control-chars'
* alias sc='source ~/.bashrc'
* alias my='mysql -h 127.0.0.1 -uroot -p -P'
* alias myp='mysql -h 127.0.0.1 -upadmin -p -P'
* 3.	修改root密码
* mysqladmin -h127.0.0.1 -P3401 -u root password
*/