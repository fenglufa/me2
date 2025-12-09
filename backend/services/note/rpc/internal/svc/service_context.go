package svc

import (
	"github.com/me2/ai/rpc/ai_client"
	"github.com/me2/note/rpc/internal/config"
	"github.com/me2/note/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config           config.Config
	NoteModel        model.NoteModel
	NoteTodoModel    model.NoteTodoModel
	NoteExpenseModel model.NoteExpenseModel
	AiRpc            ai_client.Ai
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:           c,
		NoteModel:        model.NewNoteModel(conn),
		NoteTodoModel:    model.NewNoteTodoModel(conn),
		NoteExpenseModel: model.NewNoteExpenseModel(conn),
		AiRpc:            ai_client.NewAi(zrpc.MustNewClient(c.AiRpc)),
	}
}
