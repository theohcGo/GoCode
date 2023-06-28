package rotation

import (
	"context"

	"goframe-shop/internal/service"
	"github.com/gogf/gf/v2/encoding/ghtml"
	"goframe-shop/internal/dao"
	"goframe-shop/internal/model"
)

type sRotation struct{}

func init() {
	// 在logic模块注册轮播图服务
	service.RegisterRotation(New())
}

func New() *sRotation {
	return &sRotation{}
}

// Create 创建内容
func (s *sRotation) Create(ctx context.Context, in model.RotationCreateInput) (out model.RotationCreateOutput, err error) {
	// 不允许HTML代码
	if err = ghtml.SpecialCharsMapOrStruct(in); err != nil {
		return out, err
	}
	lastInsertID, err := dao.RotationInfo.Ctx(ctx).Data(in).InsertAndGetId()
	if err != nil {
		return out, err
	}
	return model.RotationCreateOutput{RotationId: int(lastInsertID)}, err
}