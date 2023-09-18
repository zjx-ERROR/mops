package models

import (
	"context"
	"mops/internal/biz"
	"mops/internal/data"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProductTag  产品标签
type ProductTag struct {
	gorm.Model
	CreateUserID int64             ``                      //创建者
	UpdateUserID int64             ``                      //最后更新者
	Name         string            `gorm:"size:20;unique"` //产品标签名称
	Type         string            `gorm:"size:20"`        //标签类型
	Products     []*ProductProduct ``                      //产品规格
}

type productTagRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewProductTagCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ProductTagRepo {
	if err := data.DB(ctx).AutoMigrate(&ProductTag{}); err != nil {
		log.Error(err)
	}
	return &productTagRepo{
		data: data,
		log:  log,
	}
}
