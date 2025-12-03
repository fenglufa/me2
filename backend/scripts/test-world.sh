#!/usr/bin/env zsh

# World 服务接口测试脚本
# 使用 grpcurl 测试 World RPC 服务

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 服务地址
WORLD_RPC_ADDR="127.0.0.1:8006"

# Proto 文件路径
PROTO_DIR="../services/world/rpc"
PROTO_FILE="world.proto"

# 打印标题
print_title() {
    echo "${BLUE}========================================${NC}"
    echo "${BLUE}$1${NC}"
    echo "${BLUE}========================================${NC}"
}

# 打印成功信息
print_success() {
    echo "${GREEN}✓ $1${NC}"
}

# 打印错误信息
print_error() {
    echo "${RED}✗ $1${NC}"
}

# 打印信息
print_info() {
    echo "${YELLOW}ℹ $1${NC}"
}

# 检查服务是否启动
check_service() {
    print_title "检查 World 服务状态"

    if grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE $WORLD_RPC_ADDR list > /dev/null 2>&1; then
        print_success "World 服务运行正常"
        echo ""
        grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE $WORLD_RPC_ADDR list
        echo ""
    else
        print_error "World 服务未启动或连接失败"
        exit 1
    fi
}

# 测试 ListMaps - 获取所有地图列表
test_list_maps() {
    print_title "测试 ListMaps - 获取所有地图列表"

    print_info "请求参数: page=1, page_size=10, only_active=true"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "page": 1,
        "page_size": 10,
        "only_active": true
    }' $WORLD_RPC_ADDR world.World/ListMaps 2>&1)

    if echo "$RESPONSE" | grep -q "maps"; then
        print_success "获取地图列表成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取地图列表失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetMap - 获取单个地图详情
test_get_map() {
    print_title "测试 GetMap - 获取地图详情"

    print_info "请求参数: map_id=1"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "map_id": "1"
    }' $WORLD_RPC_ADDR world.World/GetMap 2>&1)

    if echo "$RESPONSE" | grep -q "map"; then
        print_success "获取地图详情成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取地图详情失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 ListRegionsInMap - 获取地图的所有区域
test_list_regions_in_map() {
    print_title "测试 ListRegionsInMap - 获取地图的所有区域"

    print_info "请求参数: map_id=1, page=1, page_size=10"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "map_id": "1",
        "page": 1,
        "page_size": 10,
        "only_active": true
    }' $WORLD_RPC_ADDR world.World/ListRegionsInMap 2>&1)

    if echo "$RESPONSE" | grep -q "regions"; then
        print_success "获取区域列表成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取区域列表失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetRegion - 获取区域详情
test_get_region() {
    print_title "测试 GetRegion - 获取区域详情"

    print_info "请求参数: region_id=1 (繁华都市)"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "region_id": "1"
    }' $WORLD_RPC_ADDR world.World/GetRegion 2>&1)

    if echo "$RESPONSE" | grep -q "region"; then
        print_success "获取区域详情成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取区域详情失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 ListScenesInRegion - 列出区域内的场景
test_list_scenes_in_region() {
    print_title "测试 ListScenesInRegion - 列出区域内的场景"

    print_info "请求参数: region_id=1 (繁华都市), page=1, page_size=10"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "region_id": "1",
        "page": 1,
        "page_size": 10,
        "only_active": true
    }' $WORLD_RPC_ADDR world.World/ListScenesInRegion 2>&1)

    if echo "$RESPONSE" | grep -q "scenes"; then
        print_success "获取场景列表成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取场景列表失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetScene - 获取场景详情
test_get_scene() {
    print_title "测试 GetScene - 获取场景详情"

    print_info "请求参数: scene_id=1 (星巴克咖啡厅)"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "scene_id": "1"
    }' $WORLD_RPC_ADDR world.World/GetScene 2>&1)

    if echo "$RESPONSE" | grep -q "scene"; then
        print_success "获取场景详情成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取场景详情失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetScenesForAction - 根据行为推荐场景（学习）
