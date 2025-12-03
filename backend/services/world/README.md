# World Service

世界地图服务，提供虚拟世界的地图、区域、场景管理和行为推荐功能。

## 功能

- **世界地图管理** - 管理多个虚拟世界地图
- **区域管理** - 管理地图下的各个区域（城市、森林等）
- **场景管理** - 管理区域内的具体场景（咖啡厅、图书馆等）
- **行为推荐** - 根据分身的行为类型智能推荐适合的场景
- **场景匹配算法** - 基于场景特征计算匹配度和推荐理由

## 技术栈

- go-zero RPC
- MySQL (数据存储)
- Redis (缓存，预留)
- Etcd (服务注册与发现)

## 配置

配置文件: `rpc/etc/world-dev.yaml`

```yaml
Name: world.rpc
ListenOn: 0.0.0.0:8006

Etcd:
  Hosts:
    - 127.0.0.1:2379
  Key: world.rpc

Mysql:
  DataSource: root:123456@tcp(127.0.0.1:3306)/me2?charset=utf8mb4&parseTime=true&loc=Local

RedisConf:
  Host: 127.0.0.1:6379
  Type: node
  Pass: ""
```

## 数据库

### 创建数据库表

```bash
mysql -u root -p < rpc/world.sql
```

### 表结构

#### world_maps - 世界地图表
```sql
CREATE TABLE world_maps (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    cover_image VARCHAR(500),
    level_required INT DEFAULT 0,
    is_active TINYINT(1) DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

#### world_regions - 区域表
```sql
CREATE TABLE world_regions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    map_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    cover_image VARCHAR(500),
    atmosphere VARCHAR(50),
    tags JSON,
    level_required INT DEFAULT 0,
    is_active TINYINT(1) DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (map_id) REFERENCES world_maps(id) ON DELETE CASCADE
);
```

#### world_scenes - 场景表
```sql
CREATE TABLE world_scenes (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    region_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    cover_image VARCHAR(500),
    atmosphere VARCHAR(50),
    tags JSON,
    suitable_actions JSON,
    capacity INT DEFAULT 100,
    level_required INT DEFAULT 0,
    is_active TINYINT(1) DEFAULT 1,
    -- 场景特征
    has_wifi TINYINT(1) DEFAULT 0,
    has_food TINYINT(1) DEFAULT 0,
    has_seating TINYINT(1) DEFAULT 1,
    is_indoor TINYINT(1) DEFAULT 1,
    is_quiet TINYINT(1) DEFAULT 0,
    comfort_level INT DEFAULT 5,
    social_level INT DEFAULT 5,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (region_id) REFERENCES world_regions(id) ON DELETE CASCADE
);
```

### 初始数据

数据库脚本中已包含初始数据：
- 1个世界地图：初始世界
- 4个区域：繁华都市、宁静森林、学术区、艺术街区
- 15个场景：包含各种类型的场所（咖啡厅、公园、图书馆等）

## 使用

```bash
# 1. 生成 RPC 代码
make gen-rpc

# 2. 初始化依赖
make init

# 3. 创建数据库表
mysql -u root -p < rpc/world.sql

# 4. 编译
make build

# 5. 运行（开发模式）
make run-dev
```

## API

### GetMap - 获取世界地图

请求:
```protobuf
message GetMapRequest {
  int64 map_id = 1;
}
```

响应:
```protobuf
message GetMapResponse {
  WorldMap map = 1;
}
```

### ListMaps - 获取所有地图列表

请求:
```protobuf
message ListMapsRequest {
  int32 page = 1;
  int32 page_size = 2;
  bool only_active = 3;
}
```

响应:
```protobuf
message ListMapsResponse {
  repeated WorldMap maps = 1;
  int64 total = 2;
}
```

### GetRegion - 获取区域详情

请求:
```protobuf
message GetRegionRequest {
  int64 region_id = 1;
}
```

响应:
```protobuf
message GetRegionResponse {
  WorldRegion region = 1;
}
```

### ListRegionsInMap - 获取地图的所有区域

请求:
```protobuf
message ListRegionsInMapRequest {
  int64 map_id = 1;
  int32 page = 2;
  int32 page_size = 3;
  bool only_active = 4;
}
```

响应:
```protobuf
message ListRegionsInMapResponse {
  repeated WorldRegion regions = 1;
  int64 total = 2;
}
```

### GetScene - 获取场景详情

请求:
```protobuf
message GetSceneRequest {
  int64 scene_id = 1;
}
```

响应:
```protobuf
message GetSceneResponse {
  WorldScene scene = 1;
}
```

### ListScenesInRegion - 列出区域内的场景

请求:
```protobuf
message ListScenesInRegionRequest {
  int64 region_id = 1;
  int32 page = 2;
  int32 page_size = 3;
  bool only_active = 4;
  repeated string tags = 5;
}
```

响应:
```protobuf
message ListScenesInRegionResponse {
  repeated WorldScene scenes = 1;
  int64 total = 2;
}
```

### GetScenesForAction - 根据行为类型推荐场景 ⭐

**用途**: 供 Action Service 调用，为分身的行为推荐合适的场景

请求:
```protobuf
message GetScenesForActionRequest {
  string action_type = 1;       // 行为类型: exploration, social, study, creative, rest, play
  int64 map_id = 2;             // 指定地图ID（可选）
  int64 region_id = 3;          // 指定区域ID（可选）
  int32 limit = 4;              // 返回数量限制，默认10
  int32 avatar_level = 5;       // 分身等级，用于过滤等级要求
}
```

响应:
```protobuf
message GetScenesForActionResponse {
  repeated SceneRecommendation recommendations = 1;
}

