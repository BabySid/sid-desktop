package common

import (
	"github.com/BabySid/proto/sodor"
	"sync"
)

type sodorCache struct {
	thomasInfos          *sodor.ThomasInfos
	alertPluginInstances *sodor.AlertPluginInstances
	alertGroups          *sodor.AlertGroups
}

var (
	sodorCacheHandle *sodorCache
	sodorCacheOnce   sync.Once
)

func GetSodorCache() *sodorCache {
	sodorOnce.Do(func() {
		sodorCacheHandle = &sodorCache{}
	})
	return sodorCacheHandle
}

func (s *sodorCache) LoadThomasInfos() error {
	return nil
}

func (s *sodorCache) GetThomasInfos() *sodor.ThomasInfos {
	return s.thomasInfos
}

func (s *sodorCache) LoadAlertPluginInstances() error {
	return nil
}

func (s *sodorCache) GetAlertPluginInstances() *sodor.AlertPluginInstances {
	return s.alertPluginInstances
}

func (s *sodorCache) LoadAlertGroups() error {
	return nil
}

func (s *sodorCache) GetAlertGroups() *sodor.AlertGroups {
	return s.alertGroups
}
