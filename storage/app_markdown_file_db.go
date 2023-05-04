package storage

import (
	"github.com/BabySid/gobase"
	"gorm.io/gorm"
	"path/filepath"
	"sid-desktop/common"
	"sync"
)

type AppMarkDownDB struct {
	handle *gorm.DB
	dbName string
}

var (
	appMarkDownDB   *AppMarkDownDB
	appMarkDownOnce sync.Once
)

func GetAppMarkDownDB() *AppMarkDownDB {
	appMarkDownOnce.Do(func() {
		appMarkDownDB = &AppMarkDownDB{
			dbName: "sid_app_markdown.db",
		}
	})
	return appMarkDownDB
}

func (db *AppMarkDownDB) Open(rootPath string) error {
	sqlDB, err := initGorm(filepath.Join(rootPath, db.dbName))
	if err != nil {
		return err
	}

	db.handle = sqlDB
	return nil
}

func (db *AppMarkDownDB) Close() {
	if db.handle != nil {
		if d, _ := db.handle.DB(); d != nil {
			_ = d.Close()
		}
	}
}

func (db *AppMarkDownDB) Init() error {
	if err := db.createMarkDown(); err != nil {
		return err
	}
	return nil
}

func (db *AppMarkDownDB) createMarkDown() error {
	return db.handle.AutoMigrate(&common.MarkDownFile{})
}

func (db *AppMarkDownDB) NeedInit() (bool, error) {
	if db.handle.Migrator().HasTable(&common.MarkDownFile{}) {
		return false, nil
	}
	return true, nil
}

func (db *AppMarkDownDB) LoadMarkDownFiles() (*common.MarkDownFileList, error) {
	var files common.MarkDownFileList
	if rs := db.handle.Order("updated_at desc").Find(&files); rs.Error != nil {
		return nil, rs.Error
	}
	return &files, nil
}

func (db *AppMarkDownDB) AddMarkDownFile(file *common.MarkDownFile) error {
	if rs := db.handle.Create(file); rs.Error != nil {
		return rs.Error
	}
	return nil
}

func (db *AppMarkDownDB) DelMarkDownFile(file *common.MarkDownFile) error {
	if rs := db.handle.Delete(&common.MarkDownFile{
		TableModel: gobase.TableModel{ID: file.ID},
	}); rs.Error != nil {
		return rs.Error
	}
	return nil
}

func (db *AppMarkDownDB) UpdateMarkDownFile(file *common.MarkDownFile) error {
	if rst := db.handle.Model(file).Select(file.UpdateFields()).Updates(file); rst.Error != nil {
		return rst.Error
	}
	return nil
}
