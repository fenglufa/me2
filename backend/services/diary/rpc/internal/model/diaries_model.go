package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ DiariesModel = (*customDiariesModel)(nil)

type (
	// DiariesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customDiariesModel.
	DiariesModel interface {
		diariesModel
		withSession(session sqlx.Session) DiariesModel
		FindByAvatarAndType(ctx context.Context, avatarId int64, tp string, page, pageSize int32, startDate, endDate string) ([]*Diaries, error)
		CountByAvatarAndType(ctx context.Context, avatarId int64, tp string, startDate, endDate string) (int64, error)
		GetDiaryStats(ctx context.Context, avatarId int64) (map[string]interface{}, error)
	}

	customDiariesModel struct {
		*defaultDiariesModel
	}
)

// NewDiariesModel returns a model for the database table.
func NewDiariesModel(conn sqlx.SqlConn) DiariesModel {
	return &customDiariesModel{
		defaultDiariesModel: newDiariesModel(conn),
	}
}

func (m *customDiariesModel) withSession(session sqlx.Session) DiariesModel {
	return NewDiariesModel(sqlx.NewSqlConnFromSession(session))
}
