FROM golang:1.12
COPY keys /keys
COPY . $GOPATH/src/etcdcc/
RUN cd $GOPATH/src/etcdcc && GO111MODULE=off go build -o etcdcc
WORKDIR $GOPATH/src/etcdcc
ENV ETCD_HOSTS "http://127.0.0.1:2379"
ENTRYPOINT ["./etcdcc"]
CMD []
