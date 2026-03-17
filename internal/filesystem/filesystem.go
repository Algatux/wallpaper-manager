package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/Algatux/wallpaper-manager/internal/log"
)

func ListSupportedImages(dir string) []string {
	log.LogDebug(fmt.Sprintf("Building file list for directory `%s`", dir))
	supportedFiles := []string{
		".png",
		".jpg",
		".jpeg",
		".bmp",
		".pnm",
		".tga",
		".tiff",
		".webp",
		".gif",
	}

	list := []string{}

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.LogError(fmt.Sprintf("error reading path `%s`, %s", dir, err.Error()))
	}

	pathSeparator := ""
	if dir[len(dir)-1] != os.PathSeparator {
		pathSeparator = string(os.PathSeparator)
	}

	for _, entry := range entries {
		if !entry.IsDir() && slices.Contains(supportedFiles, filepath.Ext(entry.Name())) {
			list = append(list, dir+pathSeparator+entry.Name())
		}
	}

	return list
}
