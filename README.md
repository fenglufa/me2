# AI 分身宇宙 (AI Avatar Universe)

<div align="center">

**让每个人都拥有一个在另一个世界生活的数字生命体**

[产品设计](./产品设计文档.md) • [技术架构](./服务架构说明.md) • [开发计划](./TODO.md)

</div>

---

## 📖 项目简介

AI 分身宇宙是一款创新的 AI 陪伴产品，融合了 **AI 人格系统**、**元宇宙行为引擎** 和 **情感陪伴** 三大核心能力。用户可以创建属于自己的 AI 分身，这个分身不仅能与用户对话交流，还能在一个名为"**第二空间 (Second Sphere)**"的虚拟世界中自主生活、探索、社交和成长。

### ✨ 核心特性

- 🧠 **6维人格系统**: 情绪温度、冒险倾向、人际能量、创造性、情绪稳定性、生活动力
- 🌍 **自主行为引擎**: AI 分身每天自动探索、旅行、交友、创作
- 📝 **双向日记系统**: 分身自动生成日记，用户日记获得 AI 回应
- 💬 **深度对话系统**: 基于人格和记忆的个性化对话
- 🎨 **AI 内容生成**: 自动生成事件故事、插画、情绪图
- 📈 **成长追踪**: 性格演化、成长节点、关系进度可视化
- 🔮 **元宇宙世界**: 4大区域自动事件生成系统

### 🎯 目标用户

18-35岁城市青年，寻求情绪陪伴、自我探索和独特数字体验的用户。

---

## 🏗️ 技术架构

### 后端技术栈

- **语言**: Go 1.21+
- **框架**: go-zero (微服务框架)
- **数据库**: PostgreSQL 14+
- **缓存**: Redis 7+
- **向量数据库**: Milvus / Pinecone (RAG 长期记忆)
- **AI 引擎**: Deepseek API (统一 AI 调用)
- **云服务**: 阿里云 (短信 SMS + 对象存储 OSS)
- **消息队列**: Asynq / NSQ
- **图像生成**: Stable Diffusion / DALL-E

### 前端技术栈

- **移动端**: 待定 (React Native / Flutter / 原生)
- **Web 端**: 待定 (React / Vue)

### 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        API Gateway                          │
│          (认证、限流、路由、A/B测试)                           │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
┌───────▼────────┐   ┌────────▼────────┐   ┌───────▼────────┐
│  基础设施服务   │   │   核心业务服务   │   │   AI & 内容    │
├────────────────┤   ├─────────────────┤   ├────────────────┤
│ SMS Service    │   │ User Service    │   │ AI Service  ⭐ │
│ OSS Service    │   │ Avatar Service  │   │ Dialogue       │
└────────────────┘   │ World Service   │   │ Event Gen      │
                     │ Action Service  │   │ ImageGen       │
                     │ Diary Service   │   │ Memory (RAG)   │
                     │ Notify Service  │   └────────────────┘
                     └─────────────────┘
```

### 核心服务列表 (15个微服务)

1. **API Gateway** - 统一网关
2. **SMS Service** - 短信验证 (阿里云)
3. **OSS Service** - 对象存储 (阿里云)
4. **User Service** - 用户账号管理
5. **Avatar Service** - 分身人格与成长 ⭐
6. **Action Service** - 行动逻辑调度 ⭐
7. **World Service** - 世界地图与场景
8. **Event Service** - 事件生成引擎 ⭐
9. **Dialogue Service** - 对话系统
10. **Diary Service** - 日记系统
11. **AI Service** - 统一 AI 调用 (Deepseek) ⭐⭐⭐
12. **Memory Service** - 长期记忆 (向量数据库)
13. **ImageGen Service** - 图像生成
14. **Notify Service** - 推送通知
15. ~~StoryGen Service~~ (已整合到 AI Service)

---

## 📂 项目结构

```
me2/
├── backend/                 # Go 后端服务
│   ├── services/           # 微服务目录
│   │   ├── user/          # 用户服务
│   │   ├── avatar/        # 分身服务
│   │   ├── ai/            # AI 统一服务
│   │   ├── event/         # 事件服务
│   │   ├── action/        # 行动逻辑服务
│   │   ├── world/         # 世界服务
│   │   ├── dialogue/      # 对话服务
│   │   ├── diary/         # 日记服务
│   │   ├── memory/        # 记忆服务
│   │   ├── imagegen/      # 图像生成服务
│   │   ├── notify/        # 通知服务
│   │   ├── sms/           # 短信服务
│   │   └── oss/           # OSS服务
│   ├── gateway/            # API 网关
│   ├── pkg/                # 公共库
│   │   ├── llmclient/     # LLM 客户端封装
│   │   ├── vectorclient/  # 向量数据库客户端
│   │   └── middleware/    # 中间件
│   └── deployments/        # 部署配置
│       ├── docker/
│       └── k8s/
├── app/                     # 移动端应用
├── web/                     # Web 应用
├── docs/                    # 文档目录
├── 产品设计文档.md           # 产品设计文档
├── 服务架构说明.md           # 技术架构文档
├── TODO.md                  # 开发任务清单
└── README.md                # 本文件
```

---

## 🚀 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 14+
- Redis 7+
- Docker & Docker Compose (可选)
- Node.js 18+ (前端开发)

### 后端开发环境设置

```bash
# 1. 克隆项目
cd /Users/flf/Desktop/ProjectCode/me2

# 2. 安装 go-zero 工具
go install github.com/zeromicro/go-zero/tools/goctl@latest

# 3. 初始化后端项目 (待执行)
cd backend
# goctl api new gateway
# goctl rpc new services/user
# ...

# 4. 配置数据库
# 创建 PostgreSQL 数据库
createdb me2_dev

