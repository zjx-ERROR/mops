package data

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type contextTxKey struct{}

type Data struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db
}

func (d *Data) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

func NewDB(dialector, dns string, conf *gorm.Config) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)
	switch dialector {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dns), conf)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dns), conf)
	}
	if db == nil || err != nil {
		return nil
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(6 * time.Hour)
	return db
}

func NewData(db *gorm.DB, log *zap.SugaredLogger) (*Data, func(), error) {
	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
	return &Data{db: db, log: log}, cleanup, nil
}
