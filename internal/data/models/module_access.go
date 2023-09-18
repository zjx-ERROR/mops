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

// ModelAccess 模块(表)操作权限
type ModelAccess struct {
	gorm.Model
	CreateUserID int64         ``                     //创建者
	UpdateUserID int64         ``                     //最后更新者
	Module       *ModuleModule ``                     //模块(表)
	Group        *BaseGroup    ``                     //权限组
	PermCreate   bool          `gorm:"default:true"`  //创建权限
	PermUnlink   bool          `gorm:"default:false"` //删除权限
	PermWrite    bool          `gorm:"default:true"`  //修改权限
	PermRead     bool          `gorm:"default:true"`  //读权限
	Domain       string        ``                     //过滤条件，只在本级有效(权限组直属访问权限)
}

type modelAccessRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewModelAccessCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ModelAccessRepo {
	if err := data.DB(ctx).AutoMigrate(&ModelAccess{}); err != nil {
		log.Error(err)
	}
	return &modelAccessRepo{
		data: data,
		log:  log,
	}
}

func (repo *modelAccessRepo) Create(ctx context.Context, modelAccess *ModelAccess) (uint, error) {
	err := repo.data.DB(ctx).Create(modelAccess).Error
	return modelAccess.ID, err
}

func (repo *modelAccessRepo) Update(ctx context.Context, modelAccess *ModelAccess) (uint, error) {
	err := repo.data.DB(ctx).Updates(modelAccess).Error
	return modelAccess.ID, err
}

func (repo *modelAccessRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []ModelAccess, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []ModelAccess
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
