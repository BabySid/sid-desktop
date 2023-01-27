package common

import (
	"github.com/BabySid/proto/sodor"
	"sync"
)

type sodorCache struct {
	thomasInfos          *sodor.ThomasInfos
	alertPluginInstances *sodor.AlertPluginInstances
	alertGroups          *sodor.AlertGroups
	jobs                 *sodor.Jobs
}

var (
	sodorCacheHandle *sodorCache
	sodorCacheOnce   sync.Once
)

func GetSodorCache() *sodorCache {
	sodorCacheOnce.Do(func() {
		sodorCacheHandle = &sodorCache{}
	})
	return sodorCacheHandle
}

func (s *sodorCache) LoadThomasInfos() error {
	resp := sodor.ThomasInfos{}
	err := GetSodorClient().Call(ListThomas, nil, &resp)
	if err != nil {
		return err
	}
	s.thomasInfos = &resp
	return nil
}

func (s *sodorCache) GetThomasInfos() *sodor.ThomasInfos {
	return s.thomasInfos
}

func (s *sodorCache) LoadAlertPluginInstances() error {
	resp := sodor.AlertPluginInstances{}
	err := GetSodorClient().Call(ListAlertPluginInstances, nil, &resp)
	if err != nil {
		return err
	}
	s.alertPluginInstances = &resp
	return nil
}

func (s *sodorCache) GetAlertPluginInstances(pluginID ...int32) *sodor.AlertPluginInstances {
	if s.alertPluginInstances == nil {
		return nil
	}
	if len(pluginID) == 0 {
		return s.alertPluginInstances
	}

	var rs sodor.AlertPluginInstances
	rs.AlertPluginInstances = make([]*sodor.AlertPluginInstance, 0)
	for _, id := range pluginID {
		for _, plugin := range s.alertPluginInstances.AlertPluginInstances {
			if id == plugin.Id {
				rs.AlertPluginInstances = append(rs.AlertPluginInstances, plugin)
				break
			}
		}
	}

	return &rs
}

func (s *sodorCache) LoadAlertGroups() error {
	resp := sodor.AlertGroups{}
	err := GetSodorClient().Call(ListAlertGroup, nil, &resp)
	if err != nil {
		return err
	}
	s.alertGroups = &resp
	return nil
}

func (s *sodorCache) GetAlertGroups() *sodor.AlertGroups {
	return s.alertGroups
}

func (s *sodorCache) LoadJobs() error {
	resp := sodor.Jobs{}
	err := GetSodorClient().Call(ListJobs, nil, &resp)
	if err != nil {
		return err
	}
	s.jobs = &resp
	return nil
}

func (s *sodorCache) GetJobs() *sodor.Jobs {
	return s.jobs
}
