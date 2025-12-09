#!/usr/bin/env zsh

# Note 服务接口测试脚本
# 使用 grpcurl 测试 Note RPC 服务

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 服务地址
NOTE_RPC_ADDR="127.0.0.1:8010"

# 测试用的用户ID
TEST_USER_ID=1
TEST_AVATAR_ID=5

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
    print_title "检查 Note 服务状态"

    if grpcurl -plaintext $NOTE_RPC_ADDR list > /dev/null 2>&1; then
        print_success "Note 服务运行正常"
        echo ""
        grpcurl -plaintext $NOTE_RPC_ADDR list
        echo ""
    else
        print_error "Note 服务未启动或连接失败"
        exit 1
    fi
}

# 测试 CreateNote - 创建普通笔记
test_create_note() {
    print_title "测试 CreateNote - 创建普通笔记"

    print_info "请求参数: user_id=$TEST_USER_ID, raw_text='今天天气很好，心情很愉快'"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"avatar_id\": $TEST_AVATAR_ID,
        \"raw_text\": \"今天天气很好，心情很愉快\"
    }" $NOTE_RPC_ADDR note.NoteService/CreateNote 2>&1)

    if echo "$RESPONSE" | grep -q "noteId"; then
        print_success "创建笔记成功"
        echo "$RESPONSE" | jq '.'

        # 提取 note_id 供后续测试使用
        export NOTE_ID=$(echo "$RESPONSE" | jq -r '.noteId')
        echo ""
        print_info "笔记ID: $NOTE_ID"
    else
        print_error "创建笔记失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 CreateNote - 创建情绪笔记
test_create_emotion_note() {
    print_title "测试 CreateNote - 创建情绪笔记"

    print_info "请求参数: user_id=$TEST_USER_ID, raw_text='今天工作压力很大，感到很焦虑'"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"raw_text\": \"今天工作压力很大，感到很焦虑\"
    }" $NOTE_RPC_ADDR note.NoteService/CreateNote 2>&1)

    if echo "$RESPONSE" | grep -q "noteId"; then
        print_success "创建情绪笔记成功"
        echo "$RESPONSE" | jq '.'

        export EMOTION_NOTE_ID=$(echo "$RESPONSE" | jq -r '.noteId')
        print_info "情绪笔记ID: $EMOTION_NOTE_ID"
    else
        print_error "创建情绪笔记失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 CreateNote - 创建TODO笔记
test_create_todo_note() {
    print_title "测试 CreateNote - 创建TODO笔记"

    print_info "请求参数: user_id=$TEST_USER_ID, raw_text='明天要完成项目文档'"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"raw_text\": \"明天要完成项目文档\"
    }" $NOTE_RPC_ADDR note.NoteService/CreateNote 2>&1)

    if echo "$RESPONSE" | grep -q "noteId"; then
        print_success "创建TODO笔记成功"
        echo "$RESPONSE" | jq '.'

        export TODO_NOTE_ID=$(echo "$RESPONSE" | jq -r '.noteId')
        print_info "TODO笔记ID: $TODO_NOTE_ID"
    else
        print_error "创建TODO笔记失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 CreateNote - 创建记账笔记
test_create_expense_note() {
    print_title "测试 CreateNote - 创建记账笔记"

    print_info "请求参数: user_id=$TEST_USER_ID, raw_text='今天吃午饭花了35块钱'"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"raw_text\": \"今天吃午饭花了35块钱，买咖啡花了25元\"
    }" $NOTE_RPC_ADDR note.NoteService/CreateNote 2>&1)

    if echo "$RESPONSE" | grep -q "noteId"; then
        print_success "创建记账笔记成功"
        echo "$RESPONSE" | jq '.'

        export EXPENSE_NOTE_ID=$(echo "$RESPONSE" | jq -r '.noteId')
        print_info "记账笔记ID: $EXPENSE_NOTE_ID"
    else
        print_error "创建记账笔记失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetNotes - 获取笔记列表
test_get_notes() {
    print_title "测试 GetNotes - 获取笔记列表"

    print_info "请求参数: user_id=$TEST_USER_ID, page=1, page_size=10"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"page\": 1,
        \"page_size\": 10
    }" $NOTE_RPC_ADDR note.NoteService/GetNotes 2>&1)

    if echo "$RESPONSE" | grep -q "notes"; then
        print_success "获取笔记列表成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取笔记列表失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetNotes - 按类型过滤
test_get_notes_by_type() {
    print_title "测试 GetNotes - 按类型过滤"

    print_info "请求参数: user_id=$TEST_USER_ID, types=['emotion']"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"page\": 1,
        \"page_size\": 10,
        \"types\": [\"emotion\"]
    }" $NOTE_RPC_ADDR note.NoteService/GetNotes 2>&1)

    if echo "$RESPONSE" | grep -q "notes"; then
        print_success "按类型过滤笔记成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "按类型过滤笔记失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetNoteDetail - 获取笔记详情
