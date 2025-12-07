package logic

import (
	"context"
	"fmt"
	"strings"

	"github.com/me2/dialogue/rpc/internal/model"
	"github.com/me2/dialogue/rpc/internal/svc"

	"github.com/me2/avatar/rpc/avatar_client"
	"github.com/me2/event/rpc/event_client"
)

type ContextBuilder struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
}

func NewContextBuilder(ctx context.Context, svcCtx *svc.ServiceContext) *ContextBuilder {
	return &ContextBuilder{
		svcCtx: svcCtx,
		ctx:    ctx,
	}
}

func (b *ContextBuilder) BuildSystemPrompt(avatarID, sessionID int64) (string, error) {
	avatarResp, err := b.svcCtx.AvatarRpc.GetAvatarInfo(b.ctx, &avatar_client.GetAvatarInfoRequest{
		AvatarId: avatarID,
	})
	if err != nil {
		return "", err
	}

	personality := b.formatPersonality(avatarResp.Avatar.Personality)

	recentEvents, err := b.getRecentEvents(avatarID)
	if err != nil {
		return "", err
	}

	history, err := b.getConversationHistory(sessionID)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`你是用户的AI分身,具有以下性格特征:
%s

最近发生的事件:
%s

对话历史:
%s

请以分身的口吻回复用户,体现你的性格特点。回复要自然、简洁,不要过于正式。`,
		personality, recentEvents, history)

	return prompt, nil
}

func (b *ContextBuilder) formatPersonality(p *avatar_client.PersonalityInfo) string {
	if p == nil {
		return "暂无人格信息"
	}

	traits := []string{
		fmt.Sprintf("情绪温度: %d/100 (%s)", p.Warmth, b.describeLevel(p.Warmth)),
		fmt.Sprintf("冒险倾向: %d/100 (%s)", p.Adventurous, b.describeLevel(p.Adventurous)),
		fmt.Sprintf("人际能量: %d/100 (%s)", p.Social, b.describeLevel(p.Social)),
		fmt.Sprintf("创造性: %d/100 (%s)", p.Creative, b.describeLevel(p.Creative)),
		fmt.Sprintf("情绪稳定性: %d/100 (%s)", p.Calm, b.describeLevel(p.Calm)),
		fmt.Sprintf("生活动力: %d/100 (%s)", p.Energetic, b.describeLevel(p.Energetic)),
	}

	return strings.Join(traits, "\n")
}

func (b *ContextBuilder) describeLevel(value int32) string {
	if value >= 80 {
		return "非常高"
	} else if value >= 60 {
		return "较高"
	} else if value >= 40 {
		return "中等"
	} else if value >= 20 {
		return "较低"
	}
	return "很低"
}

func (b *ContextBuilder) getRecentEvents(avatarID int64) (string, error) {
	eventsResp, err := b.svcCtx.EventRpc.GetEventTimeline(b.ctx, &event_client.GetEventTimelineRequest{
		AvatarId: avatarID,
		Page:     1,
		PageSize: b.svcCtx.Config.Context.MaxRecentEvents,
	})
	if err != nil || len(eventsResp.Events) == 0 {
		return "暂无最近事件", nil
	}

	var events []string
	for i, e := range eventsResp.Events {
		events = append(events, fmt.Sprintf("%d. %s - %s", i+1, e.EventTitle, e.EventText))
	}

	return strings.Join(events, "\n"), nil
}

func (b *ContextBuilder) getConversationHistory(sessionID int64) (string, error) {
	messages, err := b.svcCtx.MessageModel.FindRecentBySession(sessionID, b.svcCtx.Config.Context.MaxHistoryMessages)
	if err != nil || len(messages) == 0 {
		return "这是你们的第一次对话", nil
	}

	var history []string
	for _, m := range messages {
		role := "用户"
		if m.Role == "assistant" {
			role = "分身"
		}
		history = append(history, fmt.Sprintf("%s: %s", role, m.Content))
	}

	return strings.Join(history, "\n"), nil
}

func (b *ContextBuilder) SaveMessage(sessionID int64, role, content string) (int64, error) {
	msg := &model.DialogueMessage{
		SessionId: sessionID,
		Role:      role,
		Content:   content,
		CreatedAt: getCurrentTimestamp(),
	}

	result, err := b.svcCtx.MessageModel.Insert(msg)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return id, nil
}

func (b *ContextBuilder) UpdateSessionLastMessage(sessionID int64, message string) error {
	truncated := message
	if len(message) > 50 {
		truncated = message[:50] + "..."
	}
	return b.svcCtx.SessionModel.UpdateLastMessage(sessionID, truncated)
}
