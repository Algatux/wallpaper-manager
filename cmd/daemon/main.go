package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Algatux/wallpaper-manager/internal/filesystem"
	"github.com/Algatux/wallpaper-manager/internal/log"
	"github.com/Algatux/wallpaper-manager/internal/setup"
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

	monitorWallpapers := map[string][]string{}

	for _, monitor := range config.Monitors {
		monitorWallpapers[monitor.Name] = filesystem.ListSupportedImages(monitor.Path)
		log.LogInfo(fmt.Sprintf("Found %d wallpapers for device `%s`", len(monitorWallpapers[monitor.Name]), monitor.Name))
		log.LogDebug(fmt.Sprintf("[%s]`", strings.Join(monitorWallpapers[monitor.Name], ",")))
	}

}
