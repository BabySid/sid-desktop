package main

import (
	"encoding/json"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"os"
	"sid-desktop/proto"
)

type Option struct {
	ScriptFile string `short:"f" long:"script" description:"script file to be run" required:"true"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	var opt Option
	parser := flags.NewParser(&opt, flags.Default)
	if _, err := parser.Parse(); err != nil {
		log.Println("flags parser failed.", err)
		return
	}

	defer os.Remove(opt.ScriptFile)

	data, err := ioutil.ReadFile(opt.ScriptFile)
	if err != nil {
		log.Printf("read file[%s] failed. %s\n", opt.ScriptFile, err)
		return
	}

	var script proto.ScriptRunnerRequest
	err = json.Unmarshal(data, &script)
	if err != nil {
		log.Printf("unmarshal file[%s] failed. %s\n", opt.ScriptFile, err)
		return
	}

	log.Printf("script [%s] begin run\n", script.Title)

	lua := NewLuaRunner()
	defer lua.Close()

	err = lua.RunScript(script.Content)
	if err != nil {
		log.Printf("script[%s] run failed. %s\n", script.Title, err)
		return
	}

	log.Printf("script [%s] run finished\n", script.Title)
}
