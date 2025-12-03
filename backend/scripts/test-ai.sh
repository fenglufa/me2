#!/usr/bin/env zsh

# AI 服务接口测试脚本
# 使用 grpcurl 测试 AI RPC 服务

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 服务地址
AI_RPC_ADDR="127.0.0.1:8005"

# Proto 文件路径
PROTO_DIR="../services/ai/rpc"
PROTO_FILE="ai.proto"

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
    print_title "检查 AI 服务状态"

    if grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE $AI_RPC_ADDR list > /dev/null 2>&1; then
        print_success "AI 服务运行正常"
        echo ""
        grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE $AI_RPC_ADDR list
        echo ""
    else
        print_error "AI 服务未启动或连接失败"
        exit 1
    fi
}

# 测试 Chat - 分身对话
test_chat_avatar() {
    print_title "测试 Chat - 分身对话"

    print_info "请求参数: prompt_template=avatar_chat, user_id=100000001, avatar_id=1"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "prompt_template": "avatar_chat",
        "variables": {
            "personality": "我是一个友好、幽默的AI助手",
            "recent_events": "今天天气很好，用户刚完成了一个重要项目",
            "user_message": "嗨，最近怎么样？"
        },
        "model_config": {
            "temperature": 0.7,
            "max_tokens": 200
        },
        "user_id": "100000001",
        "avatar_id": "1"
    }' $AI_RPC_ADDR ai.Ai/Chat 2>&1)

    if echo "$RESPONSE" | grep -q "content"; then
        print_success "对话生成成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "对话生成失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 Generate - 事件故事生成
test_generate_event_story() {
    print_title "测试 Generate - 事件故事生成"

    print_info "请求参数: prompt_template=event_story, user_id=100000001, avatar_id=1"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "prompt_template": "event_story",
        "variables": {
            "event_type": "完成项目",
            "event_detail": "成功上线了一个新功能",
            "time": "今天下午3点",
            "avatar_personality": "积极向上、充满活力"
        },
        "model_config": {
            "temperature": 0.8,
            "max_tokens": 300
        },
        "user_id": "100000001",
        "avatar_id": "1"
    }' $AI_RPC_ADDR ai.Ai/Generate 2>&1)

    # 检查是否有 gRPC 错误
    if echo "$RESPONSE" | grep -qE "(Code:|ERROR:)"; then
        print_error "事件故事生成失败"
        echo "$RESPONSE"
    elif echo "$RESPONSE" | grep -q "content"; then
        print_success "事件故事生成成功"
        # 尝试格式化 JSON，如果失败则直接输出原始内容
        if echo "$RESPONSE" | jq '.' 2>/dev/null; then
            : # jq 成功，已经输出了
        else
            echo "$RESPONSE"
        fi
    else
        print_error "事件故事生成失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 Generate - 分身日记生成
test_generate_avatar_diary() {
    print_title "测试 Generate - 分身日记生成"

    print_info "请求参数: prompt_template=avatar_diary, user_id=100000001, avatar_id=1"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "prompt_template": "avatar_diary",
        "variables": {
            "avatar_name": "小智",
            "date": "2025-12-03",
            "events": "今天完成了AI服务的开发，解决了Redis配置问题",
            "emotion": "开心",
            "personality": "认真负责、注重细节"
        },
        "model_config": {
            "temperature": 0.7,
            "max_tokens": 400
        },
        "user_id": "100000001",
        "avatar_id": "1"
    }' $AI_RPC_ADDR ai.Ai/Generate 2>&1)

    # 检查是否有 gRPC 错误
    if echo "$RESPONSE" | grep -qE "(Code:|ERROR:)"; then
        print_error "分身日记生成失败"
        echo "$RESPONSE"
    elif echo "$RESPONSE" | grep -q "content"; then
        print_success "分身日记生成成功"
        # 尝试格式化 JSON，如果失败则直接输出原始内容
        if echo "$RESPONSE" | jq '.' 2>/dev/null; then
            : # jq 成功，已经输出了
        else
            echo "$RESPONSE"
        fi
    else
        print_error "分身日记生成失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 AnalyzeEmotion - 情绪分析
