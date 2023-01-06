package storage

import (
	"gorm.io/gorm"
	"path/filepath"
	"sid-desktop/common"
	"sync"
)

type AppSodorDB struct {
	dbName string

	handle *gorm.DB
}

var (
	appSodorDB   *AppSodorDB
	appSodorOnce sync.Once
)

func GetAppSodorDB() *AppSodorDB {
	appSodorOnce.Do(func() {
		appSodorDB = &AppSodorDB{
			dbName: "sid_app_sodor.db",
		}
	})
	return appSodorDB
}

func (db *AppSodorDB) Open(rootPath string) error {
	sqlDB, err := initGorm(filepath.Join(rootPath, db.dbName))
	if err != nil {
		return err
	}

	db.handle = sqlDB
	return nil
}

func (db *AppSodorDB) Close() {
	if db.handle != nil {
		if d, _ := db.handle.DB(); d != nil {
			_ = d.Close()
		}
	}
}

func (db *AppSodorDB) Init() error {
	if err := db.createFatCtrlTbl(); err != nil {
		return err
	}
	return nil
}

func (db *AppSodorDB) createFatCtrlTbl() error {
	return db.handle.AutoMigrate(&common.FatCtrl{})
}

func (db *AppSodorDB) NeedInit() (bool, error) {
	if db.handle.Migrator().HasTable(&common.FatCtrl{}) {
		return false, nil
	}
	return true, nil
}

func (db *AppSodorDB) SetFatCtrl(ctrl common.FatCtrl) error {
	if ctrl.ID == 0 {
		if rs := db.handle.Create(&ctrl); rs.Error != nil {
			return rs.Error
		}
		return nil
	}

	if rs := db.handle.Model(&ctrl).Updates(ctrl); rs.Error != nil {
		return rs.Error
	}

	return nil
}

func (db *AppSodorDB) LoadFatCtl() (*common.FatCtrl, error) {
	var ctrl common.FatCtrl
	rs := db.handle.Find(&ctrl)
	if rs.Error != nil {
		return nil, rs.Error
	}

	if rs.RowsAffected == 0 {
		return nil, nil
	}

	return &ctrl, nil
}
