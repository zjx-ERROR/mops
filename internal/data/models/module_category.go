package models

import (
	"context"
	"mops/internal/biz"
	"mops/internal/data"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ModuleCategory 模块(表)分类
type ModuleCategory struct {
	gorm.Model
	CreateModuleCategoryID int64  ``                          //创建者
	UpdateModuleCategoryID int64  ``                          //最后更新者
	Name                   string `gorm:"size:50" xml:"name"` //模块分类名称
	Description            string `xml:"description"`         //说明
}
type moduleCategoryRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewModuleCategoryCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ModuleCategoryRepo {
	if err := data.DB(ctx).AutoMigrate(&ModuleCategory{}); err != nil {
		log.Error(err)
	}
	return &moduleCategoryRepo{
		data: data,
		log:  log,
	}
}

func (repo *moduleCategoryRepo) Create(ctx context.Context, moduleCategory *ModuleCategory) (uint, error) {
	err := repo.data.DB(ctx).Create(moduleCategory).Error
	return moduleCategory.ID, err
}

func (repo *moduleCategoryRepo) CreateBatch(ctx context.Context, moduleCategory []*ModuleCategory) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(moduleCategory, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *moduleCategoryRepo) Update(ctx context.Context, moduleCategory *ModuleCategory) (uint, error) {
	err := repo.data.DB(ctx).Updates(moduleCategory).Error
	return moduleCategory.ID, err
}

func (repo *moduleCategoryRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &ModuleCategory{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *moduleCategoryRepo) Get(ctx context.Context, id uint) (*ModuleCategory, error) {
	city := &ModuleCategory{}
	city.ID = id
	err := repo.data.DB(ctx).First(city).Error
	return city, err
}

func (repo *moduleCategoryRepo) GetByName(ctx context.Context, name string) (*ModuleCategory, error) {
	city := &ModuleCategory{}
	err := repo.data.DB(ctx).Where("name = ?", name).First(city).Error
	return city, err
}