test_analyze_emotion() {
    print_title "测试 AnalyzeEmotion - 情绪分析"

    print_info "请求参数: text=今天真是太开心了..., user_id=100000001"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "text": "今天真是太开心了！项目终于上线了，团队的努力得到了回报，感觉充满了成就感。",
        "user_id": "100000001"
    }' $AI_RPC_ADDR ai.Ai/AnalyzeEmotion 2>&1)

    # 检查是否有 gRPC 错误
    if echo "$RESPONSE" | grep -qE "(Code:|ERROR:)"; then
        print_error "情绪分析失败"
        echo "$RESPONSE"
    elif echo "$RESPONSE" | grep -q "emotion"; then
        print_success "情绪分析成功"
        # 尝试格式化 JSON，如果失败则直接输出原始内容
        if echo "$RESPONSE" | jq '.' 2>/dev/null; then
            : # jq 成功，已经输出了
        else
            echo "$RESPONSE"
        fi
    else
        print_error "情绪分析失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 AnalyzeEmotion - 负面情绪
test_analyze_emotion_negative() {
    print_title "测试 AnalyzeEmotion - 负面情绪"

    print_info "请求参数: text=今天很沮丧..., user_id=100000001"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "text": "今天很沮丧，遇到了很多困难，感觉压力很大，不知道该怎么办。",
        "user_id": "100000001"
    }' $AI_RPC_ADDR ai.Ai/AnalyzeEmotion 2>&1)

    # 检查是否有 gRPC 错误
    if echo "$RESPONSE" | grep -qE "(Code:|ERROR:)"; then
        print_error "负面情绪分析失败"
        echo "$RESPONSE"
    elif echo "$RESPONSE" | grep -q "emotion"; then
        print_success "负面情绪分析成功"
        # 尝试格式化 JSON，如果失败则直接输出原始内容
        if echo "$RESPONSE" | jq '.' 2>/dev/null; then
            : # jq 成功，已经输出了
        else
            echo "$RESPONSE"
        fi
    else
        print_error "负面情绪分析失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 Embedding - 文本向量化
test_embedding() {
    print_title "测试 Embedding - 文本向量化"

    print_info "请求参数: text=这是一段测试文本"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "text": "这是一段测试文本，用于生成向量表示。"
    }' $AI_RPC_ADDR ai.Ai/Embedding 2>&1)

    # 检查是否有错误（gRPC 错误会包含 "Code:" 或 "ERROR:"）
    if echo "$RESPONSE" | grep -qE "(Code:|ERROR:)"; then
        print_error "文本向量化失败（功能待实现）"
        echo "$RESPONSE"
    elif echo "$RESPONSE" | grep -q "dimension"; then
        print_success "文本向量化接口调用成功（返回空向量，功能待实现）"
        # 尝试格式化 JSON，如果失败则直接输出原始内容
        if echo "$RESPONSE" | jq '.' 2>/dev/null; then
            : # jq 成功，已经输出了
        else
            echo "$RESPONSE"
        fi
    else
        print_error "文本向量化失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetStats - 获取使用统计
test_get_stats_user() {
    print_title "测试 GetStats - 获取用户统计"

    print_info "请求参数: user_id=100000001"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "user_id": "100000001"
    }' $AI_RPC_ADDR ai.Ai/GetStats 2>&1)

    if echo "$RESPONSE" | grep -qE "(Code:|ERROR:)"; then
        print_error "获取统计失败"
        echo "$RESPONSE"
    else
        print_success "获取用户统计成功"
        echo "$RESPONSE" | jq '.'
    fi
    echo ""
}

