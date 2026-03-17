package setup

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
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

	log.LogInfo("loading configuration file `" + i.ConfigFile + "`")

	data, err := os.ReadFile(i.ConfigFile)
	if err != nil {
		log.LogError("Error reading configuration file: `" + i.ConfigFile + "` - " + err.Error())
		return nil, err
	}

	list := MonitorList{}
	err = json.Unmarshal(data, &list)
	if err != nil {
		log.LogError("Error decoding configuration: " + err.Error())
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
		log.LogError("could not pretty print config: " + err.Error())
		log.LogDebug(fmt.Sprintf("configuration (raw): %+v", list))
		return
	}
	log.LogDebug("loaded configuration:\n" + string(prettyJSON))
}

func buildMonitorwallpaperList(config *MonitorList) {
	for i, monitor := range config.Monitors {
		monitor.Wallpapers = filesystem.ListSupportedImages(monitor.Path)
		log.LogInfo("Found " + strconv.Itoa(len(monitor.Wallpapers)) + " wallpapers for device `" + monitor.Name + "`")
		log.LogDebug("[" + strings.Join(monitor.Wallpapers, ",") + "]")
		config.Monitors[i] = monitor
	}
}
