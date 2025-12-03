package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/me2/ai/rpc/ai"
	"github.com/me2/ai/rpc/internal/deepseek"
	"github.com/me2/ai/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CalculateActionIntentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCalculateActionIntentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CalculateActionIntentLogic {
	return &CalculateActionIntentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 行动意图计算（供 Action Service 调用）
func (l *CalculateActionIntentLogic) CalculateActionIntent(in *ai.ActionIntentRequest) (*ai.ActionIntentResponse, error) {
	startTime := time.Now()

	// 构建模板变量
	variables := l.buildTemplateVariables(in)

	// 渲染 Prompt
	systemPrompt, userPrompt, err := l.svcCtx.PromptRenderer.Render("action_intent", variables)
	if err != nil {
		return nil, fmt.Errorf("渲染Prompt失败: %w", err)
	}

	// 构建 Deepseek 请求
	req := &deepseek.ChatRequest{
		Model: l.svcCtx.Config.Deepseek.Model,
		Messages: []deepseek.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.7,
		MaxTokens:   500,
	}

	// 调用 Deepseek
	resp, err := l.svcCtx.DeepseekClient.Chat(l.ctx, req)
	if err != nil {
		return nil, fmt.Errorf("调用Deepseek失败: %w", err)
	}

	// 解析 JSON 响应
	result, err := l.parseActionIntentResponse(resp.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	duration := time.Since(startTime)

	// 记录日志
	logx.Infof("行动意图计算完成: avatar_id=%d, action=%s, duration=%dms",
		in.AvatarId, result.RecommendedAction, duration.Milliseconds())

	return &ai.ActionIntentResponse{
		ActionScores:      result.ActionScores,
		RecommendedAction: result.RecommendedAction,
		Reason:            result.Reason,
		InputTokens:       resp.Usage.PromptTokens,
		OutputTokens:      resp.Usage.CompletionTokens,
	}, nil
}

// 构建模板变量
func (l *CalculateActionIntentLogic) buildTemplateVariables(in *ai.ActionIntentRequest) map[string]string {
	variables := make(map[string]string)

	// 人格向量
	if in.Personality != nil {
		variables["warmth"] = strconv.Itoa(int(in.Personality.Warmth))
		variables["adventurous"] = strconv.Itoa(int(in.Personality.Adventurous))
		variables["social"] = strconv.Itoa(int(in.Personality.Social))
		variables["creative"] = strconv.Itoa(int(in.Personality.Creative))
		variables["calm"] = strconv.Itoa(int(in.Personality.Calm))
		variables["energetic"] = strconv.Itoa(int(in.Personality.Energetic))
	}

	// 时间上下文
	if in.TimeContext != nil {
		variables["current_time"] = in.TimeContext.CurrentTime
		variables["time_period"] = in.TimeContext.TimePeriod
		variables["day_of_week"] = in.TimeContext.DayOfWeek
	}

	// 当前状态
	if in.CurrentState != nil {
		variables["energy"] = strconv.Itoa(int(in.CurrentState.Energy))
		variables["emotion_state"] = in.CurrentState.EmotionState
		variables["mood"] = in.CurrentState.Mood
	}

	// 最近用户互动
	if len(in.RecentInteractions) > 0 {
		variables["recent_interactions"] = strings.Join(in.RecentInteractions, "\n- ")
	} else {
		variables["recent_interactions"] = "暂无用户互动"
	}

	// 最近事件
	if len(in.RecentEvents) > 0 {
		variables["recent_events"] = strings.Join(in.RecentEvents, "\n- ")
	} else {
		variables["recent_events"] = "暂无最近事件"
	}

	return variables
}

// ActionIntentResult 行动意图结果
type ActionIntentResult struct {
	ActionScores      map[string]float32 `json:"action_scores"`
	RecommendedAction string             `json:"recommended_action"`
	Reason            string             `json:"reason"`
}

// 解析行动意图响应
func (l *CalculateActionIntentLogic) parseActionIntentResponse(content string) (*ActionIntentResult, error) {
	// 去除可能的 markdown 代码块标记
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var result ActionIntentResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w, 原始内容: %s", err, content)
	}

	// 验证必需字段
	if len(result.ActionScores) == 0 {
		return nil, fmt.Errorf("action_scores 为空")
	}
	if result.RecommendedAction == "" {
		return nil, fmt.Errorf("recommended_action 为空")
	}

	return &result, nil
}

