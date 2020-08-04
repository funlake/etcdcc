#!/bin/bash
go build . && ./etcdcc server.start \
	--c=/mnt/d/develop/etcd-operator/example/tls/certs/etcd-client.crt \
	--k=/mnt/d/develop/etcd-operator/example/tls/certs/etcd-client.key \
	--ca=/mnt/d/develop/etcd-operator/example/tls/certs/etcd-client-ca.crt \
	--hosts=https://example-client.default.svc:2379
