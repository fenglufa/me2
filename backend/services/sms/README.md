# SMS Service

短信服务，提供验证码发送和验证功能。

## 功能

- 发送验证码 (SendCode)
- 验证验证码 (VerifyCode)

## 技术栈

- go-zero RPC
- 阿里云短信服务
- Redis (验证码存储)
- Etcd (服务注册与发现)

## 配置

配置文件: `rpc/etc/sms.yaml`

环境变量:
- `ETCD_HOST` - Etcd 地址 (默认: 127.0.0.1:2379)
- `ALIYUN_ACCESS_KEY_ID` - 阿里云 AccessKeyId
- `ALIYUN_ACCESS_KEY_SECRET` - 阿里云 AccessKeySecret
- `ALIYUN_SMS_SIGN_NAME` - 短信签名
- `ALIYUN_SMS_TEMPLATE_CODE` - 短信模板代码
- `REDIS_HOST` - Redis 地址
- `REDIS_PASS` - Redis 密码

## 使用

```bash
# 1. 生成 RPC 代码 (必须先执行)
make gen-rpc

# 2. 初始化依赖
make init

# 3. 编译
make build

# 4. 运行
make run

# 5. 测试
make test
```

**注意**: 必须先执行 `make gen-rpc` 生成 proto 代码，再执行 `make init` 初始化依赖。

## API

### SendCode - 发送验证码

请求:
```protobuf
message SendCodeRequest {
  string phone = 1;      // 手机号
  string scene = 2;      // 场景: login, register
}
```

响应:
```protobuf
message SendCodeResponse {
  bool success = 1;      // 是否成功
  string message = 2;    // 消息
}
```

### VerifyCode - 验证验证码

请求:
```protobuf
message VerifyCodeRequest {
  string phone = 1;      // 手机号
  string code = 2;       // 验证码
}
```

响应:
```protobuf
message VerifyCodeResponse {
  bool valid = 1;        // 是否有效
  string message = 2;    // 消息
}
```

## 注意事项

- 验证码有效期: 5分钟
- 验证码长度: 6位数字
- 验证成功后验证码自动删除

# SMS 服务限制配置说明

## 配置项

在 `sms.yaml` 或 `sms-dev.yaml` 中添加 `SmsLimit` 配置段：

```yaml
# 短信限制配置
SmsLimit:
  CodeExpire: 300      # 验证码过期时间（秒）
  SendInterval: 60     # 发送间隔时间（秒）
  DailyLimit: 10       # 每日发送次数限制
```

## 配置参数说明

### CodeExpire - 验证码过期时间

**作用**: 控制验证码在 Redis 中的有效期

**单位**: 秒

**推荐值**:
- **生产环境**: 300 秒（5分钟）
- **开发环境**: 300 秒（5分钟）

**影响**:
- 过短：用户可能来不及输入验证码
- 过长：安全性降低，占用更多 Redis 空间

**Redis Key**: `sms:code:{phone}`

### SendInterval - 发送间隔时间

**作用**: 控制同一手机号两次发送验证码的最小间隔

**单位**: 秒

**推荐值**:
- **生产环境**: 60 秒（1分钟）
- **开发环境**: 30 秒（方便测试）
- **严格模式**: 120 秒（2分钟）

**影响**:
- 过短：容易被恶意刷短信
- 过长：用户体验差，可能等不及

**Redis Key**: `sms:limit:{phone}`

**错误提示**: "发送过于频繁，请XX秒后重试"（XX 为配置的秒数）

### DailyLimit - 每日发送次数限制

**作用**: 控制同一手机号每天最多发送验证码的次数

**单位**: 次

**推荐值**:
- **生产环境**: 10 次/天
- **开发环境**: 100 次/天（方便测试）
- **严格模式**: 5 次/天

**影响**:
- 过小：可能限制正常用户
- 过大：防刷效果差，短信费用高

**Redis Key**: `sms:daily:{phone}`

**过期时间**: 当天结束（每天 23:59:59 自动清零）

**错误提示**: "今日发送次数已达上限"

## 环境差异配置

### 生产环境 (sms.yaml)

```yaml
SmsLimit:
  CodeExpire: 300      # 5分钟
  SendInterval: 60     # 1分钟
  DailyLimit: 10       # 10次/天
```

**适用场景**: 正式环境，注重安全和成本控制

### 开发环境 (sms-dev.yaml)

