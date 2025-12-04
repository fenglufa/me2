#!/bin/bash

# Diary Service 测试脚本
# 测试分身 ID: 4947066875

AVATAR_ID=4947066875

# grpcurl 超时时间（秒）
TIMEOUT=35

echo "=========================================="
echo "Diary Service 测试脚本"
echo "=========================================="
echo ""

# 1. 测试生成分身日记
echo "1. 测试生成分身日记 (GenerateAvatarDiary)"
echo "------------------------------------------"
echo "注意：此接口需要先有事件数据，并且依赖 AI 服务生成日记内容"
grpcurl -plaintext -max-time $TIMEOUT -d "{\"avatar_id\": $AVATAR_ID, \"date\": \"2025-12-04\"}" \
  localhost:8011 diary.Diary/GenerateAvatarDiary
echo ""
echo ""

# 2. 测试创建用户日记
echo "2. 测试创建用户日记 (CreateUserDiary)"
echo "------------------------------------------"
echo "注意：此接口依赖 AI 服务进行情感分析和生成回复"
grpcurl -plaintext -max-time $TIMEOUT -d "{\"avatar_id\": $AVATAR_ID, \"title\": \"今天的心情\", \"content\": \"今天天气很好，心情也很愉快。和朋友一起去了公园，度过了美好的一天。\", \"tags\": [\"日常\", \"朋友\"]}" \
  localhost:8011 diary.Diary/CreateUserDiary
echo ""
echo ""

# 3. 测试获取分身日记列表
echo "3. 测试获取分身日记列表 (GetAvatarDiaryList)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID, \"page\": 1, \"page_size\": 10}" \
  localhost:8011 diary.Diary/GetAvatarDiaryList
echo ""
echo ""

# 4. 测试获取分身日记列表（带日期过滤）
echo "4. 测试获取分身日记列表 - 带日期过滤"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID, \"page\": 1, \"page_size\": 10, \"start_date\": \"2025-12-01\", \"end_date\": \"2025-12-31\"}" \
  localhost:8011 diary.Diary/GetAvatarDiaryList
echo ""
echo ""

# 5. 测试获取用户日记列表
echo "5. 测试获取用户日记列表 (GetUserDiaryList)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID, \"page\": 1, \"page_size\": 10}" \
  localhost:8011 diary.Diary/GetUserDiaryList
echo ""
echo ""

# 6. 测试获取用户日记列表（带日期过滤）
echo "6. 测试获取用户日记列表 - 带日期过滤"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID, \"page\": 1, \"page_size\": 10, \"start_date\": \"2025-12-01\", \"end_date\": \"2025-12-31\"}" \
  localhost:8011 diary.Diary/GetUserDiaryList
echo ""
echo ""

# 7. 测试获取日记统计
echo "7. 测试获取日记统计 (GetDiaryStats)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID}" \
  localhost:8011 diary.Diary/GetDiaryStats
echo ""
echo ""

# 8. 测试重复创建用户日记（应该失败）
echo "8. 测试重复创建用户日记（预期失败：今天已经写过日记）"
echo "------------------------------------------"
grpcurl -plaintext -max-time $TIMEOUT -d "{\"avatar_id\": $AVATAR_ID, \"title\": \"糟糕的一天\", \"content\": \"今天遇到了很多不顺心的事情，工作上出了问题，心情很糟糕。\"}" \
  localhost:8011 diary.Diary/CreateUserDiary
echo ""
echo ""

# 9. 测试获取日记详情（通过查询列表验证）
echo "9. 测试验证日记创建成功"
echo "------------------------------------------"
echo "查询今天的用户日记，验证是否创建成功"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID, \"page\": 1, \"page_size\": 1, \"start_date\": \"2025-12-04\", \"end_date\": \"2025-12-04\"}" \
  localhost:8011 diary.Diary/GetUserDiaryList
echo ""
echo ""

echo "=========================================="
echo "测试完成"
echo "=========================================="
echo ""
echo "提示："
echo "1. GenerateAvatarDiary 需要确保有事件数据，可以先运行 test-event.sh 生成事件"
echo "2. CreateUserDiary 会调用 AI 服务进行情感分析和生成回复，请确保 AI 服务正常运行"
echo "3. 如果某些测试失败，请检查相关依赖服务是否正常运行（Event、AI、Avatar）"
