package common

import (
	"bytes"
	luaHttp "github.com/cjoudrey/gluahttp"
	lua "github.com/yuin/gopher-lua"
	"io"
	luaJson "layeh.com/gopher-json"
	"net/http"
	"os"
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

func (l *LuaRunner) RunScript(cont string) (<-chan string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()
	if err := l.handle.DoString(cont); err != nil {
		return outC, err
	}

	_ = w.Close()
	os.Stdout = old // restoring the real stdout

	return outC, nil
}
