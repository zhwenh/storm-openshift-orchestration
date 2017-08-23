运行在k8s（openshift）上的storm集群
====================

由于storm是流计算形式的，所以特别适合构建在无状态的k8s容器集群里面。
这个编排需要用到之前的zookeeper编排，先生成一个3节点的zk集群，利用这个集群storm来保存状态信息和集群间通信。
*当前的简单实现中，并未对zk的鉴权和storm鉴权进行设置*

构建镜像
------------------------

```oc new-build https://github.com/asiainfoLDP/storm-openshift-orchestration.git --context-dir='image' ```
*注意修改编排文件中的镜像名*

创建编排
------------

参数意义:

* STORM_CMD - The storm command to run. Currently can be "nimbus", "supervisor", or "ui"
* CONFIGURE_ZOOKEEPER - Set this to signal that we want to automatically set up the servers in the storm.yaml file. If this is set, the following variables are checked.
* ZK_SERVER_1_SERVICE_HOST - Zookeeper Server 1
* ZK_SERVER_2_SERVICE_HOST - Zookeeper Server 2
* ZK_SERVER_${N}_SERVICE_HOST - Will keep searching for Zookeeper servers
* APACHE_STORM_NIMBUS_SERVICE_HOST - Nimbus server IP

To run this in Kubernetes:

```
oc create -f kubernetes-nimbus-service.yaml
oc create -f kubernetes-ui-service.yaml
```

Create the replication controllers:

```
oc create -f kubernetes-nimbus-rc.yaml
oc create -f kubernetes-supervisor-rc.yaml
oc create -f kubernetes-ui-rc.yaml
```
需要替换其中的*instanceid* *zookeeper服务名*

暴露dashboard ui
```oc expose svc sb-instanceid-su```


绑定和解除绑定
------------
由于没有用户名和密码的认证，简单返回storm集群的地址就好了。*kubernetes-nimbus-service.yaml* 中的内容。

