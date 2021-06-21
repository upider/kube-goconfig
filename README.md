# goconfig

## 用途

将nacos的配置同步到k8s，同步方式为：

```mermaid
graph LR
A[nacos namespace] --> B[k8s namespace];
C[nacos dataId] --> D[k8s configMap]
```
