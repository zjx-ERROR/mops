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

// BaseGroup  权限组
type BaseGroup struct {
	gorm.Model
	CreateUserID  int64          ``                      //创建者
	UpdateUserID  int64          ``                      //最后更新者
	Name          string         `gorm:"unique;size:50"` //权限组名称
	ModelAccesses []*ModelAccess ``                      //模块(表)
	Childs        []*BaseGroup   ``                      //下级
	Parent        *BaseGroup     ``                      //上级
	ParentLeft    int64          `gorm:"unique"`         //左边界
	ParentRight   int64          `gorm:"unique"`         //右边界
	Category      string         ``                      //分类
	Description   string         ``                      //说明
	Menus         []*BaseMenu    ``                      //菜单
}

type baseGroupRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewBaseGroupRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.BaseGroupRepo {
	if err := data.DB(ctx).AutoMigrate(&BaseGroup{}); err != nil {
		log.Error(err)
	}
	return &baseGroupRepo{
		data: data,
		log:  log,
	}
}

func (repo *baseGroupRepo) Create(ctx context.Context, baseGroup *BaseGroup) (uint, error) {
	err := repo.data.DB(ctx).Create(baseGroup).Error
	return baseGroup.ID, err
}

func (repo *baseGroupRepo) CreateBatch(ctx context.Context, baseGroup []*BaseGroup) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(baseGroup, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *baseGroupRepo) Update(ctx context.Context, baseGroup *BaseGroup) (uint, error) {
	err := repo.data.DB(ctx).Updates(baseGroup).Error
	return baseGroup.ID, err
}

func (repo *baseGroupRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &BaseGroup{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *baseGroupRepo) Get(ctx context.Context, id uint) (*BaseGroup, error) {
	ml := &BaseGroup{}
	ml.ID = id
	err := repo.data.DB(ctx).First(ml).Error
	return ml, err
}

func (repo *baseGroupRepo) GetByName(ctx context.Context, name string) (*BaseGroup, error) {
	ml := &BaseGroup{}
	err := repo.data.DB(ctx).Where("name = ?", name).First(ml).Error
	return ml, err
}

func (repo *baseGroupRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []BaseGroup, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []BaseGroup
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
