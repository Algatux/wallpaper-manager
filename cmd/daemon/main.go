package main

import (
	"flag"

	"github.com/Algatux/wallpaper-manager/internal"
)

var version = "dev"

func main() {

	internal.InitHelp(version)
	internal.InitInputs()

	flag.Parse()

}
