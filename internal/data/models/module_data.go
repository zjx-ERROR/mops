package models

import (
	"context"
	"mops/internal/biz"
	"mops/internal/data"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ModuleData xml初始化数据记录
type ModuleData struct {
	gorm.Model
	CreateUserID int64  ``                                  //创建者
	UpdateUserID int64  ``                                  //最后更新者
	XMLID        string `gorm:"column:xml_id;unique;index"` //xml文件中的id
	Data         string `orm:"null"`                        //数据内容
	Descrition   string `orm:"null"`                        //记录描述
	InsertID     int64  `orm:"column(insert_id)"`           //插入记录的ID
	ModuleName   string `orm:""`                            //模块(表)的名称
}

type moduleDataRepo struct {
	data *data.Data
	log  *zap.SugaredLogger
}

func NewModuleDataCityRepo(ctx context.Context, data *data.Data, log *zap.SugaredLogger) biz.ModuleDataRepo {
	if err := data.DB(ctx).AutoMigrate(&ModuleData{}); err != nil {
		log.Error(err)
	}
	return &moduleDataRepo{
		data: data,
		log:  log,
	}
}

func (repo *moduleDataRepo) Create(ctx context.Context, moduleData *ModuleData) (uint, error) {
	err := repo.data.DB(ctx).Create(moduleData).Error
	return moduleData.ID, err
}

func (repo *moduleDataRepo) GetByXMLID(ctx context.Context, xml_id string) (*ModuleData, error) {
	city := &ModuleData{}
	err := repo.data.DB(ctx).Where("xml_id = ?", xml_id).First(city).Error
	return city, err
}
