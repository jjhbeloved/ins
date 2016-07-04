package redis

import (
	"asiainfo.com/ins/logs"
	"asiainfo.com/ins/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"asiainfo.com/ins/cli"
)

/*
调优如下
vm.overcommit_memory = 1
/etc/sysctl.conf
sysctl vm.overcommit_memory=1
设置系统参数vm.overcommit_memory=1, 可以避免bgsave/aofrewrite失败.

The TCP backlog setting of 511 cannot be enforced because /proc/sys/net/core/somaxconn is set to the lower value of 128

WARNING you have Transparent Huge Pages (THP) support enabled in your kernel. This will create latency and memory usage issues with Redis.
To fix this issue run the command 'echo never > /sys/kernel/mm/transparent_hugepage/enabled' as root,
and add it to your /etc/rc.local in order to retain the setting after a reboot.
Redis must be restarted after THP is disabled.
*/

type Redis struct {
	Redis_PKG  string `json:"redis_pkg"`
	Redis_HOME string `json:"redis_home"`
	Ruby_PKG   string `json:"ruby_pkg"`
	Ruby_HOME  string `json:"ruby_home"`
	GEMS_PKG   string `json:"gems_pkg"`
	IsCluster  bool   `json:"isCluster"`
}

func (redis *Redis) Json(bs []byte) error {
	return json.Unmarshal(bs, &redis)
}

const installRedis = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
make PREFIX=%s install
rm -rf %s-SRC
cd %s
mkdir -p conf logs
`

const installRedisCluster = `#!/bin/bash
########################################
# AUTO CREATE BY XIAOXIAO INS %s
########################################
mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
sh configure -prefix=%s --disable-install-doc --disable-install-rdoc --disable-install-capi --enable-rubygems
make -j4 && make install
rm -rf %s-SRC

cd %s
./gem install %s --local

mkdir -p %s-SRC
tar xzf %s -C %s-SRC --strip-components=1
cd %s-SRC
make PREFIX=%s install
cd %s
mkdir -p conf logs
cd bin
cp %s .
sed -i "1s#ruby#%s#" redis-trib.rb
rm -rf %s-SRC
`

func (redis *Redis) Install() error {
	pkg, err := utils.DownloadToDir(redis.Redis_PKG, cli.PKG_PATH)
	if err != nil {
		return err
	}
	redis.Redis_PKG = pkg

	_, err = os.Stat(redis.Redis_PKG)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	now := time.Now().String()
	redisSh := filepath.Join(utils.TMPD, "redisCreate.sh")
	defer os.Remove(redisSh)
	if !redis.IsCluster {
		logs.Print(ioutil.WriteFile(
			redisSh,
			[]byte(fmt.Sprintf(installRedis,
				now,
				redis.Redis_HOME,
				redis.Redis_PKG, redis.Redis_HOME,
				redis.Redis_HOME,
				redis.Redis_HOME,
				redis.Redis_HOME,
				redis.Redis_HOME,
			)),
			0750,
		))
	} else {
		pkg, err = utils.DownloadToDir(redis.Ruby_PKG, cli.PKG_PATH)
		if err != nil {
			return err
		}
		redis.Ruby_PKG = pkg
		pkg, err = utils.DownloadToDir(redis.GEMS_PKG, cli.PKG_PATH)
		if err != nil {
			return err
		}
		redis.GEMS_PKG = pkg

		logs.Print(ioutil.WriteFile(
			redisSh,
			[]byte(fmt.Sprintf(installRedisCluster,
				now,
				redis.Ruby_HOME,
				redis.Ruby_PKG, redis.Ruby_HOME,
				redis.Ruby_HOME,
				redis.Ruby_HOME,
				redis.Ruby_HOME,
				filepath.Join(redis.Ruby_HOME, "bin"), redis.GEMS_PKG,
				redis.Redis_HOME,
				redis.Redis_PKG, redis.Redis_HOME,
				redis.Redis_HOME,
				redis.Redis_HOME,
				redis.Redis_HOME,
				filepath.Join(redis.Redis_HOME+"-SRC", "src", "redis-trib.rb"),
				filepath.Join(redis.Ruby_HOME, "bin", "ruby"),
				redis.Redis_HOME,
			)),
			0750,
		))
	}

	// 根据模板生成domain
	cmd = exec.Command("sh", redisSh)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return err
}
