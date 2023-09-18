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

// GroupUser 城市
type GroupUser struct {
	gorm.Model
	CreateUserID int64      `` //创建者
	UpdateUserID int64      `` //最后更新者
	Group        *BaseGroup `` //权限组
	User         *User      `` //用户
}

type groupUserRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewGroupUserRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.GroupUserRepo {
	if err := data.DB(ctx).AutoMigrate(&GroupUser{}); err != nil {
		log.Error(err)
	}
	return &groupUserRepo{
		data: data,
		log:  log,
	}
}

func (repo *groupUserRepo) Create(ctx context.Context, groupUser *GroupUser) (uint, error) {
	err := repo.data.DB(ctx).Create(groupUser).Error
	return groupUser.ID, err
}

func (repo *groupUserRepo) CreateBatch(ctx context.Context, groupUser []*GroupUser) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(groupUser, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *groupUserRepo) Update(ctx context.Context, groupUser *GroupUser) (uint, error) {
	err := repo.data.DB(ctx).Updates(groupUser).Error
	return groupUser.ID, err
}

func (repo *groupUserRepo) Get(ctx context.Context, id uint) (*GroupUser, error) {
	ml := &GroupUser{}
	ml.ID = id
	err := repo.data.DB(ctx).First(ml).Error
	return ml, err
}

func (repo *groupUserRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []GroupUser, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []GroupUser
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
