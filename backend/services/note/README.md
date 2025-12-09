# Note Service (笔记服务)

## 服务概述

Note Service 是 Me2 项目的笔记管理服务，提供笔记的创建、查询、分类、检索等功能。本服务是用户生活记录的核心入口，支持 AI 自动分类、TODO 管理、记账统计等功能。

## 功能特性

### MVP Phase 1 （已完成基础框架）
- ✅ 笔记 CRUD（创建、查询、更新、删除）
- ✅ 数据库模型（notes, note_todos, note_expenses）
- ✅ RPC 服务框架
- ✅ 简单分类（基于关键词）

### Phase 2 （待实现）
- ⏸ AI 智能分类（调用 AI Service）
- ⏸ 情绪分析
- ⏸ TODO 管理完整功能
- ⏸ 记账统计功能
- ⏸ 笔记检索（用于对话集成）

## 目录结构

```
note/
├── Makefile              # 构建和管理命令
├── README.md            # 本文档
├── rpc/                 # RPC 服务
│   ├── note.proto       # Proto 定义
│   ├── note.sql         # 数据库脚本
│   ├── note.go          # 服务入口
│   ├── etc/             # 配置文件
│   │   └── note-dev.yaml
│   ├── internal/
│   │   ├── config/      # 配置结构
│   │   ├── logic/       # 业务逻辑
│   │   │   ├── create_note_logic.go         ✅ 已实现
│   │   │   ├── get_notes_logic.go           ⏸ 待实现
│   │   │   ├── get_note_detail_logic.go     ⏸ 待实现
│   │   │   ├── update_note_logic.go         ⏸ 待实现
│   │   │   ├── delete_note_logic.go         ⏸ 待实现
│   │   │   ├── search_notes_logic.go        ⏸ 待实现
│   │   │   ├── get_todos_logic.go           ⏸ 待实现
│   │   │   ├── update_todo_status_logic.go  ⏸ 待实现
│   │   │   ├── get_expenses_logic.go        ⏸ 待实现
│   │   │   └── get_expense_stats_logic.go   ⏸ 待实现
│   │   ├── model/       # 数据模型
│   │   │   ├── note_model.go                ✅ 已完成
│   │   │   ├── note_todo_model.go           ✅ 已完成
│   │   │   └── note_expense_model.go        ✅ 已完成
│   │   ├── server/      # gRPC Server
│   │   └── svc/         # Service Context      ✅ 已完成
│   └── note_client/     # RPC Client
└── go.mod
```

## 数据库表

### notes (笔记主表)
```sql
CREATE TABLE `notes` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `user_id` BIGINT NOT NULL,
  `avatar_id` BIGINT DEFAULT NULL,
  `raw_text` TEXT NOT NULL,
  `ai_summary` VARCHAR(500) DEFAULT '',
  `types` JSON DEFAULT NULL,
  `emotion_data` JSON DEFAULT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_avatar_id` (`avatar_id`),
  INDEX `idx_created_at` (`created_at`),
  FULLTEXT INDEX `idx_fulltext_raw_text` (`raw_text`)
);
```

### note_todos (TODO表)
```sql
CREATE TABLE `note_todos` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `note_id` BIGINT NOT NULL,
  `user_id` BIGINT NOT NULL,
  `title` VARCHAR(500) NOT NULL,
  `due_date` DATE DEFAULT NULL,
  `status` TINYINT DEFAULT 0,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `completed_at` TIMESTAMP NULL DEFAULT NULL,
  INDEX `idx_note_id` (`note_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_status` (`status`)
);
```

### note_expenses (记账表)
```sql
CREATE TABLE `note_expenses` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `note_id` BIGINT NOT NULL,
  `user_id` BIGINT NOT NULL,
  `item` VARCHAR(200) NOT NULL,
  `amount` DECIMAL(10,2) NOT NULL,
  `category` VARCHAR(50) DEFAULT 'other',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_note_id` (`note_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_category` (`category`),
  INDEX `idx_created_at` (`created_at`)
);
```

## 配置说明

### note-dev.yaml
```yaml
Name: note.rpc
ListenOn: 0.0.0.0:8010

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: note.rpc

Mysql:
  DataSource: root:123456@tcp(127.0.0.1:3306)/me2?charset=utf8mb4&parseTime=true&loc=Local

# 依赖的 RPC 服务
AiRpc:
  Etcd:
    Hosts:
      - 127.0.0.1:2379
    Key: ai.rpc
  Timeout: 30000

# AI 配置
AI:
  Enabled: true  # MVP 阶段可设为 false
```

## 快速开始

### 1. 初始化数据库

```bash
mysql -u root -p123456 me2 < rpc/note.sql
```

### 2. 初始化依赖

```bash
make init

# 或手动执行
go mod edit -require=github.com/me2/ai@v0.0.0
go mod edit -replace=github.com/me2/ai=../ai
go mod tidy
```

### 3. 生成代码（如果修改了 proto）

```bash
make gen-rpc
```

### 4. 运行服务

开发模式:
```bash
make run-dev
```

生产模式:
```bash
make run
```

### 5. 编译

```bash
make build
```

## API 接口

### CreateNote - 创建笔记 ✅
```protobuf
rpc CreateNote(CreateNoteRequest) returns (CreateNoteResponse);
```

**已实现功能**：
- 笔记内容保存
- 简单的关键词分类（MVP）
- 数据库存储

**待完善**：
- AI 智能分类
- 自动提取 TODO
- 自动提取记账信息
- 生成摘要

### GetNotes - 获取笔记列表 ⏸
```protobuf
rpc GetNotes(GetNotesRequest) returns (GetNotesResponse);
```

### GetNoteDetail - 获取笔记详情 ⏸
```protobuf
rpc GetNoteDetail(GetNoteDetailRequest) returns (GetNoteDetailResponse);
```