test_get_scenes_for_study() {
    print_title "测试 GetScenesForAction - 推荐学习场景"

    print_info "请求参数: action_type=study, limit=5"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "action_type": "study",
        "limit": 5
    }' $WORLD_RPC_ADDR world.World/GetScenesForAction 2>&1)

    if echo "$RESPONSE" | grep -q "recommendations"; then
        print_success "获取学习场景推荐成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取学习场景推荐失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetScenesForAction - 根据行为推荐场景（社交）
test_get_scenes_for_social() {
    print_title "测试 GetScenesForAction - 推荐社交场景"

    print_info "请求参数: action_type=social, limit=5"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "action_type": "social",
        "limit": 5
    }' $WORLD_RPC_ADDR world.World/GetScenesForAction 2>&1)

    if echo "$RESPONSE" | grep -q "recommendations"; then
        print_success "获取社交场景推荐成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取社交场景推荐失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetScenesForAction - 根据行为推荐场景（休息）
test_get_scenes_for_rest() {
    print_title "测试 GetScenesForAction - 推荐休息场景"

    print_info "请求参数: action_type=rest, region_id=2 (宁静森林), limit=3"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "action_type": "rest",
        "region_id": "2",
        "limit": 3
    }' $WORLD_RPC_ADDR world.World/GetScenesForAction 2>&1)

    if echo "$RESPONSE" | grep -q "recommendations"; then
        print_success "获取休息场景推荐成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取休息场景推荐失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetScenesForAction - 根据行为推荐场景（创作）
test_get_scenes_for_creative() {
    print_title "测试 GetScenesForAction - 推荐创作场景"

    print_info "请求参数: action_type=creative, region_id=4 (艺术街区), limit=3"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "action_type": "creative",
        "region_id": "4",
        "limit": 3
    }' $WORLD_RPC_ADDR world.World/GetScenesForAction 2>&1)

    if echo "$RESPONSE" | grep -q "recommendations"; then
        print_success "获取创作场景推荐成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取创作场景推荐失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 主测试流程
main() {
    print_title "World 服务接口测试"
    echo ""

    # 检查 jq 是否安装
    if ! command -v jq &> /dev/null; then
        print_error "jq 未安装，请先安装: brew install jq"
        exit 1
    fi

    # 检查 grpcurl 是否安装
    if ! command -v grpcurl &> /dev/null; then
        print_error "grpcurl 未安装，请先安装: brew install grpcurl"
        exit 1
    fi

    # 检查服务状态
    check_service

    # 运行测试
    test_list_maps
    test_get_map
    test_list_regions_in_map
    test_get_region
    test_list_scenes_in_region
    test_get_scene
    test_get_scenes_for_study
    test_get_scenes_for_social
    test_get_scenes_for_rest
    test_get_scenes_for_creative

    print_title "测试完成"
}

# 如果有参数，执行特定测试
if [ $# -eq 0 ]; then
    main
else
    case "$1" in
        "check")
            check_service
            ;;
        "maps")
            test_list_maps
            test_get_map
            ;;
        "regions")
            test_list_regions_in_map
            test_get_region
            ;;
        "scenes")
            test_list_scenes_in_region
            test_get_scene
            ;;
        "actions")
            test_get_scenes_for_study
            test_get_scenes_for_social
            test_get_scenes_for_rest
            test_get_scenes_for_creative
            ;;
        *)
            echo "用法: $0 [check|maps|regions|scenes|actions]"
            echo ""
            echo "参数说明:"
            echo "  check     - 检查服务状态"
            echo "  maps      - 测试地图相关接口"
            echo "  regions   - 测试区域相关接口"
            echo "  scenes    - 测试场景相关接口"
            echo "  actions   - 测试行为推荐接口"
            echo "  (无参数)  - 运行所有测试"
            ;;
    esac
fi
