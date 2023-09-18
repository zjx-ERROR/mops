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

// SaleOrder 销售订单
type SaleOrder struct {
	gorm.Model
	CreateUserID int64            ``                          //创建者
	UpdateUserID int64            ``                          //最后更新者
	Name         string           `gorm:"unique" json:"name"` //订单号
	Partner      *Partner         ``                          //客户
	SalesMan     *User            ``                          //业务员
	Company      *Company         ``                          //公司
	Country      *AddressCountry  `json:"-"`                  //国家
	Province     *AddressProvince `json:"-"`                  //省份
	City         *SaleOrder       `json:"-"`                  //城市
	District     *AddressDistrict `json:"-"`                  //区县
	Street       string           `json:"Street"`             //街道
	OrderLine    []*SaleOrderLine ``                          //订单明细
	State        string           `gorm:"default:draft"`      //状态draft/confirm/process/done/cancel
}

type saleOrderRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewSaleOrderCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.SaleOrderRepo {
	if err := data.DB(ctx).AutoMigrate(&SaleOrder{}); err != nil {
		log.Error(err)
	}
	return &saleOrderRepo{
		data: data,
		log:  log,
	}
}

func (repo *saleOrderRepo) Create(ctx context.Context, saleOrder *SaleOrder) (uint, error) {
	err := repo.data.DB(ctx).Create(saleOrder).Error
	return saleOrder.ID, err
}

func (repo *saleOrderRepo) CreateBatch(ctx context.Context, saleOrder []*SaleOrder) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(saleOrder, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *saleOrderRepo) Update(ctx context.Context, saleOrder *SaleOrder) (uint, error) {
	err := repo.data.DB(ctx).Updates(saleOrder).Error
	return saleOrder.ID, err
}

func (repo *saleOrderRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &SaleOrder{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *saleOrderRepo) Get(ctx context.Context, id uint) (*SaleOrder, error) {
	city := &SaleOrder{}
	city.ID = id
	err := repo.data.DB(ctx).First(city).Error
	return city, err
}

func (repo *saleOrderRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []SaleOrder, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []SaleOrder
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
