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

// ProductTemplate 产品款式
type ProductTemplate struct {
	gorm.Model
	CreateUserID        int64                   ``                     //创建者
	UpdateUserID        int64                   ``                     //最后更新者
	Name                string                  ``                     //款式名称
	Description         string                  `gorm:"type:text"`     //描述
	DescriptionSale     string                  `gorm:"type:text"`     //销售描述
	DescriptionPurchase string                  `gorm:"type:text"`     //采购描述
	Rental              bool                    `gorm:"default:false"` //代售品
	Category            *ProductCategory        ``                     //产品类别
	Price               float64                 ``                     //款式价格
	StandardPrice       float64                 ``                     //成本价格
	StandardWeight      float64                 ``                     //标准重量
	SaleOk              bool                    `gorm:"default:true"`  //可销售
	Active              bool                    `gorm:"default:true"`  //有效
	IsProductVariant    bool                    `gorm:"default:true"`  //是规格产品
	FirstSaleUom        *ProductUom             ``                     //第一销售单位
	SecondSaleUom       *ProductUom             ``                     //第二销售单位
	FirstPurchaseUom    *ProductUom             ``                     //第一采购单位
	SecondPurchaseUom   *ProductUom             ``                     //第二采购单位
	AttributeLines      []*ProductAttributeLine ``                     //属性明细
	ProductVariants     []*ProductProduct       ``                     //产品规格明细
	VariantCount        int32                   ``                     //产品规格数量
	Barcode             string                  ``                     //条码,如ean13
	DefaultCode         string                  ``                     //产品编码
	BigImages           []*ProductImage         ``                     //产品款式图片
	MidImages           []*ProductImage         ``                     //产品款式图片
	SmallImages         []*ProductImage         ``                     //产品款式图片
	ProductType         string                  `gorm:"default:stock"` //产品类型 stock consume service
	ProductMethod       string                  `gorm:"default:hand"`  //产品规格创建方式 auto hand
	// TemplatePackagings  []*ProductPackaging     `orm:"reverse(many)"`               //打包方式
}

type productTemplateRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewProductTemplateCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ProductTemplateRepo {
	if err := data.DB(ctx).AutoMigrate(&ProductTemplate{}); err != nil {
		log.Error(err)
	}
	return &productTemplateRepo{
		data: data,
		log:  log,
	}
}

func (repo *productTemplateRepo) Create(ctx context.Context, productTemplate *ProductTemplate) (uint, error) {
	err := repo.data.DB(ctx).Create(productTemplate).Error
	return productTemplate.ID, err
}

func (repo *productTemplateRepo) CreateBatch(ctx context.Context, productTemplate []*ProductTemplate) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(productTemplate, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *productTemplateRepo) Update(ctx context.Context, productTemplate *ProductTemplate) (uint, error) {
	err := repo.data.DB(ctx).Updates(productTemplate).Error
	return productTemplate.ID, err
}

func (repo *productTemplateRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &ProductTemplate{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *productTemplateRepo) Get(ctx context.Context, id uint) (*ProductTemplate, error) {
	city := &ProductTemplate{}
	city.ID = id
	err := repo.data.DB(ctx).Preload(clause.Associations).First(city).Error
	return city, err
}

func (repo *productTemplateRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []ProductTemplate, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []ProductTemplate
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
