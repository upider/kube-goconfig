# Build the manager binary
FROM golang:1.16 as builder
WORKDIR /workspace/
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
COPY ./ src/
RUN cd src && go build -o main cmd/main.go

FROM centos:7
WORKDIR /
COPY --from=builder /workspace/src/main .
USER 65532:65532
ENTRYPOINT ./main --autoCreatek8sNs $AUTO_CREATE_K8S_NS --logLevel $LOG_LEVEL --configScanTime $CONFIG_SCAN_TIME --nacosPort $NACOS_PORT --nacosIPs $NACOS_IPS --namespaces $NAMESPACES