[![Build Status](https://travis-ci.org/funlake/etcdcc.svg?branch=master)](https://travis-ci.org/funlake/etcdcc)
[![Go Report Card](https://goreportcard.com/badge/github.com/funlake/etcdcc)](https://goreportcard.com/report/github.com/funlake/etcdcc)
# Etcdcc
#### What's this
Restful/grpc service for config center base on etcd which's distributed,stable,high performance k/v store storage.

#### How it works
Receive http/grpc request,and save/del/update specific k/v record in etcd3

#### Performace
Hihgly depends on  etcd ,see how etcd proved it [here](https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/performance.md) ,
100k write request / 50k qps,this's way beyond what we expected

#### Set up
1. Start up mysql, see install/config.sql , select a database,put it in it.
2. Add env variables for mysql connections , they are `MYSQL_HOST`,`MYSQL_USERNAME`,`MYSQL_PASSWORD`,`MYSQL_DB`
3. Start up etcd, run it with TLS mode,so you need to have key,cert,ca files

#### Run & Test
```
GO111MODULE=on go test ./...
GO111MODULE=on go build -o etcdcc
./etcdcc --h #Ssl verification is needed,so you must specify --k,--c,--ca for etcd connection
#server
./etcdcc server.start --hosts=https://127.0.0.1:2379
#file client,see files in /opt/dev/abc
./etcdcc client.file --hosts=https://127.0.0.1:2379 --prefix=dev/abc --storeDir=/opt
#unix socket client,serve unix socket for application
#command for application to request unix socket is just : get [config type]/[config name] [specific key]
./etcdcc client.sock --hosts=https://127.0.0.1:2379 --prefix=dev/abc --sock=/run/etcdcc.sock
```

###### Example
> Put data into mysql & etcd,mysql install
```
$ echo -en "get toml/tm" | nc  -U /run/etcdcc.sock
ok,{"clients":{"data":[["gamma","delta"],[1,2]]},"database":{"connection_max":5000,"enabled":true,"ports":[8001,8001,8002],"server":"192.168.1.1"},"owner":{"bio":"GitHub Cofounder \u0026 CEO\nLikes tater tots and beer.","dob":"1979-05-27T07:32:00Z","name":"Tom Preston-Werner","organization":"GitHub"},"servers":{"alpha":{"dc":"eqdc10","ip":"10.0.0.1"},"beta":{"dc":"eqdc10","ip":"10.0.0.2"}},"title":"TOML Example"}
$ echo -en "get toml/tm owner.name" | nc  -U /run/etcdcc.sock
ok,Tom Preston-Werner
```
##### By Docker
##### Compile and run