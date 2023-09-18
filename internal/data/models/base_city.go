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

// AddressCity 城市
type AddressCity struct {
	gorm.Model
	CreateUserID int64              ``               //创建者
	UpdateUserID int64              ``               //最后更新者
	Name         string             `gorm:"size:50"` //城市名称
	Province     *AddressProvince   ``               //省份
	Districts    []*AddressDistrict ``               //区县
}

type addressCityRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewAddressCityCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.AddressCityRepo {
	if err := data.DB(ctx).AutoMigrate(&AddressCity{}); err != nil {
		log.Error(err)
	}
	return &addressCityRepo{
		data: data,
		log:  log,
	}
}

func (repo *addressCityRepo) Create(ctx context.Context, addressCity *AddressCity) (uint, error) {
	err := repo.data.DB(ctx).Create(addressCity).Error
	return addressCity.ID, err
}

func (repo *addressCityRepo) CreateBatch(ctx context.Context, addressCity []*AddressCity) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(addressCity, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *addressCityRepo) Update(ctx context.Context, addressCity *AddressCity) (uint, error) {
	err := repo.data.DB(ctx).Updates(addressCity).Error
	return addressCity.ID, err
}

func (repo *addressCityRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &AddressCity{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *addressCityRepo) Get(ctx context.Context, id uint) (*AddressCity, error) {
	city := &AddressCity{}
	city.ID = id
	err := repo.data.DB(ctx).First(city).Error
	return city, err
}

func (repo *addressCityRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []AddressCity, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []AddressCity
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
