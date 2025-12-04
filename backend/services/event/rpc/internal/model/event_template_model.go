package model

import (
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// EventTemplate 事件模板
type EventTemplate struct {
	Id                 int64          `db:"id"`
	Category           string         `db:"category"`
	Name               string         `db:"name"`
	Description        sql.NullString `db:"description"`
	TriggerConditions  sql.NullString `db:"trigger_conditions"` // JSON
	Rarity             string         `db:"rarity"`
	CooldownHours      int            `db:"cooldown_hours"`
	ContentTemplate    string         `db:"content_template"`
	PersonalityImpact  sql.NullString `db:"personality_impact"` // JSON
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
}

// EventTemplateModel 事件模板模型
type EventTemplateModel struct {
	conn sqlx.SqlConn
}

// NewEventTemplateModel 创建事件模板模型
func NewEventTemplateModel(conn sqlx.SqlConn) *EventTemplateModel {
	return &EventTemplateModel{
		conn: conn,
	}
}

// FindByCategory 根据分类查找模板
func (m *EventTemplateModel) FindByCategory(category string) ([]*EventTemplate, error) {
	query := `SELECT id, category, name, description, trigger_conditions, rarity,
	          cooldown_hours, content_template, personality_impact, created_at, updated_at
	          FROM event_templates WHERE category = ? ORDER BY rarity, id`

	var templates []*EventTemplate
	err := m.conn.QueryRows(&templates, query, category)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// FindByCategoryAndRarity 根据分类和稀有度查找模板
func (m *EventTemplateModel) FindByCategoryAndRarity(category, rarity string) ([]*EventTemplate, error) {
	query := `SELECT id, category, name, description, trigger_conditions, rarity,
	          cooldown_hours, content_template, personality_impact, created_at, updated_at
	          FROM event_templates WHERE category = ? AND rarity = ? ORDER BY id`

	var templates []*EventTemplate
	err := m.conn.QueryRows(&templates, query, category, rarity)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// FindAll 查找所有模板（可选过滤）
func (m *EventTemplateModel) FindAll(category, rarity string) ([]*EventTemplate, error) {
	var query string
	var args []interface{}

	if category != "" && rarity != "" {
		query = `SELECT id, category, name, description, trigger_conditions, rarity,
		         cooldown_hours, content_template, personality_impact, created_at, updated_at
		         FROM event_templates WHERE category = ? AND rarity = ? ORDER BY category, rarity, id`
		args = []interface{}{category, rarity}
	} else if category != "" {
		query = `SELECT id, category, name, description, trigger_conditions, rarity,
		         cooldown_hours, content_template, personality_impact, created_at, updated_at
		         FROM event_templates WHERE category = ? ORDER BY category, rarity, id`
		args = []interface{}{category}
	} else if rarity != "" {
		query = `SELECT id, category, name, description, trigger_conditions, rarity,
		         cooldown_hours, content_template, personality_impact, created_at, updated_at
		         FROM event_templates WHERE rarity = ? ORDER BY category, rarity, id`
		args = []interface{}{rarity}
	} else {
		query = `SELECT id, category, name, description, trigger_conditions, rarity,
		         cooldown_hours, content_template, personality_impact, created_at, updated_at
		         FROM event_templates ORDER BY category, rarity, id`
	}

	var templates []*EventTemplate
	err := m.conn.QueryRows(&templates, query, args...)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// FindOne 查找单个模板
func (m *EventTemplateModel) FindOne(id int64) (*EventTemplate, error) {
	query := `SELECT id, category, name, description, trigger_conditions, rarity,
	          cooldown_hours, content_template, personality_impact, created_at, updated_at
	          FROM event_templates WHERE id = ?`

	var template EventTemplate
	err := m.conn.QueryRow(&template, query, id)
	if err != nil {
		return nil, err
	}

	return &template, nil
}
