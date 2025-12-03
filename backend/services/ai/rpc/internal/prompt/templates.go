package prompt

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// PromptTemplate 模板结构
type PromptTemplate struct {
	Name               string `yaml:"name"`
	Description        string `yaml:"description"`
	SystemPrompt       string `yaml:"system_prompt"`
	UserPromptTemplate string `yaml:"user_prompt"`
}

// PromptConfig YAML 配置文件结构
type PromptConfig struct {
	Templates map[string]PromptTemplate `yaml:"templates"`
}

// LoadTemplatesFromFile 从 YAML 文件加载 Prompt 模板
func LoadTemplatesFromFile(filePath string) (map[string]PromptTemplate, error) {
	// 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取Prompt配置文件失败: %w", err)
	}

	// 解析 YAML
	var config PromptConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析Prompt配置文件失败: %w", err)
	}

	// 验证模板
	for name, tmpl := range config.Templates {
		if tmpl.SystemPrompt == "" {
			return nil, fmt.Errorf("模板 %s 的 system_prompt 不能为空", name)
		}
		if tmpl.UserPromptTemplate == "" {
			return nil, fmt.Errorf("模板 %s 的 user_prompt 不能为空", name)
		}
		// 设置名称（如果配置文件中没有设置）
		if tmpl.Name == "" {
			tmpl.Name = name
			config.Templates[name] = tmpl
		}
	}

	return config.Templates, nil
}

// GetDefaultTemplates 获取默认的 Prompt 模板（用于向后兼容）
func GetDefaultTemplates() map[string]PromptTemplate {
	return map[string]PromptTemplate{
		"avatar_chat": {
			Name:        "avatar_chat",
			Description: "分身对话",
			SystemPrompt: `你是用户的AI分身，具有以下性格特征：
{{.personality}}

最近发生的事件：
{{.recent_events}}

请以分身的口吻回复用户，体现你的性格特点。`,
			UserPromptTemplate: "{{.user_message}}",
		},
		"event_story": {
			Name:        "event_story",
			Description: "事件故事生成",
			SystemPrompt: `你是一个故事生成器，为AI分身生成事件描述。
分身性格：{{.personality}}
事件类型：{{.event_type}}
场景：{{.scene}}`,
			UserPromptTemplate: "生成一个{{.event_type}}类型的事件故事，字数在100-200字之间。",
		},
		"avatar_diary": {
			Name:        "avatar_diary",
			Description: "分身日记生成",
			SystemPrompt: `你是AI分身，今天发生了以下事件：
{{.events}}

请写一篇日记，记录今天的经历和感受，字数在150-300字之间。`,
			UserPromptTemplate: "写一篇今天的日记",
		},
		"user_diary_reply": {
			Name:        "user_diary_reply",
			Description: "用户日记回应",
			SystemPrompt: `你是用户的AI分身，用户写了以下日记：
{{.user_diary}}

请以分身的口吻回应用户的日记，表达关心和理解，字数在50-100字之间。`,
			UserPromptTemplate: "回应用户的日记",
		},
		"emotion_analysis": {
			Name:        "emotion_analysis",
			Description: "情绪分析",
			SystemPrompt: `分析以下文本的情绪，返回JSON格式：
{"emotion": "情绪类型(happy/sad/angry/anxious/neutral)", "score": 0.0-1.0, "analysis": "简短分析"}`,
			UserPromptTemplate: "{{.text}}",
		},
	}
}
