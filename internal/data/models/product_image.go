package models

import (
	"context"
	"mops/internal/biz"
	"mops/internal/data"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProductImage 产品图片
type ProductImage struct {
	gorm.Model
	CreateUserID    int64            ``                             //创建者
	UpdateUserID    int64            ``                             //最后更新者
	Name            string           `gorm:"unique" form:"name"`    //图片名称
	ProductTemplate *ProductTemplate ``                             //款式图片
	ProductProduct  *ProductProduct  ``                             //规格图片
	FormAction      string           `gorm:"-" json:"FormAction"`   //非数据库字段，用于表示记录的增加，修改
	ActionFields    []string         `gorm:"-" json:"ActionFields"` //需要操作的字段,用于update时
}

type productImageRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewProductImageCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ProductImageRepo {
	if err := data.DB(ctx).AutoMigrate(&ProductImage{}); err != nil {
		log.Error(err)
	}
	return &productImageRepo{
		data: data,
		log:  log,
	}
}
