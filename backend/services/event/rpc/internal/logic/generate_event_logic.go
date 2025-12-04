package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/me2/ai/rpc/ai"
	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/event/rpc/event"
	"github.com/me2/event/rpc/internal/model"
	"github.com/me2/event/rpc/internal/svc"
	"github.com/me2/world/rpc/world"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateEventLogic {
	return &GenerateEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GenerateEvent 生成事件
func (l *GenerateEventLogic) GenerateEvent(in *event.GenerateEventRequest) (*event.GenerateEventResponse, error) {
	// 1. 参数验证
	if in.AvatarId == 0 {
		return nil, fmt.Errorf("avatar_id 不能为空")
	}
	if in.ActionType == "" {
		return nil, fmt.Errorf("action_type 不能为空")
	}
	if in.SceneId == 0 {
		return nil, fmt.Errorf("scene_id 不能为空")
	}

	// 2. 获取分身信息
	avatarInfo, err := l.svcCtx.AvatarRpc.GetAvatarInfo(l.ctx, &avatar.GetAvatarInfoRequest{
		AvatarId: in.AvatarId,
	})
	if err != nil {
		l.Errorf("获取分身信息失败: %v", err)
		return nil, fmt.Errorf("获取分身信息失败")
	}

	// 3. 获取场景详细信息
	sceneInfo, err := l.svcCtx.WorldRpc.GetScene(l.ctx, &world.GetSceneRequest{
		SceneId: in.SceneId,
	})
	if err != nil {
		l.Errorf("获取场景信息失败: %v", err)
		return nil, fmt.Errorf("获取场景信息失败")
	}

	// 4. 查询该行为类型的所有模板
	templates, err := l.svcCtx.EventTemplateModel.FindByCategory(in.ActionType)
	if err != nil {
		l.Errorf("查询事件模板失败: %v", err)
		return nil, fmt.Errorf("查询事件模板失败")
	}

	if len(templates) == 0 {
		l.Errorf("未找到类型为 %s 的事件模板", in.ActionType)
		return nil, fmt.Errorf("未找到合适的事件模板")
	}

	// 5. 选择模板
	selector := NewTemplateSelector(templates)
	selectedTemplate := selector.SelectTemplate(in.ActionType)
	if selectedTemplate == nil {
		return nil, fmt.Errorf("未能选择合适的事件模板")
	}

	// 6. 使用 AI Service 生成事件内容
	// 创建一个新的 context，使用配置的超时时间（默认 30 秒）
	timeout := l.svcCtx.Config.AIGenerateTimeout
	if timeout == 0 {
		timeout = 30000 // 默认 30 秒
	}
	aiCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	aiResp, err := l.svcCtx.AIRpc.Generate(aiCtx, &ai.GenerateRequest{
		PromptTemplate: "event_generation",
		Variables: map[string]string{
			"avatar_name":       avatarInfo.Avatar.Nickname,
			"scene_name":        sceneInfo.Scene.Name,
			"scene_description": sceneInfo.Scene.Description,
			"action_type":       in.ActionType,
			"warmth":            fmt.Sprintf("%d", avatarInfo.Avatar.Personality.Warmth),
			"adventurous":       fmt.Sprintf("%d", avatarInfo.Avatar.Personality.Adventurous),
			"social":            fmt.Sprintf("%d", avatarInfo.Avatar.Personality.Social),
			"creative":          fmt.Sprintf("%d", avatarInfo.Avatar.Personality.Creative),
			"calm":              fmt.Sprintf("%d", avatarInfo.Avatar.Personality.Calm),
		},
		AvatarId: in.AvatarId,
	})
	if err != nil {
		l.Errorf("AI生成事件失败: %v", err)
		return nil, fmt.Errorf("AI生成事件失败")
	}

	eventTitle, eventText := parseAIResponse(aiResp.Content)

	// 7. 保存事件历史
	eventHistory := &model.EventHistory{
		AvatarId:    in.AvatarId,
		TemplateId:  selectedTemplate.Id,
		EventType:   in.ActionType,
		EventTitle:  eventTitle,
		EventText:   eventText,
		ImageUrl:    "", // MVP阶段暂不生成图片
		SceneId:     in.SceneId,
		SceneName:   sceneInfo.Scene.Name,
		OccurredAt:  time.Now(),
	}

	// PersonalityChanges 暂时为空，后续实现
	eventHistory.PersonalityChanges = sql.NullString{Valid: false}

	result, err := l.svcCtx.EventHistoryModel.Insert(eventHistory)
	if err != nil {
		l.Errorf("保存事件历史失败: %v", err)
		return nil, fmt.Errorf("保存事件历史失败")
	}

	eventId, err := result.LastInsertId()
	if err != nil {
		l.Errorf("获取事件ID失败: %v", err)
		return nil, fmt.Errorf("获取事件ID失败")
	}

	l.Infof("成功生成事件 (event_id=%d, avatar_id=%d, type=%s, scene=%s)",
		eventId, in.AvatarId, in.ActionType, sceneInfo.Scene.Name)

	// 8. 返回结果
	return &event.GenerateEventResponse{
		EventId:    eventId,
		EventType:  in.ActionType,
		EventTitle: eventTitle,
		EventText:  eventText,
		ImageUrl:   "", // MVP阶段暂不生成图片
	}, nil
}

// parseAIResponse 解析 AI 返回的文本，提取标题和内容
func parseAIResponse(text string) (title, content string) {
	lines := []rune(text)
	titleStart := -1
	contentStart := -1

	// 查找标题和内容
	for i := 0; i < len(lines)-3; i++ {
		if titleStart == -1 && i+3 < len(lines) {
			if string(lines[i:i+3]) == "标题：" || string(lines[i:i+6]) == "标题:" {
				titleStart = i + 3
				if string(lines[i:i+6]) == "标题:" {
					titleStart = i + 3
				}
			}
		}
		if contentStart == -1 && i+3 < len(lines) {
			if string(lines[i:i+3]) == "内容：" || string(lines[i:i+6]) == "内容:" {
				contentStart = i + 3
				if string(lines[i:i+6]) == "内容:" {
					contentStart = i + 3
				}
				break
			}
		}
	}

	// 提取标题
	if titleStart != -1 {
		end := titleStart
		for end < len(lines) && lines[end] != '\n' {
			end++
		}
		title = string(lines[titleStart:end])
	}

	// 提取内容
	if contentStart != -1 {
		content = string(lines[contentStart:])
	}

	// 如果解析失败，使用默认值
	if title == "" {
		title = "精彩的一天"
	}
	if content == "" {
		content = text
	}

	return title, content
}
