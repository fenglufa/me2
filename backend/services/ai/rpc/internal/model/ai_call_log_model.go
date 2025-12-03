package model

import (
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type AiCallLog struct {
	Id           int64     `db:"id"`
	ServiceName  string    `db:"service_name"`
	SceneType    string    `db:"scene_type"`
	UserId       int64     `db:"user_id"`
	AvatarId     int64     `db:"avatar_id"`
	InputTokens  int32     `db:"input_tokens"`
	OutputTokens int32     `db:"output_tokens"`
	Cost         int64     `db:"cost"`
	DurationMs   int64     `db:"duration_ms"`
	Status       string    `db:"status"`
	ErrorMessage string    `db:"error_message"`
	CreatedAt    time.Time `db:"created_at"`
}

type AiCallLogModel struct {
	conn sqlx.SqlConn
}

func NewAiCallLogModel(conn sqlx.SqlConn) *AiCallLogModel {
	return &AiCallLogModel{conn: conn}
}

func (m *AiCallLogModel) Insert(log *AiCallLog) (sql.Result, error) {
	query := `INSERT INTO ai_call_logs (service_name, scene_type, user_id, avatar_id, 
		input_tokens, output_tokens, cost, duration_ms, status, error_message)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return m.conn.Exec(query, log.ServiceName, log.SceneType, log.UserId, log.AvatarId,
		log.InputTokens, log.OutputTokens, log.Cost, log.DurationMs, log.Status, log.ErrorMessage)
}

func (m *AiCallLogModel) GetStats(userId int64, startDate, endDate string) (*Stats, error) {
	query := `SELECT 
		COUNT(*) as total_calls,
		SUM(input_tokens) as total_input_tokens,
		SUM(output_tokens) as total_output_tokens,
		SUM(cost) as total_cost,
		SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as success_calls,
		SUM(CASE WHEN status = 'error' THEN 1 ELSE 0 END) as error_calls,
		AVG(duration_ms) as avg_duration_ms
		FROM ai_call_logs
		WHERE (? = 0 OR user_id = ?)
		AND DATE(created_at) >= ? AND DATE(created_at) <= ?`

	var stats Stats
	err := m.conn.QueryRow(&stats, query, userId, userId, startDate, endDate)
	return &stats, err
}

type Stats struct {
	TotalCalls        int64   `db:"total_calls"`
	TotalInputTokens  int64   `db:"total_input_tokens"`
	TotalOutputTokens int64   `db:"total_output_tokens"`
	TotalCost         int64   `db:"total_cost"`
	SuccessCalls      int64   `db:"success_calls"`
	ErrorCalls        int64   `db:"error_calls"`
	AvgDurationMs     float64 `db:"avg_duration_ms"`
}
