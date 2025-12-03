#!/usr/bin/env zsh

# OSS 服务接口测试脚本
# 使用 grpcurl 测试 OSS RPC 服务

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 服务地址
OSS_RPC_ADDR="127.0.0.1:8002"

# Proto 文件路径（相对于 scripts 目录）
PROTO_DIR="../services/oss/rpc"
PROTO_FILE="oss.proto"

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
    print_title "检查 OSS 服务状态"

    if grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE $OSS_RPC_ADDR list > /dev/null 2>&1; then
        print_success "OSS 服务运行正常"
        echo ""
        grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE $OSS_RPC_ADDR list
        echo ""
    else
        print_error "OSS 服务未启动或连接失败"
        exit 1
    fi
}

# 测试 GetUploadToken - 头像上传
test_get_upload_token_avatar() {
    print_title "测试 GetUploadToken - 头像上传 (avatar)"

    print_info "请求参数: service_name=avatar, file_name=profile.jpg, user_id=12345"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "service_name": "avatar",
        "file_name": "profile.jpg",
        "user_id": 12345
    }' $OSS_RPC_ADDR oss.Oss/GetUploadToken 2>&1)

    if echo "$RESPONSE" | grep -q "host"; then
        print_success "获取上传令牌成功"
        echo "$RESPONSE" | jq '.'

        # 提取 complete_token 供后续测试使用
        AVATAR_COMPLETE_TOKEN=$(echo "$RESPONSE" | jq -r '.completeToken')
        AVATAR_DIR=$(echo "$RESPONSE" | jq -r '.dir')

        echo ""
        print_info "完成令牌 (complete_token): $AVATAR_COMPLETE_TOKEN"
        print_info "文件目录 (dir): $AVATAR_DIR"
    else
        print_error "获取上传令牌失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetUploadToken - Banner 上传
test_get_upload_token_banner() {
    print_title "测试 GetUploadToken - Banner 上传 (banner)"

    print_info "请求参数: service_name=banner, file_name=banner.png, user_id=12345"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "service_name": "banner",
        "file_name": "banner.png",
        "user_id": 12345
    }' $OSS_RPC_ADDR oss.Oss/GetUploadToken 2>&1)

    if echo "$RESPONSE" | grep -q "host"; then
        print_success "获取上传令牌成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取上传令牌失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetUploadToken - 帖子上传 (视频)
test_get_upload_token_post_video() {
    print_title "测试 GetUploadToken - 帖子上传 (post - 视频)"

    print_info "请求参数: service_name=post, file_name=video.mp4, user_id=12345"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "service_name": "post",
        "file_name": "video.mp4",
        "user_id": 12345
    }' $OSS_RPC_ADDR oss.Oss/GetUploadToken 2>&1)

    if echo "$RESPONSE" | grep -q "host"; then
        print_success "获取上传令牌成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取上传令牌失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetUploadToken - 不支持的服务类型
test_get_upload_token_invalid_service() {
    print_title "测试 GetUploadToken - 不支持的服务类型"

    print_info "请求参数: service_name=invalid_service (应该失败)"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "service_name": "invalid_service",
        "file_name": "test.jpg",
        "user_id": 12345
    }' $OSS_RPC_ADDR oss.Oss/GetUploadToken 2>&1)

    if echo "$RESPONSE" | grep -q "不支持的服务类型"; then
        print_success "正确返回错误: 不支持的服务类型"
        echo "$RESPONSE"
    else
        print_error "应该返回错误但没有"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetUploadToken - 不支持的文件类型
test_get_upload_token_invalid_ext() {
    print_title "测试 GetUploadToken - 不支持的文件类型"

    print_info "请求参数: service_name=avatar, file_name=test.pdf (应该失败)"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "service_name": "avatar",
        "file_name": "test.pdf",
        "user_id": 12345
    }' $OSS_RPC_ADDR oss.Oss/GetUploadToken 2>&1)

    if echo "$RESPONSE" | grep -q "不支持的文件类型"; then
        print_success "正确返回错误: 不支持的文件类型"
        echo "$RESPONSE"
    else
        print_error "应该返回错误但没有"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 CompleteUpload - 成功场景
