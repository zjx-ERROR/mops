package models

import (
	"context"
	"mops/internal/biz"
	"mops/internal/data"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ModuleModule 模块(表)名称
type ModuleModule struct {
	gorm.Model
	CreateUserID int64           ``                                 //创建者
	UpdateUserID int64           ``                                 //最后更新者
	Name         string          `gorm:"unique;size:50" xml:"name"` //表名称
	Description  string          `xml:"description"`                //说明
	Category     *ModuleCategory ``                                 //模块分类
}

type moduleModuleRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewModuleModuleRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ModuleModuleRepo {
	if err := data.DB(ctx).AutoMigrate(&ModuleModule{}); err != nil {
		log.Error(err)
	}
	return &moduleModuleRepo{
		data: data,
		log:  log,
	}
}

func (repo *moduleModuleRepo) Create(ctx context.Context, moduleModule *ModuleModule) (uint, error) {
	err := repo.data.DB(ctx).Create(moduleModule).Error
	return moduleModule.ID, err
}

func (repo *moduleModuleRepo) Update(ctx context.Context, moduleModule *ModuleModule) (uint, error) {
	err := repo.data.DB(ctx).Updates(moduleModule).Error
	return moduleModule.ID, err
}

func (repo *moduleModuleRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &ModuleModule{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *moduleModuleRepo) Get(ctx context.Context, id uint) (*ModuleModule, error) {
	ml := &ModuleModule{}
	ml.ID = id
	err := repo.data.DB(ctx).First(ml).Error
	return ml, err
}

func (repo *moduleModuleRepo) GetByName(ctx context.Context, name string) (*ModuleModule, error) {
	ml := &ModuleModule{}
	err := repo.data.DB(ctx).Where("name = ?", name).First(ml).Error
	return ml, err
}
