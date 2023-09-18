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

// ProductUomCateg 产品单位类别
type ProductUomCateg struct {
	gorm.Model
	CreateUserID int64         ``              //创建者
	UpdateUserID int64         ``              //最后更新者
	Name         string        `gorm:"unique"` //计量单位分类
	Uoms         []*ProductUom ``              //计量单位
}

type productUomCategRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewProductUomCategCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ProductUomCategRepo {
	if err := data.DB(ctx).AutoMigrate(&ProductUomCateg{}); err != nil {
		log.Error(err)
	}
	return &productUomCategRepo{
		data: data,
		log:  log,
	}
}

func (repo *productUomCategRepo) Create(ctx context.Context, ProductUomCateg *ProductUomCateg) (uint, error) {
	err := repo.data.DB(ctx).Create(ProductUomCateg).Error
	return ProductUomCateg.ID, err
}

func (repo *productUomCategRepo) CreateBatch(ctx context.Context, ProductUomCateg []*ProductUomCateg) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(ProductUomCateg, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *productUomCategRepo) Update(ctx context.Context, ProductUomCateg *ProductUomCateg) (uint, error) {
	err := repo.data.DB(ctx).Updates(ProductUomCateg).Error
	return ProductUomCateg.ID, err
}

func (repo *productUomCategRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &ProductUomCateg{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *productUomCategRepo) Get(ctx context.Context, id uint) (*ProductUomCateg, error) {
	city := &ProductUomCateg{}
	city.ID = id
	err := repo.data.DB(ctx).First(city).Error
	return city, err
}

func (repo *productUomCategRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []ProductUomCateg, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []ProductUomCateg
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
