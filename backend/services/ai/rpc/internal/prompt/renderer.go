package prompt

import (
	"bytes"
	"fmt"
	"text/template"
)

type Renderer struct {
	templates map[string]PromptTemplate
}

// NewRenderer 创建一个新的 Renderer，使用指定文件的模板
func NewRenderer(promptFile string) (*Renderer, error) {
	var templates map[string]PromptTemplate
	var err error

	// 如果指定了配置文件，从文件加载
	if promptFile != "" {
		templates, err = LoadTemplatesFromFile(promptFile)
		if err != nil {
			return nil, fmt.Errorf("加载Prompt配置失败: %w", err)
		}
	} else {
		// 否则使用默认模板
		templates = GetDefaultTemplates()
	}

	return &Renderer{
		templates: templates,
	}, nil
}

// Render 渲染 Prompt 模板
func (r *Renderer) Render(templateName string, variables map[string]string) (string, string, error) {
	tmpl, ok := r.templates[templateName]
	if !ok {
		return "", "", fmt.Errorf("模板不存在: %s", templateName)
	}

	systemPrompt, err := r.renderTemplate(tmpl.SystemPrompt, variables)
	if err != nil {
		return "", "", fmt.Errorf("渲染系统提示词失败: %w", err)
	}

	userPrompt, err := r.renderTemplate(tmpl.UserPromptTemplate, variables)
	if err != nil {
		return "", "", fmt.Errorf("渲染用户提示词失败: %w", err)
	}

	return systemPrompt, userPrompt, nil
}

func (r *Renderer) renderTemplate(tmplStr string, variables map[string]string) (string, error) {
	tmpl, err := template.New("prompt").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, variables); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GetTemplates 获取所有加载的模板（用于调试）
func (r *Renderer) GetTemplates() map[string]PromptTemplate {
	return r.templates
}
