package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Server struct {
	Endpoint string
	Key      string
}

func JudgeOS() {
	switch runtime.GOOS {
	case "windows":
		progData := os.Getenv("PROGRAMDATA")
		if progData == "" {
			progData = "."
		}
		viper.AddConfigPath(filepath.Join(progData, "magenet"))
		fmt.Println("Detected OS: windows, config path:", filepath.Join(progData, "magenet"))
	case "darwin":
		viper.AddConfigPath("/usr/local/etc/magenet")
		viper.AddConfigPath("/etc/magenet")
		fmt.Println("Detected OS: darwin, config paths: /usr/local/etc/magenet, /etc/magenet")
	default: // linux, freebsd, etc.
		viper.AddConfigPath("/etc/magenet")
		viper.AddConfigPath(".")
		fmt.Println("Detected OS:", runtime.GOOS, "config paths: /etc/magenet, .")
	}
}

func init() {
	JudgeOS()

	viper.OnConfigChange(func(e fsnotify.Event) {

	})

	viper.WatchConfig()
}

func updateConfig() {

}