test_complete_upload_success() {
    print_title "测试 CompleteUpload - 成功场景"

    # 先获取上传令牌
    print_info "步骤 1: 获取上传令牌"
    TOKEN_RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "service_name": "avatar",
        "file_name": "avatar.jpg",
        "user_id": 99999
    }' $OSS_RPC_ADDR oss.Oss/GetUploadToken 2>&1)

    if echo "$TOKEN_RESPONSE" | grep -q "host"; then
        print_success "获取上传令牌成功"

        COMPLETE_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.completeToken')
        DIR=$(echo "$TOKEN_RESPONSE" | jq -r '.dir')

        print_info "完成令牌: ${COMPLETE_TOKEN:0:50}..."
        print_info "文件目录: $DIR"

        # 模拟上传完成，使用正确的文件路径
        FILE_KEY="${DIR}20231203_12345678.jpg"

        print_info "步骤 2: 完成上传验证"
        print_info "文件 Key: $FILE_KEY"

        COMPLETE_RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
            \"service_name\": \"avatar\",
            \"key\": \"$FILE_KEY\",
            \"user_id\": 99999,
            \"complete_token\": \"$COMPLETE_TOKEN\"
        }" $OSS_RPC_ADDR oss.Oss/CompleteUpload 2>&1)

        if echo "$COMPLETE_RESPONSE" | grep -q "url"; then
            print_success "上传完成验证成功"
            echo "$COMPLETE_RESPONSE" | jq '.'

            FILE_URL=$(echo "$COMPLETE_RESPONSE" | jq -r '.url')
            print_info "文件访问地址: $FILE_URL"
        else
            print_error "上传完成验证失败"
            echo "$COMPLETE_RESPONSE"
        fi
    else
        print_error "获取上传令牌失败，跳过完成上传测试"
    fi
    echo ""
}

# 测试 CompleteUpload - 用户 ID 不匹配
test_complete_upload_user_mismatch() {
    print_title "测试 CompleteUpload - 用户 ID 不匹配"

    # 获取用户 12345 的令牌
    TOKEN_RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "service_name": "avatar",
        "file_name": "test.jpg",
        "user_id": 12345
    }' $OSS_RPC_ADDR oss.Oss/GetUploadToken 2>&1)

    COMPLETE_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.completeToken')
    DIR=$(echo "$TOKEN_RESPONSE" | jq -r '.dir')
    FILE_KEY="${DIR}test.jpg"

    print_info "使用用户 12345 的令牌，但用 54321 完成上传 (应该失败)"

    # 尝试使用用户 54321 完成上传
    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"service_name\": \"avatar\",
        \"key\": \"$FILE_KEY\",
        \"user_id\": 54321,
        \"complete_token\": \"$COMPLETE_TOKEN\"
    }" $OSS_RPC_ADDR oss.Oss/CompleteUpload 2>&1)

    if echo "$RESPONSE" | grep -q "用户ID不匹配"; then
        print_success "正确返回错误: 用户ID不匹配"
        echo "$RESPONSE"
    else
        print_error "应该返回用户ID不匹配错误"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 CompleteUpload - 服务名称不匹配
test_complete_upload_service_mismatch() {
    print_title "测试 CompleteUpload - 服务名称不匹配"

    # 获取 avatar 服务的令牌
    TOKEN_RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "service_name": "avatar",
        "file_name": "test.jpg",
        "user_id": 12345
    }' $OSS_RPC_ADDR oss.Oss/GetUploadToken 2>&1)

    COMPLETE_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.completeToken')
    DIR=$(echo "$TOKEN_RESPONSE" | jq -r '.dir')
    FILE_KEY="${DIR}test.jpg"

    print_info "使用 avatar 的令牌，但用 banner 完成上传 (应该失败)"

    # 尝试使用 banner 完成上传
    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"service_name\": \"banner\",
        \"key\": \"$FILE_KEY\",
        \"user_id\": 12345,
        \"complete_token\": \"$COMPLETE_TOKEN\"
    }" $OSS_RPC_ADDR oss.Oss/CompleteUpload 2>&1)

    if echo "$RESPONSE" | grep -q "服务名称不匹配"; then
        print_success "正确返回错误: 服务名称不匹配"
        echo "$RESPONSE"
    else
        print_error "应该返回服务名称不匹配错误"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 CompleteUpload - 文件路径不合法
