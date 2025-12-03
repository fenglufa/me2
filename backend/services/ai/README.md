# AI Service

AI 统一服务，封装 Deepseek API，提供统一的 AI 调用接口、Prompt 管理、缓存和监控功能。

## 功能

- **统一 AI 调用** - 封装 Deepseek API，提供统一接口
- **Prompt 模板管理** - 集中管理不同场景的 Prompt 模板
- **智能缓存** - Redis 缓存相同请求，节省成本
- **调用监控** - 记录所有调用日志，统计 Token 使用和成本
- **情绪分析** - 分析文本情绪
- **行动意图计算** - 基于性格向量、时间、状态计算分身下一步行动
- **文本向量化** - 支持文本向量化（待实现）

## 技术栈

- go-zero RPC
- Deepseek API (deepseek-chat)
- MySQL (调用日志存储)
- Redis (缓存)
- Etcd (服务注册与发现)

## 配置

配置文件: `rpc/etc/ai-dev.yaml`

```yaml
Name: ai.rpc
ListenOn: 0.0.0.0:8005

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: ai.rpc

Mysql:
  DataSource: root:123456@tcp(127.0.0.1:3306)/me?charset=utf8mb4&parseTime=true&loc=Local

RedisConf:
  Host: 127.0.0.1:6379
  Type: node
  Pass: ""

Deepseek:
  ApiKey: "your-deepseek-api-key"
  BaseURL: "https://api.deepseek.com"
  Model: "deepseek-chat"
  Timeout: 60
```

## 数据库

### 创建数据库表

```bash
mysql -u root -p < rpc/ai.sql
```

### 表结构

```sql
CREATE TABLE ai_call_logs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    service_name VARCHAR(50) NOT NULL,
    scene_type VARCHAR(50) NOT NULL,
    user_id BIGINT DEFAULT 0,
    avatar_id BIGINT DEFAULT 0,
    input_tokens INT DEFAULT 0,
    output_tokens INT DEFAULT 0,
    cost BIGINT DEFAULT 0,
    duration_ms BIGINT DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## 使用

```bash
# 1. 生成 RPC 代码
make gen-rpc

# 2. 初始化依赖
make init

# 3. 创建数据库表
mysql -u root -p < rpc/ai.sql

# 4. 配置 Deepseek API Key
# 编辑 rpc/etc/ai-dev.yaml，填入你的 API Key

# 5. 编译
make build

# 6. 运行（开发模式）
make run-dev
```

## API

### CalculateActionIntent - 行动意图计算（P0 - 已实现）⭐

**用途**: 供 Action Service 调用，决定分身下一步应该做什么

请求:
```protobuf
message ActionIntentRequest {
  int64 avatar_id = 1;
  PersonalityVector personality = 2;  // 6维人格向量
  TimeContext time_context = 3;        // 时间上下文
  repeated string recent_interactions = 4;
  repeated string recent_events = 5;
  AvatarState current_state = 6;       // 当前状态
}
```

响应:
```protobuf
message ActionIntentResponse {
  map<string, float> action_scores = 1;  // 各行为得分 0-100
  string recommended_action = 2;          // 推荐行为
  string reason = 3;                      // 推荐理由
  int32 input_tokens = 4;
  int32 output_tokens = 5;
}
```

**行为类别**:
- `exploration` - 探索新地方
- `social` - 社交互动
- `study` - 学习知识
- `creative` - 艺术创作
- `rest` - 休息放松
- `play` - 游戏娱乐
- `respond` - 回应用户

### Chat - 对话生成

请求:
```protobuf
message ChatRequest {
  string prompt_template = 1;  // Prompt 模板名称
  map<string, string> variables = 2;  // 变量
  ModelConfig model_config = 3;  // 模型配置
  int64 user_id = 4;
  int64 avatar_id = 5;
}
```

响应:
```protobuf
message ChatResponse {
  string content = 1;  // 生成的内容
  int32 input_tokens = 2;
  int32 output_tokens = 3;
  int64 duration_ms = 4;
}
```

### Generate - 文本生成

与 Chat 接口相同，用于生成事件、日记等文本。

### AnalyzeEmotion - 情绪分析

请求:
```protobuf
message AnalyzeEmotionRequest {
  string text = 1;  // 要分析的文本
  int64 user_id = 2;
}
```

响应:
```protobuf
message AnalyzeEmotionResponse {
  string emotion = 1;  // 情绪类型
  float score = 2;  // 情绪强度 0-1
  string analysis = 3;  // 分析说明
}
```

### GetStats - 获取使用统计

请求:
```protobuf
message GetStatsRequest {
  int64 user_id = 1;  // 用户ID（0表示全部）
  string start_date = 2;  // 开始日期
  string end_date = 3;  // 结束日期
}
```

响应:
```protobuf
message GetStatsResponse {
  int64 total_calls = 1;
  int64 total_input_tokens = 2;
  int64 total_output_tokens = 3;
  int64 total_cost = 4;  // 总成本（分）
  int64 success_calls = 5;
  int64 error_calls = 6;
  int64 avg_duration_ms = 7;
}
```

## Prompt 模板

Prompt 模板现在支持从配置文件加载，方便随时调整而无需重新编译代码。

### 配置文件

模板配置文件：`rpc/etc/prompts.yaml`

在服务配置文件中指定 Prompt 配置文件路径：
```yaml
# rpc/etc/ai-dev.yaml
PromptFile: "etc/prompts.yaml"
```

### 模板格式

```yaml
templates:
  template_name:
    name: "template_name"
    description: "模板描述"
    system_prompt: |
      系统提示词，支持变量替换
      使用 {{.variable_name}} 语法
    user_prompt: "用户提示词 {{.input}}"
