package logic

import (
	"github.com/me2/action/rpc/action"
	"github.com/me2/action/rpc/internal/model"
)

// ModelToProtoActionLog 将 Model 的 ActionLog 转换为 Proto 的 ActionLog
func ModelToProtoActionLog(m *model.ActionLog) *action.ActionLog {
	return &action.ActionLog{
		Id:            m.Id,
		AvatarId:      m.AvatarId,
		ActionType:    m.ActionType,
		SceneId:       m.SceneId,
		SceneName:     m.SceneName,
		IntentScore:   m.IntentScore,
		TriggerReason: m.TriggerReason,
		EventId:       m.EventId,
		CreatedAt:     m.CreatedAt.Unix(),
	}
}
