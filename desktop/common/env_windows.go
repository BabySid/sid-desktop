package common

import "path/filepath"

var (
	luaRunner = "lua_runner.exe"
)

func GetLuaRunner() string {
	return filepath.Join(GetBinPath(), luaRunner)
}
