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

// User 登录用户
type User struct {
	gorm.Model
	CreateUserID    int64        ``                                        //创建者
	UpdateUserID    int64        ``                                        //最后更新者
	Name            string       `gorm:"size:20;unique" xml:"name"`        //用户名
	Company         *Company     ``                                        //公司
	NameZh          string       `gorm:"size:20"  xml:"NameZh"`            //中文用户名
	Email           string       `gorm:"size:20;unique" xml:"email"`       //邮箱
	Mobile          string       `gorm:"size:20;unique" xml:"mobile"`      //手机号码
	Tel             string       `gorm:"size:20"`                          //固定号码
	Password        string       `xml:"password"`                          //密码
	ConfirmPassword string       `gorm:"-" xml:"ConfirmPassword"`          //确认密码,数据库中不保存
	IsAdmin         bool         `gorm:"default:false" xml:"isAdmin"`      //是否为超级用户
	Active          bool         `gorm:"default:true" xml:"active"`        //有效
	Qq              string       `xml:"qq"`                                //QQ
	WeChat          string       `xml:"wechat"`                            //微信
	Groups          []*BaseGroup `gorm:"many2many:group_user;"`            //权限组
	IsBackground    bool         `gorm:"defalut:false" xml:"isbackground"` //后台用户可以登录后台
}

type userRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewUserCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.UserRepo {
	if err := data.DB(ctx).AutoMigrate(&User{}); err != nil {
		log.Error(err)
	}
	return &userRepo{
		data: data,
		log:  log,
	}
}

func (repo *userRepo) Create(ctx context.Context, user *User) (uint, error) {
	err := repo.data.DB(ctx).Create(user).Error
	return user.ID, err
}

func (repo *userRepo) CreateBatch(ctx context.Context, user []*User) (int64, error) {
	tx := repo.data.DB(ctx).CreateInBatches(user, 100)
	err := tx.Error
	return tx.RowsAffected, err
}

func (repo *userRepo) Update(ctx context.Context, user *User) (uint, error) {
	err := repo.data.DB(ctx).Updates(user).Error
	return user.ID, err
}

func (repo *userRepo) Delete(ctx context.Context, id uint) (int64, error) {
	ml := &User{}
	ml.ID = id
	tx := repo.data.DB(ctx).Delete(ml)
	return tx.RowsAffected, tx.Error
}

func (repo *userRepo) Get(ctx context.Context, id uint) (*User, error) {
	city := &User{}
	city.ID = id
	err := repo.data.DB(ctx).Preload("Groups").First(city).Error
	return city, err
}

func (repo *userRepo) List(ctx context.Context, exclude map[string]interface{}, condMap map[string]map[interface{}][]interface{},
	fields []string, orderBy string, page int64, limit int64) (pkg.Paginator, []User, error) {
	if limit == 0 {
		limit = 200
	}
	var (
		objArrs   []User
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
