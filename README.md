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

#### Run it
```
GO111MODULE=on go build -o etcdcc
./etcdcc --h #Ssl verification is needed,so you must specify --k,--c,--ca for etcd connection
#server
./etcdcc server.start --hosts=https://127.0.0.1:2479
#file client,see files in /opt/dev/abc
./etcdcc client.file --hosts=https://127.0.0.1:2379 --prefix=dev/abc --storeDir=/opt
#unix socket client,serve unix socket server for application
#command for application to request unix socket is just : get [config type]/[config name] [specific key]
./etcdcc client.sock --hosts=https://127.0.0.1:2379 --prefix=dev/abc --sock=/run/etcdcc.sock

```
##### By Docker
##### Compile and run