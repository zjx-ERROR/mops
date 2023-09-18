package models

import (
	"context"
	"mops/internal/biz"
	"mops/internal/data"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Company 公司
type Company struct {
	gorm.Model
	CreateUserID int64            ``              //创建者
	UpdateUserID int64            ``              //最后更新者
	Name         string           `gorm:"unique"` //公司名称
	Code         string           `gorm:"unique"` //公司编码
	Children     []*Company       ``              //子公司
	Parent       *Company         ``              //上级公司
	Country      *AddressCountry  ``              //国家
	Province     *AddressProvince ``              //省份
	City         *AddressCity     ``              //城市
	District     *AddressDistrict ``              //区县
	Street       string           ``              //街道
}

type companyRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewCompanyRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.CompanyRepo {
	if err := data.DB(ctx).AutoMigrate(&Company{}); err != nil {
		log.Error(err)
	}
	return &companyRepo{
		data: data,
		log:  log,
	}
}

func (repo *companyRepo) Create(ctx context.Context, company *Company) (uint, error) {
	err := repo.data.DB(ctx).Create(company).Error
	return company.ID, err
}

func (repo *companyRepo) CreateBatch(ctx context.Context, company []*Company) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(company, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *companyRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &Company{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}
