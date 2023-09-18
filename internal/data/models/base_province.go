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

// AddressProvince 省份
type AddressProvince struct {
	gorm.Model
	CreateUserID int64           `gorm:""`                               //创建者
	UpdateUserID int64           `gorm:""`                               //最后更新者
	Name         string          `gorm:"unique" xml:"ProvinceName,attr"` //省份名称
	Country      *AddressCountry ``                                      //国家
	Citys        []*AddressCity  ``                                      //城市
}

type addressProvinceRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewAddressProvinceCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.AddressProvinceRepo {
	if err := data.DB(ctx).AutoMigrate(&AddressProvince{}); err != nil {
		log.Error(err)
	}
	return &addressProvinceRepo{
		data: data,
		log:  log,
	}
}

func (repo *addressProvinceRepo) Create(ctx context.Context, addressProvince *AddressProvince) (uint, error) {
	err := repo.data.DB(ctx).Create(addressProvince).Error
	return addressProvince.ID, err
}

func (repo *addressProvinceRepo) CreateBatch(ctx context.Context, addressProvince []*AddressProvince) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(addressProvince, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *addressProvinceRepo) Update(ctx context.Context, addressProvince *AddressProvince) (uint, error) {
	err := repo.data.DB(ctx).Updates(addressProvince).Error
	return addressProvince.ID, err
}

func (repo *addressProvinceRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &AddressProvince{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *addressProvinceRepo) Get(ctx context.Context, id uint) (*AddressProvince, error) {
	city := &AddressProvince{}
	city.ID = id
	err := repo.data.DB(ctx).Preload("Country").First(city).Error
	return city, err
}

func (repo *addressProvinceRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []AddressProvince, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []AddressProvince
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
