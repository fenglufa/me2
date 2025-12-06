# ME2 Common Package

ME2 项目的公共库，包含所有服务共享的代码。

## 目录结构

```
pkg/
├── errcode/          # 统一错误码定义
│   └── code.go
├── response/         # 统一响应格式
│   └── response.go
├── middleware/       # HTTP 中间件
│   ├── auth.go      # JWT 认证
│   ├── logger.go    # 请求日志
│   └── recovery.go  # Panic 恢复
├── utils/           # 工具函数
│   ├── crypto.go    # 加密工具 (bcrypt, JWT, MD5)
│   ├── time.go      # 时间工具
│   └── validator.go # 验证工具
├── go.mod
└── README.md
```

## 模块说明

### 1. errcode - 错误码定义

统一的错误码和错误信息定义，覆盖所有业务场景。

```go
import "github.com/me2/pkg/errcode"

// 使用预定义错误
return nil, errcode.ErrUserNotFound

// 自定义错误
err := errcode.NewError(30001, "分身不存在")

// 错误转换
customErr := errcode.FromError(err)
```

**错误码分类:**
- 1xxxx: 通用错误
- 2xxxx: 用户相关
- 3xxxx: 分身相关
- 4xxxx: AI 相关
- 5xxxx: 事件相关
- 6xxxx: 世界相关
- 7xxxx: 对话相关
- 8xxxx: 日记相关
- 9xxxx: OSS 相关
- 11xxxx: 短信相关
- 12xxxx: 调度相关

### 2. response - 统一响应格式

标准的 HTTP JSON 响应封装。

```go
import "github.com/me2/pkg/response"

// 成功响应
response.HttpSuccess(w, data)

// 错误响应
response.HttpError(w, errcode.ErrInvalidParam)

// 分页响应
response.PageSuccess(total, page, pageSize, list)
```

**响应格式:**
```json
{
  "code": 0,
  "msg": "success",
  "data": {...}
}
```

### 3. middleware - HTTP 中间件

#### auth.go - JWT 认证中间件
```go
import "github.com/me2/pkg/middleware"

// 使用认证中间件
server.Use(middleware.JWTAuth(secret))

// 获取用户 ID
userID := middleware.GetUserID(r.Context())
```

#### logger.go - 请求日志中间件
```go
// 记录请求日志 (方法、路径、状态码、耗时等)
server.Use(middleware.Logger())
```

#### recovery.go - Panic 恢复中间件
```go
// 捕获 panic 并记录堆栈
server.Use(middleware.Recovery())
```

### 4. utils - 工具函数

#### crypto.go - 加密工具
```go
import "github.com/me2/pkg/utils"

// 密码加密
hash, _ := utils.HashPassword("password123")

// 密码验证
ok := utils.CheckPassword("password123", hash)

// 生成 JWT Token
token, _ := utils.GenerateToken(userID, secret, 7) // 7天有效期

// 解析 JWT Token
userID, _ := utils.ParseToken(token, secret)

// MD5 加密
hash := utils.MD5("hello")
```

#### time.go - 时间工具
```go
// 格式化
dateStr := utils.FormatDate(time.Now())        // 2006-01-02
dateTimeStr := utils.FormatDateTime(time.Now()) // 2006-01-02 15:04:05

// 解析
t, _ := utils.ParseDate("2024-12-06")
t, _ := utils.ParseDateTime("2024-12-06 10:30:00")

// 快捷方法
today := utils.GetToday()  // 今天日期字符串
now := utils.GetNow()      // 当前时间字符串
isToday := utils.IsToday(t) // 判断是否今天
```

#### validator.go - 验证工具
```go
// 手机号验证
ok := utils.IsValidPhone("13800138000")

// 验证码验证
ok := utils.IsValidCode("123456")

// 邮箱验证
ok := utils.IsValidEmail("user@example.com")

// 手机号脱敏
masked := utils.MaskPhone("13800138000") // 138****8000
```

## 使用方式

### 在服务中引入

在各个服务的 go.mod 中添加：

```go
require (
    github.com/me2/pkg v0.0.1
)

// 本地开发使用 replace
replace github.com/me2/pkg => ../../pkg
```

### 导入使用

```go
import (
    "github.com/me2/pkg/errcode"
    "github.com/me2/pkg/response"
    "github.com/me2/pkg/middleware"
    "github.com/me2/pkg/utils"
)
```

## 依赖

- `github.com/zeromicro/go-zero` - go-zero 框架
- `github.com/golang-jwt/jwt/v4` - JWT 认证
- `golang.org/x/crypto` - 密码加密

## 开发规范

- 保持向后兼容
- 充分测试
- 添加文档注释
- 遵循 Go 代码规范
- 不要在 pkg 中引入业务逻辑
