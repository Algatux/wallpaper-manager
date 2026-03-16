package setup

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/Algatux/wallpaper-manager/internal/log"
)

type MonitorList struct {
	Monitors []MonitorConfig `json:"monitors"`
}
type MonitorConfig struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	Flags string `json:"flags"`
}

type Inputs struct {
	Interval   int
	LockFile   string
	Debug      bool
	ConfigFile string
}

func (i *Inputs) GetConfiguration() (*MonitorList, error) {

	log.LogInfo(fmt.Sprintf("Loading configuration file `%s`", i.ConfigFile))

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

	logConfiguration(list)

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

func logConfiguration(list MonitorList) {
	prettyJSON, err := json.MarshalIndent(list, "", "  ")

	if err != nil {
		log.LogError(fmt.Sprintf("Could not pretty print config: %v", err))
		log.LogDebug(fmt.Sprintf("Configuration (raw): %+v", list))
		return
	}
	log.LogDebug(fmt.Sprintf("Loaded Configuration:\n%s", string(prettyJSON)))
}
