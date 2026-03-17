package setup

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Algatux/wallpaper-manager/internal/filesystem"
	"github.com/Algatux/wallpaper-manager/internal/log"
)

type MonitorList struct {
	Monitors []MonitorConfig `json:"monitors"`
}
type MonitorConfig struct {
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	Flags      string   `json:"flags"`
	Wallpapers []string `json:"wallpapers"`
}

type Inputs struct {
	Interval   int
	LockFile   string
	Debug      bool
	ConfigFile string
}

func (i *Inputs) GetConfiguration() (*MonitorList, error) {

	log.LogInfo(fmt.Sprintf("loading configuration file `%s`", i.ConfigFile))

	data, err := os.ReadFile(i.ConfigFile)
	if err != nil {
		log.LogError(fmt.Sprintf("Error reading configuration file: `%s` - %s", i.ConfigFile, err.Error()))
		return nil, err
	}

	list := MonitorList{}
	err = json.Unmarshal(data, &list)
	if err != nil {
		log.LogError(fmt.Sprintf("Error decoding configuration: %s", err.Error()))
		return nil, err
	}

	buildMonitorwallpaperList(&list)
	logConfiguration(&list)

	return &list, nil

}

func InitInputs() Inputs {
	inputs := Inputs{}

	flag.IntVar(&inputs.Interval, "i", 3600, "Interval in seconds between wallpaper changes")
	flag.StringVar(&inputs.LockFile, "l", "/tmp/wallpaper_cycle.lock", "Path to the wallpaper cycle lock file")
	flag.BoolVar(&inputs.Debug, "d", false, "Set to enable debug mode and print verbose messages")

	flag.Parse()

	inputs.ConfigFile = flag.Arg(0)

	return inputs
}

func InitHelp(version string) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "Wallpaper Manager Daemon Tool %s\n", version)
		fmt.Fprintf(os.Stdout, "Usage: %s [options] <config.json>\n\n", os.Args[0])
		fmt.Fprintln(os.Stdout, "Options:")
		flag.PrintDefaults()
	}
}

func logConfiguration(list *MonitorList) {
	prettyJSON, err := json.MarshalIndent(list, "", "  ")

	if err != nil {
		log.LogError(fmt.Sprintf("could not pretty print config: %s", err.Error()))
		log.LogDebug(fmt.Sprintf("configuration (raw): %+v", list))
		return
	}
	log.LogDebug(fmt.Sprintf("loaded configuration:\n%s", string(prettyJSON)))
}

func buildMonitorwallpaperList(config *MonitorList) {
	for i, monitor := range config.Monitors {
		monitor.Wallpapers = filesystem.ListSupportedImages(monitor.Path)
		log.LogInfo(fmt.Sprintf("Found %d wallpapers for device `%s`", len(monitor.Wallpapers), monitor.Name))
		log.LogDebug(fmt.Sprintf("[%s]", strings.Join(monitor.Wallpapers, ",")))
		config.Monitors[i] = monitor
	}
}
