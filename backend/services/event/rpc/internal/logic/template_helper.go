package logic

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/event/rpc/internal/model"
)

// TemplateSelector 模板选择器
type TemplateSelector struct {
	templates map[string][]*model.EventTemplate
}

// NewTemplateSelector 创建模板选择器
func NewTemplateSelector(templates []*model.EventTemplate) *TemplateSelector {
	selector := &TemplateSelector{
		templates: make(map[string][]*model.EventTemplate),
	}

	// 按分类分组
	for _, t := range templates {
		selector.templates[t.Category] = append(selector.templates[t.Category], t)
	}

	return selector
}

// SelectTemplate 选择模板
func (s *TemplateSelector) SelectTemplate(actionType string) *model.EventTemplate {
	templates, ok := s.templates[actionType]
	if !ok || len(templates) == 0 {
		return nil
	}

	// MVP阶段简单随机选择
	// 后续可以根据稀有度、冷却时间等进行加权选择
	rand.Seed(time.Now().UnixNano())
	return templates[rand.Intn(len(templates))]
}

// RenderTemplate 渲染模板（替换变量）
func RenderTemplate(template string, variables map[string]interface{}) string {
	result := template

	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}

	return result
}

// BuildTemplateVariables 构建模板变量
func BuildTemplateVariables(avatarInfo *avatar.AvatarInfo, sceneName string) map[string]interface{} {
	return map[string]interface{}{
		"avatar_name": avatarInfo.Nickname,
		"scene_name":  sceneName,
		"warmth":      avatarInfo.Personality.Warmth,
		"adventurous": avatarInfo.Personality.Adventurous,
		"social":      avatarInfo.Personality.Social,
		"creative":    avatarInfo.Personality.Creative,
		"calm":        avatarInfo.Personality.Calm,
		"energetic":   avatarInfo.Personality.Energetic,
	}
}

// GenerateEventText MVP阶段的简单文本生成
// 后续会调用 AI Service 生成更丰富的内容
func GenerateEventText(prompt string, avatarInfo *avatar.AvatarInfo, actionType, sceneName string) (title string, text string) {
	// MVP阶段使用模板直接生成
	// 这里模拟AI生成的效果

	// 生成标题
	title = generateEventTitle(actionType, sceneName)

	// 生成内容（使用模板提示词作为基础）
	text = generateEventContent(actionType, avatarInfo, sceneName)

	return title, text
}

// generateEventTitle 生成事件标题
func generateEventTitle(actionType, sceneName string) string {
	titles := map[string][]string{
		"exploration": {
			fmt.Sprintf("在%s的探索", sceneName),
			fmt.Sprintf("%s的发现之旅", sceneName),
			fmt.Sprintf("漫步%s", sceneName),
		},
		"social": {
			fmt.Sprintf("在%s的相遇", sceneName),
			fmt.Sprintf("%s的交流时光", sceneName),
			fmt.Sprintf("邂逅于%s", sceneName),
		},
		"study": {
			fmt.Sprintf("在%s的学习", sceneName),
			fmt.Sprintf("%s的思考时刻", sceneName),
			fmt.Sprintf("知识探索：%s", sceneName),
		},
		"creative": {
			fmt.Sprintf("在%s的创作", sceneName),
			fmt.Sprintf("%s的灵感时刻", sceneName),
			fmt.Sprintf("创意迸发于%s", sceneName),
		},
		"rest": {
			fmt.Sprintf("在%s的休憩", sceneName),
			fmt.Sprintf("%s的宁静时光", sceneName),
			fmt.Sprintf("静心于%s", sceneName),
		},
		"play": {
			fmt.Sprintf("在%s的欢乐", sceneName),
			fmt.Sprintf("%s的趣味时光", sceneName),
			fmt.Sprintf("玩耍于%s", sceneName),
		},
	}

	options, ok := titles[actionType]
	if !ok || len(options) == 0 {
		return fmt.Sprintf("在%s的时光", sceneName)
	}

	rand.Seed(time.Now().UnixNano())
	return options[rand.Intn(len(options))]
}

