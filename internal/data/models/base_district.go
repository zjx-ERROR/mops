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

// AddressDistrict 区县
type AddressDistrict struct {
	gorm.Model
	CreateUserID int64        `` //创建者
	UpdateUserID int64        `` //最后更新者
	Name         string       `` //区县名称
	City         *AddressCity `` //城市
}

type addressDistrictRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewAddressDistrictRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.AddressDistrictRepo {
	if err := data.DB(ctx).AutoMigrate(&AddressDistrict{}); err != nil {
		log.Error(err)
	}
	return &addressDistrictRepo{
		data: data,
		log:  log,
	}
}

func (repo *addressDistrictRepo) Create(ctx context.Context, addressDistrict *AddressDistrict) (uint, error) {
	err := repo.data.DB(ctx).Create(addressDistrict).Error
	return addressDistrict.ID, err
}

func (repo *addressDistrictRepo) CreateBatch(ctx context.Context, addressDistrict []*AddressDistrict) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(addressDistrict, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *addressDistrictRepo) Update(ctx context.Context, addressDistrict *AddressDistrict) (uint, error) {
	err := repo.data.DB(ctx).Updates(addressDistrict).Error
	return addressDistrict.ID, err
}

func (repo *addressDistrictRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &AddressDistrict{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *addressDistrictRepo) Get(ctx context.Context, id uint) (*AddressDistrict, error) {
	ml := &AddressDistrict{}
	ml.ID = id
	err := repo.data.DB(ctx).First(ml).Error
	return ml, err
}

func (repo *addressDistrictRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []AddressDistrict, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []AddressDistrict
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
