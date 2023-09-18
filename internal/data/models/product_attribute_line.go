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

// ProductAttributeLine 产品属性明细
type ProductAttributeLine struct {
	gorm.Model
	CreateUserID    int64                    `` //创建者
	UpdateUserID    int64                    `` //最后更新者
	Attribute       *ProductAttribute        `` //属性
	ProductTemplate *ProductTemplate         `` //产品模版
	AttributeValues []*ProductAttributeValue `` //属性值
}

type productAttributeLineRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewProductAttributeLineCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ProductAttributeLineRepo {
	if err := data.DB(ctx).AutoMigrate(&ProductAttributeLine{}); err != nil {
		log.Error(err)
	}
	return &productAttributeLineRepo{
		data: data,
		log:  log,
	}
}

func (repo *productAttributeLineRepo) Create(ctx context.Context, productAttributeLine *ProductAttributeLine) (uint, error) {
	err := repo.data.DB(ctx).Create(productAttributeLine).Error
	return productAttributeLine.ID, err
}

func (repo *productAttributeLineRepo) CreateBatch(ctx context.Context, productAttributeLine []*ProductAttributeLine) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(productAttributeLine, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *productAttributeLineRepo) Update(ctx context.Context, productAttributeLine *ProductAttributeLine) (uint, error) {
	err := repo.data.DB(ctx).Updates(productAttributeLine).Error
	return productAttributeLine.ID, err
}

func (repo *productAttributeLineRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &ProductAttributeLine{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *productAttributeLineRepo) Get(ctx context.Context, id uint) (*ProductAttributeLine, error) {
	city := &ProductAttributeLine{}
	city.ID = id
	err := repo.data.DB(ctx).Preload("AttributeValues").First(city).Error
	return city, err
}

func (repo *productAttributeLineRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []ProductAttributeLine, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []ProductAttributeLine
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