test_get_note_detail() {
    print_title "测试 GetNoteDetail - 获取笔记详情"

    if [ -z "$NOTE_ID" ]; then
        NOTE_ID=1
    fi

    print_info "请求参数: note_id=$NOTE_ID, user_id=$TEST_USER_ID"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"note_id\": $NOTE_ID,
        \"user_id\": $TEST_USER_ID
    }" $NOTE_RPC_ADDR note.NoteService/GetNoteDetail 2>&1)

    if echo "$RESPONSE" | grep -q "note"; then
        print_success "获取笔记详情成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取笔记详情失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 SearchNotes - 搜索笔记
test_search_notes() {
    print_title "测试 SearchNotes - 搜索笔记"

    print_info "请求参数: user_id=$TEST_USER_ID, query='天气', limit=10"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"query\": \"天气\",
        \"limit\": 10
    }" $NOTE_RPC_ADDR note.NoteService/SearchNotes 2>&1)

    if echo "$RESPONSE" | grep -q "notes"; then
        print_success "搜索笔记成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "搜索笔记失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 SearchNotes - 按时间范围搜索
test_search_notes_by_date() {
    print_title "测试 SearchNotes - 按时间范围搜索"

    print_info "请求参数: user_id=$TEST_USER_ID, date_range='today', limit=10"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"query\": \"\",
        \"date_range\": \"today\",
        \"limit\": 10
    }" $NOTE_RPC_ADDR note.NoteService/SearchNotes 2>&1)

    if echo "$RESPONSE" | grep -q "notes"; then
        print_success "按时间范围搜索成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "按时间范围搜索失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 UpdateNote - 更新笔记
test_update_note() {
    print_title "测试 UpdateNote - 更新笔记"

    if [ -z "$NOTE_ID" ]; then
        NOTE_ID=1
    fi

    print_info "请求参数: note_id=$NOTE_ID, user_id=$TEST_USER_ID, raw_text='更新后的笔记内容'"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"note_id\": $NOTE_ID,
        \"user_id\": $TEST_USER_ID,
        \"raw_text\": \"今天天气很好，心情非常愉快，更新了笔记内容\"
    }" $NOTE_RPC_ADDR note.NoteService/UpdateNote 2>&1)

    if echo "$RESPONSE" | grep -q "success"; then
        print_success "更新笔记成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "更新笔记失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetTodos - 获取TODO列表
test_get_todos() {
    print_title "测试 GetTodos - 获取TODO列表"

    print_info "请求参数: user_id=$TEST_USER_ID, status=-1 (全部), page=1, page_size=10"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"status\": -1,
        \"page\": 1,
        \"page_size\": 10
    }" $NOTE_RPC_ADDR note.NoteService/GetTodos 2>&1)

    if echo "$RESPONSE" | grep -q "todos"; then
        print_success "获取TODO列表成功"
        echo "$RESPONSE" | jq '.'

        # 提取第一个 todo_id 供后续测试使用
        export TODO_ID=$(echo "$RESPONSE" | jq -r '.todos[0].id // empty')
        if [ -n "$TODO_ID" ]; then
            print_info "TODO ID: $TODO_ID"
        fi
    else
        print_error "获取TODO列表失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetTodos - 只获取未完成的TODO
test_get_todos_incomplete() {
    print_title "测试 GetTodos - 只获取未完成的TODO"

    print_info "请求参数: user_id=$TEST_USER_ID, status=0 (未完成)"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"status\": 0,
        \"page\": 1,
        \"page_size\": 10
    }" $NOTE_RPC_ADDR note.NoteService/GetTodos 2>&1)

    if echo "$RESPONSE" | grep -q "todos"; then
        print_success "获取未完成TODO列表成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取未完成TODO列表失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 UpdateTodoStatus - 更新TODO状态
test_update_todo_status() {
    print_title "测试 UpdateTodoStatus - 更新TODO状态"

    if [ -z "$TODO_ID" ]; then
        print_error "没有可用的 TODO_ID，跳过此测试"
        echo ""
        return
    fi

    print_info "请求参数: todo_id=$TODO_ID, user_id=$TEST_USER_ID, status=1 (已完成)"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"todo_id\": $TODO_ID,
        \"user_id\": $TEST_USER_ID,
        \"status\": 1
    }" $NOTE_RPC_ADDR note.NoteService/UpdateTodoStatus 2>&1)

    if echo "$RESPONSE" | grep -q "success"; then
        print_success "更新TODO状态成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "更新TODO状态失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetExpenses - 获取记账列表
test_get_expenses() {
    print_title "测试 GetExpenses - 获取记账列表"

    print_info "请求参数: user_id=$TEST_USER_ID, page=1, page_size=10"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"page\": 1,
        \"page_size\": 10
    }" $NOTE_RPC_ADDR note.NoteService/GetExpenses 2>&1)

    if echo "$RESPONSE" | grep -q "expenses"; then
        print_success "获取记账列表成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取记账列表失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetExpenses - 按日期范围过滤
