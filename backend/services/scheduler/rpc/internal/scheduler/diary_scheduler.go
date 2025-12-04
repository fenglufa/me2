package scheduler

import (
	"context"
	"time"

	"github.com/me2/diary/rpc/diary_client"
	"github.com/me2/scheduler/rpc/internal/svc"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

// DiaryScheduler 日记生成调度器
type DiaryScheduler struct {
	svcCtx *svc.ServiceContext
	cron   *cron.Cron
}

// NewDiaryScheduler 创建日记调度器
func NewDiaryScheduler(svcCtx *svc.ServiceContext) *DiaryScheduler {
	return &DiaryScheduler{
		svcCtx: svcCtx,
		cron:   cron.New(),
	}
}

// Start 启动日记调度器
func (s *DiaryScheduler) Start() {
	cronExpr := s.svcCtx.Config.DiarySchedule.CronExpression
	if cronExpr == "" {
		cronExpr = "0 22 * * *" // 默认每天 22:00
	}

	_, err := s.cron.AddFunc(cronExpr, s.generateAvatarDiaries)
	if err != nil {
		logx.Errorf("[DiaryScheduler] 添加定时任务失败: %v", err)
		return
	}

	s.cron.Start()
	logx.Infof("[DiaryScheduler] 日记调度器启动，Cron 表达式: %s", cronExpr)
}

// Stop 停止日记调度器
func (s *DiaryScheduler) Stop() {
	s.cron.Stop()
	logx.Info("[DiaryScheduler] 日记调度器已停止")
}

// generateAvatarDiaries 生成所有分身的日记
func (s *DiaryScheduler) generateAvatarDiaries() {
	ctx := context.Background()

	// 查询所有启用的分身
	schedules, err := s.svcCtx.AvatarScheduleModel.FindActiveSchedules(ctx)
	if err != nil {
		logx.Errorf("[DiaryScheduler] 查询启用的分身失败: %v", err)
		return
	}

	if len(schedules) == 0 {
		logx.Info("[DiaryScheduler] 没有启用的分身，跳过日记生成")
		return
	}

	logx.Infof("[DiaryScheduler] 开始为 %d 个分身生成日记", len(schedules))

	today := time.Now().Format("2006-01-02")
	successCount := 0
	failCount := 0

	for _, schedule := range schedules {
		diaryCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		_, err := s.svcCtx.DiaryRpc.GenerateAvatarDiary(diaryCtx, &diary_client.GenerateAvatarDiaryRequest{
			AvatarId: schedule.AvatarId,
			Date:     today,
		})
		cancel()

		if err != nil {
			logx.Errorf("[DiaryScheduler] 分身 %d 日记生成失败: %v", schedule.AvatarId, err)
			failCount++
		} else {
			logx.Infof("[DiaryScheduler] 分身 %d 日记生成成功", schedule.AvatarId)
			successCount++
		}
	}

	logx.Infof("[DiaryScheduler] 日记生成完成，成功: %d，失败: %d", successCount, failCount)
}
