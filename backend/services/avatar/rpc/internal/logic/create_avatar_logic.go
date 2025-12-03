package logic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/avatar/rpc/internal/model"
	"github.com/me2/avatar/rpc/internal/personality"
	"github.com/me2/avatar/rpc/internal/svc"
	"github.com/me2/scheduler/rpc/scheduler"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAvatarLogic {
	return &CreateAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateAvatarLogic) CreateAvatar(in *avatar.CreateAvatarRequest) (*avatar.CreateAvatarResponse, error) {
	if in.UserId == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 检查是否已存在分身
	existing, err := l.svcCtx.AvatarModel.FindByUserId(in.UserId)
	if err != nil && err != sql.ErrNoRows {
		l.Errorf("查询分身失败: %v", err)
		return nil, fmt.Errorf("查询分身失败")
	}
	if existing != nil {
		return nil, fmt.Errorf("分身已存在")
	}

	// 生成分身ID
	avatarId := l.svcCtx.IDGen.NextID()

	// 生成人格
	p := personality.GeneratePersonality(in.Gender, in.BirthDate, in.Occupation, in.MaritalStatus)

	// 创建分身
	av := &model.Avatar{
		AvatarId:      avatarId,
		UserId:        in.UserId,
		Nickname:      in.Nickname,
		AvatarUrl:     in.AvatarUrl,
		Gender:        in.Gender,
		BirthDate:     in.BirthDate,
		Occupation:    in.Occupation,
		MaritalStatus: in.MaritalStatus,
		Warmth:        p.Warmth,
		Adventurous:   p.Adventurous,
		Social:        p.Social,
		Creative:      p.Creative,
		Calm:          p.Calm,
		Energetic:     p.Energetic,
	}

	_, err = l.svcCtx.AvatarModel.Insert(av)
	if err != nil {
		l.Errorf("创建分身失败: %v", err)
		return nil, fmt.Errorf("创建分身失败")
	}

	// 自动启用调度（异步调用，失败不影响分身创建）
	go func() {
		scheduleCtx := context.Background()
		_, scheduleErr := l.svcCtx.SchedulerRpc.EnableAvatarSchedule(scheduleCtx, &scheduler.EnableAvatarScheduleRequest{
			AvatarId: avatarId,
		})
		if scheduleErr != nil {
			l.Errorf("自动启用分身调度失败 (avatar_id=%d): %v", avatarId, scheduleErr)
		} else {
			l.Infof("已为分身 %d 自动启用调度", avatarId)
		}
	}()

	return &avatar.CreateAvatarResponse{
		AvatarId: avatarId,
		Personality: &avatar.PersonalityInfo{
			Warmth:      p.Warmth,
			Adventurous: p.Adventurous,
			Social:      p.Social,
			Creative:    p.Creative,
			Calm:        p.Calm,
			Energetic:   p.Energetic,
		},
	}, nil
}
