package main

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/gueldenstone/vogonix/pkg/config"
	"github.com/gueldenstone/vogonix/pkg/jira"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configDir := filepath.Join(homeDir, ".vogonix")
	os.Mkdir(configDir, 0755)
	configFile := filepath.Join(configDir, "config.yml")
	cfg, err := config.ReadConfig(configFile)
	if err != nil {
		panic(err)
	}

	jira, err := jira.NewJiraInstance(cfg.Url, cfg.User, cfg.Token, filepath.Join(configDir, "data.db"))
	if err != nil {
		panic(err)
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:     "vogonix",
		Width:     700,
		Height:    800,
		MinWidth:  500,
		MinHeight: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        jira.Startup,
		OnShutdown:       jira.Shutdown,
		LogLevel:         logger.DEBUG,
		Bind: []interface{}{
			jira,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
