# ME2 Common Package

ME2 项目的公共库，包含所有服务共享的代码。

## 模块

- `errcode` - 统一错误码定义
- `response` - 统一响应格式
- `middleware` - 中间件 (认证、日志、恢复等)
- `utils` - 工具函数
- `llmclient` - LLM 客户端封装
- `vectorclient` - 向量数据库客户端
- `ossclient` - OSS 客户端封装
- `smsclient` - 短信客户端封装

## 使用方式

在各个服务的 go.mod 中引入：

```go
require (
    github.com/me2/pkg v0.0.1
)
```

如果是本地开发，使用 replace：

```go
replace github.com/me2/pkg => ../../../pkg
```

## 发布

当 pkg 有更新时：

1. 提交代码到 pkg 仓库
2. 打 tag: `git tag v0.0.x`
3. 推送: `git push origin v0.0.x`
4. 各服务更新依赖: `go get github.com/me2/pkg@v0.0.x`

## 开发规范

- 保持向后兼容
- 充分测试
- 添加文档注释
- 遵循 Go 代码规范