```

### 当前支持的模板

1. **avatar_chat** - 分身对话
2. **event_story** - 事件故事生成
3. **avatar_diary** - 分身日记生成
4. **user_diary_reply** - 用户日记回应
5. **emotion_analysis** - 情绪分析

### 修改 Prompt

1. 编辑 `rpc/etc/prompts.yaml` 文件
2. 修改对应模板的 `system_prompt` 或 `user_prompt`
3. 重启服务使配置生效

注意：变量名需要与调用时传入的 variables 对应。

## 成本计算

Deepseek 价格：
- 输入: 0.001元/1K tokens
- 输出: 0.002元/1K tokens

成本以分为单位存储在数据库中。

## 依赖服务

- **Deepseek API** - AI 模型服务
- **MySQL** - 调用日志存储
- **Redis** - 缓存
- **Etcd** - 服务注册与发现

## 注意事项

- 需要配置有效的 Deepseek API Key
- 所有 AI 调用都会记录日志用于监控和成本统计
- 建议使用 Redis 缓存减少重复调用
- Prompt 模板支持变量替换，使用 `{{.variable_name}}` 语法

## 目录结构

```
ai/
├── Makefile
├── README.md
├── go.mod
└── rpc/
    ├── ai.proto
    ├── ai.sql
    ├── etc/
    │   └── ai-dev.yaml
    └── internal/
        ├── config/
        │   └── config.go
        ├── deepseek/
        │   └── client.go
        ├── logic/
        │   ├── chat_logic.go
        │   ├── generate_logic.go
        │   ├── analyze_emotion_logic.go
        │   ├── embedding_logic.go
        │   └── get_stats_logic.go
        ├── model/
        │   └── ai_call_log_model.go
        ├── prompt/
        │   ├── templates.go
        │   └── renderer.go
        └── svc/
            └── service_context.go
```

## 扩展

### 添加新的 Prompt 模板

编辑 `rpc/internal/prompt/templates.go`，添加新模板：

```go
"new_template": {
    Name:        "new_template",
    Description: "新模板描述",
    SystemPrompt: "系统提示词 {{.variable}}",
    UserPromptTemplate: "用户提示词 {{.input}}",
},
```

### 支持其他 AI 服务商

修改 `rpc/internal/deepseek/client.go`，实现新的客户端接口。

## 文本向量化（Embedding）实现规划

### 什么是文本向量化

文本向量化（Embedding）是将文本转换为数值向量的过程，用于：
- **语义捕捉** - 相似语义的文本具有相似的向量表示
- **相似度计算** - 通过余弦相似度等方法计算文本相关性
- **降维表示** - 将复杂文本用固定维度向量表示

### 应用场景

在 AI 分身系统中，Embedding 用于：

1. **语义搜索** - Memory Service 中根据语义匹配相关记忆
2. **事件匹配** - Action Service 中找到相似的历史事件
3. **上下文理解** - Dialogue Service 中理解对话主题和上下文
4. **聚类分析** - 将相似事件/对话分组归类

### 实现方案

#### 方案 1: 使用 Deepseek Embedding API

如果 Deepseek 提供 embedding 接口，直接集成：

**优点:**
- 与现有 Chat API 一致
- 无需额外服务部署

**缺点:**
- 需要确认 Deepseek 是否支持
- 每次调用产生成本

#### 方案 2: 使用第三方 Embedding 服务

集成主流 Embedding API：
- **OpenAI** - text-embedding-3-small (1536维) / text-embedding-3-large (3072维)
- **智谱 AI** - embedding-2 模型
- **通义千问** - text-embedding-v1/v2
- **百度文心** - Embedding-V1

**优点:**
- 成熟稳定，效果好
- 维度可选，性能可控

**缺点:**
- 增加服务依赖
- 持续产生调用成本

#### 方案 3: 本地部署 Embedding 模型（推荐）

使用开源模型本地部署：
- **sentence-transformers** - 如 `paraphrase-multilingual-mpnet-base-v2` (768维)
- **ONNX Runtime** - 部署优化后的轻量级模型
- **text2vec** - 中文优化的向量化模型

**优点:**
- 零调用成本
- 响应速度快（本地计算）
- 数据隐私保护

**缺点:**
- 需要服务器资源（CPU/GPU）
- 需要维护模型服务

### 实现步骤

1. **选择 Embedding 服务** - 根据成本、性能、隐私要求选择方案
2. **扩展 Client** - 在 `internal/deepseek/` 或新建 `internal/embedding/` 实现客户端
3. **实现 Logic** - 完善 `embedding_logic.go` 的实际调用逻辑
4. **添加缓存** - 相同文本的向量结果缓存到 Redis
5. **性能监控** - 记录向量化耗时和成本

### 配置示例

未来配置文件可能需要添加：

```yaml
# rpc/etc/ai-dev.yaml
Embedding:
  Provider: "openai"  # openai/zhipu/local
  ApiKey: "your-api-key"
  Model: "text-embedding-3-small"
  Dimension: 1536

  # 或本地部署
  # Provider: "local"
  # Endpoint: "http://localhost:8080/embed"
```

### 当前状态

- ✅ Proto 定义已完成（`EmbeddingRequest` / `EmbeddingResponse`）
- ✅ RPC 接口已生成
- ⏳ 实现逻辑待完成（目前返回空向量）
- ⏳ 缓存机制待实现
- ⏳ 成本统计待实现

**优先级**: P2（中等优先级）- 在核心功能完成后实现
