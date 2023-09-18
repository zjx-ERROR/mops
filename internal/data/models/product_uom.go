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

// ProductUom 产品单位
type ProductUom struct {
	gorm.Model
	CreateUserID int64            ``                         //创建者
	UpdateUserID int64            ``                         //最后更新者
	Name         string           `gorm:"unique"`            //计量单位名称
	Category     *ProductUomCateg ``                         //计量单位类别
	Factor       float64          ``                         //比率
	FactorInv    float64          ``                         //更大比率
	Rounding     float64          ``                         //舍入精度
	Type         string           `gorm:"default:reference"` //类型：参考单位:reference;大于参考单位:bigger;小于参考单位:smaller
	Symbol       bool             ``                         //符号位置
}

type productUomRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewProductUomCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ProductUomRepo {
	if err := data.DB(ctx).AutoMigrate(&ProductUom{}); err != nil {
		log.Error(err)
	}
	return &productUomRepo{
		data: data,
		log:  log,
	}
}

func (repo *productUomRepo) Create(ctx context.Context, productUom *ProductUom) (uint, error) {
	err := repo.data.DB(ctx).Create(productUom).Error
	return productUom.ID, err
}

func (repo *productUomRepo) CreateBatch(ctx context.Context, productUom []*ProductUom) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(productUom, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *productUomRepo) Update(ctx context.Context, productUom *ProductUom) (uint, error) {
	err := repo.data.DB(ctx).Updates(productUom).Error
	return productUom.ID, err
}

func (repo *productUomRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &ProductUom{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *productUomRepo) Get(ctx context.Context, id uint) (*ProductUom, error) {
	city := &ProductUom{}
	city.ID = id
	err := repo.data.DB(ctx).First(city).Error
	return city, err
}

func (repo *productUomRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []ProductUom, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []ProductUom
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
