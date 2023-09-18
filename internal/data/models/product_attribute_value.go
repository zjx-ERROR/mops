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

// ProductAttributeValue 产品属性
type ProductAttributeValue struct {
	gorm.Model
	CreateUserID int64             ``               //创建者
	UpdateUserID int64             ``               //最后更新者
	Name         string            `gorm:"size:50"` //属性值
	Attribute    *ProductAttribute ``               //属性
	Products     []*ProductProduct ``               //产品规格
}

type productAttributeValueRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewProductAttributeValueCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ProductAttributeValueRepo {
	if err := data.DB(ctx).AutoMigrate(&ProductAttributeValue{}); err != nil {
		log.Error(err)
	}
	return &productAttributeValueRepo{
		data: data,
		log:  log,
	}
}

func (repo *productAttributeValueRepo) Create(ctx context.Context, productAttributeValue *ProductAttributeValue) (uint, error) {
	err := repo.data.DB(ctx).Create(productAttributeValue).Error
	return productAttributeValue.ID, err
}

func (repo *productAttributeValueRepo) CreateBatch(ctx context.Context, productAttributeValue []*ProductAttributeValue) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(productAttributeValue, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *productAttributeValueRepo) Update(ctx context.Context, productAttributeValue *ProductAttributeValue) (uint, error) {
	err := repo.data.DB(ctx).Updates(productAttributeValue).Error
	return productAttributeValue.ID, err
}

func (repo *productAttributeValueRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &ProductAttributeValue{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *productAttributeValueRepo) Get(ctx context.Context, id uint) (*ProductAttributeValue, error) {
	city := &ProductAttributeValue{}
	city.ID = id
	err := repo.data.DB(ctx).Preload("Attribute").First(city).Error
	return city, err
}

func (repo *productAttributeValueRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []ProductAttributeValue, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []ProductAttributeValue
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
