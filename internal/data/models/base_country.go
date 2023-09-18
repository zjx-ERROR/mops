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

// AddressCountry 国家
type AddressCountry struct {
	gorm.Model
	CreateUserID int64              ``                                 //创建者
	UpdateUserID int64              ``                                 //最后更新者
	Name         string             `gorm:"unique;size:50" xml:"name"` //国家名称
	Provinces    []*AddressProvince ``                                 //省份
}

type addressCountryRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewAddressCountryRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.AddressCountryRepo {
	if err := data.DB(ctx).AutoMigrate(&AddressCountry{}); err != nil {
		log.Error(err)
	}
	return &addressCountryRepo{
		data: data,
		log:  log,
	}
}

func (repo *addressCountryRepo) Create(ctx context.Context, addressCountry *AddressCountry) (uint, error) {
	err := repo.data.DB(ctx).Create(addressCountry).Error
	return addressCountry.ID, err
}

func (repo *addressCountryRepo) CreateBatch(ctx context.Context, addressCountry []*AddressCountry) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(addressCountry, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *addressCountryRepo) Update(ctx context.Context, addressCountry *AddressCountry) (uint, error) {
	err := repo.data.DB(ctx).Updates(addressCountry).Error
	return addressCountry.ID, err
}

func (repo *addressCountryRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &AddressCountry{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *addressCountryRepo) Get(ctx context.Context, id uint) (*AddressCountry, error) {
	ml := &AddressCountry{}
	ml.ID = id
	err := repo.data.DB(ctx).First(ml).Error
	return ml, err
}

func (repo *addressCountryRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []AddressCountry, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []AddressCountry
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
