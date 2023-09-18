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

// GroupMenu 城市
type GroupMenu struct {
	gorm.Model
	CreateUserID int64      `` //创建者
	UpdateUserID int64      `` //最后更新者
	Group        *BaseGroup `` //权限组
	Menu         *BaseMenu  `` //菜单
}

type groupMenuRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewGroupMenuRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.GroupMenuRepo {
	if err := data.DB(ctx).AutoMigrate(&GroupMenu{}); err != nil {
		log.Error(err)
	}
	return &groupMenuRepo{
		data: data,
		log:  log,
	}
}

func (repo *groupMenuRepo) Create(ctx context.Context, groupMenu *GroupMenu) (uint, error) {
	err := repo.data.DB(ctx).Create(groupMenu).Error
	return groupMenu.ID, err
}

func (repo *groupMenuRepo) CreateBatch(ctx context.Context, groupMenu []*GroupMenu) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(groupMenu, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *groupMenuRepo) Update(ctx context.Context, groupMenu *GroupMenu) (uint, error) {
	err := repo.data.DB(ctx).Updates(groupMenu).Error
	return groupMenu.ID, err
}

func (repo *groupMenuRepo) Get(ctx context.Context, id uint) (*GroupMenu, error) {
	ml := &GroupMenu{}
	ml.ID = id
	err := repo.data.DB(ctx).First(ml).Error
	return ml, err
}

func (repo *groupMenuRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []GroupMenu, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []GroupMenu
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
