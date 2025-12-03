#!/usr/bin/env zsh

# User 服务接口测试脚本
# 使用 grpcurl 测试 User RPC 服务

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 服务地址
USER_RPC_ADDR="127.0.0.1:8003"

# Proto 文件路径
PROTO_DIR="../services/user/rpc"
PROTO_FILE="user.proto"

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
    print_title "检查 User 服务状态"

    if grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE $USER_RPC_ADDR list > /dev/null 2>&1; then
        print_success "User 服务运行正常"
        echo ""
        grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE $USER_RPC_ADDR list
        echo ""
    else
        print_error "User 服务未启动或连接失败"
        exit 1
    fi
}

# 测试 LoginOrRegister - 新用户注册
test_login_or_register_new() {
    print_title "测试 LoginOrRegister - 新用户注册"

    print_info "请求参数: phone=18610665422, code=888888"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "phone": "18610665422",
        "code": "888888"
    }' $USER_RPC_ADDR user.User/LoginOrRegister 2>&1)

    if echo "$RESPONSE" | grep -q "userId"; then
        print_success "登录/注册成功"
        echo "$RESPONSE" | jq '.'

        # 提取 user_id 供后续测试使用
        export USER_ID=$(echo "$RESPONSE" | jq -r '.userId')
        echo ""
        print_info "用户ID: $USER_ID"
    else
        print_error "登录/注册失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetUserInfo
test_get_user_info() {
    print_title "测试 GetUserInfo"

    # 使用之前注册的用户ID
    if [ -z "$USER_ID" ]; then
        USER_ID="100000001"
    fi

    print_info "请求参数: user_id=$USER_ID"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"user_id\": $USER_ID
    }" $USER_RPC_ADDR user.User/GetUserInfo 2>&1)

    if echo "$RESPONSE" | grep -q "userId"; then
        print_success "获取用户信息成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取用户信息失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetSubscription
test_get_subscription() {
    print_title "测试 GetSubscription"

    if [ -z "$USER_ID" ]; then
        USER_ID="100000001"
    fi

    print_info "请求参数: user_id=$USER_ID"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"user_id\": $USER_ID
    }" $USER_RPC_ADDR user.User/GetSubscription 2>&1)

    # 检查是否有错误（gRPC 错误会包含 "Code:" 或 "ERROR:"）
    if echo "$RESPONSE" | grep -qE "(Code:|ERROR:)"; then
        print_error "获取订阅信息失败"
        echo "$RESPONSE"
    else
        print_success "获取订阅信息成功"
        # 即使是空对象 {} 也表示成功（免费用户的默认值）
        echo "$RESPONSE" | jq '.'
    fi
    echo ""
}

# 测试 UpdateSubscription
test_update_subscription() {
    print_title "测试 UpdateSubscription"

    if [ -z "$USER_ID" ]; then
        USER_ID="100000001"
    fi

    # 设置为付费会员，过期时间为30天后
    EXPIRE_TIME=$(($(date +%s) + 2592000))

    print_info "请求参数: user_id=$USER_ID, subscription_type=1, expire_time=$EXPIRE_TIME"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"user_id\": $USER_ID,
        \"subscription_type\": 1,
        \"subscription_expire_time\": $EXPIRE_TIME
    }" $USER_RPC_ADDR user.User/UpdateSubscription 2>&1)

    if echo "$RESPONSE" | grep -q "success"; then
        print_success "更新订阅状态成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "更新订阅状态失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 UpdateUserStatus
test_update_user_status() {
    print_title "测试 UpdateUserStatus"

    if [ -z "$USER_ID" ]; then
        USER_ID="100000001"
    fi

    print_info "请求参数: user_id=$USER_ID, status=1 (正常)"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"user_id\": $USER_ID,
        \"status\": 1
    }" $USER_RPC_ADDR user.User/UpdateUserStatus 2>&1)

    if echo "$RESPONSE" | grep -q "success"; then
        print_success "更新用户状态成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "更新用户状态失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试错误场景 - 无效用户ID
test_invalid_user_id() {
    print_title "测试错误场景 - 无效用户ID"

    print_info "请求参数: user_id=0 (应该失败)"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "user_id": 0
    }' $USER_RPC_ADDR user.User/GetUserInfo 2>&1)

    if echo "$RESPONSE" | grep -q "用户ID不能为空"; then
        print_success "正确返回错误: 用户ID不能为空"
        echo "$RESPONSE"
    else
        print_error "应该返回错误但没有"
        echo "$RESPONSE"
    fi
    echo ""
}

# 主测试流程
main() {
    print_title "User 服务接口测试"
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
    test_login_or_register_new
    test_get_user_info
    test_get_subscription
    test_update_subscription
    test_get_subscription  # 再次查询验证更新
    test_update_user_status
    test_invalid_user_id

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
        "login")
            test_login_or_register_new
            ;;
        "info")
            test_get_user_info
            ;;
        "subscription")
            test_get_subscription
            ;;
        *)
            echo "用法: $0 [check|login|info|subscription]"
            echo ""
            echo "参数说明:"
            echo "  check        - 检查服务状态"
            echo "  login        - 测试登录/注册"
            echo "  info         - 测试获取用户信息"
            echo "  subscription - 测试获取订阅信息"
            echo "  (无参数)     - 运行所有测试"
            ;;
    esac
fi