// generateEventContent 生成事件内容
func generateEventContent(actionType string, avatarInfo *avatar.AvatarInfo, sceneName string) string {
	// MVP阶段使用模板生成
	// 后续替换为 AI Service 调用

	templates := map[string][]string{
		"exploration": {
			fmt.Sprintf("%s来到了%s，四周的景色让ta感到新奇。随着脚步的深入，一些有趣的细节吸引了ta的注意。这次探索虽然短暂，却让ta对这个地方有了更深的了解。", avatarInfo.Nickname, sceneName),
			fmt.Sprintf("在%s漫步时，%s发现了一些平时容易被忽略的美好。阳光透过树叶洒下斑驳的光影，微风带来远处的声响，一切都显得那么和谐自然。", sceneName, avatarInfo.Nickname),
		},
		"social": {
			fmt.Sprintf("%s在%s遇到了一个有趣的灵魂。他们聊起各自的经历和想法，时间在不知不觉中流逝。这次相遇虽然简短，却给%s留下了深刻的印象。", avatarInfo.Nickname, sceneName, avatarInfo.Nickname),
			fmt.Sprintf("在%s，%s参与了一场温暖的交流。大家分享着彼此的故事，笑声此起彼伏。这种人与人之间真诚的连接，让ta感到格外珍贵。", sceneName, avatarInfo.Nickname),
		},
		"study": {
			fmt.Sprintf("%s在%s沉浸在学习之中。新的知识像拼图一样逐渐组合成完整的画面，那种豁然开朗的感觉让ta兴奋不已。学习的过程虽然需要专注，但收获的满足感让一切都值得。", avatarInfo.Nickname, sceneName),
			fmt.Sprintf("在%s，%s陷入了深度思考。一个个想法在脑海中碰撞、融合，最终形成了独特的见解。这种思维的跃升让ta对世界有了新的认识。", sceneName, avatarInfo.Nickname),
		},
		"creative": {
			fmt.Sprintf("%s在%s进行了一次创作。灵感如泉涌般不断涌现，ta的双手随着心中的想法舞动。最终完成的作品虽不完美，却充满了独特的个人风格。", avatarInfo.Nickname, sceneName),
			fmt.Sprintf("在%s，%s突然产生了一个绝妙的创意。这个想法新颖而大胆，让ta迫不及待想要将它实现。创造的激情在心中燃烧，照亮了前行的道路。", sceneName, avatarInfo.Nickname),
		},
		"rest": {
			fmt.Sprintf("%s在%s找到了一处安静的角落，放松身心。外界的喧嚣渐渐远去，内心归于平静。这段休息时光让ta重新充满了活力。", avatarInfo.Nickname, sceneName),
			fmt.Sprintf("在%s，%s进行了一次心灵的沉淀。通过冥想和放松，ta感受到内在的宁静与和谐。这种平和的状态让ta更加了解自己。", sceneName, avatarInfo.Nickname),
		},
		"play": {
			fmt.Sprintf("%s在%s尽情玩耍，享受这份简单的快乐。笑声回荡在空气中，烦恼和压力都被抛到了脑后。这样纯粹的欢乐时光，正是生活中最珍贵的礼物。", avatarInfo.Nickname, sceneName),
			fmt.Sprintf("在%s，%s发现了一些有趣的事物。好奇心驱使ta不断探索和尝试，每一次新的发现都带来惊喜和欢笑。玩耍的过程充满了乐趣和创意。", sceneName, avatarInfo.Nickname),
		},
	}

	options, ok := templates[actionType]
	if !ok || len(options) == 0 {
		return fmt.Sprintf("%s在%s度过了一段美好的时光。", avatarInfo.Nickname, sceneName)
	}

	rand.Seed(time.Now().UnixNano())
	return options[rand.Intn(len(options))]
}
