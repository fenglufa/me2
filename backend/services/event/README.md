# Event Service (事件生成服务)

## 服务概述

Event Service 负责管理事件模板、生成个性化事件内容，并存储事件历史。是 Me2 项目中将分身行为转化为具体故事内容的核心服务。

## 核心功能

### 1. 事件生成 (GenerateEvent)
- 接收 Action Service 的行为请求
- 根据行为类型选择合适的事件模板
- 调用 World Service 获取场景详细信息
- 调用 AI Service 生成个性化事件内容
- 保存事件历史
- 返回完整事件数据

### 2. 事件时间线 (GetEventTimeline)
- 分页查询分身的事件历史
- 按时间倒序展示
- 支持总数统计

### 3. 事件详情 (GetEventDetail)
- 获取单个事件的详细信息

### 4. 模板管理 (GetTemplates)
- 查询事件模板列表
- 支持按分类和稀有度筛选

## 技术架构

### 目录结构
```
event/
├── Makefile                # 构建管理
├── go.mod                  # Go 模块定义
├── README.md               # 服务文档
└── rpc/                    # RPC 服务
    ├── event.proto        # Protobuf 定义
    ├── event.sql          # 数据库表结构及初始数据
    ├── event.go           # 服务入口
    ├── etc/               # 配置文件
    │   └── event-dev.yaml
    ├── internal/
    │   ├── config/        # 配置结构
    │   ├── logic/         # 业务逻辑
    │   │   ├── generate_event_logic.go    # 事件生成逻辑
    │   │   ├── get_event_timeline_logic.go
    │   │   ├── get_event_detail_logic.go
    │   │   ├── get_templates_logic.go
    │   │   └── template_helper.go         # 模板辅助函数
    │   ├── model/         # 数据模型
    │   │   ├── event_template_model.go
    │   │   └── events_history_model.go
    │   ├── server/        # gRPC 服务器
    │   └── svc/           # 服务上下文
    └── event_client/      # 客户端库
```

### 依赖服务
- **Avatar Service**: 获取分身信息和人格特征
- **World Service**: 场景信息查询
- **AI Service**: AI文本生成（MVP阶段使用模板生成）
- **MySQL**: 存储事件模板和历史
- **Etcd**: 服务发现

## 数据库表结构