# 测试 GetStats - 获取全部统计
test_get_stats_all() {
    print_title "测试 GetStats - 获取全部统计"

    print_info "请求参数: user_id=0 (全部用户)"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "user_id": "0"
    }' $AI_RPC_ADDR ai.Ai/GetStats 2>&1)

    if echo "$RESPONSE" | grep -qE "(Code:|ERROR:)"; then
        print_error "获取全部统计失败"
        echo "$RESPONSE"
    else
        print_success "获取全部统计成功"
        echo "$RESPONSE" | jq '.'
    fi
    echo ""
}

# 测试 CalculateActionIntent - 行动意图计算（冒险型性格，下午时段）
test_calculate_action_intent_adventurous() {
    print_title "测试 CalculateActionIntent - 行动意图计算（冒险型）"

    print_info "请求参数: 冒险型性格 + 下午时段 + 高能量"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "avatar_id": "1",
        "personality": {
            "warmth": 70,
            "adventurous": 85,
            "social": 65,
            "creative": 60,
            "calm": 50,
            "energetic": 80
        },
        "time_context": {
            "current_time": "14:30",
            "time_period": "afternoon",
            "day_of_week": "Tuesday"
        },
        "current_state": {
            "energy": 80,
            "emotion_state": "excited",
            "mood": "good"
        },
        "recent_interactions": [
            "用户昨天问候了分身",
            "用户分享了一张风景照片"
        ],
        "recent_events": [
            "昨天去了咖啡厅",
            "阅读了一本有趣的书"
        ]
    }' $AI_RPC_ADDR ai.Ai/CalculateActionIntent 2>&1)

    if echo "$RESPONSE" | grep -q "recommendedAction"; then
        print_success "行动意图计算成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "行动意图计算失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 CalculateActionIntent - 社交型性格，下午
test_calculate_action_intent_social() {
    print_title "测试 CalculateActionIntent - 行动意图计算（社交型）"

    print_info "请求参数: 社交型性格 + 下午时段 + 正常能量"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "avatar_id": "2",
        "personality": {
            "warmth": 90,
            "adventurous": 50,
            "social": 90,
            "creative": 55,
            "calm": 70,
            "energetic": 65
        },
        "time_context": {
            "current_time": "16:00",
            "time_period": "afternoon",
            "day_of_week": "Friday"
        },
        "current_state": {
            "energy": 70,
            "emotion_state": "happy",
            "mood": "cheerful"
        },
        "recent_interactions": [],
        "recent_events": [
            "上午独自在家学习"
        ]
    }' $AI_RPC_ADDR ai.Ai/CalculateActionIntent 2>&1)

    if echo "$RESPONSE" | grep -q "recommendedAction"; then
        print_success "行动意图计算成功（社交型应倾向社交行为）"
        echo "$RESPONSE" | jq '.'
    else
        print_error "行动意图计算失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 CalculateActionIntent - 低能量晚上
test_calculate_action_intent_tired() {
    print_title "测试 CalculateActionIntent - 行动意图计算（疲惫状态）"

    print_info "请求参数: 低能量 + 晚上时段"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "avatar_id": "3",
        "personality": {
            "warmth": 60,
            "adventurous": 40,
            "social": 50,
            "creative": 70,
            "calm": 80,
            "energetic": 40
        },
        "time_context": {
            "current_time": "21:30",
            "time_period": "night",
            "day_of_week": "Monday"
        },
        "current_state": {
            "energy": 30,
            "emotion_state": "tired",
            "mood": "peaceful"
        },
        "recent_interactions": [],
        "recent_events": [
            "今天工作了一整天"
        ]
    }' $AI_RPC_ADDR ai.Ai/CalculateActionIntent 2>&1)

    if echo "$RESPONSE" | grep -q "recommendedAction"; then
        print_success "行动意图计算成功（低能量应倾向休息）"
        echo "$RESPONSE" | jq '.'
    else
        print_error "行动意图计算失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 CalculateActionIntent - 创作型早晨