### SearchNotes - 搜索笔记 ⏸
```protobuf
rpc SearchNotes(SearchNotesRequest) returns (SearchNotesResponse);
```

**说明**：此接口用于对话服务集成，允许分身检索用户笔记。

### 其他接口 ⏸
- UpdateNote
- DeleteNote
- GetTodos
- UpdateTodoStatus
- GetExpenses
- GetExpenseStats

## 当前状态

### ✅ 已完成
1. **项目框架**
   - RPC 服务结构
   - Proto 定义
   - Makefile 构建脚本
   - 配置文件

2. **数据库设计**
   - 笔记主表设计和 Model
   - TODO 表设计和 Model
   - 记账表设计和 Model
   - 索引优化

3. **核心逻辑**
   - CreateNote 基础实现（无 AI）
   - Service Context 依赖注入

### ⏸ 待实现（按优先级）

#### 高优先级（MVP 完善）
1. **get_notes_logic.go** - 获取笔记列表
   - 分页查询
   - 类型过滤
   - 时间排序

2. **search_notes_logic.go** - 搜索笔记
   - 关键词搜索
   - 时间范围过滤
   - 类型过滤

3. **update_note_logic.go** - 更新笔记
4. **delete_note_logic.go** - 删除笔记

#### 中优先级（功能扩展）
5. **AI 分类集成**
   - 创建 ai_classifier.go
   - 调用 AI Service 进行分类
   - 情绪分析
   - 自动生成摘要

6. **TODO 功能**
   - get_todos_logic.go
   - update_todo_status_logic.go
   - 从笔记中提取 TODO

7. **记账功能**
   - get_expenses_logic.go
   - get_expense_stats_logic.go
   - 从笔记中提取记账信息
   - 统计分析

#### 低优先级（优化）
8. **性能优化**
   - Redis 缓存
   - 全文搜索优化
   - 批量操作

9. **高级功能**
   - 笔记标签
   - 笔记分类
   - 笔记导出

## 实现示例

### 创建笔记（已实现）

```go
// 调用示例
resp, err := noteClient.CreateNote(ctx, &note.CreateNoteRequest{
    UserId:   1,
    AvatarId: 5,
    RawText:  "今天很开心，吃了蛋炒饭花了10块钱",
})

// MVP 返回：
// {
//   note_id: 123,
//   types: ["note"],
//   ai_summary: "",
//   emotion_data: {primary: "neutral", score: 0.5}
// }

// Phase 2 期望返回：
// {
//   note_id: 123,
//   types: ["emotion", "expense"],
//   ai_summary: "记录了心情和消费",
//   emotion_data: {primary: "happy", score: 0.8}
// }
```

### 获取笔记列表（待实现）

```go
// 实现参考
func (l *GetNotesLogic) GetNotes(in *note.GetNotesRequest) (*note.GetNotesResponse, error) {
    // 1. 参数验证
    // 2. 调用 Model 查询
    notes, err := l.svcCtx.NoteModel.FindByUserId(
        in.UserId,
        in.Page,
        in.PageSize,
        in.Types,
    )
    // 3. 统计总数
    // 4. 转换并返回
}
```

## 开发规范

- 所有文件名使用下划线命名(snake_case)
- 单个文件不超过 500 行
- 使用 `--style go_zero` 生成代码
- 所有 RPC 方法都需要验证用户权限
- 数据库操作封装在 Model 层

## 测试

```bash
make test
```

## 代码检查

```bash
make lint
```

## 格式化代码

```bash
make fmt
```

## 清理构建文件

```bash
make clean
```

## 注意事项

1. **启动前准备**
   - 确保 MySQL、Etcd 已启动
   - 确保 AI Service 已启动（如果启用 AI 分类）
   - 首次运行需要执行数据库迁移

2. **配置文件**
   - 端口 8010 不要与其他服务冲突
   - 数据库连接信息需正确配置
   - AI.Enabled 可设为 false 跳过 AI 功能

3. **开发建议**
   - MVP 阶段先实现基础 CRUD，不依赖 AI
   - Phase 2 再集成 AI 智能分类
   - 保持代码简洁，单一职责

## 下一步工作

### 立即可做
1. 实现 GetNotes 接口（获取笔记列表）
2. 实现 SearchNotes 接口（搜索笔记）
3. 测试创建笔记功能

### 短期目标（1周内）
1. 完成所有 CRUD 接口
2. 集成 AI Service 进行智能分类
3. 实现 TODO 基础功能
4. 实现记账基础功能

### 中期目标（2周内）
1. 对话服务集成笔记检索
2. 记账统计和分析
3. 情绪趋势分析
4. 性能优化

## 版本信息

- **版本**: v0.1.0 (MVP Framework)
- **Go 版本**: 1.24.0+
- **go-zero 版本**: 1.9.3+
- **创建日期**: 2025-12-09
- **状态**: 开发中（基础框架已完成）

## 作者

ME2 Team

## 许可证

MIT

---

## 快速开发指南

如果你要继续开发此服务，建议按以下顺序：

1. **先实现 GetNotes** (get_notes_logic.go)
   - 复制 create_note_logic.go 的结构
   - 调用 NoteModel.FindByUserId
   - 返回笔记列表

2. **再实现 SearchNotes** (search_notes_logic.go)
   - 调用 NoteModel.Search
   - 支持关键词和时间范围

3. **集成 AI Service**
   - 在 create_note_logic.go 中调用 AI RPC
   - 实现真正的智能分类

4. **完善其他接口**
   - Update, Delete
   - TODO 管理
   - 记账统计

5. **测试和优化**
   - 编写单元测试
   - 性能优化
   - 错误处理完善
