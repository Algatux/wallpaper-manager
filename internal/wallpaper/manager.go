package wallpaper

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Algatux/wallpaper-manager/internal/log"
	"github.com/Algatux/wallpaper-manager/internal/setup"
)

type WallpaperManagerType string

const (
	NotAvailable WallpaperManagerType = ""
	Awww         WallpaperManagerType = "awww"
	Swww         WallpaperManagerType = "swww"
)

type Manager struct {
	executable WallpaperManagerType
	interval   int
	lockFile   string
}

func (m *Manager) Init() *Manager {

	m.executable = NotAvailable
	if checkExecutableExistance(string(Awww)) {
		m.executable = Awww
	} else if checkExecutableExistance(string(Swww)) {
		m.executable = Swww
	}

	if m.executable == NotAvailable {
		log.LogError("either 'awww' nor 'swww' were found in the system.")
		return nil
	}

	return m
}

func (m *Manager) WaitUntilReady() {

	for exec.Command(string(m.executable), "query").Run() != nil {
		log.LogError(fmt.Sprintf("waiting for %s-daemon ...", m.executable))
		time.Sleep(5 * time.Second)
	}

	log.LogInfo("awww-daemon is running! Proceeding...")
}

func (m *Manager) CycleWallpapers(config *setup.MonitorList) {
	ticker := time.NewTicker(time.Duration(m.interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.changeWallpapers(config)
	}
}

func (m *Manager) changeWallpapers(config *setup.MonitorList) {
	if _, err := os.Stat(m.lockFile); err == nil {
		log.LogInfo("wallpaper switch is paused")
		return
	}

	for _, monitor := range config.Monitors {

		if len(monitor.Wallpapers) < 1 {
			continue
		}

		wallpaper := monitor.Wallpapers[rand.Intn(len(monitor.Wallpapers))]

		args := []string{"img", wallpaper, "-o", monitor.Name}
		if monitor.Flags != "" {
			args = append(args, strings.Split(monitor.Flags, " ")...)
		}

		log.LogInfo(fmt.Sprintf("setting wallpaper %s on monitor %s", wallpaper, monitor.Name))
		err := exec.Command(string(m.executable), args...).Run()
		if err != nil {
			log.LogError(fmt.Sprintf("Error switching wallpaper: %s", err.Error()))
		}

	}
}

func NewManager(interval int, file string) *Manager {
	return (&Manager{
		interval: interval,
		lockFile: file,
	}).Init()
}

func checkExecutableExistance(executable string) bool {
	path, err := exec.LookPath(executable)
	if err != nil {
		log.LogDebug(fmt.Sprintf("executable `%s` not found %s", executable, err.Error()))
		return false
	}

	log.LogDebug(fmt.Sprintf("found executable `%s` in path %s", executable, path))
	return true
}
