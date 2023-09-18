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

// BaseMenu 菜单
type BaseMenu struct {
	gorm.Model
	CreateUserID  int64        ``                             //创建者
	UpdateUserID  int64        ``                             //最后更新者
	ParentLeft    int64        `gorm:"unique"`                //菜单左
	ParentRight   int64        `gorm:"unique"`                //菜单右
	Name          string       `gorm:"size:50" json:"name"`   //菜单名称
	Parent        *BaseMenu    ``                             //上级菜单
	Childs        []*BaseMenu  ``                             //子菜单
	Icon          string       ``                             //菜单图标样式
	Groups        []*BaseGroup `gorm:"many2many:group_menu;"` //权限组
	Path          string       `json:"path"`                  //菜单点击地址
	ComponentPath string       ``                             //组件名称
	Meta          string       ``                             //额外参数
	ViewType      string       ``                             //视图类型,json数据，需要提供path和component路径信息
	IsBackground  bool         `gorm:"default:true"`          //前台还是后台
	Index         string       `gorm:"unique"`                //唯一标识
}

type baseMenuRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewBaseMenuRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.BaseMenuRepo {
	if err := data.DB(ctx).AutoMigrate(&BaseMenu{}); err != nil {
		log.Error(err)
	}
	return &baseMenuRepo{
		data: data,
		log:  log,
	}
}

func (repo *baseMenuRepo) Create(ctx context.Context, baseMenu *BaseMenu) (uint, error) {
	err := repo.data.DB(ctx).Create(baseMenu).Error
	return baseMenu.ID, err
}

func (repo *baseMenuRepo) CreateBatch(ctx context.Context, baseMenu []*BaseMenu) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(baseMenu, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *baseMenuRepo) Update(ctx context.Context, baseMenu *BaseMenu) (uint, error) {
	err := repo.data.DB(ctx).Updates(baseMenu).Error
	return baseMenu.ID, err
}

func (repo *baseMenuRepo) Get(ctx context.Context, id uint) (*BaseMenu, error) {
	ml := &BaseMenu{}
	ml.ID = id
	err := repo.data.DB(ctx).Preload("Childs").First(ml).Error
	return ml, err
}

func (repo *baseMenuRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []BaseMenu, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []BaseMenu
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