# 5. 配置环境变量
cp .env.example .env
# 编辑 .env 配置数据库连接、Redis、阿里云密钥等

# 6. 运行数据库迁移
make migrate

# 7. 启动服务
make run
```

### 使用 Docker Compose (推荐)

```bash
# 启动所有服务 (待配置)
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

---

## 🎨 核心概念

### 1. 第二空间 (Second Sphere)

一个由用户情绪和 AI 驱动的平行数字世界，包含以下区域：

- **日常城市**: 咖啡厅、街区、公园
- **自然区域**: 森林、湖泊、山谷
- **社交广场**: 分身随机相遇的场所
- **神秘区域**: 触发特殊事件

### 2. 6维人格系统

每个 AI 分身拥有独特的人格向量 (0-100)：

| 维度 | 说明 | 影响 |
|------|------|------|
| 情绪温度 (Warmth) | 温柔 ↔ 理性 | 对话风格、情感表达 |
| 冒险倾向 (Adventurous) | 探索 ↔ 安全 | 旅行频率、事件类型 |
| 人际能量 (Social) | 社交 ↔ 独处 | 交友倾向、互动频率 |
| 创造性 (Creative) | 创意 ↔ 结构 | 日记风格、内容生成 |
| 情绪稳定性 (Calm) | 平稳 ↔ 敏感 | 情绪反应、剧情触发 |
| 生活动力 (Energetic) | 活跃 ↔ 悠闲 | 行动频率、节奏 |

### 3. 事件生成机制

**混合模式**: 模板 (20%) + AI 生成 (80%)

- **模板**: 提供结构和世界观一致性
- **AI**: 生成个性化内容
- **触发器**: 时间、人格、用户行为、世界事件
- **稀有度**: Common / Rare / Epic / Legendary

### 4. 行动逻辑系统

分身每日自主行动，基于：

- **时间节律**: 清晨、白天、午后、夜间、深夜
- **人格驱动**: 6维人格决定行为倾向
- **用户同步**: 根据用户状态调整行动
- **世界事件**: 节日、天气、随机事件

---

## 📚 文档索引

| 文档 | 描述 |
|------|------|
| [产品设计文档](./产品设计文档.md) | 完整的产品设计、世界观、商业模式 |
| [服务架构说明](./服务架构说明.md) | 详细的技术架构和 API 设计 |
| [开发任务清单](./TODO.md) | 按优先级排列的开发任务 |
| API 文档 | 待生成 (OpenAPI/Swagger) |
| 数据库设计 | 待整理 |

---

## 🔑 核心数据表

### 主要数据表 (14个)

1. **users** - 用户账号 (手机号、订阅)
2. **sms_codes** - 短信验证码
3. **avatars** - 分身核心数据 (6维人格 + 状态)
4. **avatar_milestones** - 成长节点
5. **avatar_stats** - 统计数据
6. **event_templates** - 事件模板
7. **events_history** - 事件历史
8. **world_regions** - 世界区域
9. **world_scenes** - 世界场景
10. **avatar_locations** - 分身位置
11. **dialogues** - 对话历史
12. **diaries** - 日记 (分身 + 用户)
13. **ai_prompts** - Prompt 模板
14. **ai_call_logs** - AI 调用日志

---

## 🔐 关键设计原则

### 架构原则

1. **统一 AI 调用**: 所有服务通过 AI Service 调用 Deepseek，不直接调用外部 API
2. **身份分离**: User 表只存账号，Avatar 表存身份 (昵称、头像)
3. **行为分离**: Action Service (决策) + Event Service (生成)
4. **MVP 限制**: 一个用户一个分身

### 数据流

```
定时任务 → Action Service (决定做什么)
         ↓
    World Service (查询/更新位置)
         ↓
    Event Service (选择事件模板)
         ↓
    AI Service (生成文本) ⭐
         ↓
    ImageGen Service (生成图片)
         ↓
    OSS Service (存储)
         ↓
    Notify Service (推送)
```

---

## 🛣️ 开发路线图

### Phase 1: MVP (3个月)

- [x] 完成架构设计
- [ ] 基础设施搭建 (数据库、Redis、向量库)
- [ ] 核心服务开发 (SMS、OSS、User、Avatar)
- [ ] AI Service 开发 ⭐ (优先级最高)
- [ ] 行动与事件系统
- [ ] 基础 App UI

### Phase 2: 完整体验 (3-6个月)

- [ ] 日记系统完善
- [ ] 对话系统优化
- [ ] 世界地图扩展
- [ ] 社交功能
- [ ] 商业化上线 (订阅、虚拟物品)

### Phase 3: 硬件扩展 (6-12个月)

- [ ] 桌面硬件原型
- [ ] MQTT/蓝牙协议
- [ ] 碰一碰社交
- [ ] 设备同步

---

## 🤝 贡献指南

### 开发规范

- 遵循 Go 代码规范
- 提交前运行 `gofmt` 和 `golint`
- 编写单元测试 (覆盖率 > 80%)
- API 变更需更新文档

### 分支管理

- `main` - 生产环境
- `develop` - 开发环境
- `feature/*` - 功能分支
- `bugfix/*` - 修复分支

### 提交规范

```
feat: 新功能
fix: 修复
docs: 文档
style: 格式
refactor: 重构
test: 测试
chore: 构建/工具
```

---

## 📄 许可证

[MIT License](./LICENSE)

---

## 📞 联系方式

- 项目负责人: [待添加]
- 技术支持: [待添加]
- 商务合作: [待添加]

---

## 🙏 致谢

感谢所有为这个项目贡献想法和代码的伙伴。

---

<div align="center">

**让每个孤独、焦虑、迷茫的人，都拥有一个永远陪伴、理解并一起成长的数字生命体**

Made with ❤️ by AI Avatar Team

</div>
