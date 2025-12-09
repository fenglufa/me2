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

	avatarInfo := b.formatAvatarInfo(avatarResp.Avatar)

	recentEvents, err := b.getRecentEvents(avatarID)
	if err != nil {
		return "", err
	}

	history, err := b.getConversationHistory(sessionID)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`你是用户的AI分身,具有以下基本信息:
%s

最近发生的事件:
%s

对话历史:
%s

请以分身的口吻回复用户。回复要自然、简洁,不要过于正式。`,
		avatarInfo, recentEvents, history)

	return prompt, nil
}

func (b *ContextBuilder) formatAvatarInfo(avatar *avatar_client.AvatarInfo) string {
	if avatar == nil {
		return "暂无分身信息"
	}

	genderMap := map[int32]string{1: "男", 2: "女", 3: "其他"}
	maritalMap := map[int32]string{1: "单身", 2: "恋爱中", 3: "已婚", 4: "其他"}

	info := []string{
		fmt.Sprintf("昵称: %s", avatar.Nickname),
		fmt.Sprintf("性别: %s", genderMap[avatar.Gender]),
		fmt.Sprintf("出生日期: %s", avatar.BirthDate),
		fmt.Sprintf("职业: %s", avatar.Occupation),
		fmt.Sprintf("婚姻状态: %s", maritalMap[avatar.MaritalStatus]),
	}

	return strings.Join(info, "\n")
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
	runes := []rune(message)
	if len(runes) > 50 {
		truncated = string(runes[:50]) + "..."
	}
	return b.svcCtx.SessionModel.UpdateLastMessage(sessionID, truncated)
}
