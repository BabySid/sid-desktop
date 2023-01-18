package common

import (
	"github.com/BabySid/proto/sodor"
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
