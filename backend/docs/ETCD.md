# Etcd 服务注册与发现

## 什么是 Etcd

Etcd 是一个分布式、可靠的键值存储系统，常用于服务注册与发现、配置管理等场景。

在微服务架构中，Etcd 作为服务注册中心：
- 服务启动时自动注册到 Etcd
- 服务健康检查和自动摘除
- 客户端通过 Etcd 发现可用服务

## 安装 Etcd

### macOS

```bash
# 使用 Homebrew 安装
brew install etcd

# 启动 Etcd
brew services start etcd

# 检查状态
brew services list | grep etcd
```

### Linux

```bash
# 下载 Etcd
ETCD_VER=v3.5.10
wget https://github.com/etcd-io/etcd/releases/download/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz

# 解压
tar xzvf etcd-${ETCD_VER}-linux-amd64.tar.gz
cd etcd-${ETCD_VER}-linux-amd64

# 复制二进制文件
sudo cp etcd etcdctl /usr/local/bin/

# 启动 Etcd
etcd
```

### Docker

```bash
# 使用 Docker 运行 Etcd
docker run -d \
  --name etcd \
  -p 2379:2379 \
  -p 2380:2380 \
  -e ALLOW_NONE_AUTHENTICATION=yes \
  bitnami/etcd:latest

# 查看日志
docker logs -f etcd
```

### Docker Compose

```yaml
version: '3.8'

services:
  etcd:
    image: bitnami/etcd:latest
    container_name: me2-etcd
    environment:
      - ALLOW_NONE_AUTHENTICATION=yes
      - ETCD_ADVERTISE_CLIENT_URLS=http://etcd:2379
    ports:
      - "2379:2379"
      - "2380:2380"
    networks:
      - me2-network
    restart: unless-stopped

networks:
  me2-network:
    driver: bridge
```

## 验证 Etcd 安装

```bash
# 检查版本
etcdctl version

# 写入测试数据
etcdctl put test "hello etcd"

# 读取测试数据
etcdctl get test

# 删除测试数据
etcdctl del test
```

## 服务注册配置

### SMS Service 配置示例

```yaml
Name: sms-rpc
ListenOn: 0.0.0.0:8001

# Etcd 服务注册配置
Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: sms.rpc
```

**配置说明:**
- `Hosts`: Etcd 服务器地址列表
- `Key`: 服务注册的键名，建议格式: `{service}.rpc`

### 服务自动注册

go-zero 服务启动时会自动：
1. 连接到 Etcd
2. 注册服务信息 (服务名、地址、端口)
3. 定期发送心跳保持注册
4. 服务停止时自动摘除

## 服务发现

### 客户端配置

其他服务调用 SMS 服务时，通过 Etcd 发现：

```go
// 客户端配置
type Config struct {
    SmsRpc zrpc.RpcClientConf
}

// 配置文件
SmsRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: sms.rpc
```

### 客户端调用

```go
// 创建客户端
smsClient := sms.NewSms(zrpc.MustNewClient(c.SmsRpc))

// 调用服务
resp, err := smsClient.SendCode(ctx, &sms.SendCodeRequest{
    Phone: "13800138000",
    Scene: "login",
})
```

## 查看注册的服务

```bash
# 查看所有服务
etcdctl get --prefix ""

# 查看 SMS 服务
etcdctl get --prefix "sms.rpc"

# 监听服务变化
etcdctl watch --prefix "sms.rpc"
```

## 常见问题

### 1. 服务无法注册到 Etcd

**检查 Etcd 是否运行:**
```bash
# 检查 Etcd 进程
ps aux | grep etcd

# 检查端口
lsof -i :2379
```

**检查网络连接:**
```bash
# 测试连接
telnet 127.0.0.1 2379
```

**查看服务日志:**
```bash
# 服务日志会显示注册信息
tail -f logs/sms-rpc.log
```

### 2. 服务注册成功但客户端无法发现

**检查服务键名:**
- 确保服务端和客户端的 `Key` 配置一致
- 建议使用统一格式: `{service}.rpc`

**查看 Etcd 中的数据:**
```bash
etcdctl get --prefix "sms.rpc"
```

### 3. Etcd 连接超时

**增加超时时间:**
```yaml
Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: sms.rpc
  Timeout: 5000  # 毫秒
```

## 生产环境建议

### 1. Etcd 集群部署

生产环境建议部署 Etcd 集群 (3 或 5 节点):

```yaml
Etcd:
  Hosts:
    - etcd1.example.com:2379
    - etcd2.example.com:2379
    - etcd3.example.com:2379
  Key: sms.rpc
```

### 2. 安全配置

启用 TLS 加密和身份认证:

```yaml
Etcd:
  Hosts:
    - https://etcd1.example.com:2379
  Key: sms.rpc
  CertFile: /path/to/cert.pem
  KeyFile: /path/to/key.pem
  CaCert: /path/to/ca.pem
  User: username
  Pass: password
```

### 3. 监控告警

使用 Prometheus 监控 Etcd:
- 节点健康状态
- 写入延迟
- 存储空间
- 客户端连接数

## 参考资料

- [Etcd 官方文档](https://etcd.io/docs/)
- [go-zero 服务注册与发现](https://go-zero.dev/docs/tutorials/service/registry)
- [Etcd 最佳实践](https://etcd.io/docs/v3.5/op-guide/performance/)
