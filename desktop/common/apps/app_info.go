package apps

import (
	"fmt"
	"github.com/BabySid/gobase"
	"github.com/sahilm/fuzzy"
	"runtime"
)

type AppList struct {
	apps []AppInfo
}

func NewAppList() *AppList {
	return &AppList{
		apps: make([]AppInfo, 0),
	}
}

func (s *AppList) Find(name string) *AppList {
	matches := fuzzy.FindFrom(name, s)

	rs := NewAppList()
	for _, match := range matches {
		rs.apps = append(rs.apps, s.apps[match.Index])
	}

	return rs
}

func (s *AppList) UpdateAppInfo(d AppInfo) {
	for i, app := range s.apps {
		if app.AppID == d.AppID {
			app.AccessTime = d.AccessTime
			s.apps[i] = app
			return
		}
	}
}

func (s *AppList) String(i int) string {
	return s.apps[i].AppName
}

func (s *AppList) Len() int {
	return len(s.apps)
}

func (s *AppList) Set(d []AppInfo) {
	if d == nil {
		return
	}
	s.apps = d
}

func (s *AppList) Append(d AppInfo) {
	s.apps = append(s.apps, d)
}

func (s *AppList) AsInterfaceArray() []interface{} {
	rs := make([]interface{}, len(s.apps), len(s.apps))
	for i := range s.apps {
		rs[i] = s.apps[i]
	}
	return rs
}

func (s *AppList) GetAppInfo() []AppInfo {
	return s.apps
}

func (s *AppList) Debug() {
	for _, app := range s.apps {
		fmt.Println(app.AppID, app.AppName, app.CreateTime, app.AccessTime)
	}
}

type AppInfo struct {
	AppID      int64
	AppName    string
	FullPath   string
	Icon       []byte
	CreateTime int64
	AccessTime int64
}

func (app *AppInfo) Exec() error {
	switch runtime.GOOS {
	case "windows":
		return gobase.ExecApp(app.FullPath)
	}
	gobase.AssertHere()
	return nil
}
