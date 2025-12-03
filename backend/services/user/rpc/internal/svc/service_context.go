package svc

import (
	"github.com/me2/user/rpc/internal/config"
	"github.com/me2/user/rpc/internal/idgen"
	"github.com/me2/user/rpc/internal/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/me2/sms/rpc/sms_client"
	"github.com/me2/oss/rpc/oss_client"
)

type ServiceContext struct {
	Config    config.Config
	DB        sqlx.SqlConn
	SmsRpc    sms_client.Sms
	OssRpc    oss_client.Oss
	IDGen     *idgen.Snowflake
	UserModel *model.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	db := sqlx.NewMysql(c.Mysql.DataSource)

	// 初始化 SMS RPC 客户端
	smsRpc := sms_client.NewSms(zrpc.MustNewClient(c.SmsRpc))

	// 初始化 OSS RPC 客户端
	ossRpc := oss_client.NewOss(zrpc.MustNewClient(c.OssRpc))

	// 初始化雪花算法 ID 生成器
	idGen, err := idgen.NewSnowflake(c.MachineID)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:    c,
		DB:        db,
		SmsRpc:    smsRpc,
		OssRpc:    ossRpc,
		IDGen:     idGen,
		UserModel: model.NewUserModel(db),
	}
}
