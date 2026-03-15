package internal

import (
	"flag"
	"fmt"
	"os"
)

type Inputs struct {
	Interval int
	LockFile string
}

func InitInputs() Inputs {
	inputs := Inputs{}

	flag.IntVar(&inputs.Interval, "interval", 3600, "Interval in seconds between wallpaper changes")
	flag.StringVar(&inputs.LockFile, "lockfile", "/tmp/wallpaper_cycle.lock", "Path to the wallpaper cycle lock file")

	return inputs
}

func InitHelp(version string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Wallpaper Manager Daemon Tool %s\n", version)
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [arguments...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
}
