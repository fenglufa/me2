# Action Service (行动逻辑服务)

## 服务概述

Action Service 是 Me2 项目的核心玩法服务，负责管理 AI 分身的行动逻辑。基于分身的六维人格特征和时间因素，计算不同行为类型的意图得分，并自动调度分身执行相应的行动。

## 核心功能

### 1. 行动调度 (ScheduleAction)
- 由定时任务调用，触发分身自动行动
- 根据人格特征计算意图得分
- 选择最佳行为类型
- 推荐适合的场景
- 记录行动日志

### 2. 意图计算 (CalculateActionIntent)
- 计算 6 种行为类型的意图得分
- 基于六维人格特征和时间因素
- 返回所有意图得分和推荐行为

### 3. 行动历史查询
- GetActionHistory: 分页查询分身的行动历史
- GetLastAction: 获取分身最近一次行动

## 行为类型

服务支持 6 种行为类型：

1. **exploration (探索)**: 冒险倾向和生活动力影响，适合上午和下午
2. **social (社交)**: 人际能量和情绪温度影响，适合下午和晚上
3. **study (学习)**: 结构化倾向和情绪稳定性影响，适合上午和下午
4. **creative (创作)**: 创造性和独处倾向影响，适合下午和晚上
5. **rest (休息)**: 温和倾向和情绪稳定性影响，适合中午和晚上
6. **play (娱乐)**: 生活动力和人际能量影响，适合晚上

## 意图计算算法

### 六维人格特征
- warmth: 情绪温度 (0-100)
- adventurous: 冒险倾向 (0-100)
- social: 人际能量 (0-100)
- creative: 创造性 (0-100)
- calm: 情绪稳定性 (0-100)
- energetic: 生活动力 (0-100)

### 计算逻辑
每种行为类型的得分由以下因素决定：
1. 相关人格特征的加权得分
2. 时间因素加成
3. 最终得分限制在 0-100 之间

## 技术架构

### 目录结构
```
action/
├── rpc/                    # RPC 服务
│   ├── action.proto       # Protobuf 定义
│   ├── action.go          # 服务入口
│   ├── action.sql         # 数据库表结构
│   ├── etc/               # 配置文件
│   │   └── action-dev.yaml
│   ├── internal/
│   │   ├── config/        # 配置结构
│   │   ├── logic/         # 业务逻辑
│   │   │   ├── intent_calculator.go          # 意图计算器
│   │   │   ├── schedule_action_logic.go      # 调度逻辑
│   │   │   ├── calculate_action_intent_logic.go
│   │   │   ├── get_action_history_logic.go
│   │   │   ├── get_last_action_logic.go
│   │   │   └── converter.go                  # 数据转换
│   │   ├── model/         # 数据模型
│   │   │   └── action_log_model.go
│   │   ├── server/        # gRPC 服务器
│   │   └── svc/           # 服务上下文
│   └── action_client/     # 客户端库
└── Makefile               # 构建管理
```

### 依赖服务
- Avatar Service: 获取分身信息和人格特征
- World Service: 获取场景推荐
- MySQL: 存储行动日志
- Etcd: 服务发现

## 数据库表结构

### action_logs
```sql
CREATE TABLE `action_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `avatar_id` bigint NOT NULL COMMENT '分身ID',
  `action_type` varchar(50) NOT NULL COMMENT '行为类型',
  `scene_id` bigint NOT NULL COMMENT '场景ID',
  `scene_name` varchar(255) NOT NULL COMMENT '场景名称',
  `intent_score` float NOT NULL COMMENT '意图得分',
  `trigger_reason` text NOT NULL COMMENT '触发原因',
  `event_id` bigint DEFAULT '0' COMMENT '关联事件ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_avatar_id` (`avatar_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='行动日志表';
```

## 配置说明

### action-dev.yaml
```yaml
Name: action.rpc
ListenOn: 0.0.0.0:8084

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: action.rpc

Mysql:
  DataSource: root:password@tcp(127.0.0.1:3306)/me2?charset=utf8mb4&parseTime=true

AvatarRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: avatar.rpc

WorldRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: world.rpc

ActionSchedule:
  MinIntervalHours: 2
  MaxIntervalHours: 6
```

## 使用方法

### 初始化
```bash
make init
```

### 构建
```bash
make build
```

### 运行开发环境
```bash
make run-dev
```

### 重新生成 RPC 代码
```bash
make gen-rpc
```

### 清理
```bash
make clean
```

## API 接口

### 1. ScheduleAction
调度分身行动
```protobuf
rpc ScheduleAction(ScheduleActionRequest) returns (ScheduleActionResponse);
```

### 2. CalculateActionIntent
计算行动意图
```protobuf
rpc CalculateActionIntent(CalculateActionIntentRequest) returns (CalculateActionIntentResponse);
```

### 3. GetActionHistory
获取行动历史
```protobuf
rpc GetActionHistory(GetActionHistoryRequest) returns (GetActionHistoryResponse);
```

### 4. GetLastAction
获取最近一次行动
```protobuf
rpc GetLastAction(GetLastActionRequest) returns (GetLastActionResponse);
```

## 开发规范

- 代码文件使用下划线命名
- 每个文件不超过 500 行
- 使用 go_zero 代码风格
- 所有 RPC 方法都有对应的 Logic 层实现

## 注意事项

1. 服务启动前需要确保 MySQL、Etcd、Avatar Service 和 World Service 已启动
2. 需要先执行 action.sql 创建数据库表
3. 行动调度建议通过定时任务（如 cron）定期调用
4. event_id 字段暂时为 0，后续由 Event Service 更新

## 后续优化

1. 添加更多行为类型
2. 优化意图计算算法
3. 支持自定义行为权重
4. 添加行为冷却时间
5. 实现行为链和组合行为
