package common

import (
	"bytes"
	"fmt"
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
	luaJson.Preload(l.handle)
	return l
}

func (l *LuaRunner) Close() {
	l.handle.Close()
}

func (l *LuaRunner) RunScript(cont string) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	if err := l.handle.DoString(`

        local http = require("http")
		local json = require("json")

        response, error_message = http.request("POST", "http://", {
            timeout="30s",
            headers={
                Accept="*/*"
            },
			body=[[{
				"sqlquery": "select * from mytable"}
			]]
        })
		
		print(response.status_code)
		print(response.body)
		print(error_message)
		
		body = json.decode(response.body)
		print(body.code)
		print(body.msg)

    `); err != nil {
		panic(err)
	}

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	fmt.Println("aha: ", out)
}
