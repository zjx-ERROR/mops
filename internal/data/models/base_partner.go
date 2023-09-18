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

// Partner 合作伙伴，包括客户和供应商，后期会为每个合作伙伴自动创建一个登录帐号
type Partner struct {
	gorm.Model
	CreateUserID int64            ``                                       //创建者
	UpdateUserID int64            ``                                       //最后更新者
	Name         string           `gorm:"unique" json:"Name"`              //合作伙伴名称
	IsCompany    bool             `gorm:"default:true" json:"IsCompany"`   //是公司
	IsSupplier   bool             `gorm:"default:false" json:"IsSupplier"` //是供应商
	IsCustomer   bool             `gorm:"default:true" json:"IsCustomer"`  //是客户
	Active       bool             `gorm:"default:true" json:"Active"`      //有效
	Country      *AddressCountry  ``                                       //国家
	Province     *AddressProvince ``                                       //省份
	City         *AddressCity     ``                                       //城市
	District     *AddressDistrict ``                                       //区县
	Street       string           `json:"Street"`                          //街道
	Parent       *Partner         ``                                       //母公司
	Childs       []*Partner       ``                                       //下级
	Mobile       string           `json:"Mobile"`                          //电话号码
	Tel          string           `json:"Tel"`                             //座机
	Email        string           `json:"Email"`                           //邮箱
	Qq           string           `json:"Qq"`                              //QQ
	WeChat       string           `json:"WeChat"`                          //微信
	Comment      string           `gorm:"type:text" json:"Comment"`        //备注

}

type partnerRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewPartnerRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.PartnerRepo {
	if err := data.DB(ctx).AutoMigrate(&Partner{}); err != nil {
		log.Error(err)
	}
	return &partnerRepo{
		data: data,
		log:  log,
	}
}

func (repo *partnerRepo) Create(ctx context.Context, partner *Partner) (uint, error) {
	err := repo.data.DB(ctx).Create(partner).Error
	return partner.ID, err
}

func (repo *partnerRepo) CreateBatch(ctx context.Context, partner []*Partner) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(partner, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *partnerRepo) Update(ctx context.Context, partner *Partner) (uint, error) {
	err := repo.data.DB(ctx).Updates(partner).Error
	return partner.ID, err
}

func (repo *partnerRepo) Get(ctx context.Context, id uint) (*Partner, error) {
	ml := &Partner{}
	ml.ID = id
	err := repo.data.DB(ctx).Preload("Province.Country").First(ml).Error
	return ml, err
}

func (repo *partnerRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []Partner, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []Partner
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
