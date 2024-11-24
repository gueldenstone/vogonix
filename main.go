package main

import (
	"changeme/pkg/jira"
	"embed"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	jira, err := jira.NewJiraInstance(os.Getenv("URL"), os.Getenv("MAIL"), os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "vogonix",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        jira.Startup,
		LogLevel:         logger.DEBUG,
		Bind: []interface{}{
			jira,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
