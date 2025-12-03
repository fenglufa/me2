package svc

import (
	"github.com/me2/ai/rpc/internal/config"
	"github.com/me2/ai/rpc/internal/deepseek"
	"github.com/me2/ai/rpc/internal/model"
	"github.com/me2/ai/rpc/internal/prompt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config         config.Config
	DB             sqlx.SqlConn
	Redis          *redis.Redis
	DeepseekClient *deepseek.Client
	PromptRenderer *prompt.Renderer
	LogModel       *model.AiCallLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := sqlx.NewMysql(c.Mysql.DataSource)
	rds := redis.MustNewRedis(c.RedisConf)

	// 初始化 Prompt 渲染器
	renderer, err := prompt.NewRenderer(c.PromptFile)
	if err != nil {
		panic("初始化Prompt渲染器失败: " + err.Error())
	}

	return &ServiceContext{
		Config:         c,
		DB:             db,
		Redis:          rds,
		DeepseekClient: deepseek.NewClient(c.Deepseek),
		PromptRenderer: renderer,
		LogModel:       model.NewAiCallLogModel(db),
	}
}
