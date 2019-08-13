#!/bin/bash
./etcdcc client.file --hosts="https://etcd1:2379,https://etcd2:2379,https://etcd3:2379" --prefix="gateway/v1/proxy"
