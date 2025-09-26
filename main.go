package main

import (
	"embed"
	"log"

	"github.com/MisakaTAT/GTerm/backend/cmd"
	"github.com/MisakaTAT/GTerm/backend/consts"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := cmd.NewApp()

	wailsApp := application.New(application.Options{
		Name:        "GTerm",
		Description: "A terminal tool developed with Wails + Vue3 + TypeScript.",
		Services:    app.Services(),
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            consts.ApplicationName,
		Width:            1200,
		Height:           800,
		MinWidth:         1024,
		MinHeight:        768,
		Frameless:        app.PreferencesSrv.IsDarwin(),
		BackgroundColour: application.NewRGBA(27, 38, 54, 0),
		Mac: application.MacWindow{
			// InvisibleTitleBarHeight: 10,
			// Backdrop: application.MacBackdropTranslucent,
			TitleBar: application.MacTitleBarHiddenInset,
		},
		Windows: application.WindowsWindow{},
		URL:     "/",
	})

	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}

}
