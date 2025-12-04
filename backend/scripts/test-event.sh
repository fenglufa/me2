#!/bin/bash

# Event Service 测试脚本
# 测试分身 ID: 4947066875

AVATAR_ID=4947066875
SCENE_ID=1

# grpcurl 超时时间（秒）
TIMEOUT=35

echo "=========================================="
echo "Event Service 测试脚本"
echo "=========================================="
echo ""

# 1. 测试生成事件
echo "1. 测试生成事件 (GenerateEvent) - exploration"
echo "------------------------------------------"
grpcurl -plaintext -max-time $TIMEOUT -d "{\"avatar_id\": $AVATAR_ID, \"action_type\": \"exploration\", \"scene_id\": $SCENE_ID}" \
  localhost:8086 event.Event/GenerateEvent
echo ""
echo ""

# 2. 测试生成事件 - social
echo "2. 测试生成事件 (GenerateEvent) - social"
echo "------------------------------------------"
grpcurl -plaintext -max-time $TIMEOUT -d "{\"avatar_id\": $AVATAR_ID, \"action_type\": \"social\", \"scene_id\": $SCENE_ID}" \
  localhost:8086 event.Event/GenerateEvent
echo ""
echo ""

# 3. 测试获取事件时间线
echo "3. 测试获取事件时间线 (GetEventTimeline)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID, \"page\": 1, \"page_size\": 10}" \
  localhost:8086 event.Event/GetEventTimeline
echo ""
echo ""

# 4. 测试获取事件详情（使用第一个事件ID）
echo "4. 测试获取事件详情 (GetEventDetail)"
echo "------------------------------------------"
EVENT_ID=$(grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID, \"page\": 1, \"page_size\": 1}" \
  localhost:8086 event.Event/GetEventTimeline 2>/dev/null | grep -o '"event_id": "[0-9]*"' | head -1 | grep -o '[0-9]*')

if [ -n "$EVENT_ID" ]; then
  grpcurl -plaintext -d "{\"event_id\": $EVENT_ID}" \
    localhost:8086 event.Event/GetEventDetail
else
  echo "未找到事件ID，跳过此测试"
fi
echo ""
echo ""

# 5. 测试获取模板列表
echo "5. 测试获取模板列表 (GetTemplates)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"category\": \"exploration\"}" \
  localhost:8086 event.Event/GetTemplates
echo ""
echo ""

echo "=========================================="
echo "测试完成"
echo "=========================================="
