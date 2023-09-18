package models

import (
	"context"
	"mops/internal/data"
	"testing"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestInit(t *testing.T) {
	db := data.NewDB("sqlite", "../../../test/test.sqlite", &gorm.Config{})
	data, cleanup, err := data.NewData(db, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	NewAddressCityCityRepo(context.TODO(), data, zap.NewNop().Sugar())
}