message SceneRecommendation {
  WorldScene scene = 1;
  float match_score = 2;        // 匹配度分数 0-100
  string reason = 3;            // 推荐理由
}
```

**行为类型说明**:
- `exploration` - 探索新地方
- `social` - 社交互动
- `study` - 学习知识
- `creative` - 艺术创作
- `rest` - 休息放松
- `play` - 游戏娱乐

**匹配算法**:
根据不同行为类型，考虑不同的场景特征：
- **study**: 安静度 + WiFi + 舒适度
- **social**: 社交度 + 非安静 + 舒适度
- **rest**: 舒适度 + 安静度 + 座位
- **creative**: 舒适度 + 安静度 + 氛围
- **exploration**: 户外 + 容量
- **play**: 社交度 + 非安静 + 食物

## 场景特征

每个场景包含以下特征，用于智能匹配：

- `has_wifi` - 是否有WiFi
- `has_food` - 是否有食物
- `has_seating` - 是否有座位
- `is_indoor` - 是否室内
- `is_quiet` - 是否安静
- `comfort_level` - 舒适度 (1-10)
- `social_level` - 社交度 (1-10)

## 测试

运行测试脚本：

```bash
# 运行所有测试
cd ../scripts && ./test-world.sh

# 检查服务状态
./test-world.sh check

# 测试地图相关接口
./test-world.sh maps

# 测试区域相关接口
./test-world.sh regions

# 测试场景相关接口
./test-world.sh scenes

# 测试行为推荐接口
./test-world.sh actions
```

## 依赖服务

- **MySQL** - 数据存储
- **Redis** - 缓存（预留）
- **Etcd** - 服务注册与发现

## 目录结构

```
world/
├── Makefile
├── README.md
├── bin/
│   └── world-rpc
└── rpc/
    ├── world.proto
    ├── world.sql
    ├── etc/
    │   └── world-dev.yaml
    └── internal/
        ├── config/
        │   └── config.go
        ├── logic/
        │   ├── converter.go
        │   ├── get_map_logic.go
        │   ├── list_maps_logic.go
        │   ├── get_region_logic.go
        │   ├── list_regions_in_map_logic.go
        │   ├── get_scene_logic.go
        │   ├── list_scenes_in_region_logic.go
        │   └── get_scenes_for_action_logic.go
        ├── model/
        │   ├── world_map_model.go
        │   ├── world_region_model.go
        │   └── world_scene_model.go
        ├── server/
        │   └── world_server.go
        └── svc/
            └── service_context.go
```

## 注意事项

- 场景的 `suitable_actions` 字段存储为 JSON 数组，使用 MySQL 的 `JSON_CONTAINS` 函数查询
- 场景推荐算法会根据行为类型动态计算匹配度
- 建议为高频查询场景添加 Redis 缓存

## 待实现功能

### 等级系统 (Level System)

当前版本**暂未实现**等级限制功能，以下功能将在后续版本中实现：

- **地图等级要求** (`level_required` in `world_maps`) - 控制地图解锁
- **区域等级要求** (`level_required` in `world_regions`) - 控制区域解锁
- **场景等级要求** (`level_required` in `world_scenes`) - 控制场景解锁
- **分身等级过滤** (`avatar_level` in `GetScenesForActionRequest`) - 根据分身等级过滤推荐场景

当等级系统实现后，需要：
1. 在 proto 文件中添加 `level_required` 字段到 `WorldMap`、`WorldRegion`、`WorldScene`
2. 在 proto 文件中添加 `avatar_level` 参数到 `GetScenesForActionRequest`
3. 在数据库模型中添加 `LevelRequired` 字段
4. 在 `FindByActionType` 方法中添加等级过滤逻辑：`AND level_required <= ?`
5. 更新 converter 添加等级字段转换
6. 更新数据库迁移脚本

**设计思路**：
- 等级值越高，解锁难度越大
- 场景推荐时自动过滤掉超出分身等级的场景
- 为新手用户提供循序渐进的探索体验

## 扩展

### 添加新的区域

```sql
INSERT INTO world_regions (map_id, name, description, cover_image, atmosphere, tags, is_active)
VALUES (1, '新区域', '区域描述', 'image_url', '氛围', '["标签1", "标签2"]', 1);
```

### 添加新的场景

```sql
INSERT INTO world_scenes (region_id, name, description, cover_image, atmosphere, tags,
    suitable_actions, capacity, is_active,
    has_wifi, has_food, has_seating, is_indoor, is_quiet, comfort_level, social_level)
