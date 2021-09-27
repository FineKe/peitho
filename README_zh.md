
# Peitho - Hyperledger Fabric chaincode 云原生托管系统


可以将Hyperledger Fabric的chaincode 交给k8s来管理，适用于在k8s上部署Hyperledger Fabric联盟链，解决chaincode游离于k8s管控问题。

## 功能特性
- k8s完全托管chaincode
- 对fabric无任何侵入，只需将CORE_VM_ENDPOINT环境变量配置成Peitho服务即可
- 支持fabric 1.4.x 2.0 以上
## 软件架构
![架构图](./docs/images/peitho-architecture.png)

## 快速开始
### 构建
1.获取源码
```shell
git clone https://github.com/tianrandailove/peitho
```
2.编译
```shell
cd peitho
make build
```
3.生成镜像
```shell
make image
```
## 使用指南
1. 配置peitho-configmap.yaml
```yaml
# Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

apiVersion: v1
data:
  kubeconfig: |-
    #填入你获取的k8s 访问凭证
  peitho.yml: |-
    k8s:
      namespace: fabric #命名空间
      kubeconfig: /root/kube/kubeconfig #k8s访问配置文件
      dns: #如果chaincode 和 peer不在同一个环境的情况下，需要配置peer地址的解析
        - 127.0.0.1:peer0.org1.example.com
        - 127.0.0.1:peer1.org1.example.com
        - 127.0.0.1:peer0.org2.example.com
        - 127.0.0.1:peer1.org2.example.com
    docker:
      endpoint: unix:///host/var/run/docker.sock # docker的端点
      registry: #镜像仓库相关
        server-address: #仓库地址
          xxx.xxx.xxx.xxx:xxxx
        project: #项目名
          chaincode
        email: #邮箱
          litesky@foxmail.com
        username: #用户名
          admin
        password: #密码
          harbor

    log:
      name: peitho # Logger的名字
      development: true # 是否是开发模式。如果是开发模式，会对DPanicLevel进行堆栈跟踪。
      level: debug # 日志级别，优先级从低到高依次为：debug, info, warn, error, dpanic, panic, fatal。
      format: console # 支持的日志输出格式，目前支持console和json两种。console其实就是text格式。
      enable-color: true # 是否开启颜色输出，true:是，false:否
      disable-caller: true # 是否开启 caller，如果开启会在日志中显示调用日志所在的文件、函数和行号
      disable-stacktrace: false # 是否再panic及以上级别禁止打印堆栈信息
kind: ConfigMap
metadata:
  name: peitho-configmap
  namespace: fabric

```
2. 配置peitho-deployment.yaml
```yaml
# Copyright 2021 Ke Fan <litesky@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

apiVersion: apps/v1
kind: Deployment
metadata:
  name: peitho
spec:
  replicas: 1
  selector:
    matchLabels:
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
    type: RollingUpdate
  template:
    metadata:
    spec:
      containers:
      - image: tianrandailoving/peitho:latest
        imagePullPolicy: Always
        name: peitho
        ports:
        - containerPort: 8080
          name: peitho
          protocol: TCP
        resources: {}
        securityContext:
          allowPrivilegeEscalation: false
          privileged: false
          readOnlyRootFilesystem: false
          runAsNonRoot: false
        stdin: true
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        tty: true
        volumeMounts:
        - mountPath: /host/var/run/
          name: vol2
        - mountPath: /root/peitho.yml
          name: vol1
          subPath: peitho.yml
        - mountPath: /root/kube/kubeconfig
          name: vol1
          subPath: kubeconfig
      dnsConfig: {}
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - hostPath:
          path: /var/run/
          type: ""
        name: vol2
      - configMap:
          defaultMode: 256
          items:
          - key: peitho.yml
            path: peitho.yml
          - key: kubeconfig
            path: kubeconfig
          name: peitho-configmap
          optional: false
        name: vol1
```
3. 创建configmap 和 deployment
```shell
kubectl apply -f peitho-configmap.yaml
kubectl apply -f peitho-deployment.yaml
```
4. 更改peer的环境变量
```yaml
- name: CORE_VM_ENDPOINT
  value: tcp://peitho:8080
```
## 关于作者

- kefan < litesky@foxmail.com >

## License
Peitho is licensed under the MIT.
## 其他
获取peitho镜像
```shell
docker pull tianrandailoving/peitho:latest
```