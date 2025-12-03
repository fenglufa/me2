package logic

import (
	"context"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/avatar/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAvatarInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarInfoLogic {
	return &GetAvatarInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAvatarInfoLogic) GetAvatarInfo(in *avatar.GetAvatarInfoRequest) (*avatar.GetAvatarInfoResponse, error) {
	av, err := l.svcCtx.AvatarModel.FindByAvatarId(in.AvatarId)
	if err != nil {
		l.Errorf("查询分身失败: %v", err)
		return nil, err
	}

	return &avatar.GetAvatarInfoResponse{
		Avatar: &avatar.AvatarInfo{
			AvatarId:      av.AvatarId,
			UserId:        av.UserId,
			Nickname:      av.Nickname,
			AvatarUrl:     av.AvatarUrl,
			Gender:        av.Gender,
			BirthDate:     av.BirthDate,
			Occupation:    av.Occupation,
			MaritalStatus: av.MaritalStatus,
			Personality: &avatar.PersonalityInfo{
				Warmth:      av.Warmth,
				Adventurous: av.Adventurous,
				Social:      av.Social,
				Creative:    av.Creative,
				Calm:        av.Calm,
				Energetic:   av.Energetic,
			},
			CreatedAt: av.CreatedAt.Unix(),
			UpdatedAt: av.UpdatedAt.Unix(),
		},
	}, nil
}
