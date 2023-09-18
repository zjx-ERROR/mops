package models

import (
	"context"
	"mops/internal/biz"
	"mops/internal/data"
	"mops/pkg"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProductAttribute 产品属性
type ProductAttribute struct {
	gorm.Model
	CreateUserID int64                    ``                    //创建者
	UpdateUserID int64                    ``                    //最后更新者
	Name         string                   `gorm:"size:50"`      //属性名称
	Code         string                   ``                    //属性编码
	ValueIds     []*ProductAttributeValue ``                    //属性值
	CreatVariant bool                     `gorm:"default:true"` //创建规格
}

type productAttributeRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewProductAttributeCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ProductAttributeRepo {
	if err := data.DB(ctx).AutoMigrate(&ProductAttribute{}); err != nil {
		log.Error(err)
	}
	return &productAttributeRepo{
		data: data,
		log:  log,
	}
}

func (repo *productAttributeRepo) Create(ctx context.Context, productAttribute *ProductAttribute) (uint, error) {
	err := repo.data.DB(ctx).Create(productAttribute).Error
	return productAttribute.ID, err
}

func (repo *productAttributeRepo) CreateBatch(ctx context.Context, productAttribute []*ProductAttribute) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(productAttribute, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *productAttributeRepo) Update(ctx context.Context, productAttribute *ProductAttribute) (uint, error) {
	err := repo.data.DB(ctx).Updates(productAttribute).Error
	return productAttribute.ID, err
}

func (repo *productAttributeRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &ProductAttribute{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *productAttributeRepo) Get(ctx context.Context, id uint) (*ProductAttribute, error) {
	city := &ProductAttribute{}
	city.ID = id
	err := repo.data.DB(ctx).Preload("ValueIds").First(city).Error
	return city, err
}

func (repo *productAttributeRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []ProductAttribute, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []ProductAttribute
		err       error
		paginator pkg.Paginator
		cnt       int64
		offset    int64 = (page - 1) * limit
	)
	tx := repo.data.DB(ctx).Preload(clause.Associations)
	if _, ok := condMap["and"]; ok {
		for k, v := range condMap["and"] {
			tx = tx.Where(k, v)
		}
	}
	if _, ok := condMap["or"]; ok {
		for k, v := range condMap["or"] {
			tx = tx.Or(k, v)
		}
	}
	for k, v := range exclude {
		tx = tx.Not(k, v)
	}
	if orderBy != "" {
		tx = tx.Order(orderBy)
	}
	tx.Count(&cnt)
	paginator = pkg.GenPaginator(limit, offset, cnt)
	err = tx.Limit(int(limit)).Offset(int(offset)).Find(&objArrs).Error
	return paginator, objArrs, err
}
