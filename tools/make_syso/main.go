package main

import (
	"flag"
	"image"
	"os"
	"strings"

	"github.com/tc-hib/winres"
)

var (
	output string
	src    string
)

func init() {
	flag.StringVar(&output, "o", "", "output file. ico or syso")
	flag.StringVar(&src, "s", "", "src png image")
}

// https://www.golangtc.com/t/50b58150320b52067e00000f
func main() {
	flag.Parse()
	if output == "" || src == "" {
		flag.PrintDefaults()
		return
	}
	if !strings.HasSuffix(output, ".syso") && !strings.HasSuffix(output, ".ico") {
		flag.PrintDefaults()
		return
	}

	file, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	icon, err := winres.NewIconFromImages([]image.Image{img})
	if err != nil {
		panic(err)
	}

	out, _ := os.Create(output)
	defer out.Close()

	if strings.HasSuffix(output, ".syso") {
		rs := winres.ResourceSet{}
		err = rs.SetIcon(winres.ID(1), icon)
		if err != nil {
			panic(err)
		}
		err = rs.WriteObject(out, winres.ArchAMD64)
		if err != nil {
			panic(err)
		}
	}
	if strings.HasSuffix(output, ".ico") {
		err = icon.SaveICO(out)
		if err != nil {
			panic(err)
		}
	}
}
