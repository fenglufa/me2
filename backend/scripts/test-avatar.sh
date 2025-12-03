#!/usr/bin/env zsh

# Avatar 服务接口测试脚本
# 使用 grpcurl 测试 Avatar RPC 服务

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 服务地址
AVATAR_RPC_ADDR="127.0.0.1:8004"

# Proto 文件路径
PROTO_DIR="../services/avatar/rpc"
PROTO_FILE="avatar.proto"

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
    print_title "检查 Avatar 服务状态"

    if grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE $AVATAR_RPC_ADDR list > /dev/null 2>&1; then
        print_success "Avatar 服务运行正常"
        echo ""
        grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE $AVATAR_RPC_ADDR list
        echo ""
    else
        print_error "Avatar 服务未启动或连接失败"
        exit 1
    fi
}

# 测试 CreateAvatar - 创建分身
test_create_avatar() {
    print_title "测试 CreateAvatar - 创建分身"

    if [ -z "$USER_ID" ]; then
        USER_ID="100000001"
    fi

    print_info "请求参数: user_id=$USER_ID, nickname=测试分身, gender=1, birth_date=1995-06-15"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"user_id\": $USER_ID,
        \"nickname\": \"测试分身\",
        \"avatar_url\": \"\",
        \"gender\": 1,
        \"birth_date\": \"1995-06-15\",
        \"occupation\": \"工程师\",
        \"marital_status\": 1
    }" $AVATAR_RPC_ADDR avatar.Avatar/CreateAvatar 2>&1)

    if echo "$RESPONSE" | grep -q "avatarId"; then
        print_success "创建分身成功"
        echo "$RESPONSE" | jq '.'

        # 提取 avatar_id 供后续测试使用
        export AVATAR_ID=$(echo "$RESPONSE" | jq -r '.avatarId')
        echo ""
        print_info "分身ID: $AVATAR_ID"
    else
        print_error "创建分身失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetMyAvatar - 获取我的分身
test_get_my_avatar() {
    print_title "测试 GetMyAvatar - 获取我的分身"

    if [ -z "$USER_ID" ]; then
        USER_ID="100000001"
    fi

    print_info "请求参数: user_id=$USER_ID"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"user_id\": $USER_ID
    }" $AVATAR_RPC_ADDR avatar.Avatar/GetMyAvatar 2>&1)

    if echo "$RESPONSE" | grep -q "hasAvatar"; then
        print_success "获取我的分身成功"
        echo "$RESPONSE" | jq '.'

        # 如果有分身，提取 avatar_id
        if echo "$RESPONSE" | jq -e '.hasAvatar == true' > /dev/null; then
            export AVATAR_ID=$(echo "$RESPONSE" | jq -r '.avatar.avatarId')
            print_info "分身ID: $AVATAR_ID"
        fi
    else
        print_error "获取我的分身失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetAvatarInfo - 获取分身详情
test_get_avatar_info() {
    print_title "测试 GetAvatarInfo - 获取分身详情"

    if [ -z "$AVATAR_ID" ]; then
        print_error "未找到分身ID，请先创建分身"
        return
    fi

    print_info "请求参数: avatar_id=$AVATAR_ID"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"avatar_id\": $AVATAR_ID
    }" $AVATAR_RPC_ADDR avatar.Avatar/GetAvatarInfo 2>&1)

    if echo "$RESPONSE" | grep -q "avatar"; then
        print_success "获取分身详情成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取分身详情失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 UpdateAvatarProfile - 更新分身资料
test_update_avatar_profile() {
    print_title "测试 UpdateAvatarProfile - 更新分身资料"

    if [ -z "$AVATAR_ID" ]; then
        print_error "未找到分身ID，请先创建分身"
        return
    fi

    print_info "请求参数: avatar_id=$AVATAR_ID, nickname=更新后的昵称"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"avatar_id\": $AVATAR_ID,
        \"nickname\": \"更新后的昵称\",
        \"avatar_url\": \"\"
    }" $AVATAR_RPC_ADDR avatar.Avatar/UpdateAvatarProfile 2>&1)

    if echo "$RESPONSE" | grep -q "success"; then
        print_success "更新分身资料成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "更新分身资料失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetAvatarUploadToken - 获取头像上传凭证
test_get_avatar_upload_token() {
    print_title "测试 GetAvatarUploadToken - 获取头像上传凭证"

    if [ -z "$AVATAR_ID" ]; then
        print_error "未找到分身ID，请先创建分身"
        return
    fi

    print_info "请求参数: avatar_id=$AVATAR_ID, file_name=avatar.jpg"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"avatar_id\": $AVATAR_ID,
        \"file_name\": \"avatar.jpg\"
    }" $AVATAR_RPC_ADDR avatar.Avatar/GetAvatarUploadToken 2>&1)

    if echo "$RESPONSE" | grep -q "host"; then
        print_success "获取头像上传凭证成功"
        echo "$RESPONSE" | jq '.'

        # 提取 key 和 complete_token 供后续测试使用
        export UPLOAD_KEY=$(echo "$RESPONSE" | jq -r '.key')
        export COMPLETE_TOKEN=$(echo "$RESPONSE" | jq -r '.completeToken')
    else
        print_error "获取头像上传凭证失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 主测试流程
main() {
    print_title "Avatar 服务接口测试"
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
    test_get_my_avatar           # 先查询是否已有分身

    if [ -z "$AVATAR_ID" ]; then
        test_create_avatar       # 如果没有分身，创建一个
    fi

    test_get_avatar_info         # 获取分身详情
    test_update_avatar_profile   # 更新分身资料
    test_get_avatar_info         # 再次查询验证更新
    test_get_avatar_upload_token # 获取上传凭证

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
        "create")
            test_create_avatar
            ;;
        "my")
            test_get_my_avatar
            ;;
        "info")
            test_get_avatar_info
            ;;
        "update")
            test_update_avatar_profile
            ;;
        "upload")
            test_get_avatar_upload_token
            ;;
        *)
            echo "用法: $0 [check|create|my|info|update|upload]"
            echo ""
            echo "参数说明:"
            echo "  check  - 检查服务状态"
            echo "  create - 测试创建分身"
            echo "  my     - 测试获取我的分身"
            echo "  info   - 测试获取分身详情"
            echo "  update - 测试更新分身资料"
            echo "  upload - 测试获取上传凭证"
            echo "  (无参数) - 运行所有测试"
            ;;
    esac
fi
