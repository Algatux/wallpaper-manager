package main

import (
	"os"

	"github.com/Algatux/wallpaper-manager/internal/log"
	"github.com/Algatux/wallpaper-manager/internal/setup"
	"github.com/Algatux/wallpaper-manager/internal/wallpaper"
)

var version = "dev"

func main() {
	setup.InitHelp(version)
	inputs := setup.InitInputs()

	log.DebugMode = inputs.Debug

	config, err := inputs.GetConfiguration()
	if err != nil {
		os.Exit(1)
	}

	manager := wallpaper.NewManager(inputs.Interval, inputs.LockFile)
	if manager == nil {
		log.LogError("Wallpaper manager was unable to start.")
		os.Exit(1)
	}

	manager.WaitUntilReady()
	manager.CycleWallpapers(config)
}
