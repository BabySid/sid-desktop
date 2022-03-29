package main

import (
	luaHttp "github.com/cjoudrey/gluahttp"
	lua "github.com/yuin/gopher-lua"
	luaJson "layeh.com/gopher-json"
	"net/http"
)

type LuaRunner struct {
	handle *lua.LState
}

func NewLuaRunner() *LuaRunner {
	l := &LuaRunner{
		handle: lua.NewState(),
	}

	l.handle.PreloadModule("http", luaHttp.NewHttpModule(&http.Client{}).Loader)
	l.handle.PreloadModule("json", luaJson.Loader)

	return l
}

func (l *LuaRunner) Close() {
	l.handle.Close()
}

func (l *LuaRunner) RunScript(cont string) error {
	if err := l.handle.DoString(cont); err != nil {
		return err
	}

	return nil
}
