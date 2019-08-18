package client

import (
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
)

var watcher = &EtcdUdsWatcher{}

func initData() {
	watcher.SaveLocal("json/lake", base64.StdEncoding.EncodeToString([]byte(`{"a":{"b":{"c":"hello"}},"d":"world"}`)))
	watcher.SaveLocal("toml/lake", base64.StdEncoding.EncodeToString([]byte(`
# This is a TOML document. Boom.
title = "TOML Example"
[owner]
name = "Tom Preston-Werner"
organization = "GitHub"
bio = "GitHub Cofounder & CEO\nLikes tater tots and beer."
dob = 1979-05-27T07:32:00Z # First class dates? Why not?
[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true
[servers]
  # You can indent as you please. Tabs or spaces. TOML don't care.
  [servers.alpha]
  ip = "10.0.0.1"
  dc = "eqdc10"
  [servers.beta]
  ip = "10.0.0.2"
  dc = "eqdc10"
[clients]
data = [ ["gamma", "delta"], [1, 2] ]
# Line breaks are OK when inside arrays
hosts = [
  "alpha",
  "omega"
]
	`)))
	watcher.SaveLocal("yaml/lake", base64.StdEncoding.EncodeToString([]byte(`
service:
  image: registry.cn/abcd
  restart: always
  hostname: a.etcchebao.com
  expose:
    - 80/tcp
  environment:
    - 'MYSQL_ETC1_MASTER_HOST=hosts.com'
    - 'MYSQL_ETC1_PORT=3306'
    - 'MYSQL_ETC1_USERNAME=etc'
    - 'MYSQL_ETC1_PASSWORD=abcd'
    - 'MYSQL_ETC1_SLAVE_HOST=hosts.com'
    - 'MYSQL_ETC1_SLAVE2_HOST=hosts.com'
    - 'MYSQL_ETC2_MASTER_HOST=hosts.com'
    - 'MYSQL_ETC2_PORT=3306'
    - 'MYSQL_ETC2_USERNAME=etc'
    - 'MYSQL_ETC2_PASSWORD=abcd'
    - 'MYSQL_ETC3_MASTER_HOST=hosts.com'
    - 'MYSQL_ETC3_PORT=3306'
    - 'MYSQL_ETC3_USERNAME=etc'
    - 'MYSQL_ETC3_PASSWORD=abcd'
    - 'MYSQL_UNITOLL_MASTER_HOST=hosts.com'
    - 'MYSQL_UNITOLL_PORT=3306'
    - 'MYSQL_UNITOLL_PASSWORD=abcd'
    - 'MONGO_HOST=10.30.222.42'
    - 'MONGO_USER=root'
    - 'MONGO_PASSWD=hosts.com'
    - 'MONGO_PORT=27017'
    - 'MONGO_HOST2=hosts.com'
    - 'MONGO_USER2=root'
    - 'MONGO_PASSWD2=abcd'
    - 'MONGO_PORT2=27017'
    - 'REDIS_ETC1_HOST=hosts.com'
    - 'REDIS_ETC1_PORT=6380'
    - 'REDIS_ETC1_PASSWORD=asdf'
    - TZ=Asia/Shanghai
    - ENV=test
    - 'PHONE_DATA_DIR=/'
    - Run_Mode=pro
  labels:
    aliyun.scale: 1
    aliyun.routing.port_80: a.etcchebao.com
`)))
}

type testCase struct {
	cmd  []string
	want string
}

func TestEtcdUdsWatcher_Find(t *testing.T) {
	initData()

	for _, ts := range []testCase{
		{
			cmd:  []string{"get", "json/lake", "a.b.c"},
			want: "hello",
		},
		{
			cmd:  []string{"get", "toml/lake", "servers.alpha.ip"},
			want: "10.0.0.1",
		},
		{
			cmd:  []string{"get", "toml/lake", "owner.bio"},
			want: "GitHub Cofounder & CEO\nLikes tater tots and beer.",
		},
		{
			cmd:  []string{"get", "yaml/lake", "service.image"},
			want: "registry.cn/abcd",
		},
		{
			cmd:  []string{"get", "yaml/lake", "service.environment.0"},
			want: "MYSQL_ETC1_MASTER_HOST=hosts.com",
		},
		{
			cmd:  []string{"get", "yaml/lake", "service.labels.aliyun\\.routing\\.port_80"},
			want: "a.etcchebao.com",
		},
	} {
		val, err := watcher.Find(ts.cmd)
		if err != nil {
			t.Error(err.Error())
		} else {
			if val != ts.want {
				t.Error(errors.New(fmt.Sprintf("Expect %s,get %s", ts.want, val)))
			}
		}
	}
}

func BenchmarkEtcdUdsWatcher_Find(b *testing.B) {
	initData()
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = watcher.Find([]string{"get", "yaml/lake", "service.labels.aliyun\\.routing\\.port_80"})
		}
	})
}
