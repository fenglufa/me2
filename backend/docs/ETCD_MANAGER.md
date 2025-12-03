# Etcd 管理脚本使用说明

## 脚本位置

`backend/scripts/etcd-manager.sh`

## 快速使用

```bash
# 从项目根目录运行
cd backend

# 列出所有已注册的服务
./scripts/etcd-manager.sh list

# 查看服务详情
./scripts/etcd-manager.sh show sms.rpc

# 监控服务注册变化
./scripts/etcd-manager.sh watch

# 清理离线服务
./scripts/etcd-manager.sh cleanup
```

## 所有命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `list` | 列出所有已注册的服务 | `./scripts/etcd-manager.sh list` |
| `show <service>` | 查看指定服务的详细信息 | `./scripts/etcd-manager.sh show sms.rpc` |
| `watch` | 监控服务注册变化（实时） | `./scripts/etcd-manager.sh watch` |
| `delete <service>` | 删除指定服务的注册信息 | `./scripts/etcd-manager.sh delete sms.rpc` |
| `cleanup` | 清理所有离线的服务实例 | `./scripts/etcd-manager.sh cleanup` |
| `help` | 显示帮助信息 | `./scripts/etcd-manager.sh help` |

## 示例输出

### list 命令

```
=== 已注册的服务 ===

服务名称: sms.rpc
实例数量: 1
服务地址:
  [1] 192.168.3.36:8001

服务名称: oss.rpc
实例数量: 1
服务地址:
  [1] 192.168.3.36:8002
```

### show 命令

```
=== sms.rpc 详细信息 ===

实例 #1
Key: sms.rpc/7587891231985505544
地址: 192.168.3.36:8001
Host: 192.168.3.36
Port: 8001
状态: ✓ 在线
```

## 技术说明

- **Shell**: 使用 zsh（macOS 默认 shell）
- **依赖**: 需要安装 etcdctl
- **权限**: 脚本已设置为可执行（`chmod +x`）

## 故障排查

### 错误: etcdctl 未安装

```bash
# macOS
brew install etcd

# Linux
sudo apt-get install etcd-client  # Ubuntu/Debian
sudo yum install etcd              # CentOS/RHEL
```

### 错误: 权限被拒绝

```bash
chmod +x scripts/etcd-manager.sh
```

### 错误: Etcd 连接失败

```bash
# 检查 Etcd 是否运行
etcdctl endpoint health

# 启动 Etcd
brew services start etcd  # macOS
```

## 注意事项

1. **delete 命令需要确认**: 删除服务前会提示确认，输入 `y` 确认删除
2. **watch 命令持续运行**: 按 `Ctrl+C` 退出监控
3. **cleanup 命令自动检测**: 自动检测并删除离线的服务实例