test_complete_upload_invalid_path() {
    print_title "测试 CompleteUpload - 文件路径不合法"

    TOKEN_RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "service_name": "avatar",
        "file_name": "test.jpg",
        "user_id": 12345
    }' $OSS_RPC_ADDR oss.Oss/GetUploadToken 2>&1)

    COMPLETE_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.completeToken')

    print_info "使用错误的文件路径 (应该失败)"

    # 使用不在允许目录下的路径
    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d "{
        \"service_name\": \"avatar\",
        \"key\": \"wrong_dir/test.jpg\",
        \"user_id\": 12345,
        \"complete_token\": \"$COMPLETE_TOKEN\"
    }" $OSS_RPC_ADDR oss.Oss/CompleteUpload 2>&1)

    if echo "$RESPONSE" | grep -q "文件路径不合法"; then
        print_success "正确返回错误: 文件路径不合法"
        echo "$RESPONSE"
    else
        print_error "应该返回文件路径不合法错误"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 CompleteUpload - 无效令牌
test_complete_upload_invalid_token() {
    print_title "测试 CompleteUpload - 无效令牌"

    print_info "使用伪造的令牌 (应该失败)"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "service_name": "avatar",
        "key": "avatar/test.jpg",
        "user_id": 12345,
        "complete_token": "invalid.token.here"
    }' $OSS_RPC_ADDR oss.Oss/CompleteUpload 2>&1)

    if echo "$RESPONSE" | grep -q "无效的验证令牌"; then
        print_success "正确返回错误: 无效的验证令牌"
        echo "$RESPONSE"
    else
        print_error "应该返回无效令牌错误"
        echo "$RESPONSE"
    fi
    echo ""
}

# 查看支持的文件类型
show_supported_types() {
    print_title "支持的业务和文件类型"

    echo "${YELLOW}Avatar (头像):${NC}"
    echo "  - 最大文件: 5MB"
    echo "  - 支持格式: jpg, jpeg, png, gif"
    echo ""

    echo "${YELLOW}Banner (横幅):${NC}"
    echo "  - 最大文件: 10MB"
    echo "  - 支持格式: jpg, jpeg, png"
    echo ""

    echo "${YELLOW}Post (帖子):${NC}"
    echo "  - 最大文件: 100MB"
    echo "  - 支持格式: jpg, jpeg, png, gif, mp4, mov, avi"
    echo ""
}

# 主测试流程
main() {
    print_title "OSS 服务接口测试"
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

    # 显示支持的文件类型
    show_supported_types

    # GetUploadToken 测试
    test_get_upload_token_avatar
    test_get_upload_token_banner
    test_get_upload_token_post_video
    test_get_upload_token_invalid_service
    test_get_upload_token_invalid_ext

    # CompleteUpload 测试
    test_complete_upload_success
    test_complete_upload_user_mismatch
    test_complete_upload_service_mismatch
    test_complete_upload_invalid_path
    test_complete_upload_invalid_token

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
        "avatar")
            test_get_upload_token_avatar
            ;;
        "banner")
            test_get_upload_token_banner
            ;;
        "post")
            test_get_upload_token_post_video
            ;;
        "complete")
            test_complete_upload_success
            ;;
        "types")
            show_supported_types
            ;;
        *)
            echo "用法: $0 [check|avatar|banner|post|complete|types]"
            echo ""
            echo "参数说明:"
            echo "  check    - 检查服务状态"
            echo "  avatar   - 测试头像上传令牌"
            echo "  banner   - 测试横幅上传令牌"
            echo "  post     - 测试帖子上传令牌"
            echo "  complete - 测试完成上传"
            echo "  types    - 显示支持的文件类型"
            echo "  (无参数) - 运行所有测试"
            ;;
    esac
fi
