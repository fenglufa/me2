package logic

import (
	"context"

	"github.com/me2/event/rpc/event"
	"github.com/me2/event/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTemplatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTemplatesLogic {
	return &GetTemplatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetTemplates 获取模板列表（后续可用于管理后台）
func (l *GetTemplatesLogic) GetTemplates(in *event.GetTemplatesRequest) (*event.GetTemplatesResponse, error) {
	// 查询模板列表
	templates, err := l.svcCtx.EventTemplateModel.FindAll(in.Category, in.Rarity)
	if err != nil {
		l.Errorf("查询模板失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	templateInfos := make([]*event.TemplateInfo, 0, len(templates))
	for _, t := range templates {
		desc := ""
		if t.Description.Valid {
			desc = t.Description.String
		}
		templateInfos = append(templateInfos, &event.TemplateInfo{
			TemplateId:  t.Id,
			Category:    t.Category,
			Name:        t.Name,
			Description: desc,
			Rarity:      t.Rarity,
		})
	}

	return &event.GetTemplatesResponse{
		Templates: templateInfos,
	}, nil
}
