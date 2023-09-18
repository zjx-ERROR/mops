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

// ProductCategory 产品分类
type ProductCategory struct {
	gorm.Model
	CreateUserID int64              ``                           //创建者
	UpdateUserID int64              ``                           //最后更新者
	ParentLeft   int64              `gorm:"unique"`              //分类左
	ParentRight  int64              `gorm:"unique"`              //分类右
	Name         string             `gorm:"size:50" json:"name"` //分类名称
	Parent       *ProductCategory   ``                           //上级分类
	Childs       []*ProductCategory ``                           //下级分类
}

type productCategoryRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewProductCategoryCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ProductCategoryRepo {
	if err := data.DB(ctx).AutoMigrate(&ProductCategory{}); err != nil {
		log.Error(err)
	}
	return &productCategoryRepo{
		data: data,
		log:  log,
	}
}

func (repo *productCategoryRepo) Create(ctx context.Context, productCategory *ProductCategory) (uint, error) {
	err := repo.data.DB(ctx).Create(productCategory).Error
	return productCategory.ID, err
}

func (repo *productCategoryRepo) CreateBatch(ctx context.Context, productCategory []*ProductCategory) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(productCategory, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *productCategoryRepo) Update(ctx context.Context, productCategory *ProductCategory) (uint, error) {
	err := repo.data.DB(ctx).Updates(productCategory).Error
	return productCategory.ID, err
}

func (repo *productCategoryRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &ProductCategory{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *productCategoryRepo) Get(ctx context.Context, id uint) (*ProductCategory, error) {
	city := &ProductCategory{}
	city.ID = id
	err := repo.data.DB(ctx).Preload(clause.Associations).First(city).Error
	return city, err
}

func (repo *productCategoryRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []ProductCategory, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []ProductCategory
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
