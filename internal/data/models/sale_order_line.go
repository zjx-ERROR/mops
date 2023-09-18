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

// SaleOrderLine 销售订单明细
type SaleOrderLine struct {
	gorm.Model
	CreateUserID  int64           ``                                      //创建者
	UpdateUserID  int64           ``                                      //最后更新者
	Name          string          `json:"Name"`                           //订单明细号
	Company       *Company        ``                                      //公司
	SaleOrder     *SaleOrder      ``                                      //销售订单
	Partner       *Partner        ``                                      //客户
	Product       *ProductProduct ``                                      //产品
	ProductName   string          `json:"ProductName"`                    //产品名称
	ProductCode   string          `json:"ProductCode"`                    //产品编码
	FirstSaleUom  *ProductUom     ``                                      //第一销售单位
	SecondSaleUom *ProductUom     ``                                      //第二销售单位
	FirstSaleQty  float32         `gorm:"default:1" json:"FirstSaleQty"`  //第一销售单位
	SecondSaleQty float32         `gorm:"default:0" json:"SecondSaleQty"` //第二销售单位
	State         string          `gorm:"default:draft"`                  //订单明细状态:draft/confirm/process/done/cancel
	PriceUnit     float32         `gorm:"default:0" json:"PriceUnit"`     //单价
	Total         float32         `gorm:"default:0" json:"Total"`         //小计
}

type saleOrderLineRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewSaleOrderLineCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.SaleOrderLineRepo {
	if err := data.DB(ctx).AutoMigrate(&SaleOrderLine{}); err != nil {
		log.Error(err)
	}
	return &saleOrderLineRepo{
		data: data,
		log:  log,
	}
}

func (repo *saleOrderLineRepo) Create(ctx context.Context, saleOrderLine *SaleOrderLine) (uint, error) {
	err := repo.data.DB(ctx).Create(saleOrderLine).Error
	return saleOrderLine.ID, err
}

func (repo *saleOrderLineRepo) CreateBatch(ctx context.Context, saleOrderLine []*SaleOrderLine) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(saleOrderLine, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *saleOrderLineRepo) Update(ctx context.Context, saleOrderLine *SaleOrderLine) (uint, error) {
	err := repo.data.DB(ctx).Updates(saleOrderLine).Error
	return saleOrderLine.ID, err
}

func (repo *saleOrderLineRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &SaleOrderLine{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *saleOrderLineRepo) Get(ctx context.Context, id uint) (*SaleOrderLine, error) {
	city := &SaleOrderLine{}
	city.ID = id
	err := repo.data.DB(ctx).First(city).Error
	return city, err
}

func (repo *saleOrderLineRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []SaleOrderLine, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []SaleOrderLine
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
