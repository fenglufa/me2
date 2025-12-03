package scheduler

import (
	"context"
	"math/rand"
	"time"

	"github.com/me2/scheduler/rpc/internal/model"
	"github.com/me2/scheduler/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

// SchedulerCore 核心调度器
type SchedulerCore struct {
	svcCtx   *svc.ServiceContext
	ticker   *time.Ticker
	ctx      context.Context
	cancel   context.CancelFunc
	executor *TaskExecutor
}

// NewSchedulerCore 创建核心调度器
func NewSchedulerCore(svcCtx *svc.ServiceContext) *SchedulerCore {
	ctx, cancel := context.WithCancel(context.Background())
	return &SchedulerCore{
		svcCtx:   svcCtx,
		ctx:      ctx,
		cancel:   cancel,
		executor: NewTaskExecutor(svcCtx),
	}
}

// Start 启动调度器
func (s *SchedulerCore) Start() {
	interval := time.Duration(s.svcCtx.Config.Schedule.ScanInterval) * time.Second
	s.ticker = time.NewTicker(interval)

	logx.Infof("[Scheduler] 调度器启动，扫描间隔: %v", interval)

	// 立即执行一次扫描
	go s.scan()

	// 定时执行
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.scan()
			case <-s.ctx.Done():
				logx.Info("[Scheduler] 调度器停止")
				return
			}
		}
	}()
}

// Stop 停止调度器
func (s *SchedulerCore) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.cancel()
	logx.Info("[Scheduler] 调度器已停止")
}

// scan 扫描需要调度的分身
func (s *SchedulerCore) scan() {
	ctx := context.Background()

	// 1. 查询需要调度的分身
	schedules, err := s.svcCtx.AvatarScheduleModel.FindDueSchedules(ctx, time.Now())
	if err != nil {
		logx.Errorf("[Scheduler] 查询待调度分身失败: %v", err)
		return
	}

	if len(schedules) == 0 {
		logx.Infof("[Scheduler] 扫描完成，无待调度分身")
		return
	}

	logx.Infof("[Scheduler] 找到 %d 个待调度分身", len(schedules))

	// 2. 并发执行调度（使用信号量控制并发数）
	maxWorkers := s.svcCtx.Config.Schedule.MaxWorkers
	sem := make(chan struct{}, maxWorkers)

	for _, schedule := range schedules {
		sem <- struct{}{} // 获取信号量
		go func(sch *model.AvatarSchedule) {
			defer func() { <-sem }() // 释放信号量
			s.executor.Execute(sch)
		}(schedule)
	}
}

// CalculateNextScheduleTime 计算下次调度时间
func (s *SchedulerCore) CalculateNextScheduleTime() time.Time {
	minHours := s.svcCtx.Config.ActionSchedule.MinIntervalHours
	maxHours := s.svcCtx.Config.ActionSchedule.MaxIntervalHours

	// 在最小和最大间隔之间随机选择
	randomHours := minHours + rand.Intn(maxHours-minHours+1)
	return time.Now().Add(time.Duration(randomHours) * time.Hour)
}
