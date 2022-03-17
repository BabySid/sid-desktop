package main

import (
	"flag"
	"image"
	"os"

	"github.com/tc-hib/winres"
)

var (
	output string
	src    string
)

func init() {
	flag.StringVar(&output, "o", "", "output file")
	flag.StringVar(&src, "s", "", "src png image")
}

// https://www.golangtc.com/t/50b58150320b52067e00000f
func main() {
	flag.Parse()
	if output == "" || src == "" {
		flag.PrintDefaults()
		return
	}

	rs := winres.ResourceSet{}

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

	rs.SetIcon(winres.ID(1), icon)

	out, _ := os.Create(output)
	defer out.Close()
	err = rs.WriteObject(out, winres.ArchAMD64)
	if err != nil {
		panic(err)
	}
}
