package diary

import (
	"context"

	"github.com/me2/diary/rpc/diary"
	"github.com/me2/gateway/api/internal/svc"
	"github.com/me2/gateway/api/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserDiaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserDiaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserDiaryLogic {
	return &CreateUserDiaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserDiaryLogic) CreateUserDiary(req *types.CreateUserDiaryRequest) (resp *types.DiaryResponse, err error) {
	userID := l.ctx.Value("user_id").(int64)

	rpcResp, err := l.svcCtx.DiaryRpc.CreateUserDiary(l.ctx, &diary.CreateUserDiaryRequest{
		AvatarId: userID,
		Content:  req.Content,
	})
	if err != nil {
		return nil, err
	}

	return &types.DiaryResponse{
		Id:      rpcResp.DiaryId,
		Type:    "user",
		Content: req.Content,
		Mood:    req.Mood,
	}, nil
}
