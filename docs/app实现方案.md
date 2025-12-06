# ME2 App 实现方案

## 项目概述

基于 Flutter 开发的 AI 分身移动应用，实现用户与 AI 分身的互动体验。

## 技术栈

- **框架**: Flutter 3.x
- **状态管理**: Riverpod 2.x
- **网络请求**: Dio + Retrofit
- **本地存储**: SharedPreferences + Hive
- **WebSocket**: web_socket_channel
- **路由**: go_router
- **图表**: fl_chart
- **动画**: lottie

## 项目结构

```
app/
├── lib/
│   ├── core/                    # 核心基础设施
│   │   ├── api/                # API 客户端
│   │   │   ├── api_client.dart
│   │   │   ├── interceptors/
│   │   │   └── endpoints.dart
│   │   ├── auth/               # 认证管理
│   │   │   ├── auth_provider.dart
│   │   │   └── token_storage.dart
│   │   ├── storage/            # 本地存储
│   │   │   └── storage_service.dart
│   │   ├── websocket/          # WebSocket 连接
│   │   │   └── ws_client.dart
│   │   └── constants/          # 常量配置
│   │       └── app_config.dart
│   ├── features/               # 功能模块
│   │   ├── auth/              # 登录注册
│   │   │   ├── data/
│   │   │   ├── domain/
│   │   │   └── presentation/
│   │   ├── home/              # 首页（分身动态）
│   │   ├── avatar/            # 分身模块
│   │   ├── world/             # 世界模块
│   │   ├── dialogue/          # 对话功能
│   │   └── profile/           # 个人中心
│   ├── shared/                # 共享组件
│   │   ├── widgets/           # 通用组件
│   │   └── models/            # 数据模型
│   └── main.dart
└── pubspec.yaml
```

## MVP 核心功能

### 1. 首页 - 分身实时动态
- 展示分身当前状态（散步/看书/旅行等）
- 动画效果展示分身行为
- 氛围背景（昼夜/天气）
- 分身主动消息推送
- 对话入口

**API:**
- `GET /api/v1/avatar/my`
- `GET /api/v1/events/timeline`
- `GET /api/v1/action/last`

### 2. 分身模块 - 人格与成长
- 性格面板: 6 维雷达图
- 分身日记列表
- 成长节点展示
- 关系面板

**API:**
- `GET /api/v1/avatar/:id`
- `GET /api/v1/diary/avatar`
- `GET /api/v1/diary/stats`

### 3. 世界模块 - 第二空间
- 4 个区域卡片展示
- 场景列表
- 分身位置追踪
- 事件时间线

**API:**
- `GET /api/v1/world/maps`
- `GET /api/v1/world/regions`
- `GET /api/v1/world/scenes`

### 4. 我模块 - 个人中心
- 用户信息展示
- 订阅状态
- 用户日记列表
- 设置入口

**API:**
- `GET /api/v1/user/info`
- `GET /api/v1/user/subscription`
- `GET /api/v1/diary/user`

## 开发优先级

### Phase 1: 基础框架
- [x] 项目结构搭建
- [ ] API 客户端封装
- [ ] 认证流程（登录/注册）
- [ ] 路由配置

### Phase 2: 核心功能
- [ ] 首页（分身动态展示）
- [ ] 分身模块（性格面板 + 日记）
- [ ] 对话功能（WebSocket 流式对话）

### Phase 3: 完善功能
- [ ] 世界模块（区域/场景展示）
- [ ] 用户日记（创建 + 分身回应）
- [ ] 个人中心

### Phase 4: 体验优化
- [ ] 动画效果
- [ ] 离线缓存
- [ ] 推送通知
- [ ] 错误处理

## 关键技术实现

### 1. JWT 认证
```dart
// 登录后存储 token
await storage.write('token', response.token);
// 拦截器自动添加 Authorization header
```

### 2. WebSocket 流式对话
```dart
final channel = WebSocketChannel.connect(
  Uri.parse('ws://localhost:8888/api/v1/dialogue/stream')
);
```

### 3. 性格雷达图
```dart
RadarChart(data: [warmth, adventurous, social, creative, calm, energetic])
```

### 4. 图片上传流程
```dart
// 1. 获取上传凭证
final token = await api.getAvatarUploadToken();
// 2. 客户端直传 OSS
await uploadToOSS(token, imageFile);
// 3. 完成上传回调
await api.completeAvatarUpload(key);
```

## API 端点配置

- **开发环境**: `http://localhost:8888`
- **生产环境**: TBD

## 依赖包

```yaml
dependencies:
  flutter_riverpod: ^2.4.0
  dio: ^5.4.0
  retrofit: ^4.0.0
  web_socket_channel: ^2.4.0
  shared_preferences: ^2.2.0
  hive: ^2.2.3
  fl_chart: ^0.65.0
  lottie: ^2.7.0
  cached_network_image: ^3.3.0
  go_router: ^12.0.0
```

## 开发规范

参考 `app开发规范.md` 文档执行开发。
