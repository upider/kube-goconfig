#!/bin/sh

if [ -z $AUTO_CREATE_K8S_NS ];then
    echo "NO AUTO_CREATE_K8S_NS, set it to true"
    export AUTO_CREATE_K8S_NS=true
else 
    echo "AUTO_CREATE_K8S_NS = "$AUTO_CREATE_K8S_NS
fi

if [ -z $LOG_LEVEL ];then
    echo "NO LOG_LEVEL, set it to info"
    export LOG_LEVEL=info
else 
    echo "LOG_LEVEL = "$LOG_LEVEL
fi

if [ -z $NACOS_PORT ];then
    echo "NO NACOS_PORT, set it to 8848"
    export NACOS_PORT=8848
else 
    echo "NACOS_PORT = "$NACOS_PORT
fi

if [ -z $CONFIG_SCAN_TIME ];then
    echo "NO CONFIG_SCAN_TIME, set it to 10"
    export CONFIG_SCAN_TIME=10
else 
    echo "CONFIG_SCAN_TIME = "$CONFIG_SCAN_TIME
fi

if [ -z $NACOS_IPS ];then
    echo "NO NACOS_IPS, exit"
else
    echo NACOS_IPS=$NACOS_IPS
fi

if [ -z $NAMESPACES ];then
    echo "NO NAMESPACES, exit"
else
    echo NAMESPACES=$NAMESPACES
fi

./main --autoCreatek8sNs $AUTO_CREATE_K8S_NS --logLevel $LOG_LEVEL --configScanTime $CONFIG_SCAN_TIME --nacosPort $NACOS_PORT --nacosIPs $NACOS_IPS --namespaces $NAMESPACES