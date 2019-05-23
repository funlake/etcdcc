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
1. Start up mysql, see install/config.sql , select a database, source it.
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
> Set up a key value pair
```
$ cat setup.sh
#Content with url encode
toml='%23+This+is+a+TOML+document.+Boom.
title+%3d+%22TOML+Example%22
%5bowner%5d
name+%3d+%22Tom+Preston-Werner%22
organization+%3d+%22GitHub%22
bio+%3d+%22GitHub+Cofounder+%26+CEO%5cnLikes+tater+tots+and+beer.%22
dob+%3d+1979-05-27T07%3a32%3a00Z+%23+First+class+dates%3f+Why+not%3f
%5bdatabase%5d
server+%3d+%22192.168.1.1%22
ports+%3d+%5b+8001%2c+8001%2c+8002+%5d
connection_max+%3d+5000
enabled+%3d+true
%5bservers%5d
++%23+You+can+indent+as+you+please.+Tabs+or+spaces.+TOML+don%27t+care.
++%5bservers.alpha%5d
++ip+%3d+%2210.0.0.1%22
++dc+%3d+%22eqdc10%22
++%5bservers.beta%5d
++ip+%3d+%2210.0.0.2%22
++dc+%3d+%22eqdc10%22
%5bclients%5d
data+%3d+%5b+%5b%22gamma%22%2c+%22delta%22%5d%2c+%5b1%2c+2%5d+%5d'
curl -X POST -d "key=tm&mod=abc&env=dev&val=$toml&type=toml" -H "Content-Type: application/x-www-form-urlencoded" http://127.0.0.1/config
$ ./setup.sh
```
> Application search by unix socket
> you can make it run continually,shell script or something else? you decide
```
$ echo -en "get toml/tm" | nc  -U /run/etcdcc.sock
ok,{"clients":{"data":[["gamma","delta"],[1,2]]},"database":{"connection_max":5000,"enabled":true,"ports":[8001,8001,8002],"server":"192.168.1.1"},"owner":{"bio":"GitHub Cofounder \u0026 CEO\nLikes tater tots and beer.","dob":"1979-05-27T07:32:00Z","name":"Tom Preston-Werner","organization":"GitHub"},"servers":{"alpha":{"dc":"eqdc10","ip":"10.0.0.1"},"beta":{"dc":"eqdc10","ip":"10.0.0.2"}},"title":"TOML Example"}
$ echo -en "get toml/tm owner.name" | nc  -U /run/etcdcc.sock
ok,Tom Preston-Werner
```
> Update value to see if application get new value
```
$ cat update.sh
#Content with url encode
toml='%23+This+is+a+TOML+document.+Boom.
title+%3d+%22TOML+Example%22
'
curl -X PUT -d "id=1&key=tm&mod=abc&env=dev&val=$toml&type=toml" -H "Content-Type: application/x-www-form-urlencoded" http://127.0.0.1/config
$ ./update.sh
```
##### By Docker
##### Compile and run