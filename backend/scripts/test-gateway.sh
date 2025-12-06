#!/bin/bash

# Gateway 测试脚本
# 使用已有的 token 测试所有需要认证的接口

BASE_URL="http://localhost:8888"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzEwNjEzMjgsInVzZXJfaWQiOjI0ODY1ODM0M30.wxbnLJg0zD00Kpd7iGFNXvNVhI6HK4o9f3Zkezj-jQk"
USER_ID=248658343

echo "=========================================="
echo "Gateway API 测试"
echo "=========================================="
echo ""

# 测试用户服务
echo "【用户服务】"
echo "1. 获取用户信息"
curl -s -X GET "$BASE_URL/api/v1/user/info" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "2. 获取订阅信息"
curl -s -X GET "$BASE_URL/api/v1/user/subscription" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# 测试分身服务
echo "【分身服务】"
echo "3. 获取我的分身"
curl -s -X GET "$BASE_URL/api/v1/avatar/my" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# 测试世界服务
echo "【世界服务】"
echo "4. 获取地图列表"
curl -s -X GET "$BASE_URL/api/v1/world/maps?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "5. 获取区域列表"
curl -s -X GET "$BASE_URL/api/v1/world/regions?map_id=1&page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "6. 获取场景列表"
curl -s -X GET "$BASE_URL/api/v1/world/scenes?region_id=1&page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# 测试事件服务
echo "【事件服务】"
echo "7. 获取事件时间线"
curl -s -X GET "$BASE_URL/api/v1/events/timeline?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# 测试行动服务
echo "【行动服务】"
echo "8. 计算行动意图"
curl -s -X GET "$BASE_URL/api/v1/action/intent" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "9. 获取行动历史"
curl -s -X GET "$BASE_URL/api/v1/action/history?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "10. 获取最近行动"
curl -s -X GET "$BASE_URL/api/v1/action/last" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# 测试日记服务
echo "【日记服务】"
echo "11. 获取分身日记列表"
curl -s -X GET "$BASE_URL/api/v1/diary/avatar?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "12. 获取用户日记列表"
curl -s -X GET "$BASE_URL/api/v1/diary/user?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "13. 获取日记统计"
curl -s -X GET "$BASE_URL/api/v1/diary/stats" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "=========================================="
echo "测试完成"
echo "=========================================="
