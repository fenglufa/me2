package logic

import (
	"context"

	"github.com/me2/diary/rpc/diary"
	"github.com/me2/diary/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarDiaryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAvatarDiaryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarDiaryListLogic {
	return &GetAvatarDiaryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分身日记列表
func (l *GetAvatarDiaryListLogic) GetAvatarDiaryList(in *diary.GetAvatarDiaryListRequest) (*diary.GetAvatarDiaryListResponse, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	diaries, err := l.svcCtx.DiaryModel.FindByAvatarAndType(l.ctx, in.AvatarId, "avatar", page, pageSize, in.StartDate, in.EndDate)
	if err != nil {
		l.Errorf("查询分身日记列表失败: %v", err)
		return nil, err
	}

	total, err := l.svcCtx.DiaryModel.CountByAvatarAndType(l.ctx, in.AvatarId, "avatar", in.StartDate, in.EndDate)
	if err != nil {
		l.Errorf("统计分身日记数量失败: %v", err)
		return nil, err
	}

	var diaryList []*diary.DiaryInfo
	for _, d := range diaries {
		diaryList = append(diaryList, convertToDiaryInfo(d))
	}

	return &diary.GetAvatarDiaryListResponse{
		Diaries: diaryList,
		Total:   total,
	}, nil
}
