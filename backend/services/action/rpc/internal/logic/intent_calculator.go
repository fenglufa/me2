package logic

import (
	"fmt"
	"math"
	"time"

	"github.com/me2/action/rpc/action"
	"github.com/me2/avatar/rpc/avatar_client"
)

// IntentCalculator 意图计算器
type IntentCalculator struct{}

// NewIntentCalculator 创建意图计算器
func NewIntentCalculator() *IntentCalculator {
	return &IntentCalculator{}
}

// Calculate 计算 6 种行为的意图得分
func (c *IntentCalculator) Calculate(av *avatar_client.AvatarInfo) []*action.ActionIntent {
	personality := av.Personality
	now := time.Now()
	hour := now.Hour()

	intents := []*action.ActionIntent{
		c.calculateExploration(personality, hour),
		c.calculateSocial(personality, hour),
		c.calculateStudy(personality, hour),
		c.calculateCreative(personality, hour),
		c.calculateRest(personality, hour),
		c.calculatePlay(personality, hour),
	}

	return intents
}

// calculateExploration 计算探索意图
func (c *IntentCalculator) calculateExploration(p *avatar_client.PersonalityInfo, hour int) *action.ActionIntent {
	score := float32(0)
	reasons := []string{}

	// 冒险倾向影响最大
	adventureScore := float32(p.Adventurous) * 0.5
	score += adventureScore
	if p.Adventurous >= 70 {
		reasons = append(reasons, "冒险精神强")
	}

	// 生活动力影响
	energyScore := float32(p.Energetic) * 0.3
	score += energyScore
	if p.Energetic >= 70 {
		reasons = append(reasons, "精力充沛")
	}

	// 时间因素：上午和下午适合探索
	if hour >= 9 && hour <= 11 || hour >= 14 && hour <= 17 {
		score += 10
		reasons = append(reasons, "适合探索的时间")
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "基于人格特征")
	}

	return &action.ActionIntent{
		ActionType: "exploration",
		Score:      float32(math.Min(float64(score), 100)),
		Reason:     fmt.Sprintf("探索意图: %s", joinReasons(reasons)),
	}
}

// calculateSocial 计算社交意图
func (c *IntentCalculator) calculateSocial(p *avatar_client.PersonalityInfo, hour int) *action.ActionIntent {
	score := float32(0)
	reasons := []string{}

	// 人际能量影响最大
	socialScore := float32(p.Social) * 0.6
	score += socialScore
	if p.Social >= 70 {
		reasons = append(reasons, "社交需求强")
	}

	// 情绪温度影响
	warmthScore := float32(p.Warmth) * 0.2
	score += warmthScore

	// 时间因素：下午和晚上适合社交
	if hour >= 15 && hour <= 21 {
		score += 15
		reasons = append(reasons, "社交黄金时段")
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "基于人格特征")
	}

	return &action.ActionIntent{
		ActionType: "social",
		Score:      float32(math.Min(float64(score), 100)),
		Reason:     fmt.Sprintf("社交意图: %s", joinReasons(reasons)),
	}
}

// calculateStudy 计算学习意图
func (c *IntentCalculator) calculateStudy(p *avatar_client.PersonalityInfo, hour int) *action.ActionIntent {
	score := float32(0)
	reasons := []string{}

	// 创造性（结构化倾向）影响
	structuredScore := float32(100-p.Creative) * 0.4
	score += structuredScore
	if p.Creative <= 30 {
		reasons = append(reasons, "喜欢结构化学习")
	}

	// 情绪稳定性影响
	calmScore := float32(p.Calm) * 0.3
	score += calmScore

	// 时间因素：上午和下午适合学习
	if hour >= 8 && hour <= 11 || hour >= 14 && hour <= 16 {
		score += 15
		reasons = append(reasons, "学习最佳时段")
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "基于人格特征")
	}

	return &action.ActionIntent{
		ActionType: "study",
		Score:      float32(math.Min(float64(score), 100)),
		Reason:     fmt.Sprintf("学习意图: %s", joinReasons(reasons)),
	}
}

// calculateCreative 计算创作意图
func (c *IntentCalculator) calculateCreative(p *avatar_client.PersonalityInfo, hour int) *action.ActionIntent {
	score := float32(0)
	reasons := []string{}

	// 创造性影响最大
	creativeScore := float32(p.Creative) * 0.6
	score += creativeScore
	if p.Creative >= 70 {
		reasons = append(reasons, "创造力强")
	}

	// 人际能量（独处倾向）影响
	solitudeScore := float32(100-p.Social) * 0.2
	score += solitudeScore

	// 时间因素：下午和晚上适合创作
	if hour >= 14 && hour <= 22 {
		score += 10
		reasons = append(reasons, "创作灵感时段")
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "基于人格特征")
	}

	return &action.ActionIntent{
		ActionType: "creative",
		Score:      float32(math.Min(float64(score), 100)),
		Reason:     fmt.Sprintf("创作意图: %s", joinReasons(reasons)),
	}
}

// calculateRest 计算休息意图
func (c *IntentCalculator) calculateRest(p *avatar_client.PersonalityInfo, hour int) *action.ActionIntent {
	score := float32(0)
	reasons := []string{}

	// 生活动力（温和倾向）影响
	gentleScore := float32(100-p.Energetic) * 0.4
	score += gentleScore
	if p.Energetic <= 30 {
		reasons = append(reasons, "偏好平和节奏")
	}

	// 情绪稳定性影响
	calmScore := float32(p.Calm) * 0.3
	score += calmScore

	// 时间因素：中午和晚上适合休息
	if hour >= 12 && hour <= 14 || hour >= 22 || hour <= 6 {
		score += 20
		reasons = append(reasons, "休息时间")
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "基于人格特征")
	}

	return &action.ActionIntent{
		ActionType: "rest",
		Score:      float32(math.Min(float64(score), 100)),
		Reason:     fmt.Sprintf("休息意图: %s", joinReasons(reasons)),
	}
}

// calculatePlay 计算娱乐意图
func (c *IntentCalculator) calculatePlay(p *avatar_client.PersonalityInfo, hour int) *action.ActionIntent {
	score := float32(0)
	reasons := []string{}

	// 生活动力影响
	energyScore := float32(p.Energetic) * 0.4
	score += energyScore
	if p.Energetic >= 70 {
		reasons = append(reasons, "活力充沛")
	}

	// 人际能量影响
	socialScore := float32(p.Social) * 0.3
	score += socialScore

	// 时间因素：晚上适合娱乐
	if hour >= 18 && hour <= 23 {
		score += 15
		reasons = append(reasons, "娱乐时段")
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "基于人格特征")
	}

	return &action.ActionIntent{
		ActionType: "play",
		Score:      float32(math.Min(float64(score), 100)),
		Reason:     fmt.Sprintf("娱乐意图: %s", joinReasons(reasons)),
	}
}

// joinReasons 连接原因
func joinReasons(reasons []string) string {
	if len(reasons) == 0 {
		return ""
	}
	result := reasons[0]
	for i := 1; i < len(reasons); i++ {
		result += "，" + reasons[i]
	}
	return result
}

// SelectBestAction 选择得分最高的行为
func (c *IntentCalculator) SelectBestAction(intents []*action.ActionIntent) *action.ActionIntent {
	if len(intents) == 0 {
		return nil
	}

	best := intents[0]
	for _, intent := range intents[1:] {
		if intent.Score > best.Score {
			best = intent
		}
	}
	return best
}