### event_templates (事件模板表)
```sql
CREATE TABLE `event_templates` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `category` VARCHAR(50) NOT NULL,          -- 分类：exploration/social/study/creative/rest/play
    `name` VARCHAR(100) NOT NULL,             -- 模板名称
    `description` TEXT,                       -- 模板描述
    `trigger_conditions` JSON,                -- 触发条件
    `rarity` VARCHAR(20) DEFAULT 'common',    -- 稀有度: common/rare/epic
    `cooldown_hours` INT DEFAULT 0,           -- 冷却时间
    `content_template` TEXT NOT NULL,         -- AI Prompt 模板
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `idx_category` (`category`),
    INDEX `idx_rarity` (`rarity`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### events_history (事件历史表)
```sql
CREATE TABLE `events_history` (
    `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
    `avatar_id` BIGINT NOT NULL,              -- 分身ID
    `template_id` BIGINT NOT NULL,            -- 模板ID
    `event_type` VARCHAR(50) NOT NULL,        -- 事件类型
    `event_title` VARCHAR(200) NOT NULL,      -- 事件标题
    `event_text` TEXT NOT NULL,               -- 事件描述
    `image_url` VARCHAR(500) DEFAULT '',      -- 事件配图（MVP暂为空）
    `scene_id` BIGINT NOT NULL,               -- 场景ID
    `scene_name` VARCHAR(255) NOT NULL,       -- 场景名称
    `personality_changes` JSON,               -- 性格变化（后续实现）
    `occurred_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX `idx_avatar_id` (`avatar_id`),
    INDEX `idx_occurred_at` (`occurred_at`),
    INDEX `idx_event_type` (`event_type`),
    INDEX `idx_template_id` (`template_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 配置说明

### event-dev.yaml
```yaml
Name: event.rpc
ListenOn: 0.0.0.0:8086

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: event.rpc

Mysql:
  DataSource: root:123456@tcp(127.0.0.1:3306)/me?charset=utf8mb4&parseTime=true&loc=Local

AvatarRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: avatar.rpc

AIRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: ai.rpc

WorldRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: world.rpc
```

## 使用方法

### 初始化
```bash
# 1. 生成 RPC 代码（必须先执行）
make gen-rpc

# 2. 初始化依赖
make init

# 3. 创建数据库表并插入初始模板
mysql -u root -p < rpc/event.sql
```

### 构建
```bash
make build
```

### 运行
```bash
# 开发环境
make run-dev

# 生产环境
make run
```

### 测试
```bash
make test
```

## API 接口

### 1. GenerateEvent - 生成事件
```protobuf
rpc GenerateEvent(GenerateEventRequest) returns (GenerateEventResponse);
```

**请求参数**:
- `avatar_id`: 分身ID
- `action_type`: 行为类型（exploration/social/study/creative/rest/play）
- `scene_id`: 场景ID
- `scene_name`: 场景名称

**响应**:
- `event_id`: 事件ID
- `event_type`: 事件类型
- `event_title`: 事件标题
- `event_text`: 事件内容
- `image_url`: 图片URL（MVP暂为空）

### 2. GetEventTimeline - 获取事件时间线
```protobuf
rpc GetEventTimeline(GetEventTimelineRequest) returns (GetEventTimelineResponse);
```

**请求参数**:
- `avatar_id`: 分身ID
- `page`: 页码（从1开始）
- `page_size`: 每页数量（1-100，默认20）

**响应**:
- `events`: 事件列表
- `total`: 总数

### 3. GetEventDetail - 获取事件详情
```protobuf
rpc GetEventDetail(GetEventDetailRequest) returns (GetEventDetailResponse);
```

### 4. GetTemplates - 获取模板列表
```protobuf
rpc GetTemplates(GetTemplatesRequest) returns (GetTemplatesResponse);
```

## MVP 实现说明

### 当前已实现
1. ✅ 事件模板管理（数据库存储，14个初始模板）
2. ✅ 事件生成核心接口
3. ✅ 事件历史查询（时间线、详情）
4. ✅ 模板选择逻辑
5. ✅ AI Service 集成（AI生成个性化事件内容）
6. ✅ World Service 集成（获取场景详细信息）

### 待后续实现
1. ⏸ ImageGen Service 集成（生成事件配图）
2. ⏸ 事件对性格的影响计算
3. ⏸ 模板的 CRUD 管理接口
4. ⏸ 事件稀有度和冷却时间机制

## 核心工作流程

```
Action Service (决定行为：exploration)
    ↓
Event Service 接收 GenerateEvent 请求
    ↓
1. 调用 Avatar Service 获取分身信息
    ↓
2. 调用 World Service 获取场景详细信息
    ↓
3. 从数据库查询该类型的所有模板
    ↓
4. 随机选择一个模板
    ↓
5. 调用 AI Service 生成个性化事件内容
   (基于分身性格和场景描述)
    ↓
6. 保存事件历史到数据库
    ↓
7. 返回完整事件数据
    ↓
(可选) Notify Service 推送给用户
```

## 事件类型和模板

### 支持的行为类型
1. **exploration** (探索): 森林探险、城市漫步、神秘发现
2. **social** (社交): 偶遇交谈、群体活动
3. **study** (学习): 知识探索、深度思考
4. **creative** (创作): 艺术创作、灵感迸发
5. **rest** (休息): 静心休息、冥想思考
6. **play** (娱乐): 快乐玩耍、趣味发现

### 模板示例
```
分身{{avatar_name}}来到{{scene_name}}，开始了一场探索之旅。
请根据分身的性格特征（冒险倾向{{adventurous}}，生活动力{{energetic}}），
生成一段生动有趣的探索事件描述，字数控制在150-200字。
```

## 注意事项

1. **AI Service 集成**: 已集成 AI Service，根据分身性格和场景描述生成个性化事件内容
2. **World Service 集成**: 已集成 World Service，获取场景的详细信息（名称、描述）
3. **图片生成**: MVP阶段 image_url 为空，待 ImageGen Service 开发
4. **性格影响**: personality_changes 字段暂时为空，后续实现
5. **AI 响应解析**: 使用 `parseAIResponse` 函数从 AI 生成的文本中提取标题和内容
6. **并发性能**: 事件生成为同步操作，高并发场景需要优化

## 开发规范

- 代码文件使用下划线命名
- 每个文件不超过 500 行
- 使用 go_zero 代码风格
- 所有 RPC 方法都有对应的 Logic 层实现
- 数据库操作封装在 Model 层

## 编译问题说明

~~由于生成的代码中import路径存在问题，需要手动修正以下文件中的import语句~~（已修复）

goctl生成的代码中import路径问题已经修复，服务可以正常编译运行。

## 后续优化

1. 添加事件模板的版本管理
2. 实现事件对性格的影响算法
3. 添加事件推荐算法（基于分身性格和历史）
4. 支持事件链和复合事件
5. 优化事件生成性能（缓存、异步）
6. 集成 ImageGen Service 生成事件配图

## 版本信息

- **版本**: v0.1.0 (MVP)
- **创建日期**: 2025-12-04
- **状态**: 可用（核心功能已实现并通过编译）
