package logic

import (
	"context"

	"github.com/me2/diary/rpc/diary"
	"github.com/me2/diary/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserDiaryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserDiaryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDiaryListLogic {
	return &GetUserDiaryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户日记列表（用户自己写的日记）
func (l *GetUserDiaryListLogic) GetUserDiaryList(in *diary.GetUserDiaryListRequest) (*diary.GetUserDiaryListResponse, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	// 使用 user_id 查询用户日记
	diaries, err := l.svcCtx.DiaryModel.FindByUserAndType(l.ctx, in.UserId, "user", page, pageSize, in.StartDate, in.EndDate)
	if err != nil {
		l.Errorf("查询用户日记列表失败: %v", err)
		return nil, err
	}

	total, err := l.svcCtx.DiaryModel.CountByUserAndType(l.ctx, in.UserId, "user", in.StartDate, in.EndDate)
	if err != nil {
		l.Errorf("统计用户日记数量失败: %v", err)
		return nil, err
	}

	var diaryList []*diary.DiaryInfo
	for _, d := range diaries {
		diaryList = append(diaryList, convertToDiaryInfo(d))
	}

	return &diary.GetUserDiaryListResponse{
		Diaries: diaryList,
		Total:   total,
	}, nil
}
