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
GO111MODULE=on go build
./apiserver --h
#server
./apiserver --hosts=https://127.0.0.1:6379
#file client,see files in /opt/dev/abc
./apiserver file --hosts=https://127.0.0.1 --prefix=dev/abc --storeDir=/opt
#uds client,default socket is /run/etcdcc.sock
./apiserver uds --hosts=https://127.0.0.1 --prefix=dev/act
```
##### By Docker
##### Compile and run