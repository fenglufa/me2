#!/bin/bash

# SMS 服务测试脚本
# 用于测试发送验证码和验证验证码功能

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# SMS 服务地址
SMS_HOST="localhost:8001"

# 检查 grpcurl 是否安装
if ! command -v grpcurl &> /dev/null; then
    echo -e "${RED}错误: grpcurl 未安装${NC}"
    echo "安装方法: brew install grpcurl"
    exit 1
fi

# 检查 SMS 服务是否运行
if ! lsof -i :8001 > /dev/null 2>&1; then
    echo -e "${RED}错误: SMS 服务未运行${NC}"
    echo "启动方法: cd services/sms && make run"
    exit 1
fi

echo -e "${GREEN}=== SMS 服务测试 ===${NC}\n"

# 获取手机号
read -p "请输入手机号（例如：13800138000）: " PHONE

if [ -z "$PHONE" ]; then
    echo -e "${RED}错误: 手机号不能为空${NC}"
    exit 1
fi

# 发送验证码
echo -e "\n${YELLOW}1. 发送验证码到 $PHONE${NC}"
echo "请求内容: {\"phone\": \"$PHONE\", \"scene\": \"login\"}"
echo ""

SEND_RESULT=$(grpcurl -plaintext -d "{\"phone\": \"$PHONE\", \"scene\": \"login\"}" \
    -import-path /Users/flf/Desktop/ProjectCode/me2/backend/services/sms/rpc \
    -proto sms.proto \
    $SMS_HOST sms.Sms/SendCode 2>&1)

echo "$SEND_RESULT"

# 检查是否发送成功
if echo "$SEND_RESULT" | grep -q '"success": true'; then
    echo -e "\n${GREEN}✓ 验证码发送成功！${NC}"
    echo -e "${YELLOW}请查收手机短信${NC}\n"

    # 等待用户输入验证码
    read -p "请输入收到的验证码: " CODE

    if [ -z "$CODE" ]; then
        echo -e "${RED}未输入验证码，测试结束${NC}"
        exit 0
    fi

    # 验证验证码
    echo -e "\n${YELLOW}2. 验证验证码${NC}"
    echo "请求内容: {\"phone\": \"$PHONE\", \"code\": \"$CODE\"}"
    echo ""

    VERIFY_RESULT=$(grpcurl -plaintext -d "{\"phone\": \"$PHONE\", \"code\": \"$CODE\"}" \
        -import-path /Users/flf/Desktop/ProjectCode/me2/backend/services/sms/rpc \
        -proto sms.proto \
        $SMS_HOST sms.Sms/VerifyCode 2>&1)

    echo "$VERIFY_RESULT"

    # 检查验证结果
    if echo "$VERIFY_RESULT" | grep -q '"valid": true'; then
        echo -e "\n${GREEN}✓ 验证码验证成功！${NC}"
    else
        echo -e "\n${RED}✗ 验证码验证失败${NC}"
        if echo "$VERIFY_RESULT" | grep -q "验证码错误"; then
            echo -e "${YELLOW}原因: 验证码不正确${NC}"
        elif echo "$VERIFY_RESULT" | grep -q "已过期"; then
            echo -e "${YELLOW}原因: 验证码已过期（有效期5分钟）${NC}"
        fi
    fi
else
    echo -e "\n${RED}✗ 验证码发送失败${NC}"
    if echo "$SEND_RESULT" | grep -q "InvalidAccessKeyId"; then
        echo -e "${YELLOW}原因: 阿里云 AccessKey ID 无效${NC}"
    elif echo "$SEND_RESULT" | grep -q "SignatureDoesNotMatch"; then
        echo -e "${YELLOW}原因: 阿里云 AccessKey Secret 错误${NC}"
    elif echo "$SEND_RESULT" | grep -q "isv.MOBILE_NUMBER_ILLEGAL"; then
        echo -e "${YELLOW}原因: 手机号格式不正确${NC}"
    fi
fi

echo ""