VALUES (1, '新场景', '场景描述', 'image_url', '氛围', '["标签1"]',
    '["study", "social"]', 50, 1,
    1, 1, 1, 1, 0, 8, 7);
```

### 自定义匹配算法

在 `internal/model/world_scene_model.go` 的 `CalculateMatchScore` 和 `GetRecommendationReason` 方法中修改匹配逻辑。



三层架构关系

  世界地图 (WorldMap)
      └── 区域 (WorldRegion) [多个]
              └── 场景 (WorldScene) [多个]

  1️⃣ 世界地图 (WorldMap) - 顶层容器

  作用：最大的地理单位，代表一个完整的虚拟世界

  数据库关系：
  world_maps (id, name, description, cover_image, is_active)

  示例数据：
  - ID: 1, 名称: "初始世界"

  特点：
  - 可以有多个世界（比如：现代世界、古代世界、科幻世界）
  - 为不同主题或难度提供隔离

  ---
  2️⃣ 区域 (WorldRegion) - 中层分区

  作用：世界地图下的大型区域，相当于"城市"或"地区"的概念

  数据库关系：
  world_regions (id, map_id, name, description, atmosphere, tags)
                      ↑
                      └── 外键关联到 world_maps.id

  示例数据（来自初始世界）：
  - ID: 1, 区域名: "繁华都市" (map_id=1)
  - ID: 2, 区域名: "宁静森林" (map_id=1)
  - ID: 3, 区域名: "学术区" (map_id=1)
  - ID: 4, 区域名: "艺术街区" (map_id=1)

  特点：
  - 一个世界包含多个区域
  - 每个区域有独特的氛围和标签
  - 提供大范围的环境分类

  ---
  3️⃣ 场景 (WorldScene) - 底层具体地点

  作用：区域内的具体场所，用户分身实际活动的地方

  数据库关系：
  world_scenes (id, region_id, name, suitable_actions, has_wifi, comfort_level...)
                     ↑
                     └── 外键关联到 world_regions.id

  示例数据（来自"繁华都市"区域）：
  - ID: 1, 场景名: "星巴克咖啡厅" (region_id=1)
    - 适合行为: ["study", "social", "rest"]
    - 特征: 有WiFi、有食物、安静、舒适度8、社交度6
  - ID: 2, 场景名: "中央公园" (region_id=1)
    - 适合行为: ["rest", "social", "exploration"]
    - 特征: 户外、安静、舒适度9、社交度5

  特点：
  - 一个区域包含多个场景
  - 每个场景有详细的特征属性（WiFi、食物、座位等）
  - 定义了适合的行为类型（学习、社交、休息等）
  - 有匹配算法计算与行为的适配度

  ---
  实际应用流程

  场景 1：用户浏览世界

  1. 调用 ListMaps → 获取所有世界列表
  2. 用户选择 "初始世界"
  3. 调用 ListRegionsInMap(map_id=1) → 获取该世界的所有区域
  4. 用户选择 "繁华都市"
  5. 调用 ListScenesInRegion(region_id=1) → 获取该区域的所有场景

  场景 2：AI分身执行行为

  1. AI分身想要 "学习"
  2. 调用 GetScenesForAction(action_type="study")
  3. 系统返回推荐场景：
     - 星巴克咖啡厅 (匹配度: 88分) - "环境安静，有WiFi，舒适度高"
     - 市立图书馆 (匹配度: 92分) - "环境安静，有WiFi，舒适度高"
     - 大学自习室 (匹配度: 85分) - "环境安静，有WiFi"

  ---
  数据流向图

  用户视角导航:
  ListMaps → GetMap → ListRegionsInMap → GetRegion → ListScenesInRegion → GetScene
    ↓           ↓            ↓               ↓              ↓               ↓
  世界列表    世界详情      区域列表         区域详情        场景列表         场景详情


  AI行为推荐:
  GetScenesForAction(action_type)
      ↓
  直接返回所有世界中符合条件的场景（可选过滤 map_id、region_id）
      ↓
  场景 + 匹配度 + 推荐理由

  ---
  关键设计优势

  1. 层级清晰：世界 → 区域 → 场景，符合真实地理概念
  2. 灵活扩展：可以轻松添加新世界、新区域、新场景
  3. 智能推荐：跨越层级直接推荐最合适的场景
  4. 数据隔离：级联删除确保数据一致性（删除世界 → 删除所有区域和场景）

  这种三层架构既支持用户主动探索浏览，也支持AI智能推荐，非常适合虚拟世界分身系统！