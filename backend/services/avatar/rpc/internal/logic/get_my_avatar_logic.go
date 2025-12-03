package logic

import (
	"context"
	"database/sql"

	"github.com/me2/avatar/rpc/avatar"
	"github.com/me2/avatar/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyAvatarLogic {
	return &GetMyAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMyAvatarLogic) GetMyAvatar(in *avatar.GetMyAvatarRequest) (*avatar.GetMyAvatarResponse, error) {
	av, err := l.svcCtx.AvatarModel.FindByUserId(in.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return &avatar.GetMyAvatarResponse{HasAvatar: false}, nil
		}
		l.Errorf("查询分身失败: %v", err)
		return &avatar.GetMyAvatarResponse{HasAvatar: false}, nil
	}

	return &avatar.GetMyAvatarResponse{
		HasAvatar: true,
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
