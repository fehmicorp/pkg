package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/fehmicorp/go/pkg/v1/win/notify"
	"github.com/fehmicorp/go/pkg/v1/win/systray"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	Windows(runtime.GOARCH)
}

func Windows(Arch string) {
	if err := ensureWintunDLL(); err != nil {
		log.Fatalf("[CLIENT ERR] Failed to setup Wintun dependency: %v", err)
	}
	app := application.New(application.Options{
		Name:        Conf.AppName,
		Description: Conf.Description,
	})
	notify.RegisterAlert(Conf.AppName)
	Tooltip := fmt.Sprintf("%s\nVersion: %s\n%s", Conf.AppName, Conf.Version, Conf.Domain)
	customMenus := PrepareMenuItems()
	systray.NewTrayManager(app, TrayIcon, Tooltip, customMenus)
	if err := app.Run(); err != nil {
		log.Fatalf("❌ System tray service failed: %v", err)
	}
}
