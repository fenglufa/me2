package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/avatar/rpc/internal/model"
	"github.com/me2/avatar/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePersonalityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePersonalityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePersonalityLogic {
	return &UpdatePersonalityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新人格特征（由事件触发）
func (l *UpdatePersonalityLogic) UpdatePersonality(in *avatar.UpdatePersonalityRequest) (*avatar.UpdatePersonalityResponse, error) {
	// 1. 解析人格变化值
	var changes map[string]int32
	if err := json.Unmarshal([]byte(in.PersonalityChanges), &changes); err != nil {
		l.Errorf("解析人格变化失败: %v", err)
		return nil, fmt.Errorf("invalid personality_changes JSON: %w", err)
	}

	// 2. 获取当前分身数据
	av, err := l.svcCtx.AvatarModel.FindByAvatarId(in.AvatarId)
	if err != nil {
		l.Errorf("查询分身失败: %v", err)
		return nil, fmt.Errorf("avatar not found: %w", err)
	}

	// 3. 记录变化前的人格值
	before := &avatar.PersonalityInfo{
		Warmth:      av.Warmth,
		Adventurous: av.Adventurous,
		Social:      av.Social,
		Creative:    av.Creative,
		Calm:        av.Calm,
		Energetic:   av.Energetic,
	}

	// 4. 应用变化并限制范围 0-100
	newWarmth := clamp(av.Warmth+changes["warmth"], 0, 100)
	newAdventurous := clamp(av.Adventurous+changes["adventurous"], 0, 100)
	newSocial := clamp(av.Social+changes["social"], 0, 100)
	newCreative := clamp(av.Creative+changes["creative"], 0, 100)
	newCalm := clamp(av.Calm+changes["calm"], 0, 100)
	newEnergetic := clamp(av.Energetic+changes["energetic"], 0, 100)

	// 5. 更新 avatars 表
	err = l.svcCtx.AvatarModel.UpdatePersonality(in.AvatarId, newWarmth, newAdventurous, newSocial, newCreative, newCalm, newEnergetic)
	if err != nil {
		l.Errorf("更新人格失败: %v", err)
		return nil, fmt.Errorf("failed to update personality: %w", err)
	}

	// 6. 记录变化后的人格值
	after := &avatar.PersonalityInfo{
		Warmth:      newWarmth,
		Adventurous: newAdventurous,
		Social:      newSocial,
		Creative:    newCreative,
		Calm:        newCalm,
		Energetic:   newEnergetic,
	}

	// 7. 保存到 personality_history 表
	beforeJSON, _ := json.Marshal(map[string]int32{
		"warmth":      before.Warmth,
		"adventurous": before.Adventurous,
		"social":      before.Social,
		"creative":    before.Creative,
		"calm":        before.Calm,
		"energetic":   before.Energetic,
	})
	afterJSON, _ := json.Marshal(map[string]int32{
		"warmth":      after.Warmth,
		"adventurous": after.Adventurous,
		"social":      after.Social,
		"creative":    after.Creative,
		"calm":        after.Calm,
		"energetic":   after.Energetic,
	})

	history := &model.PersonalityHistory{
		AvatarId:     in.AvatarId,
		EventId:      in.EventId,
		Changes:      in.PersonalityChanges,
		BeforeValues: string(beforeJSON),
		AfterValues:  string(afterJSON),
	}

	err = l.svcCtx.PersonalityHistoryModel.Insert(history)
	if err != nil {
		l.Errorf("记录人格历史失败: %v", err)
		// 不返回错误，历史记录失败不影响主流程
	}

	l.Infof("分身 %d 人格更新成功，事件 %d", in.AvatarId, in.EventId)
	return &avatar.UpdatePersonalityResponse{
		Success: true,
		Before:  before,
		After:   after,
	}, nil
}

// clamp 限制值在指定范围内
func clamp(value, min, max int32) int32 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
