package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"reflect"
	"sid-desktop/desktop/common/apps"
	"strings"
)

type Config struct {
	Theme binding.String

	AppLaunchAppSearchPath binding.StringList

	HideWhenQuit binding.Bool
}

func NewConfig() *Config {
	c := &Config{
		Theme:                  binding.NewString(),
		AppLaunchAppSearchPath: binding.NewStringList(),
		HideWhenQuit:           binding.NewBool(),
	}

	_ = c.Theme.Set(
		fyne.CurrentApp().Preferences().StringWithFallback("theme", "__DARK__"))

	searchPath := fyne.CurrentApp().Preferences().String("app_launch_search_path")
	if searchPath == "" {
		_ = c.AppLaunchAppSearchPath.Set(apps.DefaultAppPaths)
	} else {
		_ = c.AppLaunchAppSearchPath.Set(strings.Split(searchPath, ";"))
	}

	_ = c.HideWhenQuit.Set(fyne.CurrentApp().Preferences().BoolWithFallback("hide_when_quit", true))

	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(c)
	s := reflect.ValueOf(c).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		f.MethodByName("AddListener").Call(in)
	}

	return c
}

func (c *Config) DataChanged() {
	// todo flush with field tag instead of the first param
	theme, _ := c.Theme.Get()
	fyne.CurrentApp().Preferences().SetString("theme", theme)

	searchPath, _ := c.AppLaunchAppSearchPath.Get()
	fyne.CurrentApp().Preferences().SetString("app_launch_search_path", strings.Join(searchPath, ";"))

	hideWhenQuit, _ := c.HideWhenQuit.Get()
	fyne.CurrentApp().Preferences().SetBool("hide_when_quit", hideWhenQuit)
}