test_calculate_action_intent_creative() {
    print_title "测试 CalculateActionIntent - 行动意图计算（创作型）"

    print_info "请求参数: 高创造性 + 早晨时段"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "avatar_id": "4",
        "personality": {
            "warmth": 65,
            "adventurous": 60,
            "social": 45,
            "creative": 90,
            "calm": 75,
            "energetic": 70
        },
        "time_context": {
            "current_time": "08:00",
            "time_period": "morning",
            "day_of_week": "Saturday"
        },
        "current_state": {
            "energy": 85,
            "emotion_state": "inspired",
            "mood": "focused"
        },
        "recent_interactions": [
            "用户昨晚分享了一些艺术作品"
        ],
        "recent_events": []
    }' $AI_RPC_ADDR ai.Ai/CalculateActionIntent 2>&1)

    if echo "$RESPONSE" | grep -q "recommendedAction"; then
        print_success "行动意图计算成功（创作型应倾向创作）"
        echo "$RESPONSE" | jq '.'
    else
        print_error "行动意图计算失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试错误场景 - 无效的 Prompt 模板
test_invalid_prompt_template() {
    print_title "测试错误场景 - 无效的 Prompt 模板"

    print_info "请求参数: prompt_template=invalid_template (应该失败)"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "prompt_template": "invalid_template",
        "variables": {
            "test": "value"
        },
        "user_id": "100000001"
    }' $AI_RPC_ADDR ai.Ai/Chat 2>&1)

    if echo "$RESPONSE" | grep -qE "(not found|找不到|Code:|ERROR:)"; then
        print_success "正确返回错误: Prompt 模板不存在"
        echo "$RESPONSE"
    else
        print_error "应该返回错误但没有"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试错误场景 - 缺少必需变量
test_missing_variables() {
    print_title "测试错误场景 - 缺少必需变量"

    print_info "请求参数: prompt_template=avatar_chat 但缺少变量 (应该失败)"

    RESPONSE=$(grpcurl -plaintext -import-path $PROTO_DIR -proto $PROTO_FILE -d '{
        "prompt_template": "avatar_chat",
        "variables": {},
        "user_id": "100000001"
    }' $AI_RPC_ADDR ai.Ai/Chat 2>&1)

    if echo "$RESPONSE" | grep -qE "(missing|缺少|Code:|ERROR:)"; then
        print_success "正确返回错误: 缺少必需变量"
        echo "$RESPONSE"
    else
        print_error "应该返回错误但没有"
        echo "$RESPONSE"
    fi
    echo ""
}

# 主测试流程
main() {
    print_title "AI 服务接口测试"
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
    test_chat_avatar
    test_generate_event_story
    test_generate_avatar_diary
    test_analyze_emotion
    test_analyze_emotion_negative
    test_embedding
    test_get_stats_user
    test_get_stats_all

    # 行动意图计算测试（核心功能）
    test_calculate_action_intent_adventurous
    test_calculate_action_intent_social
    test_calculate_action_intent_tired
    test_calculate_action_intent_creative

    # 错误场景测试
    test_invalid_prompt_template
    test_missing_variables

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
        "chat")
            test_chat_avatar
            ;;
        "generate")
            test_generate_event_story
            test_generate_avatar_diary
            ;;
        "emotion")
            test_analyze_emotion
            test_analyze_emotion_negative
            ;;
        "embedding")
            test_embedding
            ;;
        "stats")
            test_get_stats_user
            test_get_stats_all
            ;;
        "action")
            test_calculate_action_intent_adventurous
            test_calculate_action_intent_social
            test_calculate_action_intent_tired
            test_calculate_action_intent_creative
            ;;
        *)
            echo "用法: $0 [check|chat|generate|emotion|embedding|stats|action]"
            echo ""
            echo "参数说明:"
            echo "  check     - 检查服务状态"
            echo "  chat      - 测试对话生成"
            echo "  generate  - 测试文本生成（事件故事、日记）"
            echo "  emotion   - 测试情绪分析"
            echo "  embedding - 测试文本向量化"
            echo "  stats     - 测试使用统计"
            echo "  action    - 测试行动意图计算（核心功能）⭐"
            echo "  (无参数)  - 运行所有测试"
            ;;
    esac
fi
