package common

import "github.com/BabySid/gobase"

type FatCtrl struct {
	gobase.TableModel
	Addr string `gorm:"size:256;not null"`
}

func (t FatCtrl) TableName() string {
	return "fat_ctrl"
}
