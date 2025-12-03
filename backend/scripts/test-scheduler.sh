#!/bin/bash

# Scheduler Service 测试脚本
# 测试分身 ID: 4947066875

AVATAR_ID=4947066875
AVATAR_ID_2=1234567890
AVATAR_ID_3=9876543210

echo "=========================================="
echo "Scheduler Service 测试脚本"
echo "=========================================="
echo ""

# 1. 测试启用分身调度
echo "1. 测试启用分身调度 (EnableAvatarSchedule)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID}" \
  localhost:8008 scheduler.Scheduler/EnableAvatarSchedule
echo ""
echo ""

# 2. 测试查询调度状态
echo "2. 测试查询调度状态 (GetAvatarScheduleStatus)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID}" \
  localhost:8008 scheduler.Scheduler/GetAvatarScheduleStatus
echo ""
echo ""

# 3. 测试手动触发调度
echo "3. 测试手动触发调度 (TriggerSchedule)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID}" \
  localhost:8008 scheduler.Scheduler/TriggerSchedule
echo ""
echo ""

# 4. 测试暂停分身调度
echo "4. 测试暂停分身调度 (PauseAvatarSchedule)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID}" \
  localhost:8008 scheduler.Scheduler/PauseAvatarSchedule
echo ""
echo ""

# 5. 再次查询状态（验证暂停状态）
echo "5. 再次查询状态（验证暂停状态）"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID}" \
  localhost:8008 scheduler.Scheduler/GetAvatarScheduleStatus
echo ""
echo ""

# 6. 测试恢复分身调度
echo "6. 测试恢复分身调度 (ResumeAvatarSchedule)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID}" \
  localhost:8008 scheduler.Scheduler/ResumeAvatarSchedule
echo ""
echo ""

# 7. 再次查询状态（验证恢复状态）
echo "7. 再次查询状态（验证恢复状态）"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID}" \
  localhost:8008 scheduler.Scheduler/GetAvatarScheduleStatus
echo ""
echo ""

# 8. 启用更多分身（用于批量测试）
echo "8. 启用更多分身（用于批量测试）"
echo "------------------------------------------"
echo "启用分身 $AVATAR_ID_2:"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID_2}" \
  localhost:8008 scheduler.Scheduler/EnableAvatarSchedule
echo ""
echo "启用分身 $AVATAR_ID_3:"
grpcurl -plaintext -d "{\"avatar_id\": $AVATAR_ID_3}" \
  localhost:8008 scheduler.Scheduler/EnableAvatarSchedule
echo ""
echo ""

# 9. 测试批量查询调度状态
echo "9. 测试批量查询调度状态 (BatchGetScheduleStatus)"
echo "------------------------------------------"
grpcurl -plaintext -d "{\"avatar_ids\": [$AVATAR_ID, $AVATAR_ID_2, $AVATAR_ID_3]}" \
  localhost:8008 scheduler.Scheduler/BatchGetScheduleStatus
echo ""
echo ""

echo "=========================================="
echo "测试完成"
echo "=========================================="
echo ""
echo "注意事项："
echo "- 调度器每 60 秒扫描一次数据库"
echo "- 下次调度时间为 2-6 小时后的随机时间"
echo "- 手动触发调度会立即执行一次行为"
echo "- 连续失败 5 次会自动暂停调度"
echo ""
