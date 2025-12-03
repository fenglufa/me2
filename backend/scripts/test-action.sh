#!/bin/bash

# Action Service 测试脚本
# 测试分身 ID: 4947066875

AVATAR_ID=4947066875

echo "=========================================="
echo "Action Service 测试脚本"
echo "=========================================="
echo ""

# 1. 测试计算行动意图
echo "1. 测试计算行动意图 (CalculateActionIntent)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID}" \
  localhost:8007 action.Action/CalculateActionIntent
echo ""
echo ""

# 2. 测试调度行动
echo "2. 测试调度行动 (ScheduleAction)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID}" \
  localhost:8007 action.Action/ScheduleAction
echo ""
echo ""

# 3. 测试获取最近一次行动
echo "3. 测试获取最近一次行动 (GetLastAction)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID}" \
  localhost:8007 action.Action/GetLastAction
echo ""
echo ""

# 4. 测试获取行动历史
echo "4. 测试获取行动历史 (GetActionHistory)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID, \"page\": 1, \"page_size\": 10}" \
  localhost:8007 action.Action/GetActionHistory
echo ""
echo ""

echo "=========================================="
echo "测试完成"
echo "=========================================="