test_get_expenses_by_date() {
    print_title "测试 GetExpenses - 按日期范围过滤"

    # 获取今天的日期
    TODAY=$(date +%Y-%m-%d)

    print_info "请求参数: user_id=$TEST_USER_ID, start_date=$TODAY, end_date=$TODAY"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"start_date\": \"$TODAY\",
        \"end_date\": \"$TODAY\",
        \"page\": 1,
        \"page_size\": 10
    }" $NOTE_RPC_ADDR note.NoteService/GetExpenses 2>&1)

    if echo "$RESPONSE" | grep -q "expenses"; then
        print_success "按日期范围获取记账成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "按日期范围获取记账失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 GetExpenseStats - 获取记账统计
test_get_expense_stats() {
    print_title "测试 GetExpenseStats - 获取记账统计"

    # 获取今天的日期
    TODAY=$(date +%Y-%m-%d)
    WEEK_AGO=$(date -v-7d +%Y-%m-%d)

    print_info "请求参数: user_id=$TEST_USER_ID, start_date=$WEEK_AGO, end_date=$TODAY"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"user_id\": $TEST_USER_ID,
        \"start_date\": \"$WEEK_AGO\",
        \"end_date\": \"$TODAY\"
    }" $NOTE_RPC_ADDR note.NoteService/GetExpenseStats 2>&1)

    if echo "$RESPONSE" | grep -q "totalAmount"; then
        print_success "获取记账统计成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "获取记账统计失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试 DeleteNote - 删除笔记
test_delete_note() {
    print_title "测试 DeleteNote - 删除笔记"

    if [ -z "$EMOTION_NOTE_ID" ]; then
        print_error "没有可用的笔记ID，跳过此测试"
        echo ""
        return
    fi

    print_info "请求参数: note_id=$EMOTION_NOTE_ID, user_id=$TEST_USER_ID"

    RESPONSE=$(grpcurl -plaintext -d "{
        \"note_id\": $EMOTION_NOTE_ID,
        \"user_id\": $TEST_USER_ID
    }" $NOTE_RPC_ADDR note.NoteService/DeleteNote 2>&1)

    if echo "$RESPONSE" | grep -q "success"; then
        print_success "删除笔记成功"
        echo "$RESPONSE" | jq '.'
    else
        print_error "删除笔记失败"
        echo "$RESPONSE"
    fi
    echo ""
}

# 测试错误场景 - 无效参数
test_invalid_params() {
    print_title "测试错误场景 - 无效参数"

    print_info "请求参数: user_id=0 (应该失败)"

    RESPONSE=$(grpcurl -plaintext -d '{
        "user_id": 0,
        "raw_text": "test"
    }' $NOTE_RPC_ADDR note.NoteService/CreateNote 2>&1)

    if echo "$RESPONSE" | grep -qE "(user_id|用户)"; then
        print_success "正确返回错误: 用户ID无效"
        echo "$RESPONSE"
    else
        print_error "应该返回错误但没有"
        echo "$RESPONSE"
    fi
    echo ""
}

# 主测试流程
main() {
    print_title "Note 服务接口测试"
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

    # 运行测试 - 创建笔记
    test_create_note
    test_create_emotion_note
    test_create_todo_note
    test_create_expense_note

    # 测试查询笔记
    test_get_notes
    test_get_notes_by_type
    test_get_note_detail

    # 测试搜索笔记
    test_search_notes
    test_search_notes_by_date

    # 测试更新笔记
    test_update_note

    # 测试TODO功能
    test_get_todos
    test_get_todos_incomplete
    test_update_todo_status

    # 测试记账功能
    test_get_expenses
    test_get_expenses_by_date
    test_get_expense_stats

    # 测试删除笔记
    test_delete_note

    # 测试错误场景
    test_invalid_params

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
            test_create_note
            ;;
        "list")
            test_get_notes
            ;;
        "detail")
            test_get_note_detail
            ;;
        "search")
            test_search_notes
            ;;
        "update")
            test_update_note
            ;;
        "delete")
            test_delete_note
            ;;
        "todos")
            test_get_todos
            ;;
        "expenses")
            test_get_expenses
            ;;
        "stats")
            test_get_expense_stats
            ;;
        *)
            echo "用法: $0 [check|create|list|detail|search|update|delete|todos|expenses|stats]"
            echo ""
            echo "参数说明:"
            echo "  check    - 检查服务状态"
            echo "  create   - 测试创建笔记"
            echo "  list     - 测试获取笔记列表"
            echo "  detail   - 测试获取笔记详情"
            echo "  search   - 测试搜索笔记"
            echo "  update   - 测试更新笔记"
            echo "  delete   - 测试删除笔记"
            echo "  todos    - 测试TODO功能"
            echo "  expenses - 测试记账功能"
            echo "  stats    - 测试记账统计"
            echo "  (无参数) - 运行所有测试"
            ;;
    esac
fi
