package main

import (
	"fmt"
	"github.com/jchv/go-webview-selector"
	"github.com/jessevdk/go-flags"
	"os"
)

type Option struct {
	Url      string `short:"u" long:"url" description:"Url to open by default" default:"https://www.baidu.com"`
	DeskPort uint16 `short:"p" long:"desk_port" description:"GRPC Port of sid-desktop listening" default:"0"`
	ReqID    int    `short:"i" long:"req_id" description:"Request ID for communication with sid-desktop" default:"0"`
}

const initJS = `window.onload=function(){
    getWebInfo(document.title, window.location.href)
}`

func main() {
	var opt Option
	parser := flags.NewParser(&opt, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(0)
	}

	if opt.DeskPort == 0 {
		// local mode
	}

	w := webview.New(false)
	if w == nil {
		fmt.Println("webview error")
		panic(false)
	}
	defer w.Destroy()

	w.SetTitle("WebView")
	w.SetSize(800, 600, webview.HintNone)
	_ = w.Bind("getWebInfo", getWebInfo)

	w.Navigate(opt.Url)

	w.Init(initJS)

	w.Run()
}

func getWebInfo(title, uri string) {
	fmt.Println("getWebInfo", title, uri)
}
