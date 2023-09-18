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

// ProductProduct 产品规格
type ProductProduct struct {
	gorm.Model
	CreateUserID          int64                    ``                                  //创建者
	UpdateUserID          int64                    ``                                  //最后更新者
	Name                  string                   `gorm:"index"`                      //产品属性名称
	Company               *Company                 ``                                  //公司
	Category              *ProductCategory         ``                                  //产品类别
	IsProductVariant      bool                     `gorm:"default:true"`               //是多规格产品
	ProductTags           []*ProductTag            ``                                  //产品标签
	SaleOk                bool                     `gorm:"default:true" json:"SaleOk"` //可销售
	Active                bool                     `gorm:"default:true"`               //有效
	Barcode               string                   `json:"Barcode"`                    //条码,如ean13
	StandardPrice         float64                  `json:"StandardPrice"`              //成本价格
	DefaultCode           string                   `gorm:"unique"`                     //产品编码
	ProductTemplate       *ProductTemplate         ``                                  //产品款式
	AttributeValues       []*ProductAttributeValue ``                                  //产品属性值
	ProductType           string                   `gorm:"default:stock"`              //产品类型
	AttributeValuesString string                   `gorm:"index;"`                     //产品属性值ID编码，用于修改和增加时对应的产品是否已经存在
	FirstSaleUom          *ProductUom              ``                                  //第一销售单位
	SecondSaleUom         *ProductUom              ``                                  //第二销售单位
	FirstPurchaseUom      *ProductUom              ``                                  //第一采购单位
	SecondPurchaseUom     *ProductUom              ``                                  //第二采购单位
	PackagingDependTemp   bool                     `gorm:"default:true"`               //根据款式打包
	BigImages             []*ProductImage          ``                                  //产品款式图片
	MidImages             []*ProductImage          ``                                  //产品款式图片
	SmallImages           []*ProductImage          ``                                  //产品款式图片
	PurchaseDependTemp    bool                     `gorm:"default:true"`               //根据款式采购，ture一个供应商可以供应所有的款式
	// ProductPackagings     []*ProductPackaging      `orm:"reverse(many)"`                        //打包方式
}

type productProductRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewProductProductCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ProductProductRepo {
	if err := data.DB(ctx).AutoMigrate(&ProductProduct{}); err != nil {
		log.Error(err)
	}
	return &productProductRepo{
		data: data,
		log:  log,
	}
}

func (repo *productProductRepo) Create(ctx context.Context, productProduct *ProductProduct) (uint, error) {
	err := repo.data.DB(ctx).Create(productProduct).Error
	return productProduct.ID, err
}

func (repo *productProductRepo) CreateBatch(ctx context.Context, productProduct []*ProductProduct) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(productProduct, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *productProductRepo) Update(ctx context.Context, productProduct *ProductProduct) (uint, error) {
	err := repo.data.DB(ctx).Updates(productProduct).Error
	return productProduct.ID, err
}

func (repo *productProductRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &ProductProduct{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *productProductRepo) Get(ctx context.Context, id uint) (*ProductProduct, error) {
	city := &ProductProduct{}
	city.ID = id
	err := repo.data.DB(ctx).First(city).Error
	return city, err
}

func (repo *productProductRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []ProductProduct, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []ProductProduct
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