```yaml
SmsLimit:
  CodeExpire: 300      # 5分钟
  SendInterval: 30     # 30秒（更快测试）
  DailyLimit: 100      # 100次/天（方便测试）
```

**适用场景**: 本地开发测试，注重开发效率

### 严格模式（高安全场景）

```yaml
SmsLimit:
  CodeExpire: 180      # 3分钟（更短有效期）
  SendInterval: 120    # 2分钟（更长间隔）
  DailyLimit: 5        # 5次/天（更严格限制）
```

**适用场景**: 金融、支付等高安全要求场景

## 配置修改后的操作

### 1. 修改配置文件

编辑 `sms.yaml` 或 `sms-dev.yaml`:

```yaml
SmsLimit:
  CodeExpire: 600      # 改为10分钟
  SendInterval: 30     # 改为30秒
  DailyLimit: 20       # 改为20次
```

### 2. 重启服务

配置修改后需要**重启服务**才能生效：

```bash
# 停止当前服务 (Ctrl+C)

# 重新启动
cd services/sms
make run
```

### 3. 验证配置

发送验证码后，查看日志确认配置已生效：

```
验证码发送成功: phone=18610665421, daily_count=1/20, code_expire=600s, send_interval=30s
```

日志中会显示：
- `daily_count=1/20`: 今日发送1次，限制20次
- `code_expire=600s`: 验证码有效期600秒
- `send_interval=30s`: 发送间隔30秒

## 查看 Redis 中的限制信息

```bash
# 连接 Redis
redis-cli

# 查看某个手机号的所有限制信息
keys sms:*18610665421*

# 输出示例:
# sms:code:18610665421      # 验证码
# sms:limit:18610665421     # 发送频率限制
# sms:daily:18610665421     # 每日计数

# 查看验证码
get sms:code:18610665421
# 输出: "123456"

# 查看每日发送次数
get sms:daily:18610665421
# 输出: "3"

# 查看 TTL（剩余过期时间）
ttl sms:code:18610665421
# 输出: 287（还有287秒过期）

ttl sms:limit:18610665421
# 输出: 45（还有45秒才能再次发送）
```

## 默认值保护

如果配置文件中没有设置这些参数，代码会使用以下默认值：

```go
CodeExpire: 300       // 5分钟
SendInterval: 60      // 1分钟
DailyLimit: 10        // 10次
```

这确保了即使配置缺失，服务也能正常运行并有基本的保护。

## 配置建议

### 场景 1: 普通应用

```yaml
SmsLimit:
  CodeExpire: 300
  SendInterval: 60
  DailyLimit: 10
```

### 场景 2: 高频应用（如电商促销）

```yaml
SmsLimit:
  CodeExpire: 300
  SendInterval: 60
  DailyLimit: 20      # 增加每日限制
```

### 场景 3: 高安全应用（如金融支付）

```yaml
SmsLimit:
  CodeExpire: 180      # 缩短有效期
  SendInterval: 120    # 增加间隔
  DailyLimit: 5        # 严格限制
```

### 场景 4: 开发测试

```yaml
SmsLimit:
  CodeExpire: 300
  SendInterval: 10     # 快速测试
  DailyLimit: 999      # 不限制
```

## 监控建议

建议监控以下指标：

1. **每日发送总量** - 统计所有用户的发送次数
2. **触发限制次数** - 统计有多少请求被限制拦截
3. **验证码使用率** - 发送数量 vs 验证成功数量
4. **平均验证时间** - 从发送到验证的时间

可以通过日志分析或添加指标收集来实现。

## 常见问题

### Q: 修改配置后立即生效吗？

A: 不会。需要重启服务。未来可以考虑支持热更新（从 etcd 读取配置）。

### Q: 如果我想临时放开限制怎么办？

A: 有两种方式：
1. 修改配置文件，设置大数值（如 DailyLimit: 999）
2. 直接删除 Redis 中的限制 key

```bash
redis-cli
del sms:limit:18610665421
del sms:daily:18610665421
```

### Q: 不同环境可以用不同的配置吗？

A: 可以！这就是为什么我们有 `sms.yaml`（生产）和 `sms-dev.yaml`（开发）两个配置文件。

### Q: 如何知道用户还需要等多久才能再次发送？

A: 可以查询 Redis 中的 TTL：

```bash
redis-cli
ttl sms:limit:18610665421
# 输出剩余秒数
```

未来可以考虑在 API 响应中返回这个信息。
