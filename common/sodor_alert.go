package common

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"strings"
)

type SodorAlertPluginsWrapper struct {
	ins *sodor.AlertPluginInstances
}

func NewSodorAlertPluginsWrapper(m *sodor.AlertPluginInstances) *SodorAlertPluginsWrapper {
	return &SodorAlertPluginsWrapper{
		ins: m,
	}
}

func (s *SodorAlertPluginsWrapper) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(s.ins.AlertPluginInstances), len(s.ins.AlertPluginInstances))

	for i := 0; i < len(s.ins.AlertPluginInstances); i++ {
		rs[len(s.ins.AlertPluginInstances)-1-i] = s.ins.AlertPluginInstances[i]
	}

	return rs
}

type SodorAlertGroupsWrapper struct {
	groups *sodor.AlertGroups
}

func NewSodorAlertGroupsWrapperWrapper(m *sodor.AlertGroups) *SodorAlertGroupsWrapper {
	return &SodorAlertGroupsWrapper{
		groups: m,
	}
}

func (s *SodorAlertGroupsWrapper) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, 0)
	for i := 0; i < len(s.groups.AlertGroups); i++ {
		plugins := GetSodorCache().GetAlertPluginInstances(s.groups.AlertGroups[i].PluginInstances...)
		ag := convertToSodorAlertGroup(s.groups.AlertGroups[i], plugins)
		if ag != nil {
			rs = append(rs, *ag)
		}
	}

	return rs
}

type SodorAlertGroup struct {
	ID          string
	Name        string
	PluginNames string
	CreateTime  string
	UpdateTime  string

	GroupObj *sodor.AlertGroup
}

func convertToSodorAlertGroup(group *sodor.AlertGroup, plugins *sodor.AlertPluginInstances) *SodorAlertGroup {
	if plugins == nil {
		return nil
	}
	var ag SodorAlertGroup

	ag.ID = fmt.Sprintf("%d", group.Id)
	ag.Name = group.Name
	names := make([]string, 0)
	for _, p := range plugins.AlertPluginInstances {
		names = append(names, p.Name)
	}
	ag.PluginNames = strings.Join(names, ArraySeparator)
	ag.CreateTime = gobase.FormatTimeStamp(int64(group.CreateAt))
	ag.UpdateTime = gobase.FormatTimeStamp(int64(group.UpdateAt))
	ag.GroupObj = group
	return &ag
}

type SodorAlertGroupHistoriesWrapper struct {
	histories *sodor.AlertPluginInstanceHistories
}

func NewSodorAlertGroupHistoriesWrapperWrapper(m *sodor.AlertPluginInstanceHistories) *SodorAlertGroupHistoriesWrapper {
	return &SodorAlertGroupHistoriesWrapper{
		histories: m,
	}
}

func (s *SodorAlertGroupHistoriesWrapper) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(s.histories.AlertPluginInstanceHistory), len(s.histories.AlertPluginInstanceHistory))
	for i := 0; i < len(s.histories.AlertPluginInstanceHistory); i++ {
		rs[len(s.histories.AlertPluginInstanceHistory)-1-i] = s.histories.AlertPluginInstanceHistory[i]
	}

	return rs
}
