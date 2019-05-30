#!/bin/bash
docker run -v /keys-test:/keys -p 8081:80 \
-e MYSQL_HOST=$MYSQL_HOST \
-e MYSQL_USERNAME=$MYSQL_USERNAME \
-e MYSQL_PASSWORD=$MYSQL_PASSWORD \
-e MYSQL_DB=$MYSQL_DB \
-e ETCD_HOSTS=$ETCD_HOSTS \
--add-host etcd1:120.76.26.106 --add-host etcd2:120.76.101.254 --add-host etcd3:47.107.210.123 \
etcdcc:latest server.start \
--hosts="https://etcd1:2379,https://etcd2:2379,https://etcd3:2379"